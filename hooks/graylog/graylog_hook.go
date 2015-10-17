package graylog

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/gemnasium/logrus-graylog-hook.v1"

	"github.com/gogap/logrus_mate"
)

type GraylogHookConfig struct {
	Address  string                 `json:"address"`
	Facility string                 `json:"facility"`
	Extra    map[string]interface{} `json:"extra"`
}

func init() {
	logrus_mate.RegisterHook("graylog", NewGraylogHook)
}

func NewGraylogHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := GraylogHookConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	hook = graylog.NewGraylogHook(
		conf.Address,
		conf.Facility,
		conf.Extra)

	return
}
