package http

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/gin"
)

func NewHttpEngine(container framework.Container) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetContainer(container)
	Routes(r)
	return r, nil
}
