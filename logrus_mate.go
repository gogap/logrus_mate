package logrus_mate

import (
	"io"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	defaultMate         *LogrusMate
	defaultMateInitOnce sync.Once
)

type LogrusMate struct {
	loggersLock sync.Mutex
	initialOnce sync.Once

	loggers map[string]*logrus.Logger
}

func Logger(loggerName ...string) (logger *logrus.Logger) {
	if defaultMate == nil {
		defaultMateInitOnce.Do(func() {
			defaultMate = defaultLogrusMate()
		})
	}

	return defaultMate.Logger(loggerName...)
}

func NewLogger(name string, conf LoggerConfig) (logger *logrus.Logger, err error) {
	if defaultMate == nil {
		defaultMateInitOnce.Do(func() {
			defaultMate = defaultLogrusMate()
		})
	}

	return defaultMate.NewLogger(name, conf)
}

func (p LogrusMate) NewLogger(name string, conf LoggerConfig) (logger *logrus.Logger, err error) {
	tmpLogger := logrus.New()

	if conf.Out.Name == "" {
		conf.Out.Name = "stdout"
		conf.Out.Options = nil
	}

	var out io.Writer
	if out, err = NewWriter(conf.Out.Name, conf.Out.Options); err != nil {
		return
	}

	tmpLogger.Out = out

	if conf.Formatter.Name == "" {
		conf.Formatter.Name = "text"
		conf.Formatter.Options = nil
	}

	var formatter logrus.Formatter
	if formatter, err = NewFormatter(conf.Formatter.Name, conf.Formatter.Options); err != nil {
		return
	}

	tmpLogger.Formatter = formatter

	if conf.Hooks != nil {
		for _, hookConf := range conf.Hooks {
			var hook logrus.Hook
			if hook, err = NewHook(hookConf.Name, hookConf.Options); err != nil {
				return
			}
			tmpLogger.Hooks.Add(hook)
		}
	}

	var lvl = logrus.DebugLevel
	if lvl, err = logrus.ParseLevel(conf.Level); err != nil {
		return
	} else {
		tmpLogger.Level = lvl
	}

	logger = tmpLogger

	p.loggers[name] = logger

	return
}

func NewLogrusMate(mateConf LogrusMateConfig) (logrusMate *LogrusMate, err error) {
	mate := new(LogrusMate)

	if err = mate.initial(mateConf); err != nil {
		return
	}

	logrusMate = mate

	return
}

func (p *LogrusMate) initial(mateConf LogrusMateConfig) (err error) {
	p.loggersLock.Lock()
	defer p.loggersLock.Unlock()

	if err = mateConf.Validate(); err != nil {
		return
	}

	p.loggers = make(map[string]*logrus.Logger, len(mateConf.Loggers))

	p.initialOnce.Do(func() {

		runEnv := mateConf.RunEnv()

		for _, loggerConfs := range mateConf.Loggers {
			var conf LoggerConfig
			if loggerConf, exist := loggerConfs.Config[runEnv]; exist {
				conf = loggerConf
			} else {
				conf = defaultLoggerConfig()
			}

			if _, err = p.NewLogger(loggerConfs.Name, conf); err != nil {
				return
			}
		}

	})

	return
}

func (p *LogrusMate) Logger(loggerName ...string) (logger *logrus.Logger) {
	p.loggersLock.Lock()
	defer p.loggersLock.Unlock()

	name := ""
	if loggerName != nil && len(loggerName) == 1 {
		name = loggerName[0]
	}

	logger, _ = p.loggers[name]

	return
}

func defaultLogrusMate() (logrusMate *LogrusMate) {
	if mate, err := NewLogrusMate(defaultLogrusMateConfig()); err != nil {
		panic(err)
	} else {
		logrusMate = mate
	}
	return
}

func defaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:     "debug",
		Formatter: FormatterConfig{Name: "text", Options: nil},
	}
}

func defaultLogrusMateConfig() LogrusMateConfig {
	return LogrusMateConfig{
		EnvironmentKeys: Environments{RunEnv: "development"},
		Loggers: []LoggerItem{
			{
				Name:   "",
				Config: map[string]LoggerConfig{"development": defaultLoggerConfig()}}},
	}
}
