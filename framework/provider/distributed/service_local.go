package distributed

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

type LocalDistributedService struct {
	container framework.Container
}

func NewLocalDistributedService(params ...any) (any, error) {
	if len(params) != 1 {
		return nil, errors.New("local distributed service params error")
	}
	container := params[0].(framework.Container)
	return &LocalDistributedService{container: container}, nil
}

func (service LocalDistributedService) Select(serviceName, appID string, holdTime time.Duration) (string, error) {
	appService := service.container.MustMake(contract.AppKey).(contract.App)
	runtimeFolder := appService.RuntimeFolder()
	lockFile := filepath.Join(runtimeFolder, "distribute_"+serviceName)
	lock, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	// 尝试独占文件锁
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		selectedAppIDByt, err := ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectedAppIDByt), nil
	}

	// 一段时间内，选举有效，其它节点在这段时间内不能再抢占
	go func() {
		defer func() {
			// 释放文件锁
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			// 释放文件
			lock.Close()
			// 删除文件锁对应的文件
			os.Remove(lockFile)
		}()
		timer := time.NewTimer(holdTime)
		<-timer.C
	}()

	// 将抢占到的appID写入文件
	if _, err := lock.WriteString(appID); err != nil {
		return "", err
	}
	return appID, nil
}
