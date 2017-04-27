package bugsnag

import (
	"github.com/Shopify/logrus-bugsnag"
	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"

	"github.com/gogap/logrus_mate"
)

func init() {
	logrus_mate.RegisterHook("bugsnag", NewBugsnagHook)
}

func NewBugsnagHook(options *logrus_mate.Options) (hook logrus.Hook, err error) {

	if options != nil {
		bugsnag.Configure(
			bugsnag.Configuration{
				Endpoint:     options.GetString("endpoint"),
				ReleaseStage: options.GetString("release-stage"),
				APIKey:       options.GetString("api-key"),
				Synchronous:  options.GetBoolean("synchronous"),
			})
	}

	hook, err = logrus_bugsnag.NewBugsnagHook()
	return
}
