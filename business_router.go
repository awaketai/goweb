package goweb

import (
	"github.com/awaketai/goweb/framework2"
	"github.com/awaketai/goweb/middleware"
)

func registerRouter(c *framework2.Core) {
	c.Get("/user/login", middleware.Test3(), UserLoginController)
}
