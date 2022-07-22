package main

import (
	"goweb/framework"
	"net/http"
)

func main() {
	serve := &http.Server{
		Handler: framework.NewCore(),
		Addr:    ":8080",
	}

	serve.ListenAndServe()
}
