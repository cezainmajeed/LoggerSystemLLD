package logger

import (
	"errors"
	"io"
	"os"
	"time"
)

var (
	ErrInvalidLevel = errors.New("invalid log level")
	ErrNoWriters    = errors.New("at least one writer is required")
)

type Logger struct {
	handlers []Handler
}

func New(minLevel Level, writers ...io.Writer) (*Logger, error) {
	if !minLevel.Valid() {
		return nil, ErrInvalidLevel
	}
	if len(writers) == 0 {
		return nil, ErrNoWriters
	}

	for _, w := range writers {
		if w == nil {
			return nil, errors.New("writer cannot be nil")
		}
	}

	output := NewWriterHandler(NewMultiWriter(writers...))

	return &Logger{
		handlers: []Handler{
			&LevelFilterHandler{MinLevel: minLevel},
			output,
		},
	}, nil
}

func NewConsoleLogger(minLevel Level) (*Logger, error) {
	return New(minLevel, os.Stdout)
}

func (l *Logger) Debug(message string) error {
	return l.log(DebugLevel, message)
}

func (l *Logger) Info(message string) error {
	return l.log(InfoLevel, message)
}

func (l *Logger) Warn(message string) error {
	return l.log(WarnLevel, message)
}

func (l *Logger) Error(message string) error {
	return l.log(ErrorLevel, message)
}

func (l *Logger) log(level Level, message string) error {
	if !level.Valid() {
		return ErrInvalidLevel
	}

	entry := Entry{
		TimeStamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	return l.dispatch(0, entry)
}

func (l *Logger) dispatch(index int, entry Entry) error {
	if index >= len(l.handlers) {
		return nil
	}

	handler := l.handlers[index]
	next := func(nextEntry Entry) error {
		return l.dispatch(index+1, nextEntry)
	}

	return handler.Handle(entry, next)
}
