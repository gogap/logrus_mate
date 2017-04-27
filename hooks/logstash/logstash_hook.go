package logstash

import (
	"github.com/Sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/gogap/logrus_mate"
)

type LogstashHookConfig struct {
	AppName          string
	Protocol         string
	Address          string
	AlwaysSentFields logrus.Fields
	Prefix           string
}

func init() {
	logrus_mate.RegisterHook("logstash", NewLogstashHook)
}

func NewLogstashHook(options *logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := LogstashHookConfig{}

	if options != nil {
		conf.AppName = options.GetString("app-name")
		conf.Protocol = options.GetString("protocol")
		conf.Address = options.GetString("address")
		conf.Prefix = options.GetString("prefix")

		alwaysSentFieldsConf := options.GetConfig("always-sent-fields")
		keys := alwaysSentFieldsConf.Root().GetObject().GetKeys()
		fields := make(logrus.Fields, len(keys))

		for i := 0; i < len(keys); i++ {
			fields[keys[i]] = alwaysSentFieldsConf.GetString(keys[i])
		}

		conf.AlwaysSentFields = fields
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
