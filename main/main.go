package main

import (
	"fmt"
	"log"

	"github.com/awaketai/goweb/app/console"
	"github.com/awaketai/goweb/app/http"
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/provider/app"
	"github.com/awaketai/goweb/framework/provider/distributed"
	"github.com/awaketai/goweb/framework/provider/env"
	"github.com/awaketai/goweb/framework/provider/kernel"
)

func main() {
	container := framework.NewWebContainer()
	// service bind
	container.Bind(&app.AppProvider{})
	container.Bind(&distributed.LocalDistributedProvider{})
	err := container.Bind(&env.WebEnvProvider{})
	fmt.Println("e:", err)
	engine, err := http.NewHttpEngine(container)
	if err != nil {
		log.Fatalf("start http engine error:%v", err)
	}
	container.Bind(&kernel.WebKernelProvider{HttpEngine: engine})
	// 运行root命令
	console.RunCommand(container)
	// core := gin.New()
	// // bind service
	// core.Bind(&demo.DemoServiceProvider{})
	// core.Use(gin.Recovery())
	// core.Use(middleware.Cost())
	// goweb.RegRouter(core)
	// server := &http.Server{
	// 	Handler: core,
	// 	Addr:    ":8080",
	// }
	// go func() {
	// 	server.ListenAndServe()
	// }()
	// quit := make(chan os.Signal, 1)
	// // 监控信号：SIGINT SIGTERM SIGQUIT
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// <-quit

	// timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := server.Shutdown(timeoutCtx); err != nil {
	// 	log.Fatal("Server shutdown:", err)
	// }

	// core := framework.NewCore()
	// core.Use(middleware.Recovery())
	// goweb.RegRouter(core)
	// core.RegisterRouter("foo", controller.FooControllerHandler)
	// serve := &http.Server{
	// 	Handler: core,
	// 	Addr:    ":8080",
	// }
	// log.Printf("listen port:%s\n", serve.Addr)

	// go func() {
	// 	serve.ListenAndServe()
	// }()

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// <-quit
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := serve.Shutdown(ctx); err != nil {
	// 	log.Fatal("Server Shutdown:", err)
	// }

	// log.Println("server exiting")
}
