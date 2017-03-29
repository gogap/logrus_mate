package logrus_mate

import (
	"github.com/sirupsen/logrus"
)

type TextFormatterConfig struct {
	ForceColors      bool   `json:"force_colors"`
	DisableColors    bool   `json:"disable_colors"`
	DisableTimestamp bool   `json:"disable_timestamp"`
	FullTimestamp    bool   `json:"full_timestamp"`
	TimestampFormat  string `json:"timestamp_format"`
	DisableSorting   bool   `json:"disable_sorting"`
}

func init() {
	RegisterFormatter("text", NewTextFormatter)
}

func NewTextFormatter(options Options) (formatter logrus.Formatter, err error) {
	conf := TextFormatterConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	formatter = &logrus.TextFormatter{
		ForceColors:      conf.ForceColors,
		DisableColors:    conf.DisableColors,
		DisableTimestamp: conf.DisableTimestamp,
		FullTimestamp:    conf.FullTimestamp,
		TimestampFormat:  conf.TimestampFormat,
		DisableSorting:   conf.DisableSorting,
	}
	return
}
