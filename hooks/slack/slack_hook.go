package slack

import (
	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"

	"github.com/gogap/logrus_mate"
)

type SlackHookConfig struct {
	URL      string
	Levels   []string
	Channel  string
	Emoji    string
	Username string
}

func init() {
	logrus_mate.RegisterHook("slack", NewSlackHook)
}

func NewSlackHook(options *logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := SlackHookConfig{}

	if options != nil {
		conf.URL = options.GetString("url")
		conf.Levels = options.GetStringList("levels")
		conf.Channel = options.GetString("channel")
		conf.Emoji = options.GetString("emoji")
		conf.Username = options.GetString("username")
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

	hook = &slackrus.SlackrusHook{
		HookURL:        conf.URL,
		AcceptedLevels: levels,
		Channel:        conf.Channel,
		IconEmoji:      conf.Emoji,
		Username:       conf.Username,
	}

	return
}
