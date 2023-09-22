package main

import (
	"log"

	"github.com/awaketai/goweb/framework/provider/cache"
	"github.com/awaketai/goweb/framework/provider/redis"
	"github.com/awaketai/goweb/framework/provider/ssh"

	"github.com/awaketai/goweb/app/console"
	"github.com/awaketai/goweb/app/http"
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/provider/app"
	"github.com/awaketai/goweb/framework/provider/config"
	"github.com/awaketai/goweb/framework/provider/distributed"
	"github.com/awaketai/goweb/framework/provider/env"
	"github.com/awaketai/goweb/framework/provider/kernel"
	gwlog "github.com/awaketai/goweb/framework/provider/log"
	"github.com/awaketai/goweb/framework/provider/orm"
)

func main() {
	container := framework.NewWebContainer()
	// service bind
	var err error
	err = container.Bind(&app.WebAppProvider{})
	err = container.Bind(&distributed.LocalDistributedProvider{})
	err = container.Bind(&env.WebEnvProvider{})
	err = container.Bind(&config.WebConfigProvider{})
	err = container.Bind(&gwlog.WebLogProvider{})
	err = container.Bind(&orm.WebGormProvider{})
	err = container.Bind(&redis.WebRedisProvider{})
	err = container.Bind(&cache.WebCacheProvider{})
	err = container.Bind(&ssh.WebSSHProvider{})
	if err != nil {
		log.Fatalf("bind provider err:%v\n", err)
	}
	engine, err := http.NewHttpEngine(container)
	if err != nil {
		log.Fatalf("start http engine error:%v", err)
	}
	container.Bind(&kernel.WebKernelProvider{HttpEngine: engine})
	// 运行root命令
	console.RunCommand(container)
}
