package logrus_mate

import (
	"github.com/go-akka/configuration"
)

type Option func(*Config)

type Config struct {
	ConfigFile   string
	ConfigString string
	config       *configuration.Config
}

func newConfig(opts ...Option) *Config {
	conf := &Config{}
	conf.init(opts...)
	return conf
}

func (p *Config) init(opts ...Option) {
	for i := 0; i < len(opts); i++ {
		opts[i](p)
	}

	var confString, confFile *configuration.Config

	if len(p.ConfigFile) > 0 {
		confFile = configuration.LoadConfig(p.ConfigFile)
	}

	if len(p.ConfigString) > 0 {
		confString = configuration.ParseString(p.ConfigString)
	}

	if confFile == nil && confString == nil {
		p.config = &configuration.Config{}
		return
	}

	if confString != nil && confFile != nil {
		confString.WithFallback(confFile)
		p.config = confString
		return
	}

	if confString != nil {
		p.config = confString
	} else {
		p.config = confFile
	}

}

func ConfigFile(fn string) Option {
	return func(o *Config) {
		o.ConfigFile = fn
	}
}

func ConfigString(str string) Option {
	return func(o *Config) {
		o.ConfigString = str
	}
}
