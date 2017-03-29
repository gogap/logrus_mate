package logstash

import (
	"github.com/sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/gogap/logrus_mate"
)

type LogstashHookConfig struct {
	AppName          string        `json:"app_name"`
	Protocol         string        `json:"protocol"`
	Address          string        `json:"address"`
	AlwaysSentFields logrus.Fields `json:"always_sent_fields"`
	Prefix           string        `json:"prefix"`
}

func init() {
	logrus_mate.RegisterHook("logstash", NewLogstashHook)
}

func NewLogstashHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := LogstashHookConfig{}
	if err = options.ToObject(&conf); err != nil {
		return
	}

	hook, err = logrus_logstash.NewHookWithFieldsAndPrefix(
		conf.Protocol,
		conf.Address,
		conf.AppName,
		conf.AlwaysSentFields,
		conf.Prefix,
	)

	return
}
