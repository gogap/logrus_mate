package logrus_mate

import (
	"github.com/sirupsen/logrus"
)

type NullFormatter struct {
}

func init() {
	RegisterFormatter("null", NewNullFormatter)
}

func NewNullFormatter(options Options) (formatter logrus.Formatter, err error) {
	formatter = &NullFormatter{}
	return
}

func (NullFormatter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte{}, nil
}
