package util

import (
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Exists 目录是否存在
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
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

// CopyFile copy file to dest dir
func CopyFile(source, dest string) error {
	var data, err = os.ReadFile(source)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, data, 0777)
}

func CopyFolder(source, dest string) error {
	err := filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		relPath := strings.Replace(path,source,"",1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(dest,relPath),0755)
		}else{
			data,err := os.ReadFile(filepath.Join(source,relPath))
			if err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(dest,relPath),data,0777)
		}
	})
	
	return err
}
