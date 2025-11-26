package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	writer  io.Writer
	logFile *os.File
	FileBuf *bufio.Writer
}

func (l *Logger) Logf(format string, args ...any) {
	format = "\n" + format
	l.writer.Write([]byte(fmt.Sprintf(format, args...)))
	if l.logFile != nil {
		l.FileBuf.Flush()
		l.logFile.Sync()
	}
}

func NewLoggerFile(path string) (*Logger, func()) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("\nunable to create file: %s", err)
		return nil, nil
	}

	writer := bufio.NewWriter(file)
	lgr := NewLogger(writer)
	lgr.logFile = file
	lgr.FileBuf = writer
	return lgr, func() {
		writer.Flush()
		file.Close()
	}
}

func NewLogger(writer io.Writer) *Logger {
	return &Logger{
		writer: writer,
	}
}
