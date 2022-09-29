package formatter

import (
	"bytes"
	"fmt"
	"time"

	"github.com/awaketai/goweb/framework/contract"
)

func TextFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]any) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	separator := "\t"
	prefix := Prefix(level)
	bf.WriteString(prefix)
	bf.WriteString(separator)

	ts := t.Format(time.RFC3339)
	bf.WriteString(ts)
	bf.WriteString(separator)

	// msg
	bf.WriteString("\"")
	bf.WriteString(msg)
	bf.WriteString("\"")
	bf.WriteString(separator)

	bf.WriteString(fmt.Sprint(fields))
	return bf.Bytes(), nil
}
