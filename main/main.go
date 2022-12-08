package main

import (
	"log"

	"github.com/awaketai/goweb/app/console"
	"github.com/awaketai/goweb/app/http"
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/provider/app"
	"github.com/awaketai/goweb/framework/provider/config"
	"github.com/awaketai/goweb/framework/provider/distributed"
	"github.com/awaketai/goweb/framework/provider/env"
	"github.com/awaketai/goweb/framework/provider/kernel"
	gwlog "github.com/awaketai/goweb/framework/provider/log"
)

func main() {
	container := framework.NewWebContainer()
	// service bind
	container.Bind(&app.AppProvider{BaseFolder: "./app"})
	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&env.WebEnvProvider{})
	container.Bind(&config.WebConfigProvider{})
	container.Bind(&gwlog.WebLogServiceProvider{})

	engine, err := http.NewHttpEngine(container)
	if err != nil {
		log.Fatalf("start http engine error:%v", err)
	}
	container.Bind(&kernel.WebKernelProvider{HttpEngine: engine})
	// 运行root命令
	console.RunCommand(container)
}
