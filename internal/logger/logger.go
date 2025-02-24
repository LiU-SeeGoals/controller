package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	Logger *zap.SugaredLogger
)

func init() {
	// Configure log rotation using Lumberjack
	logFile := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "../log",
		MaxSize:    10,  // Max size in MB before rotating
		MaxBackups: 5,   // Max number of old log files to keep
		MaxAge:     30,  // Max number of days to retain old logs
		Compress:   true, // Compress old logs (gzip)
	})

	// Create JSON encoder for file logs
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-readable timestamp
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)  // JSON format for logs

	// Create console encoder (human-readable logs)
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// Set up cores for file and console logging
	fileCore := zapcore.NewCore(fileEncoder, logFile, zapcore.DebugLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	// Combine file and console logging
	logCore := zapcore.NewTee(fileCore, consoleCore)

	// Create the logger
	logger := zap.New(logCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Convert to SugaredLogger
	Logger = logger.Sugar()

	// Example of setting a dynamic log level (optional)
	// zap.ReplaceGlobals(logger) // Uncomment to replace the global Zap logger
}

/*
	// Example: Custom Log Level (e.g., "NOTICE")
	const NoticeLevel zapcore.Level = 1  // Between Debug (-1) and Info (0)
	
	// Define a custom level enabler
	func CustomLevelEnabler(level zapcore.Level) bool {
		return level == NoticeLevel || level >= zapcore.InfoLevel
	}

	// Example: Add custom level to the logger
	customCore := zapcore.NewCore(fileEncoder, logFile, zap.LevelEnablerFunc(CustomLevelEnabler))
	logCore := zapcore.NewTee(fileCore, consoleCore, customCore)

	// Log a message with the custom level (if supported in your logger setup)
	LoggerS.Desugar().Check(NoticeLevel, "This is a NOTICE level log").Write()
*/


