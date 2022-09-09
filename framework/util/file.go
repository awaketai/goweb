package util

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Exists 目录是否存在
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			// 文件不存在err
			return true, nil
		}
		return false, err
	}
	return true, nil
}

// IsHiddenDir 路径是否隐藏路径
func IsHiddenDir(path string) bool {
	return len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")
}

// SubDir 输出所有子目录
func SubDir(folder string) ([]string, error) {
	subs, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, sub := range subs {
		if sub.IsDir() {
			ret = append(ret, sub.Name())
		}
	}
	return ret, nil
}

// DownloadFile
func DownloadFile(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
