package middleware

import (
	"goweb/framework"
	"net/http"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.SetStatus(http.StatusInternalServerError).Json(err)
			}
		}()
		c.Next()
		return nil
	}

}
