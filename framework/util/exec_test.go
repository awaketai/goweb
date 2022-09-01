package util

import (
	"testing"
)

func TestGetExecDir(t *testing.T) {
	folder, err := GetExecDir()
	if err != nil {
		t.Errorf("get exec dir err:%v", err)
	}
	absoluteDir := "/Users/admin/www/goweb/framework/util"
	if folder != absoluteDir {
		t.Errorf("get exec folder error")
	}
}
