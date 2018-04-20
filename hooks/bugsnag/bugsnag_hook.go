package bugsnag

import (
	"github.com/Shopify/logrus-bugsnag"
	"github.com/bugsnag/bugsnag-go"
	"github.com/sirupsen/logrus"

	"github.com/gogap/config"
	"github.com/gogap/logrus_mate"
)

func init() {
	logrus_mate.RegisterHook("bugsnag", NewBugsnagHook)
}

func NewBugsnagHook(conf config.Configuration) (hook logrus.Hook, err error) {

	if conf != nil {
		bugsnag.Configure(
			bugsnag.Configuration{
				Endpoint:     conf.GetString("endpoint"),
				ReleaseStage: conf.GetString("release-stage"),
				APIKey:       conf.GetString("api-key"),
				Synchronous:  conf.GetBoolean("synchronous"),
			})
	}

	hook, err = logrus_bugsnag.NewBugsnagHook()
	return
}
