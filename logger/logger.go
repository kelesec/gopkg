package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Level int8

func (l Level) toZeroLogLevel() zerolog.Level {
	switch l {
	case TraceLevel:
		return zerolog.TraceLevel
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	case PanicLevel:
		return zerolog.PanicLevel
	case NoLevel:
		return zerolog.NoLevel
	case Disabled:
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel
	Disabled
	TraceLevel Level = -1
)

// Config 日志初始化配置
type Config struct {
	level           Level  // 日志等级
	logDir          string // 日志保存目录
	logFile         string // 日志文件名
	allowConsoleLog bool   // 允许控制台日志输出
	allowFileLog    bool   // 允许保存日志文件
	maxSize         int    // 每个日志文件最大尺寸（MB）
	maxBackups      int    // 保留的旧日志文件数量
	maxAge          int    // 保留的旧日志文件天数
	compress        bool   // 是否压缩旧日志文件
	jsonFormat      bool   // 是否使用JSON格式（文件日志）
	caller          bool   // 是否显示调用者信息
}

// InitLogger 初始化日志配置
func InitLogger(ops ...Option) error {
	config := DefaultConfig()
	for _, f := range ops {
		f(config)
	}

	// 全局配置：确保多个 writter 日志时间格式统一
	zerolog.TimeFieldFormat = time.DateTime

	// 允许控制台日志
	var writers []io.Writer
	if config.allowConsoleLog {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.DateTime,
		})
	}

	// 允许保存日志文件
	if config.allowFileLog {
		if config.logFile == "" {
			config.logFile = "app.log"
		}

		if config.logDir != "" {
			if err := os.MkdirAll(config.logDir, 0755); err != nil {
				return fmt.Errorf("create log directory: %w", err)
			}
		}

		logPath := filepath.Clean(filepath.Join(config.logDir, config.logFile))
		lumberLogger := lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    config.maxSize,
			MaxAge:     config.maxAge,
			MaxBackups: config.maxBackups,
			Compress:   config.compress,
		}

		if config.jsonFormat {
			writers = append(writers, &lumberLogger)
		} else {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        &lumberLogger,
				NoColor:    true,
				TimeFormat: time.DateTime,
			})
		}
	}

	var multiWriter io.Writer
	if len(writers) == 0 {
		// 禁用控制台和文件日志，说明不需要日志输出
		multiWriter = io.Discard
	} else {
		multiWriter = zerolog.MultiLevelWriter(writers...)
	}

	loggerCtx := zerolog.New(multiWriter).With().Timestamp()
	if config.caller {
		loggerCtx = loggerCtx.Caller()
	}

	logger := loggerCtx.Logger().Level(config.level.toZeroLogLevel())
	log.Logger = logger
	zerolog.DefaultContextLogger = &logger
	return nil
}

func Trace() *zerolog.Event {
	return log.Trace()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}

func Panic() *zerolog.Event {
	return log.Panic()
}
