package app

import (
	"testing"

	"github.com/awaketai/goweb/framework"
)

func TestBaseFolder(t *testing.T) {
	container := framework.NewWebContainer()
	// 1. baseFolder empty
	params := make([]any, 2)
	params[0] = container
	params[1] = ""
	appIns, err := NewApp(params...)
	if err != nil {
		t.Errorf("get app instance err:%v", err)
	}
	ins := appIns.(App)
	folder := ins.BaseFolder()
	if folder != "/Users/admin/www/goweb/framework/provider/app" {
		t.Errorf("not expected folder")
	}
}
