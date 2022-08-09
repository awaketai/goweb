package middleware

import "goweb/framework"

func RequestTime() framework.ControllerHandler {
	return func(c *framework.Context) error {
		// 统计请求时长

		return nil
	}
}
