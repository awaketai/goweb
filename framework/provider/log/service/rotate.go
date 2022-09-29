package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/util"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
)

// 单个文件输出，但是自动进行切割
type WebRotateLog struct {
	WebLog
	// 日志文件存储目录
	folder string
	// 日志文件名
	file string
}

func NewWebRotateLog(params ...any) (any, error) {
	container := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	appService := container.MustMake(contract.AppKey).(contract.App)
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	folder := appService.LogFolder()
	if configService.IsExists("log.folder") {
		folder = configService.GetString("log.folder")
	}
	exists, err := util.Exists(folder)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	file := "web.log"
	if configService.IsExists("log.file") {
		file = configService.GetString("log.file")
	}
	dateFormat := "%Y%m%d%H"
	if configService.IsExists("log.date_format") {
		dateFormat = configService.GetString("log.date_format")
	}
	linkName := rotatelogs.WithLinkName(filepath.Join(folder, file))
	options := []rotatelogs.Option{linkName}
	if configService.IsExists("log.rotate_count") {
		rotateCount := configService.GetInt("log.rotate_count")
		options = append(options, rotatelogs.WithRotationCount(uint(rotateCount)))
	}

	if configService.IsExists("log.rotate_size") {
		rotateSize := configService.GetInt("log.rotate_size")
		options = append(options, rotatelogs.WithRotationSize(int64(rotateSize)))
	}
	if configService.IsExists("log._max_age") {
		if maxAgeParse, err := time.ParseDuration(configService.GetString("log.max_age")); err == nil {
			options = append(options, rotatelogs.WithMaxAge(maxAgeParse))
		}
	}

	log := &WebRotateLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.folder = folder
	log.file = file

	w, err := rotatelogs.New(fmt.Sprintf("%s.%s", filepath.Join(log.folder, log.file), dateFormat), options...)
	if err != nil {
		return nil, errors.Wrap(err, "new rotatelogs error")
	}

	log.SetOutput(w)
	log.c = container
	return log, nil
}
