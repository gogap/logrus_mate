package logrus_mate

import (
	"github.com/sirupsen/logrus"
)

type JSONFormatterConfig struct {
	TimestampFormat string `json:"timestamp_format"`
}

func init() {
	RegisterFormatter("json", NewJSONFormatter)
}

func NewJSONFormatter(options Options) (formatter logrus.Formatter, err error) {
	conf := JSONFormatterConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	formatter = &logrus.JSONFormatter{TimestampFormat: conf.TimestampFormat}
	return
}
