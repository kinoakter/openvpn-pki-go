package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	sugared *LoggerS
	mx      sync.RWMutex
)

func init() {
	SetLogger(zap.Must(zap.NewProduction()))
}

func SetLogger(logger *Logger) {
	mx.Lock()
	defer mx.Unlock()
	sugared = logger.
		WithOptions(pkgCallerSkip).
		Sugar()
}
func logf(lvl Level, template string, args ...any) {
	mx.RLock()
	s := sugared
	mx.RUnlock()
	s.Logf(lvl, template, args...)
}

func Debugf(template string, args ...any) {
	logf(zapcore.DebugLevel, template, args...)
}

func Infof(template string, args ...any) {
	logf(zapcore.InfoLevel, template, args...)
}

func Warnf(template string, args ...any) {
	logf(zapcore.WarnLevel, template, args...)
}

func Errorf(template string, args ...any) {
	logf(zapcore.ErrorLevel, template, args...)
}

func Fatalf(template string, args ...any) {
	logf(zapcore.FatalLevel, template, args...)
}
