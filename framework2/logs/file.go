package logs

import (
	"fmt"
	"os"
	"sync"
)

type FileWriter struct {
	file *os.File
	mu   sync.Mutex
}

func NewFileWriter(fileName string) (*FileWriter, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("open log file writer [%s] failed:%w", fileName, err)
	}

	return &FileWriter{file: file, mu: sync.Mutex{}}, nil
}

func (f *FileWriter) Write(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.file.Write(p)
}

func (f *FileWriter) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.file.Close()
}
