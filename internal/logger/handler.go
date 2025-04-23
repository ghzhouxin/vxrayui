package logger

import (
	"context"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"zhouxin.learn/go/vxrayui/config"
)

type ShardFileHandler struct {
	cfg     *config.LogConfig
	writer  *os.File
	current string
	shardFn func() string
}

func NewSharedFileHandler(cfg *config.LogConfig) *ShardFileHandler {
	if err := os.MkdirAll(cfg.File.Path, 0755); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	h := &ShardFileHandler{cfg: cfg}
	h.setShardFunction()

	if err := h.rotate(); err != nil {
		log.Fatalf("FileHandler.rotate err: %v", err)
	}

	return h
}

func (h *ShardFileHandler) setShardFunction() {
	switch h.cfg.File.ShardBy {
	case "hour":
		h.shardFn = func() string {
			return time.Now().Format("2006-01-02_15")
		}
	case "minute":
		h.shardFn = func() string {
			return time.Now().Format("2006-01-02_15-04")
		}
	default: // day
		h.shardFn = func() string {
			return time.Now().Format("2006-01-02")
		}
	}
}

func (h *ShardFileHandler) rotate() error {
	shard := h.shardFn()
	if shard == h.current && h.writer != nil {
		return nil
	}

	if h.writer != nil {
		h.writer.Close()
	}

	filename := filepath.Join(
		h.cfg.File.Path,
		h.cfg.File.Filename+"."+shard,
	)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	h.writer = f
	h.current = shard
	return nil
}

func (h *ShardFileHandler) Write(p []byte) (n int, err error) {
	if err := h.rotate(); err != nil {
		return 0, err
	}
	return h.writer.Write(p)
}

func (h *ShardFileHandler) Close() error {
	if h.writer != nil {
		return h.writer.Close()
	}
	return nil
}

type multiHandler struct {
	handlers []slog.Handler
}

func newMultiHandler(handlers ...slog.Handler) *multiHandler {
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) Enabled(context context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(context, level) {
			return true
		}
	}

	return false
}

func (h *multiHandler) Handle(context context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if !handler.Enabled(context, r.Level) {
			continue
		}

		if err := handler.Handle(context, r); err != nil {
			return err
		}
	}

	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return newMultiHandler(handlers...)
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return newMultiHandler(handlers...)
}
