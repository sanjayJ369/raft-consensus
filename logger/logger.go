package utils

import (
	"fmt"
	"io"
)

type Logger struct {
	writer io.Writer
}

func (l *Logger) Logf(format string, args ...any) {
	l.writer.Write([]byte(fmt.Sprintf(format, args...)))
}

func NewLogger(writer io.Writer) Logger {
	return Logger{
		writer: writer,
	}
}
