package logrus_mate

import (
	"github.com/Sirupsen/logrus"
)

func init() {
	RegisterFormatter("text", NewTextFormatter)
}

func NewTextFormatter(options *Options) (formatter logrus.Formatter, err error) {

	f := &logrus.TextFormatter{}

	if options != nil {
		f.ForceColors = options.GetBoolean("force-colors")
		f.DisableColors = options.GetBoolean("disable-colors")
		f.DisableTimestamp = options.GetBoolean("disable-timestamp")
		f.FullTimestamp = options.GetBoolean("full-timestamp")
		f.TimestampFormat = options.GetString("timestamp-format")
		f.DisableSorting = options.GetBoolean("disable-sorting")
	}

	formatter = f

	return
}
