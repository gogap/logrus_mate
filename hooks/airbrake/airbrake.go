package airbrake

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/gemnasium/logrus-airbrake-hook.v2"

	"github.com/gogap/logrus_mate"
)

type AirbrakeHookConfig struct {
	ProjectId int    `json:"project_id"`
	APIKey    string `json:"api_key"`
	Env       string `json:"env"`
}

func init() {
	logrus_mate.RegisterHook("airbrake", NewAirbrakeHook)
}

func NewAirbrakeHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := AirbrakeHookConfig{}
	if err = options.ToObject(&conf); err != nil {
		return
	}

	hook = airbrake.NewHook(
		int64(conf.ProjectId),
		conf.APIKey,
		conf.Env)

	return
}
