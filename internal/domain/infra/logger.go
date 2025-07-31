package infra

type Logger interface {
	Debug(msg string, fields ...interface{})
	Debugf(format string, args ...interface{})
	Info(msg string, fields ...interface{})
	Infof(format string, args ...interface{})
	Warn(msg string, fields ...interface{})
	Warnf(format string, args ...interface{})
	Error(msg string, fields ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(msg string, fields ...interface{})
	Fatalf(format string, args ...interface{})

	With(fields ...interface{}) Logger
	Sync() error
}
