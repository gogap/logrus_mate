package logrus_mate

import (
	"github.com/Sirupsen/logrus"
)

type JSONFormatterConfig struct {
	TimestampFormat string `json:"timestamp_format"`
}

func init() {
	RegisterFormatter("json", NewJSONFormatter)
}

func NewJSONFormatter(options *Options) (formatter logrus.Formatter, err error) {
	var format string
	if options != nil {
		format = options.GetString("timestamp_format")
	}
	formatter = &logrus.JSONFormatter{TimestampFormat: format}
	return
}
