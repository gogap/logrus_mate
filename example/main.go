package main

import (
	"fmt"
	"os"

	"github.com/gogap/logrus_mate"
)

func main() {
	logrus_mate.Logger().Infoln("=== Using internal defualt logurs mate ===")
	logrus_mate.Logger().Debugln("Hello Default Logrus Mate")

	loggerConf := logrus_mate.LoggerConfig{
		Level: "info",
		Formatter: logrus_mate.FormatterConfig{
			Name: "json",
		},
	}

	if _, err := logrus_mate.NewLogger("jack", loggerConf); err != nil {
		logrus_mate.Logger().Error(err)
		return
	}

	logrus_mate.Logger().Warnln("*** Add Logger named jack, and it will use json format")

	logrus_mate.Logger("jack").Debugln("not print")
	logrus_mate.Logger("jack").Infoln("Hello, I am A Logger from jack")

	fmt.Println("")
	os.Setenv("RUN_MODE", "production")
	logrus_mate.Logger().Infoln("=== Load logrus mate config from mate.conf ===")

	if mateConf, err := logrus_mate.LoadLogrusMateConfig("mate.conf"); err != nil {
		logrus_mate.Logger().Error(err)
		return
	} else {
		logrus_mate.Logger().Debugf("Run mode is %s", mateConf.RunEnv())

		if newMate, err := logrus_mate.NewLogrusMate(mateConf); err != nil {
			logrus_mate.Logger().Error(err)
			return
		} else {
			newMate.Logger("mike").Errorln("I am mike in new logrus mate")
		}
	}

}
