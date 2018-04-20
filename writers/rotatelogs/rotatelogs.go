package redisio

import (
	"errors"
	"io"
	"time"

	"github.com/gogap/config"
	"github.com/gogap/logrus_mate"
	"github.com/lestrrat-go/file-rotatelogs"
)

type RedisIOConfig struct {
	Network  string
	Address  string
	Password string
	Db       int64
	ListName string
}

func init() {
	logrus_mate.RegisterWriter("rotatelogs", NewRotateLogs)
}

func NewRotateLogs(conf config.Configuration) (writer io.Writer, err error) {

	path := conf.GetString("path")

	if len(path) == 0 {
		err = errors.New("config of path is empty, e.g.: /path/to/access_log.%Y%m%d%H%M")
		return
	}

	clock := toClock(conf.GetString("clock", "Local"))
	strLoc := conf.GetString("location", time.Now().Local().Location().String())

	loc, err := time.LoadLocation(strLoc)
	if err != nil {
		return
	}

	linkName := conf.GetString("link-name")
	rotationTime := conf.GetTimeDuration("rotation-time", time.Hour*24)
	maxAge := conf.GetTimeDuration("max-age", time.Hour*24*7)

	log, err := rotatelogs.New(path,
		rotatelogs.WithClock(clock),
		rotatelogs.WithLocation(loc),
		rotatelogs.WithLinkName(linkName),
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithMaxAge(maxAge),
	)

	if err != nil {
		return
	}

	writer = log

	return
}

func toClock(name string) rotatelogs.Clock {
	if name == "UTC" || name == "utc" {
		return rotatelogs.UTC
	}

	return rotatelogs.Local
}
