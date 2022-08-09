package middleware

import "goweb/framework"

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.Json(500, err)
			}
		}()
		c.Next()
		return nil
	}

}
