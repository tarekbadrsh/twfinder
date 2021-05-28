package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	*zap.Logger
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("|%s|", t.Format("2006-01-02T15:04:05")))
}

func filelogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05"))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("[%v]", level.CapitalString()))
}

func getZapLogLevelFunc(level string) zap.LevelEnablerFunc {
	switch level {
	case "DEBUG":
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.DebugLevel
		})
	case "INFO":
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel
		})
	case "WARN":
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.WarnLevel
		})
	case "ERROR":
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
	case "FATAL":
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.FatalLevel
		})
	}
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
}

// NewZapLogger :
func NewZapLogger(fileLogLevel, terminalLogLevel string) ILogger {
	terminalEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:     "M",
		LevelKey:       "L",
		TimeKey:        "T",
		NameKey:        "N",
		CallerKey:      "C",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     syslogTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	terminalOutput := zapcore.AddSync(os.Stderr)

	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:     "M",
		LevelKey:       "L",
		TimeKey:        "T",
		NameKey:        "N",
		CallerKey:      "C",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     filelogTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	fileOutput := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/logs.json",
		MaxSize:    10, // megabytes
		MaxBackups: 100,
		MaxAge:     28, // days
		Compress:   true,
	})

	core := zapcore.NewTee()
	if fileLogLevel != "" && terminalLogLevel != "" {
		core = zapcore.NewTee(
			zapcore.NewCore(terminalEncoder, terminalOutput, getZapLogLevelFunc(terminalLogLevel)),
			zapcore.NewCore(fileEncoder, fileOutput, getZapLogLevelFunc(fileLogLevel)),
		)
	} else if fileLogLevel != "" {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, fileOutput, getZapLogLevelFunc(fileLogLevel)),
		)
	} else if terminalLogLevel != "" {
		core = zapcore.NewTee(
			zapcore.NewCore(terminalEncoder, terminalOutput, getZapLogLevelFunc(terminalLogLevel)),
		)
	}
	return &zapLogger{zap.New(core)}
}

func (l *zapLogger) writer(lvl zapcore.Level, a ...interface{}) {
	var msg = fmt.Sprint(a...)
	if ce := l.Check(lvl, msg); ce != nil {
		ce.Entry.Caller = zapcore.NewEntryCaller(runtime.Caller(3))
		ce.Write()
	}
}

func (l *zapLogger) writerf(lvl zapcore.Level, format string, prm ...interface{}) {
	var msg = fmt.Sprintf(format, prm...)
	if ce := l.Check(lvl, msg); ce != nil {
		ce.Entry.Caller = zapcore.NewEntryCaller(runtime.Caller(3))
		ce.Write()
	}
}

// Debug :
func (l *zapLogger) Debug(a ...interface{}) {
	l.writer(zap.DebugLevel, a...)
}

// Debugf :
func (l *zapLogger) Debugf(format string, prm ...interface{}) {
	l.writerf(zap.DebugLevel, format, prm...)
}

// Info :
func (l *zapLogger) Info(a ...interface{}) {
	l.writer(zap.InfoLevel, a...)
}

// Infof :
func (l *zapLogger) Infof(format string, prm ...interface{}) {
	l.writerf(zap.InfoLevel, format, prm...)
}

// Warn :
func (l *zapLogger) Warn(a ...interface{}) {
	l.writer(zap.WarnLevel, a...)
}

// Warnf :
func (l *zapLogger) Warnf(format string, prm ...interface{}) {
	l.writerf(zap.WarnLevel, format, prm...)
}

// Error :
func (l *zapLogger) Error(a ...interface{}) {
	l.writer(zap.ErrorLevel, a...)
}

// Errorf :
func (l *zapLogger) Errorf(format string, prm ...interface{}) {
	l.writerf(zap.ErrorLevel, format, prm...)
}

// Fatal :
func (l *zapLogger) Fatal(a ...interface{}) {
	l.writer(zap.FatalLevel, a...)
}

// Fatalf :
func (l *zapLogger) Fatalf(format string, prm ...interface{}) {
	l.writerf(zap.FatalLevel, format, prm...)
}

// Sync :
func (l *zapLogger) Close() {
	l.Sync()
}
