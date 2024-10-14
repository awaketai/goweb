package config

import (
	"bufio"
	"bytes"
	"context"
	"strings"
	"sync"
)

// core code from beego
var (
	defaultSection = "default"
	bNumComment    = []byte{'#'}
	bSemComment    = []byte{';'}
	bEmpty         = []byte{}
	bEqual         = []byte{'='}
	bDQuote        = []byte{'"'}
	sectionStart   = []byte{'['}
	sectionEnd     = []byte{']'}
	lineBreak      = "\n"
)

type InitConfig struct {
}

func (i *InitConfig) Parse(name string) (Configer, error) {
	return i.parseFile(name)
}

type IniConfigContainer struct {
	BaseConfiger
	data           map[string]map[string]string
	sectionComment map[string]string
	keyComment     map[string]string
	sync.RWMutex
}

func (i *InitConfig) parseFile(name string) (*IniConfigContainer, error) {
	return nil, nil
}

func (i *InitConfig) parseData(dir string, data []byte) (*IniConfigContainer, error) {
	cfg := &IniConfigContainer{
		data:           make(map[string]map[string]string),
		sectionComment: make(map[string]string),
		keyComment:     make(map[string]string),
	}
	cfg.BaseConfiger = NewBaseConfiger(func(ctx context.Context, key string) (string, error) {
		return cfg.getData(key), nil
	})
	cfg.Lock()
	defer cfg.Unlock()
	//var comment bytes.Buffer
	buf := bufio.NewReader(bytes.NewReader(data))
	head, err := buf.Peek(3)
	// check UTF-8 byte order mark的特定组合
	// 0xEF 0xBB 0xBF
	if err == nil && head[0] == 239 && head[1] == 187 && head[2] == 191 {
		for i := 1; i <= 3; i++ {
			buf.ReadByte()
		}
	}
	//section := defaultSection
	tmpBuf := bytes.NewBuffer(nil)
	for {
		tmpBuf.Reset()
		//shouldBreak := false
	}
}

func (i *IniConfigContainer) getData(key string) string {
	if key == "" {
		return key
	}
	i.RLock()
	defer i.RUnlock()

	var (
		section, k string
		sectionKey = strings.Split(strings.ToLower(key), "::")
	)
	if len(sectionKey) >= 2 {
		section = sectionKey[0]
		k = sectionKey[1]
	} else {
		section = defaultSection
		k = sectionKey[0]
	}
	if v, ok := i.data[section]; ok {
		if vv, ok := v[k]; ok {
			return vv
		}
	}

	return ""
}
