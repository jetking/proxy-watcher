package instance

import (
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

var loggerOnce sync.Once

func Logger() *zap.Logger {
	path := viper.GetString("logger.path")
	level := zapcore.DebugLevel
	levelStr := viper.GetString("logger.level")
	switch levelStr {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}
	loggerOnce.Do(func() {
		logger = zap.New(zapcore.NewCore(jsonEncoder(), writeSyncer(path), level)) //.Sugar()
	})

	// return logger
	fs := zap.Fields()
	return logger.WithOptions(fs)
}
func writeSyncer(f string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   f,
		MaxSize:    512,
		MaxBackups: 120,
		MaxAge:     60,
		Compress:   true,
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}

func jsonEncoder() zapcore.Encoder {
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("20060102 15:04:05.000"))
	}
	return zapcore.NewJSONEncoder(encoder)
}
