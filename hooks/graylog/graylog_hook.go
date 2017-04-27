package graylog

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/gemnasium/logrus-graylog-hook.v2"

	"github.com/gogap/logrus_mate"
)

type GraylogHookConfig struct {
	Address string
	Extra   map[string]interface{}
}

func init() {
	logrus_mate.RegisterHook("graylog", NewGraylogHook)
}

func NewGraylogHook(options *logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := GraylogHookConfig{}

	if options != nil {
		conf.Address = options.GetString("address")

		extraConf := options.GetConfig("extra")
		keys := extraConf.Root().GetObject().GetKeys()
		extra := make(map[string]interface{}, len(keys))

		for i := 0; i < len(keys); i++ {
			extra[keys[i]] = extraConf.GetString(keys[i])
		}

		conf.Extra = extra
	}

	hook = graylog.NewAsyncGraylogHook(
		conf.Address,
		conf.Extra)

	return
}
