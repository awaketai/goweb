package http

import (
	"github.com/awaketai/goweb/app/http/module/demo"
	"github.com/awaketai/goweb/framework/gin"
)

func Routes(r *gin.Engine) {
	r.Static("/dist", "./dist")
	demo.Register(r)
}
