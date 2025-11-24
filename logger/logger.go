package utils

import (
	"fmt"
	"io"
)

type Logger struct {
	writer io.Writer
}

func (l *Logger) Logf(format string, args ...any) {
	format = "\n" + format
	l.writer.Write([]byte(fmt.Sprintf(format, args...)))
}

func NewLogger(writer io.Writer) Logger {
	return Logger{
		writer: writer,
	}
}
