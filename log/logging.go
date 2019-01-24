package log

import (
	"bytes"
	"encoding/json"
	"github.com/alecthomas/kingpin"
	"github.com/sirupsen/logrus"
)

var DefaultLogger RootLogger = &LogrusLogger{
	Level:     LogrusLevel{logrus.InfoLevel},
	Format:    LogrusFormat("text"),
	ColorMode: LogrusColorMode("auto"),
	Delegate:  logrus.New(),
}

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithDeepField(key string, value interface{}) Logger
	WithDeepFieldOn(key string, value interface{}, on func() bool) Logger
	WithError(err error) Logger

	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})

	IsTraceEnabled() bool
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool
	IsFatalEnabled() bool
}

type HasFlags interface {
	Flag(name, help string) *kingpin.FlagClause
}

type RootLogger interface {
	Logger

	Init() error
	ConfigureFlags(HasFlags)
}

func WithField(key string, value interface{}) Logger {
	return DefaultLogger.WithField(key, value)
}

func WithDeepField(key string, value interface{}) Logger {
	return DefaultLogger.WithDeepField(key, value)
}

func WithDeepFieldOn(key string, value interface{}, on func() bool) Logger {
	return DefaultLogger.WithDeepFieldOn(key, value, on)
}

func WithError(err error) Logger {
	return DefaultLogger.WithError(err)
}

func Trace(msg string, args ...interface{}) {
	DefaultLogger.Trace(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	DefaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	DefaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	DefaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	DefaultLogger.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	DefaultLogger.Fatal(msg, args...)
}

func IsTraceEnabled() bool {
	return DefaultLogger.IsTraceEnabled()
}

func IsDebugEnabled() bool {
	return DefaultLogger.IsDebugEnabled()
}

func IsInfoEnabled() bool {
	return DefaultLogger.IsInfoEnabled()
}

func IsWarnEnabled() bool {
	return DefaultLogger.IsWarnEnabled()
}

func IsErrorEnabled() bool {
	return DefaultLogger.IsErrorEnabled()
}

func IsFatalEnabled() bool {
	return DefaultLogger.IsFatalEnabled()
}

type JsonValue struct {
	Value       interface{}
	PrettyPrint bool
}

func (instance JsonValue) String() string {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if instance.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(instance.Value); err != nil {
		panic(err)
	}
	return buf.String()
}
