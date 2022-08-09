package main

import (
	"goweb"
	"goweb/controller"
	"goweb/framework"
	"goweb/framework/middleware"
	"log"
	"net/http"
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
	serve.ListenAndServe()
}
