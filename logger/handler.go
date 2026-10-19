package logger

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

type Next func(Entry) error

type Handler interface {
	Handle(entry Entry, next Next) error
}

type LevelFilterHandler struct {
	MinLevel Level
}

func (h *LevelFilterHandler) Handle(
	entry Entry,
	next Next,
) error {
	if entry.Level < h.MinLevel {
		return nil
	}
	return next(entry)
}

type WriterHandler struct {
	writer io.Writer
	mu     sync.Mutex
}

func NewWriterHandler(w io.Writer) *WriterHandler {
	return &WriterHandler{writer: w}
}

func (h *WriterHandler) Handle(
	entry Entry,
	_ Next,
) error {
	line := fmt.Sprintf(
		"%s [%s] %s\n",
		entry.TimeStamp.Format("2006-01-02 15:04:05"),
		entry.Level,
		strings.ReplaceAll(entry.Message, "\n", "\\n"),
	)

	h.mu.Lock()
	defer h.mu.Unlock()

	n, err := io.WriteString(h.writer, line)
	if err != nil {
		return err
	}
	if n != len(line) {
		return io.ErrShortWrite
	}

	return nil
}

func NewMultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}
