package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awaketai/goweb"
	"github.com/awaketai/goweb/framework/gin"
	"github.com/awaketai/goweb/framework/middleware"
	"github.com/awaketai/goweb/app/provider/demo"
)

func main() {
	core := gin.New()
	// bind service
	core.Bind(&demo.DemoServiceProvider{})
	core.Use(gin.Recovery())
	core.Use(middleware.Cost())
	goweb.RegRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}
	go func() {
		server.ListenAndServe()
	}()
	quit := make(chan os.Signal, 1)
	// 监控信号：SIGINT SIGTERM SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server shutdown:", err)
	}

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
