package logrus_mate

import (
	"github.com/Sirupsen/logrus"
)

func init() {
	RegisterFormatter("json", NewJSONFormatter)
}

func NewJSONFormatter(options Options) (formatter logrus.Formatter, err error) {
	formatter = &logrus.JSONFormatter{}
	return
}
