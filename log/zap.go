package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Level   = zapcore.Level
	Logger  = zap.Logger
	LoggerS = zap.SugaredLogger
)

type (
	Option = zap.Option
)

var pkgCallerSkip = zap.AddCallerSkip(2)
