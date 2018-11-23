package sls

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gogap/config"
	"github.com/gogap/logrus_mate"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var allLevels = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

type SLSHookConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string

	Project string
	Store   string

	StoreField   string
	TopicField   string
	SourceField  string
	TagsField    string
	ContextField string

	Levels []string
	Async  bool
}

func init() {
	logrus_mate.RegisterHook("sls", NewSLSHook)
}

func NewSLSHook(config config.Configuration) (hook logrus.Hook, err error) {

	conf := SLSHookConfig{}

	if config != nil {
		conf.Endpoint = config.GetString("endpoint")
		conf.AccessKeyID = config.GetString("access-key-id")
		conf.AccessKeySecret = config.GetString("access-key-secret")

		conf.Project = config.GetString("project")
		conf.Store = config.GetString("store")

		conf.StoreField = config.GetString("fields.store")
		conf.TopicField = config.GetString("fields.topic", "topic")
		conf.SourceField = config.GetString("fields.source", "source")
		conf.ContextField = config.GetString("fields.context")
		conf.TagsField = config.GetString("fields.tags")

		conf.Levels = config.GetStringList("levels")

		conf.Async = config.GetBoolean("async")
	}

	levels := []logrus.Level{}

	if conf.Levels != nil {
		for _, level := range conf.Levels {
			if lv, e := logrus.ParseLevel(level); e != nil {
				err = e
				return
			} else {
				levels = append(levels, lv)
			}
		}
	}

	if len(levels) == 0 && conf.Levels != nil {
		levels = append(levels, logrus.ErrorLevel, logrus.PanicLevel, logrus.FatalLevel)
	}

	slsHook := &SLSHook{
		AcceptedLevels: levels,
		Config:         conf,
	}

	proj, err := sls.NewLogProject(conf.Project, conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)

	if err != nil {
		return
	}

	slsHook.Config = conf
	slsHook.project = proj

	hook = slsHook

	return
}

type SLSHook struct {
	AcceptedLevels []logrus.Level
	Config         SLSHookConfig

	mapStores sync.Map
	project   *sls.LogProject
}

// Levels sets which levels to sent to sls
func (p *SLSHook) Levels() []logrus.Level {
	if p.AcceptedLevels == nil {
		return allLevels
	}
	return p.AcceptedLevels
}

// Fire -  Sent event to slack
func (p *SLSHook) Fire(entry *logrus.Entry) (err error) {

	storeName := p.Config.Store

	if len(p.Config.StoreField) > 0 {
		if specStoreName, ok := entry.Data[p.Config.StoreField].(string); ok {
			storeName = specStoreName
		}
	}

	store, storeExist := p.mapStores.Load(storeName)
	if !storeExist {

		newStore, newStoreErr := sls.NewLogStore(storeName, p.project)
		if newStoreErr != nil {
			return newStoreErr
		}
		store, _ = p.mapStores.LoadOrStore(storeName, newStore)
	}

	delete(entry.Data, p.Config.Store)

	logStore, ok := store.(*sls.LogStore)
	if !ok {
		err = fmt.Errorf("could not convert %s into *sls.LogStore in log store maps", storeName)
		return
	}

	slsLog := &sls.Log{Time: proto.Uint32(uint32(entry.Time.Unix()))}

	var strCtx string

	if len(p.Config.ContextField) > 0 {
		if ctxData, exist := entry.Data[p.Config.ContextField]; exist {

			if ctx, ok := ctxData.(string); ok {
				strCtx = ctx
			} else {
				jsonData, errJson := json.Marshal(ctxData)
				if errJson != nil {
					return errJson
				}

				strCtx = string(jsonData)
			}

			delete(entry.Data, p.Config.ContextField)

			slsLog.Contents = append(slsLog.Contents,
				&sls.LogContent{
					Key:   proto.String(p.Config.ContextField),
					Value: proto.String(strCtx),
				})

		}
	}

	topic := proto.String(fieldToString(entry, p.Config.TopicField))
	source := proto.String(fieldToString(entry, p.Config.SourceField))

	delete(entry.Data, p.Config.TopicField)
	delete(entry.Data, p.Config.SourceField)

	if len(entry.Data) > 0 {
		for k, v := range entry.Data {

			content := &sls.LogContent{
				Key:   proto.String(k),
				Value: proto.String(fmt.Sprintf("%v", v)),
			}

			slsLog.Contents = append(slsLog.Contents, content)
		}
	}

	slsLog.Contents = append(slsLog.Contents,
		&sls.LogContent{
			Key:   proto.String("message"),
			Value: proto.String(entry.Message),
		})

	slsLog.Contents = append(slsLog.Contents,
		&sls.LogContent{
			Key:   proto.String("log_level"),
			Value: proto.String(entry.Level.String()),
		})

	var slsTags []*sls.LogTag

	if len(p.Config.TagsField) > 0 {
		logrusTags, ok := entry.Data[p.Config.TagsField].(map[string]string)

		if ok {
			for k, v := range logrusTags {
				slsTags = append(slsTags,
					&sls.LogTag{
						Key:   proto.String(k),
						Value: proto.String(v),
					})
			}
		}
	}

	logGroup := &sls.LogGroup{
		Topic:   topic,
		Source:  source,
		Logs:    []*sls.Log{slsLog},
		LogTags: slsTags,
	}

	if p.Config.Async {
		go logStore.PutLogs(logGroup)
		return
	}

	err = logStore.PutLogs(logGroup)

	return
}

func fieldToString(entry *logrus.Entry, key string) string {

	fieldValue, exist := entry.Data[key]

	if !exist {
		return ""
	}

	return fmt.Sprintf("%v", fieldValue)
}
