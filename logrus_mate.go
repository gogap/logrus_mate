package logrus_mate

import (
	"errors"
	"io"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/go-akka/configuration"
	"github.com/orcaman/concurrent-map"
)

type Options struct {
	*configuration.Config
}

var (
	ErrLoggerNotExist = errors.New("logger not exist")
)

type LogrusMate struct {
	loggersConf cmap.ConcurrentMap //map[string]*Config
	loggers     cmap.ConcurrentMap //map[string]*logrus.Logger
}

func NewLogger(opts ...Option) (logger *logrus.Logger, err error) {
	l := logrus.New()
	if err = Hijack(l, opts...); err != nil {
		return
	}

	return l, nil
}

func Hijack(logger *logrus.Logger, opts ...Option) (err error) {
	hijackConf := newConfig(opts...)
	conf := hijackConf.config

	return hijackByConfig(logger, conf)
}

func hijackByConfig(logger *logrus.Logger, conf *configuration.Config) (err error) {

	outConf := conf.GetConfig("out")
	formatterConf := conf.GetConfig("formatter")

	outName := "stdout"
	formatterName := "text"

	var outOptionsConf, formatterOptionsConf *configuration.Config

	if outConf != nil {
		outName = outConf.GetString("name", "stdout")
		outOptionsConf = outConf.GetConfig("options")
	}

	if formatterConf != nil {
		formatterName = formatterConf.GetString("name", "text")
		formatterOptionsConf = formatterConf.GetConfig("options")
	}

	var out io.Writer
	var outOptions *Options
	if outOptionsConf != nil {
		outOptions = &Options{outOptionsConf}
	}
	if out, err = NewWriter(outName, outOptions); err != nil {
		return
	}

	var formatter logrus.Formatter
	var formatterOptions *Options
	if formatterOptionsConf != nil {
		formatterOptions = &Options{formatterOptionsConf}
	}
	if formatter, err = NewFormatter(formatterName, formatterOptions); err != nil {
		return
	}

	var hooks []logrus.Hook

	confHooks := conf.GetConfig("hooks")

	if confHooks != nil {
		hookNames := confHooks.Root().GetObject().GetKeys()

		for i := 0; i < len(hookNames); i++ {
			var hook logrus.Hook
			if hook, err = NewHook(hookNames[i], confHooks.GetConfig(hookNames[i])); err != nil {
				return
			}
			hooks = append(hooks, hook)
		}
	}

	level := conf.GetString("level", "debug")

	var lvl = logrus.DebugLevel
	if lvl, err = logrus.ParseLevel(level); err != nil {
		return
	}

	l := logrus.New()

	l.Level = lvl
	l.Out = out
	l.Formatter = formatter
	for i := 0; i < len(hooks); i++ {
		l.Hooks.Add(hooks[i])
	}

	*logger = *l

	return
}

func NewLogrusMate(opts ...Option) (logrusMate *LogrusMate, err error) {
	mate := &LogrusMate{
		loggersConf: cmap.New(),
		loggers:     cmap.New(),
	}

	hijackConf := newConfig(opts...)
	conf := hijackConf.config

	if conf.Root() == nil {
		logrusMate = mate
		return
	}

	loggerNames := conf.Root().GetObject().GetKeys()

	for i := 0; i < len(loggerNames); i++ {
		mate.loggersConf.SetIfAbsent(loggerNames[i], conf.GetConfig(loggerNames[i]))
	}

	logrusMate = mate

	return
}

func (p *LogrusMate) Hijack(logger *logrus.Logger, loggerName string, opts ...Option) (err error) {
	confV, exist := p.loggersConf.Get(loggerName)
	if !exist {
		err = ErrLoggerNotExist
		return
	}

	conf := confV.(*configuration.Config)

	if len(opts) > 0 {
		conf2 := newConfig(opts...)
		conf2.config.WithFallback(conf)
		err = hijackByConfig(logger, conf2.config)
		return
	}

	err = hijackByConfig(logger, confV.(*configuration.Config))

	return
}

func (p *LogrusMate) Logger(loggerName ...string) (logger *logrus.Logger) {
	name := "default"

	if len(loggerName) > 0 {
		name = strings.TrimSpace(loggerName[0])
		if len(name) == 0 {
			name = "default"
		}
	}

	lv, exist := p.loggers.Get(name)

	if exist {
		return lv.(*logrus.Logger)
	}

	confV, exist := p.loggersConf.Get(name)
	if !exist {
		return nil
	}

	l := logrus.New()

	if err := hijackByConfig(l, confV.(*configuration.Config)); err != nil {
		return nil
	}

	p.loggers.SetIfAbsent(name, l)

	return l
}
