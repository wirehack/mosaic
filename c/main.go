package c

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logger := zap.Must(config.Build())
	return logger.Sugar()
}

func PrintRoutes(router chi.Router) {
	err := chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Println(method, route)
		return nil
	})
	if err != nil {
		fmt.Println("Error walking routes:", err)
	}
}

func PrintFS(fsys fs.FS, path string) {
	fs.WalkDir(fsys, path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("error:", err)
			return nil
		}
		fmt.Println(path)
		return nil
	})
}
