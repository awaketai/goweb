package main

import (
	"context"
	"goweb"
	"goweb/controller"
	"goweb/framework"
	"goweb/framework/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	core := framework.NewCore()
	core.Use(middleware.Recovery())
	goweb.RegRouter(core)
	core.RegisterRouter("foo", controller.FooControllerHandler)
	serve := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}
	log.Printf("listen port:%s\n", serve.Addr)

	go func() {
		serve.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serve.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("server exiting")
}
