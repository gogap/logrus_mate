package logstash

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"

	"github.com/gogap/logrus_mate"
)

func init() {
	logrus_mate.RegisterFormatter("logstash", NewLogStashFormatter)
}

func NewLogStashFormatter(options logrus_mate.Options) (formatter logrus.Formatter, err error) {
	formatter = &logstash.LogstashFormatter{}
	return
}
