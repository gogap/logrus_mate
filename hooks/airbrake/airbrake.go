package airbrake

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/gemnasium/logrus-airbrake-hook.v2"

	"github.com/gogap/logrus_mate"
)

type AirbrakeHookConfig struct {
	ProjectId int64
	APIKey    string
	Env       string
}

func init() {
	logrus_mate.RegisterHook("airbrake", NewAirbrakeHook)
}

func NewAirbrakeHook(options *logrus_mate.Options) (hook logrus.Hook, err error) {

	conf := AirbrakeHookConfig{}
	if options != nil {
		conf.ProjectId = options.GetInt64("project-id")
		conf.APIKey = options.GetString("api-key")
		conf.Env = options.GetString("env")
	}

	hook = airbrake.NewHook(
		conf.ProjectId,
		conf.APIKey,
		conf.Env,
	)

	return
}
