package formatter

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/awaketai/goweb/framework/contract"
	"github.com/pkg/errors"
)

func JsonFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]any) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	fields["msg"] = msg
	fields["level"] = level
	fields["timestamp"] = t.Format(time.RFC3339)
	c, err := json.Marshal(fields)
	if err != nil {
		return bf.Bytes(), errors.Wrap(err, "json foramt error")
	}
	bf.Write(c)
	return bf.Bytes(), nil
}
