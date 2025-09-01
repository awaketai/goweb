package logs

import (
	"fmt"
	"strconv"
	"time"
	jsoniter "github.com/json-iterator/go"
)

type JSONFormatter struct{
	IgnoreBasicFields bool
}

func (j *JSONFormatter) Format(e *Entry) error {
	if !j.IgnoreBasicFields {
		e.Map["level"] = LevelNameMapping[e.Level]
		e.Map["time"] = e.Time.Format(time.RFC3339)
		if e.File != "" {
			e.Map["file"] = e.File + ":" + strconv.Itoa(e.Line)
			e.Map["func"] = e.Func
		}
		switch e.Format {
			case FmtEmptySeparate:
			e.Map["message"] = fmt.Sprint(e.Args...)
			default:
			e.Map["message"] = fmt.Sprintf(e.Format, e.Args...)
		}
		
		return jsoniter.NewEncoder(e.Buffer).Encode(e.Map)
	}
	
	switch e.Format {
		case FmtEmptySeparate:
		for _,arg := range e.Args {
			if err := jsoniter.NewEncoder(e.Buffer).Encode(arg);err != nil {
				return fmt.Errorf("log json formatter encode err:%w",err)
			}
		}
	}
	
	return nil
}