package logs

import (
	"os"
	"sync"
)

type ConsoleWriter struct {
	mu sync.Mutex
}

func (c *ConsoleWriter) Write(p []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return os.Stdout.Write(p)
}

func (c *ConsoleWriter) Close() error {
	return nil
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{
		mu: sync.Mutex{},
	}
}
