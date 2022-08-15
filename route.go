package goweb

import (
	"goweb/framework"
	"goweb/framework/middleware"
	"net/http"
)

func RegRouter(core *framework.Core) {
	core.Get("/user/login", middleware.Test1(), func(c *framework.Context) error {
		c.SetStatus(http.StatusOK).Json("login success")
		return nil
	})
	sub := core.Group("/subject")
	sub.Delete("/:id", func(c *framework.Context) error {
		c.SetStatus(http.StatusOK).Json("delete")
		return nil
	})
	sub.Put("/:id", func(c *framework.Context) error {
		c.SetStatus(http.StatusOK).Json("put :id")
		return nil
	})
	sub.Get("/:id", middleware.Test2(), func(c *framework.Context) error {
		c.SetStatus(http.StatusOK).Json("get :id")
		return nil
	})
	sub.Get("/list/all", func(c *framework.Context) error {
		c.SetStatus(http.StatusOK).Json("get list all")
		return nil
	})

}
