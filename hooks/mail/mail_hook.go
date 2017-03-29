package mail

import (
	"github.com/sirupsen/logrus"
	"github.com/zbindenren/logrus_mail"

	"github.com/gogap/logrus_mate"
)

type MailHookConfig struct {
	AppName  string `json:"app_name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	From     string `json:"from"`
	To       string `json:"to"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	logrus_mate.RegisterHook("mail", NewMailHook)
}

func NewMailHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := MailHookConfig{}
	if err = options.ToObject(&conf); err != nil {
		return
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
