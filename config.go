package logrus_mate

import (
	"github.com/gogap/config"
)

type Option func(*Config)

type Config struct {
	configOpts []config.Option
}

func ConfigFile(fn string) Option {
	return func(o *Config) {
		o.configOpts = append(o.configOpts, config.ConfigFile(fn))
	}
}

func ConfigString(str string) Option {
	return func(o *Config) {
		o.configOpts = append(o.configOpts, config.ConfigString(str))
	}
}

func WithConfig(conf config.Configuration) Option {
	return func(o *Config) {
		o.configOpts = append(o.configOpts, config.WithConfig(conf))
	}
}

func ConfigProvider(provider config.ConfigurationProvider) Option {
	return func(o *Config) {
		o.configOpts = append(o.configOpts, config.ConfigProvider(provider))
	}
}
