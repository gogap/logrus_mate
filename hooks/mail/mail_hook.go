package mail

import (
	"github.com/Sirupsen/logrus"
	"github.com/zbindenren/logrus_mail"

	"github.com/gogap/logrus_mate"
)

type MailHookConfig struct {
	AppName  string
	Host     string
	Port     int
	From     string
	To       string
	Username string
	Password string
}

func init() {
	logrus_mate.RegisterHook("mail", NewMailHook)
}

func NewMailHook(options *logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := MailHookConfig{}
	if options != nil {
		conf.AppName = options.GetString("app-name")
		conf.Host = options.GetString("host")
		conf.Port = int(options.GetInt32("port"))
		conf.From = options.GetString("from")
		conf.To = options.GetString("to")
		conf.Username = options.GetString("username")
		conf.Password = options.GetString("password")
	}

	hook, err = logrus_mail.NewMailAuthHook(
		conf.AppName,
		conf.Host,
		conf.Port,
		conf.From,
		conf.To,
		conf.Username,
		conf.Password)

	return
}
