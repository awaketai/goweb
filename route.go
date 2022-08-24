package goweb

import (
	"net/http"
	"time"

	"github.com/awaketai/goweb/framework/gin"
	"github.com/awaketai/goweb/framework/middleware"
)

func RegRouter(core *gin.Engine) {
	core.GET("/user/login", middleware.Test1(), func(c *gin.Context) {
		c.ISetStatus(http.StatusOK).IJson("login success")
	})
	core.GET("/shutdown", func(c *gin.Context) {
		foo, _ := c.DefaultQueryString("foo", "def")
		time.Sleep(10 * time.Second)
		c.ISetOkStatus().IJson("ok,UserloginController:" + foo)
	})

	sub := core.Group("/subject")
	sub.DELETE("/:id", func(c *gin.Context) {
		id, ret := c.Params.Get("id")
		if ret {
			c.ISetStatus(http.StatusOK).IJson("delete-" + id)
		} else {
			c.ISetStatus(http.StatusOK).IJson("delete-" + id)
		}
	})
	sub.PUT("/:id", func(c *gin.Context) {
		c.ISetStatus(http.StatusOK).IJson("put :id")
	})
	sub.GET("/:id", middleware.Test2(), func(c *gin.Context) {
		id, ret := c.Params.Get("id")
		if ret {
			c.ISetStatus(http.StatusOK).IJson("get :id" + id)
		} else {
			c.ISetStatus(http.StatusOK).IJson("get-" + id)
		}
	})
	sub.GET("/list/all", func(c *gin.Context) {
		c.ISetStatus(http.StatusOK).IJson("get list all")
	})

}
