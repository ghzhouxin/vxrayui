package log

import (
	"io"
	"log"
	"log/slog"
	"os"
	"sync"

	"zhouxin.learn/go/vxrayui/config"
)

var (
	Logger   *slog.Logger
	initOnce sync.Once
)

func Init() {
	initOnce.Do(func() {
		Logger = initLogger(config.Config.Logger)
	})
}

func initLogger(cfg *config.LogConfig) *slog.Logger {
	var handlers []slog.Handler
	var level slog.Level
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		log.Fatalf("invalid log level: %v", err)
	}

	// Console handler
	if cfg.Console.Enabled {
		var w io.Writer = os.Stdout
		opts := &slog.HandlerOptions{Level: level}

		if cfg.Console.Format == "json" {
			handlers = append(handlers, slog.NewJSONHandler(w, opts))
		} else {
			handlers = append(handlers, slog.NewTextHandler(w, opts))
		}
	}

	// File handler
	if cfg.File.Enabled {
		fileHandler := NewSharedFileHandler(cfg)
		opts := &slog.HandlerOptions{Level: level}
		if cfg.Console.Format == "json" {
			handlers = append(handlers, slog.NewJSONHandler(fileHandler, opts))
		} else {
			handlers = append(handlers, slog.NewTextHandler(fileHandler, opts))
		}
	}

	// Multi-handler
	var handler slog.Handler
	switch len(handlers) {
	case 0:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	case 1:
		handler = handlers[0]
	default:
		handler = newMultiHandler(handlers...)
	}

	return slog.New(handler)
}
