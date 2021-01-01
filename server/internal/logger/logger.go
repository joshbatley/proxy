package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Setup logger
func Setup() *zap.SugaredLogger {
	// writerSyncer := getLogWriter()
	// core := zapcore.NewCore(
	// 	zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
	// 	writerSyncer,
	// 	zapcore.DebugLevel,
	// )
	logger, _ := zap.NewDevelopment()
	sugarLogger := logger.Sugar()

	return sugarLogger
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
