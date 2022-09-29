package service

import (
	"os"
	"path/filepath"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/util"
	"github.com/pkg/errors"
)

// 单个文件输出
type WebSingleLog struct {
	WebLog
	folder string
	file   string
	fd     *os.File
}

func NewWebSingleLog(params ...any) (any, error) {
	container := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	appService := container.MustMake(contract.AppKey).(contract.App)
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	folder := appService.LogFolder()

	log := &WebSingleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	if configService.IsExists("log.folder") {
		folder = configService.GetString("log.folder")
	}
	log.folder = folder
	exists, err := util.Exists(folder)
	if err != nil {
		return nil, err
	}
	if !exists {
		os.MkdirAll(folder, os.ModePerm)
	}
	log.file = "web.log"
	if configService.IsExists("log.file") {
		log.file = configService.GetString("log.file")
	}
	fd, err := os.OpenFile(filepath.Join(log.folder, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open log file err")
	}
	log.SetOutput(fd)
	log.c = container
	return log, nil
}
