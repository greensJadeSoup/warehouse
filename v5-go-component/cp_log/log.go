package cp_log

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_dc"
)

var L *Logger

type Logger struct {
	*zap.Logger
}

func NewLogger(logConf *cp_dc.DcLogConfig) *Logger {
	level := zapcore.DebugLevel
	coreList := make([]zapcore.Core, 0)

	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// ================= ConsoleWriter ===================
	if logConf.ConsoleWriter.Enable {
		switch logConf.ConsoleWriter.Level {
			case cp_constant.LevelPanic: 		level = zapcore.DebugLevel
			case cp_constant.LevelFatal: 		level = zapcore.FatalLevel
			case cp_constant.LevelError: 		level = zapcore.ErrorLevel
			case cp_constant.LevelWarning: 		level = zapcore.WarnLevel
			case cp_constant.LevelInformational: 	level = zapcore.InfoLevel
			case cp_constant.LevelDebug: 		level = zapcore.DebugLevel

			default: panic(errors.New("logConf.ConsoleWriter.level 非法:" + logConf.ConsoleWriter.Level))
		}
		coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	// ================= InfoFileWriter ===================
	if logConf.InfoFileWriter.Enable {
		infoJackLogger := &lumberjack.Logger{
			Filename:   logConf.InfoFileWriter.Path,
			MaxSize:    50,
			MaxBackups: 30,
			MaxAge:     30,
			LocalTime:  true,
			Compress:   false,
		}
		switch logConf.InfoFileWriter.Level {
			case cp_constant.LevelPanic: 		level = zapcore.DebugLevel
			case cp_constant.LevelFatal: 		level = zapcore.FatalLevel
			case cp_constant.LevelError: 		level = zapcore.ErrorLevel
			case cp_constant.LevelWarning: 		level = zapcore.WarnLevel
			case cp_constant.LevelInformational: 	level = zapcore.InfoLevel
			case cp_constant.LevelDebug: 		level = zapcore.DebugLevel

			default: panic(errors.New("logConf.InfoFileWriter.level 非法:" + logConf.InfoFileWriter.Level))
		}
		coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(infoJackLogger), level))
	}

	// ================= ErrorFileWriter ===================
	if logConf.ErrorFileWriter.Enable {
		errJackLogger := &lumberjack.Logger{
			Filename:  logConf.ErrorFileWriter.Path,
			MaxSize:    50,
			MaxBackups: 50,
			MaxAge:     30,
			LocalTime:  true,
			Compress:   false,
		}
		switch logConf.ErrorFileWriter.Level {
			case cp_constant.LevelPanic: 		level = zapcore.DebugLevel
			case cp_constant.LevelFatal: 		level = zapcore.FatalLevel
			case cp_constant.LevelError: 		level = zapcore.ErrorLevel
			case cp_constant.LevelWarning: 		level = zapcore.WarnLevel
			case cp_constant.LevelInformational: 	level = zapcore.InfoLevel
			case cp_constant.LevelDebug: 		level = zapcore.DebugLevel

			default: panic(errors.New("logConf.ErrorFileWriter.level 非法:" + logConf.ErrorFileWriter.Level))
		}
		coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(errJackLogger), level))
	}

	core := zapcore.NewTee(coreList...)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.WarnLevel))
	L = &Logger{
		Logger: zapLogger,
	}

	return L
}

func Error(msg string, fields ...zap.Field) {
	L.Error(msg, fields...)
}

func Warning(msg string, fields ...zap.Field) {
	L.Warn(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	L.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	L.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	L.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	L.Panic(msg, fields...)
}