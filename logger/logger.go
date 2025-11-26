package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	sync.Mutex
	writer  io.Writer
	logFile *os.File
	FileBuf *bufio.Writer
	strict  bool // if strict will log to file and flush
}

func (l *Logger) Sync() {
	if l.logFile != nil {
		l.FileBuf.Flush()
		l.logFile.Sync()
	}
}

func (l *Logger) Logf(format string, args ...any) {
	l.Lock()
	defer l.Unlock()

	timestamp := time.Now().UTC()
	message := fmt.Sprintf(format, args...)
	log := fmt.Sprintf("\n[%s]%s", timestamp, message)
	l.writer.Write([]byte(log))
	if l.logFile != nil && l.strict {
		l.FileBuf.Flush()
		l.logFile.Sync()
	}
}

func NewLoggerFile(path string, strict bool) (*Logger, func()) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("\n[]unable to create file: %s", err)
		return nil, nil
	}

	writer := bufio.NewWriter(file)
	lgr := NewLogger(writer)
	lgr.logFile = file
	lgr.FileBuf = writer
	lgr.strict = strict
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
