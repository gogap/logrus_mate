package logrus_mate

import (
	"github.com/Sirupsen/logrus"
)

func init() {
	RegisterFormatter("text", NewTextFormatter)
}

func NewTextFormatter(options Options) (formatter logrus.Formatter, err error) {
	formatter = &logrus.TextFormatter{}
	return
}
