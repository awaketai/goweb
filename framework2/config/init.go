package config

// from beego config
import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

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
	var comment bytes.Buffer
	buf := bufio.NewReader(bytes.NewReader(data))
	head, err := buf.Peek(3)
	// check UTF-8 byte order mark的特定组合
	// 0xEF 0xBB 0xBF
	if err == nil && head[0] == 239 && head[1] == 187 && head[2] == 191 {
		for i := 1; i <= 3; i++ {
			buf.ReadByte()
		}
	}
	section := defaultSection
	tmpBuf := bytes.NewBuffer(nil)
	for {
		tmpBuf.Reset()
		shouldBreak := false
		for {
			tmp, isPrefix, err := buf.ReadLine()
			if err == io.EOF {
				shouldBreak = true
				break
			}
			var pathError *os.PathError
			if errors.As(err, &pathError) {
				return nil, err
			}
			tmpBuf.Write(tmp)
			if isPrefix {
				continue
			}
			if !isPrefix {
				break
			}

		}
		if shouldBreak {
			break
		}
		line := tmpBuf.Bytes()
		line = bytes.TrimSpace(line)
		if bytes.Equal(line, bEmpty) {
			continue
		}
		var bComment []byte
		switch {
		case bytes.HasPrefix(line, bNumComment):
			bComment = bNumComment
		case bytes.HasPrefix(line, bSemComment):
			bComment = bSemComment

		}
		if bComment != nil {
			line = bytes.TrimLeft(line, string(bComment))
			if comment.Len() > 0 {
				comment.WriteByte('\n')
			}
			comment.Write(line)
			continue
		}
		if bytes.HasPrefix(line, sectionStart) && bytes.HasSuffix(line, sectionEnd) {
			section = strings.ToLower(string(line[1 : len(line)-1]))
			if comment.Len() > 0 {
				cfg.sectionComment[section] = comment.String()
			}
			if _, ok := cfg.data[section]; !ok {
				cfg.data[section] = make(map[string]string)
			}
			continue
		}
		if _, ok := cfg.data[section]; !ok {
			cfg.data[section] = make(map[string]string)
		}
		// key is insensitive
		keyValue := bytes.SplitN(line, bEqual, 2)
		key := string(bytes.TrimSpace(keyValue[0]))
		key = strings.ToLower(key)

		if len(keyValue) == 1 && strings.HasPrefix(key, "include") {
			includeFiles := strings.Fields(key)
			if includeFiles[0] == "include" && len(includeFiles) == 2 {
				otherFile := strings.Trim(includeFiles[1], "\"")
				if !filepath.IsAbs(otherFile) {
					otherFile = filepath.Join(dir, otherFile)
				}
				ifg, err := i.parseFile(otherFile)
				if err != nil {
					return nil, err
				}
				for sec, dt := range ifg.data {
					if _, ok := cfg.data[sec]; !ok {
						cfg.data[section] = make(map[string]string)
					}
					for k, v := range dt {
						cfg.data[sec][k] = v
					}
				}
				for sec, comm := range ifg.sectionComment {
					cfg.sectionComment[sec] = comm
				}
				for k, comm := range ifg.keyComment {
					cfg.keyComment[k] = comm
				}
				continue
			}
		}
		if len(keyValue) != 2 {
			return nil, errors.New("read the content error: \"" + string(line) + "\", should key = val")
		}
		val := bytes.TrimSpace(keyValue[1])
		if bytes.HasPrefix(val, bDQuote) {
			val = bytes.Trim(val, `"`)
		}
		// ExpandValueEnv this method waiting finished
		cfg.data[section][key] = ExpandValueEnv(string(val))
		if comment.Len() > 0 {
			cfg.keyComment[section+"."+key] = comment.String()
			comment.Reset()
		}
	}

	return cfg, nil
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
