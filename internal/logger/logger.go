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
type WebWriter struct{}

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
	multiWriter := zapcore.NewMultiWriteSyncer(writeSyncer, os.Stdout)

	lvl := zapcore.DebugLevel
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, multiWriter, lvl)
	Logger = zap.New(core, zap.AddCaller())

	// Sugar wraps the Logger to provide a more ergonomic, but slightly slower, API.
	LoggerS = Logger.Sugar()
}
