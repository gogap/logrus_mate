package main

import (
	"fmt"
	"github.com/gogap/errors"
	"os"
	"time"

	"github.com/gogap/logrus_mate"

	_ "github.com/gogap/logrus_mate/hooks/bugsnag"
	_ "github.com/gogap/logrus_mate/hooks/mail"
	_ "github.com/gogap/logrus_mate/hooks/slack"
	_ "github.com/gogap/logrus_mate/hooks/syslog"

	_ "github.com/gogap/logrus_mate/writers/redisio"
)

func main() {
	if _, err := os.Stat("mate.conf"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Please copy mate.conf.example to mate.conf, and configure this file.")
			return
		}
		fmt.Println(err)
		return
	}

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

			ErrTest := errors.TN("GOGAP", 1000, "hello {{.param}}")
			ErrTest2 := errors.TN("GOGAP", 1002, "hello")
			e := ErrTest.New(errors.Params{"param": "world"}).Append("append error").Append(ErrTest2).WithContext("key", "Value")
			newMate.Logger("mike").WithError(e).Error(e)

			// This sleep is for output of redisio to write data to redis
			time.Sleep(time.Second)
		}
	}

}
