package goweb

import "github.com/awaketai/goweb/framework2"

func UserLoginController(c *framework2.Context) error {
	c.JSON(200,"ok,UserLoginController")
	return nil
}