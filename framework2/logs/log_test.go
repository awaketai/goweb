package logs

import (
	"os"
	"testing"
)

// go test -timeout 30s -run ^TestFileLogWriter$
func TestFileLogWriter(t *testing.T) {
	file, err := os.OpenFile("goweb.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	fileWriter, err := NewFileWriter("goweb.log")
	if err != nil {
		t.Error(err)
		return
	}

	defer fileWriter.Close()
	logger := NewLog(LevelInfo, Ldate|Ltime|Lmicrosecond|Llongfile, fileWriter)
	logger.Debug("this is debug log")
	logger.Info("this is info log")
	logger.Warn("this is warn log")
	logger.Error("this is error log")
	// custom
	var custom = struct {
		Name  string
		Value int
	}{
		Name:  "test",
		Value: 1,
	}
	logger.Debug(custom)
	logger.Fatal("this is fatal log")
}

func TestConsoleLogWriter(t *testing.T) {
	consoleWriter := NewConsoleWriter()
	logger := NewLog(LevelInfo, Ldate|Ltime|Lmicrosecond, consoleWriter)
	logger.Debug("this is console debug log")
	logger.Info("this is console  info log")
	logger.Warn("this is console  warn log")
	logger.Error("this is console  error log")
	// custom
	var custom = struct {
		Name  string
		Value int
	}{
		Name:  "test",
		Value: 1,
	}
	logger.Debug(custom)
	logger.Fatal("this is console  fatal log")
}
