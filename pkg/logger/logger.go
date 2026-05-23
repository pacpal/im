package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

// Init 初始化全局 logger，level: debug/info/warn/error，format: json/console
func Init(level string, format string) error {
	var cfg zap.Config
	if strings.ToLower(format) == "console" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	switch strings.ToLower(level) {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn", "warning":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	l, err := cfg.Build()
	if err != nil {
		return err
	}
	sugar = l.Sugar()
	return nil
}

// Sync flushes any buffered logs
func Sync() {
	if sugar != nil {
		_ = sugar.Sync()
	}
}

func Debugf(format string, args ...interface{}) {
	if sugar != nil {
		sugar.Debugf(format, args...)
		return
	}
	fmt.Printf(format+"\n", args...)
}

func Infof(format string, args ...interface{}) {
	if sugar != nil {
		sugar.Infof(format, args...)
		return
	}
	fmt.Printf(format+"\n", args...)
}

func Warnf(format string, args ...interface{}) {
	if sugar != nil {
		sugar.Warnf(format, args...)
		return
	}
	fmt.Printf(format+"\n", args...)
}

func Errorf(format string, args ...interface{}) {
	if sugar != nil {
		sugar.Errorf(format, args...)
		return
	}
	fmt.Printf(format+"\n", args...)
}

func Fatalf(format string, args ...interface{}) {
	if sugar != nil {
		sugar.Fatalf(format, args...)
		return
	}
	fmt.Printf(format+"\n", args...)
}

func Debug(args ...interface{}) {
	if sugar != nil {
		sugar.Debug(args...)
		return
	}
	fmt.Println(args...)
}

func Info(args ...interface{}) {
	if sugar != nil {
		sugar.Info(args...)
		return
	}
	fmt.Println(args...)
}

func Warn(args ...interface{}) {
	if sugar != nil {
		sugar.Warn(args...)
		return
	}
	fmt.Println(args...)
}

func Error(args ...interface{}) {
	if sugar != nil {
		sugar.Error(args...)
		return
	}
	fmt.Println(args...)
}

func Fatal(args ...interface{}) {
	if sugar != nil {
		sugar.Fatal(args...)
		return
	}
	fmt.Println(args...)
}

// With returns a sugared logger with added context fields
func With(args ...interface{}) *zap.SugaredLogger {
	if sugar != nil {
		return sugar.With(args...)
	}
	return zap.NewNop().Sugar()
}
