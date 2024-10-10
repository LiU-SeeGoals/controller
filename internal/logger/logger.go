package logger

import (
	"os"

	"github.com/LiU-SeeGoals/controller/internal/client"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger  *zap.Logger
	LoggerS *zap.SugaredLogger
)

// Implements io.Writer
type WebWriter struct {}

func (w *WebWriter) Write(p []byte) (n int, err error) {
	client.UpdateWebLog(p)
	return len(p), nil
}

func init() {
	
	// AddSync converts an io.Writer to a WriteSyncer. 
	writeSyncer := zapcore.AddSync(&WebWriter{})
	// Lock wraps a WriteSyncer in a mutex to make it safe for concurrent use. 
	// In particular, *os.Files must be locked before use. 
	writeSyncer = zapcore.Lock(writeSyncer)
	
	// NewMultiWriteSyncer creates a WriteSyncer that duplicates its writes 
	// and sync calls, much like io.MultiWriter. 
	mw := zapcore.NewMultiWriteSyncer(writeSyncer, os.Stdout)
	
	encoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()) // NewProductionEncoderConfig() also available
	core := zapcore.NewCore(encoder, mw, zapcore.DebugLevel) // DebugLevel, InfoLevel, WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, FatalLevel
	Logger = zap.New(core)
	LoggerS = Logger.Sugar() // Sugar wraps the Logger to provide a more ergonomic, but slightly slower, API.
}
