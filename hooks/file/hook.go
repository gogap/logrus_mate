package logrus_file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gogap/config"
	"github.com/gogap/logrus_mate"
	"github.com/gogap/logrus_mate/hooks/utils/caller"
)

type fileHookConfig struct {
	Filename   string `json:"filename"`
	MaxLines   int64  `json:"maxLines"`
	MaxSize    int64  `json:"maxsize"`
	Daily      bool   `json:"daily"`
	MaxDays    int64  `json:"maxDays"`
	Rotate     bool   `json:"rotate"`
	Perm       string `json:"perm"`
	RotatePerm string `json:"rotateperm"`
	Level      int32  `json:"level"`
}

func init() {
	logrus_mate.RegisterHook("file", NewFileHook)
}

func NewFileHook(config config.Configuration) (hook logrus.Hook, err error) {

	filename := config.GetString("filename", "logs/logrus.log")

	dir := filepath.Dir(filename)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	hookConf := fileHookConfig{
		Filename:   filename,
		Daily:      config.GetBoolean("daily", true),
		MaxDays:    config.GetInt64("max-days", 7),
		Rotate:     config.GetBoolean("rotate", true),
		MaxLines:   config.GetInt64("max-lines", 10000),
		MaxSize:    config.GetInt64("max-size", 1024),
		RotatePerm: config.GetString("rotate-perm", "0440"),
		Perm:       config.GetString("perm", "0660"),
		Level:      config.GetInt32("level"),
	}

	w := newFileWriter()

	confData, err := json.Marshal(hookConf)

	if err != nil {
		return
	}

	err = w.Init(string(confData))
	if err != nil {
		return
	}

	hook = &FileHook{W: w}

	return
}

type FileHook struct {
	W *fileLogWriter
}

func (p *FileHook) Fire(entry *logrus.Entry) (err error) {
	message, err := getMessage(entry)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	now := time.Now()

	switch entry.Level {
	case logrus.PanicLevel:
		fallthrough
	case logrus.FatalLevel:
		fallthrough
	case logrus.ErrorLevel:
		return p.W.WriteMsg(now, fmt.Sprintf("[ERROR] %s", message), LevelError)
	case logrus.WarnLevel:
		return p.W.WriteMsg(now, fmt.Sprintf("[WARN] %s", message), LevelWarn)
	case logrus.InfoLevel:
		return p.W.WriteMsg(now, fmt.Sprintf("[INFO] %s", message), LevelInfo)
	case logrus.DebugLevel:
		return p.W.WriteMsg(now, fmt.Sprintf("[DEBUG] %s", message), LevelDebug)
	default:
		return nil
	}

	return
}

func (p *FileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func getMessage(entry *logrus.Entry) (message string, err error) {
	message = message + fmt.Sprintf("%s\n", entry.Message)
	for k, v := range entry.Data {
		if !strings.HasPrefix(k, "err_") {
			message = message + fmt.Sprintf("%v:%v\n", k, v)
		}
	}
	if errCode, exist := entry.Data["err_code"]; exist {

		ns, _ := entry.Data["err_ns"]
		ctx, _ := entry.Data["err_ctx"]
		id, _ := entry.Data["err_id"]
		tSt, _ := entry.Data["err_stack"]
		st, _ := tSt.(string)
		st = strings.Replace(st, "\n", "\n\t\t", -1)

		buf := bytes.NewBuffer(nil)
		buf.WriteString(fmt.Sprintf("\tid:\n\t\t%s#%d:%s\n", ns, errCode, id))
		buf.WriteString(fmt.Sprintf("\tcontext:\n\t\t%s\n", ctx))
		buf.WriteString(fmt.Sprintf("\tstacktrace:\n\t\t%s", st))

		message = message + fmt.Sprintf("%v", buf.String())
	} else {
		file, lineNumber := caller.GetCallerIgnoringLogMulti(2)
		if file != "" {
			sep := fmt.Sprintf("%s/src/", os.Getenv("GOPATH"))
			fileName := strings.Split(file, sep)
			if len(fileName) >= 2 {
				file = fileName[1]
			}
		}
		message = message + fmt.Sprintf("%s:%d", file, lineNumber)
	}

	return
}
