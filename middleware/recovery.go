package middleware

import "github.com/awaketai/goweb/framework2"

// Recovery capture exception
func Recovery() framework2.ControllerHandler {
	return func(c *framework2.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(500, err)
			}
		}()
		c.Next()
		
		return nil
	}
}
