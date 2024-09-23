package logs

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type LogType int8

const (
	LogTypeConsole LogType = iota
	LogTypeFile
)

type Logger interface {
	SetLevel(level Level)
	SetType(t LogType)
	Flush()
}
