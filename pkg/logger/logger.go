// Package logger 提供项目内的轻量日志封装，基于 zap 实现简易的全局日志函数。
package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugar *zap.SugaredLogger
var lgr *zap.Logger

// Init 初始化全局 logger，level: debug/info/warn/error，format: json/console
// 输出包含 timestamp、caller、stacktrace（error 及以上级别）
func Init(level string, format string) error {
	var cfg zap.Config
	if strings.ToLower(format) == "console" {
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
	} else {
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
	}

	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.StacktraceKey = "stacktrace"

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

	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = []string{"stdout"}
	}

	lg, err := cfg.Build(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		return err
	}

	lgr = lg
	sugar = lgr.Sugar()
	return nil
}

// Sync 刷新并关闭底层日志缓冲（需要在程序退出时调用）。
func Sync() {
	if lgr != nil {
		_ = lgr.Sync()
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

// 结构化日志便捷方法，方便在代码中使用 key-value 的形式记录额外字段
func Infow(msg string, keysAndValues ...interface{}) {
	if sugar != nil {
		sugar.Infow(msg, keysAndValues...)
		return
	}
	if len(keysAndValues) > 0 {
		fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
	} else {
		fmt.Println(msg)
	}
}

func Errorw(msg string, keysAndValues ...interface{}) {
	if sugar != nil {
		sugar.Errorw(msg, keysAndValues...)
		return
	}
	if len(keysAndValues) > 0 {
		fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
	} else {
		fmt.Println(msg)
	}
}

func Debugw(msg string, keysAndValues ...interface{}) {
	if sugar != nil {
		sugar.Debugw(msg, keysAndValues...)
		return
	}
	if len(keysAndValues) > 0 {
		fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
	} else {
		fmt.Println(msg)
	}
}

func Warnw(msg string, keysAndValues ...interface{}) {
	if sugar != nil {
		sugar.Warnw(msg, keysAndValues...)
		return
	}
	if len(keysAndValues) > 0 {
		fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
	} else {
		fmt.Println(msg)
	}
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	if sugar != nil {
		sugar.Fatalw(msg, keysAndValues...)
		return
	}
	if len(keysAndValues) > 0 {
		fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
	} else {
		fmt.Println(msg)
	}
}

// StartStep 用于记录某个操作的开始与结束，返回一个在操作结束时调用的函数。
// 用法：done := logger.StartStep("UserService.AddFriend", "from", fromUID, "to", toUID)
//
//	defer done(err) 或在每个分支返回前调用 done(err)
func StartStep(name string, fields ...interface{}) func(error) {
	start := time.Now()
	if sugar != nil {
		sugar.Infow(name+" started", append(fields, "phase", "start")...)
	} else {
		if len(fields) > 0 {
			fmt.Println(append([]interface{}{name + " started"}, fields...)...)
		} else {
			fmt.Println(name + " started")
		}
	}

	return func(err error) {
		duration := time.Since(start)
		if err != nil {
			if sugar != nil {
				sugar.Errorw(name+" failed", append(fields, "phase", "end", "duration", duration.String(), "error", err)...)
			} else {
				fmt.Println(append([]interface{}{name + " failed", "duration", duration.String(), "error", err}, fields...)...)
			}
			return
		}
		if sugar != nil {
			sugar.Infow(name+" completed", append(fields, "phase", "end", "duration", duration.String())...)
		} else {
			fmt.Println(append([]interface{}{name + " completed", "duration", duration.String()}, fields...)...)
		}
	}
}

// With 返回一个带有附加上下文字段的 SugaredLogger，方便链式调用。
func With(args ...interface{}) *zap.SugaredLogger {
	if sugar != nil {
		return sugar.With(args...)
	}
	return zap.NewNop().Sugar()
}
