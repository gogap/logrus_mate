package lfshook

import (
	"fmt"

	"github.com/gogap/config"
	"github.com/gogap/logrus_mate"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus_mate.RegisterHook("lfshook", NewLFShook)
}

func NewLFShook(conf config.Configuration) (hook logrus.Hook, err error) {

	pathMapConf := conf.GetConfig("path-map")

	pathMap := lfshook.PathMap{}

	for _, key := range pathMapConf.Keys() {
		var lvl logrus.Level
		lvl, err = logrus.ParseLevel(key)
		if err != nil {
			return
		}

		filename := pathMapConf.GetString(key)

		if len(filename) == 0 {
			err = fmt.Errorf("log level of %s did not assign filename", key)
			return
		}

		pathMap[lvl] = filename
	}

	hook = lfshook.NewHook(pathMap, nil)

	return
}
