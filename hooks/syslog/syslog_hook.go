package syslog

import (
	"log/syslog"

	"github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"

	"github.com/gogap/logrus_mate"
)

type SyslogHookConfig struct {
	Network  string
	Address  string
	Priority string
	Tag      string
}

func init() {
	logrus_mate.RegisterHook("syslog", NewSyslogHook)
}

func NewSyslogHook(options *logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := SyslogHookConfig{}

	if options != nil {
		conf.Network = options.GetString("network")
		conf.Address = options.GetString("address")
		conf.Priority = options.GetString("priority")
		conf.Tag = options.GetString("tag")
	}

	return logrus_syslog.NewSyslogHook(
		conf.Network,
		conf.Address,
		toPriority(conf.Priority),
		conf.Tag)
}

func toPriority(priority string) syslog.Priority {
	switch priority {
	case "LOG_EMERG":
		return syslog.LOG_EMERG
	case "LOG_ALERT":
		return syslog.LOG_ALERT
	case "LOG_CRIT":
		return syslog.LOG_CRIT
	case "LOG_ERR":
		return syslog.LOG_ERR
	case "LOG_WARNING":
		return syslog.LOG_WARNING
	case "LOG_NOTICE":
		return syslog.LOG_NOTICE
	case "LOG_INFO":
		return syslog.LOG_INFO
	case "LOG_DEBUG":
		return syslog.LOG_DEBUG
	}
	return syslog.LOG_DEBUG
}
