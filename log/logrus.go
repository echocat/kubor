package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type LogrusLevel struct {
	logrus.Level
}

func (instance *LogrusLevel) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

type LogrusFormat string

func (instance *LogrusFormat) Set(plain string) error {
	if plain != "json" && plain != "text" {
		return fmt.Errorf("unsupported log format: %s", plain)
	}
	*instance = LogrusFormat(plain)
	return nil
}

func (instance LogrusFormat) String() string {
	return string(instance)
}

type LogrusColorMode string

func (instance *LogrusColorMode) Set(plain string) error {
	if plain != "auto" && plain != "never" && plain != "always" {
		return fmt.Errorf("unsupported log color mode: %s", plain)
	}
	*instance = LogrusColorMode(plain)
	return nil
}

func (instance LogrusColorMode) String() string {
	return string(instance)
}

type LogrusLogger struct {
	Level              LogrusLevel
	Format             LogrusFormat
	ColorMode          LogrusColorMode
	ReportCaller       bool
	Delegate           *logrus.Logger
	EntryLoggerFactory func(*logrus.Logger) Logger
}

func (instance *LogrusLogger) CreateEntryLogger() Logger {
	if instance.EntryLoggerFactory == nil {
		return &LogrusEntry{
			Root:     instance,
			Delegate: logrus.NewEntry(instance.Delegate),
		}
	}
	return instance.EntryLoggerFactory(instance.Delegate)
}

func (instance *LogrusLogger) WithField(key string, value interface{}) Logger {
	return instance.CreateEntryLogger().WithField(key, value)
}

func (instance *LogrusLogger) WithDeepField(key string, value interface{}) Logger {
	return instance.CreateEntryLogger().WithDeepField(key, value)
}

func (instance *LogrusLogger) WithDeepFieldOn(key string, value interface{}, on func() bool) Logger {
	return instance.CreateEntryLogger().WithDeepFieldOn(key, value, on)
}

func (instance *LogrusLogger) WithError(err error) Logger {
	return instance.CreateEntryLogger().WithError(err)
}

func (instance *LogrusLogger) Trace(msg string, args ...interface{}) {
	instance.CreateEntryLogger().Trace(msg, args...)
}

func (instance *LogrusLogger) Debug(msg string, args ...interface{}) {
	instance.CreateEntryLogger().Debug(msg, args...)
}

func (instance *LogrusLogger) Info(msg string, args ...interface{}) {
	instance.CreateEntryLogger().Info(msg, args...)
}

func (instance *LogrusLogger) Warn(msg string, args ...interface{}) {
	instance.CreateEntryLogger().Warn(msg, args...)
}

func (instance *LogrusLogger) Error(msg string, args ...interface{}) {
	instance.CreateEntryLogger().Error(msg, args...)
}

func (instance *LogrusLogger) Fatal(msg string, args ...interface{}) {
	instance.CreateEntryLogger().Fatal(msg, args...)
}

func (instance *LogrusLogger) IsTraceEnabled() bool {
	return instance.Delegate.Level >= logrus.TraceLevel
}

func (instance *LogrusLogger) IsDebugEnabled() bool {
	return instance.Delegate.Level >= logrus.DebugLevel
}

func (instance *LogrusLogger) IsInfoEnabled() bool {
	return instance.Delegate.Level >= logrus.InfoLevel
}

func (instance *LogrusLogger) IsWarnEnabled() bool {
	return instance.Delegate.Level >= logrus.WarnLevel
}

func (instance *LogrusLogger) IsErrorEnabled() bool {
	return instance.Delegate.Level >= logrus.ErrorLevel
}

func (instance *LogrusLogger) IsFatalEnabled() bool {
	return instance.Delegate.Level >= logrus.FatalLevel
}

func (instance *LogrusLogger) ConfigureFlags(hf HasFlags) {
	hf.Flag("logLevel", "Specifies the minimum required log level.").
		Envar("KUBOR_LOG_LEVEL").
		Default(instance.Level.String()).
		SetValue(&instance.Level)
	hf.Flag("logFormat", "Specifies format output (text or json).").
		Envar("KUBOR_LOG_FORMAT").
		Default(instance.Format.String()).
		SetValue(&instance.Format)
	hf.Flag("logColorMode", "Specifies if the output is in colors or not (auto, never or always).").
		Envar("KUBOR_LOG_COLOR_MODE").
		Default(instance.ColorMode.String()).
		SetValue(&instance.ColorMode)
	hf.Flag("logCaller", "If true the caller details will be logged too.").
		Envar("KUBOR_LOG_CALLER").
		Default("false").
		BoolVar(&instance.ReportCaller)
}

func (instance *LogrusLogger) Init() error {
	instance.Delegate.Level = instance.Level.Level
	instance.Delegate.ReportCaller = instance.ReportCaller

	textFormatter := &logrus.TextFormatter{
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	}
	switch instance.ColorMode {
	case "always":
		textFormatter.ForceColors = true
	case "never":
		textFormatter.DisableColors = true
	}

	instance.Delegate.Formatter = textFormatter
	switch instance.Format {
	case "json":
		instance.Delegate.Formatter = &logrus.JSONFormatter{}
	}
	return nil
}

type LogrusEntry struct {
	Root     *LogrusLogger
	Delegate *logrus.Entry
}

func (instance *LogrusEntry) WithField(key string, value interface{}) Logger {
	return &LogrusEntry{
		Root:     instance.Root,
		Delegate: instance.Delegate.WithField(key, value),
	}
}

func (instance *LogrusEntry) WithDeepField(key string, value interface{}) Logger {
	return instance.WithField(key, JsonValue{
		Value: value,
	})
}

func (instance *LogrusEntry) WithDeepFieldOn(key string, value interface{}, on func() bool) Logger {
	if on() {
		return instance.WithDeepField(key, value)
	}
	return instance
}

func (instance *LogrusEntry) WithError(err error) Logger {
	return &LogrusEntry{
		Root:     instance.Root,
		Delegate: instance.Delegate.WithError(err),
	}
}

func (instance *LogrusEntry) Trace(msg string, args ...interface{}) {
	instance.Delegate.Tracef(msg, args...)
}

func (instance *LogrusEntry) Debug(msg string, args ...interface{}) {
	instance.Delegate.Debugf(msg, args...)
}

func (instance *LogrusEntry) Info(msg string, args ...interface{}) {
	instance.Delegate.Infof(msg, args...)
}

func (instance *LogrusEntry) Warn(msg string, args ...interface{}) {
	instance.Delegate.Warnf(msg, args...)
}

func (instance *LogrusEntry) Error(msg string, args ...interface{}) {
	instance.Delegate.Errorf(msg, args...)
}

func (instance *LogrusEntry) Fatal(msg string, args ...interface{}) {
	instance.Delegate.Fatalf(msg, args...)
}

func (instance *LogrusEntry) IsTraceEnabled() bool {
	return instance.Root.IsTraceEnabled()
}

func (instance *LogrusEntry) IsDebugEnabled() bool {
	return instance.Root.IsDebugEnabled()
}

func (instance *LogrusEntry) IsInfoEnabled() bool {
	return instance.Root.IsInfoEnabled()
}

func (instance *LogrusEntry) IsWarnEnabled() bool {
	return instance.Root.IsWarnEnabled()
}

func (instance *LogrusEntry) IsErrorEnabled() bool {
	return instance.Root.IsErrorEnabled()
}

func (instance *LogrusEntry) IsFatalEnabled() bool {
	return instance.Root.IsFatalEnabled()
}
