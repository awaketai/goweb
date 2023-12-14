package goweb

import (
	"github.com/awaketai/goweb/framework2"
	"github.com/awaketai/goweb/middleware"
)

func registerRouter(c *framework2.Core) {
	c.Get("/user/login", middleware.Test3(), UserLoginController)

	// batch common
	subjectApi := c.Group("/subject")
	{
		subjectApi.Use(middleware.Test1())
		subjectApi.Get("/", SubjectListController)
		subjectApi.Post("/", SubjectAddController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Put("/:id", SubjectUpdateController)
		subjectApi.Delete("/:id", SubjectDeleteController)
		subjectApi.Get("/list/all", SubjectListController)

		sbjectInnerApi := subjectApi.Group("/inner")
		{
			sbjectInnerApi.Get("/name",SubjectNameController)
		}
	}
}
