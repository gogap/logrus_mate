package slack

import (
	"github.com/sirupsen/logrus"
	"github.com/johntdyer/slackrus"

	"github.com/gogap/logrus_mate"
)

type SlackHookConfig struct {
	URL      string   `json:"url"`
	Levels   []string `json:"levels"`
	Channel  string   `json:"channel"`
	Emoji    string   `json:"emoji"`
	Username string   `json:"username"`
}

func init() {
	logrus_mate.RegisterHook("slack", NewSlackHook)
}

func NewSlackHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := SlackHookConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
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
