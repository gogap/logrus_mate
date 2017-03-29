package graylog

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/gemnasium/logrus-graylog-hook.v2"

	"github.com/gogap/logrus_mate"
)

type GraylogHookConfig struct {
	Address string                 `json:"address"`
	Extra   map[string]interface{} `json:"extra"`
}

func init() {
	logrus_mate.RegisterHook("graylog", NewGraylogHook)
}

func NewGraylogHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := GraylogHookConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	hook = graylog.NewAsyncGraylogHook(
		conf.Address,
		conf.Extra)

	return
}
