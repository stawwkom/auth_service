package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var globalLogger *zap.Logger

func InitLogger(logFilePath string, isDev bool) error {
	// Настройка ротации файла логов
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // МБ
		MaxBackups: 7,
		MaxAge:     30, // дней
		Compress:   true,
	})

	// Энкодер для консоли и файла (читаемый в dev, JSON в prod)
	var encoderConfig zapcore.EncoderConfig
	if isDev {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Консольный энкодер (читаемый)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	// Файловый энкодер (JSON)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Потоки вывода: файл + консоль
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, w, zapcore.InfoLevel), // В файл пишем только Info+
	)

	globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Можно заменить глобальный логгер zap (если хочешь)
	zap.ReplaceGlobals(globalLogger)

	return nil
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
