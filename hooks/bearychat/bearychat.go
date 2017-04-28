package bearychat

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gogap/bearychat/incoming"

	"github.com/gogap/logrus_mate"
)

var allLevels = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

type BearyChatHookConfig struct {
	RobotId  string
	Token    string
	Levels   []string
	Channel  string
	User     string
	Markdown bool
	Async    bool
}

func init() {
	logrus_mate.RegisterHook("bearychat", NewBearyChatHook)
}

func NewBearyChatHook(options *logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := BearyChatHookConfig{}

	if options != nil {
		conf.RobotId = options.GetString("robot-id")
		conf.Token = options.GetString("token")
		conf.Levels = options.GetStringList("levels")
		conf.Channel = options.GetString("channel")
		conf.User = options.GetString("user")
		conf.Markdown = options.GetBoolean("markdown")
		conf.Async = options.GetBoolean("async", true)
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

	hook = &BearyChatHook{
		RobotId:        conf.RobotId,
		Token:          conf.Token,
		AcceptedLevels: levels,
		Channel:        conf.Channel,
		User:           conf.User,
		Markdown:       conf.Markdown,
		Async:          conf.Async,
		cli:            incoming.NewClient(),
	}

	return
}

type BearyChatHook struct {
	AcceptedLevels []logrus.Level
	RobotId        string
	Token          string
	Channel        string
	User           string
	Markdown       bool
	Async          bool

	cli *incoming.Client
}

// Levels sets which levels to sent to slack
func (p *BearyChatHook) Levels() []logrus.Level {
	if p.AcceptedLevels == nil {
		return allLevels
	}
	return p.AcceptedLevels
}

// Fire -  Sent event to slack
func (p *BearyChatHook) Fire(e *logrus.Entry) (err error) {
	color := ""
	switch e.Level {
	case logrus.DebugLevel:
		color = "#FDFEFE"
	case logrus.InfoLevel:
		color = "#5DADE2"
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		color = "#FF0000"
	default:
		color = "#FFFF00"
	}

	req := &incoming.Request{
		Text:     e.Message,
		Markdown: p.Markdown,
		Channel:  p.Channel,
		User:     p.User,
	}

	var attachs []incoming.Attachment

	if len(e.Data) > 0 {

		for k, v := range e.Data {
			attach := incoming.Attachment{}

			attach.Title = k
			attach.Text = fmt.Sprint(v)
			attach.Color = color

			attachs = append(attachs, attach)
		}
	}

	req.Attachments = attachs

	if p.Async {
		go p.cli.Send(p.RobotId, p.Token, req)
		return
	}

	_, err = p.cli.Send(p.RobotId, p.Token, req)
	return
}
