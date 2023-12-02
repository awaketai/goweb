package middleware

import (
	"fmt"

	"github.com/awaketai/goweb/framework2"
)

func Test1() framework2.ControllerHandler {
	return func(c *framework2.Context) error {
		fmt.Println("middleware pre test1")
		c.Next()
		fmt.Println("middleware post test1")
		return nil
	}
}

func Test2() framework2.ControllerHandler {
	return func(c *framework2.Context) error {
		fmt.Println("middleware pre test2")
		c.Next()
		fmt.Println("middleware post test2")
		return nil
	}
}

func Test3() framework2.ControllerHandler {
	return func(c *framework2.Context) error {
		fmt.Println("middleware pre test3")
		c.Next()
		fmt.Println("middleware post test3")
		return nil
	}
}
