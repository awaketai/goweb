package goweb

import (
	"goweb/framework"
)

func RegRouter(core *framework.Core) {
	core.Get("/user/login", func(c *framework.Context) error {
		c.Json(200, "login success")
		return nil
	})
	sub := core.Group("/subject")
	sub.Delete("/:id", func(c *framework.Context) error {
		c.Json(200, "delete")
		return nil
	})
	sub.Put("/:id", func(c *framework.Context) error {
		c.Json(200, "put :id")
		return nil
	})
	sub.Get("/:id", func(c *framework.Context) error {
		c.Json(200, "get :id")
		return nil
	})
	sub.Get("/list/all", func(c *framework.Context) error {
		c.Json(200, "get list all")
		return nil
	})

}
