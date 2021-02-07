package logrus_mate

import (
	"github.com/gogap/config"
	"github.com/sirupsen/logrus"
)

type JSONFormatterConfig struct {
	TimestampFormat  string `json:"timestamp_format"`
	EnableHTMLEscape bool   `json:"enable_html_escape"`
	PrettyPrint      bool   `json:"pretty_print"`
	DisableTimestamp bool   `json:"enable_timestamp"`
}

func init() {
	RegisterFormatter("json", NewJSONFormatter)
}

func NewJSONFormatter(config config.Configuration) (logrus.Formatter, error) {
	jsonFormatter := &logrus.JSONFormatter{}
	if config != nil {
		jsonFormatter.TimestampFormat = config.GetString("timestamp_format", "")
		jsonFormatter.DisableTimestamp = !config.GetBoolean("enable_timestamp", true)
		jsonFormatter.DisableHTMLEscape = !config.GetBoolean("enable_html_escape", false)
		jsonFormatter.PrettyPrint = config.GetBoolean("pretty_print", false)
	}
	return jsonFormatter, nil
}
