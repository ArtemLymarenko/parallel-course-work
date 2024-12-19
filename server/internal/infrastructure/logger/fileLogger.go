package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type fileLogger struct {
	file       *os.File
	writer     *bufio.Writer
	logsBuffer []string
	maxBatch   int
}

func NewFileLogger(filepath string, maxBatch int) (*fileLogger, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)

	return &fileLogger{
		file:       file,
		writer:     writer,
		logsBuffer: []string{},
		maxBatch:   maxBatch,
	}, nil
}

func (l *fileLogger) LogUnsafe(v ...interface{}) {
	toLog := fmt.Sprint(v...)
	l.logsBuffer = append(l.logsBuffer, toLog)
	if len(l.logsBuffer) >= l.maxBatch {
		l.flush()
	}
}

func (l *fileLogger) flush() {
	for _, toLog := range l.logsBuffer {
		_, err := l.writer.WriteString(toLog + "\n")
		if err != nil {
			log.Println("error writing to file:", err)
			return
		}
	}

	err := l.writer.Flush()
	if err != nil {
		log.Println("error flushing buffer:", err)
	}

	l.logsBuffer = l.logsBuffer[:0]
}

func (l *fileLogger) Close() {
	l.flush()
	_ = l.file.Close()
}
