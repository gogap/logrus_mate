package logrus_mate

import (
	"github.com/gogap/config"
	"github.com/sirupsen/logrus"
)

type JSONFormatterConfig struct {
	TimestampFormat   string `json:"timestamp_format"`
	DisableHTMLEscape bool   `json:"disable_html_escape"`
	DisableTimestamp  bool   `json:"disable_timestamp"`
	PrettyPrint       bool   `json:"pretty_print"`
}

func init() {
	RegisterFormatter("json", NewJSONFormatter)
}

func NewJSONFormatter(config config.Configuration) (logrus.Formatter, error) {
	jsonFormatter := &logrus.JSONFormatter{}
	if config != nil {
		jsonFormatter.TimestampFormat = config.GetString("timestamp_format", "")
		jsonFormatter.DisableTimestamp = config.GetBoolean("disable_timestamp", false)
		jsonFormatter.DisableHTMLEscape = config.GetBoolean("disable_html_escape", true)
		jsonFormatter.PrettyPrint = config.GetBoolean("pretty_print", false)
	}
	return jsonFormatter, nil
}
