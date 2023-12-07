package middleware

import (
	"log"
	"time"

	"github.com/awaketai/goweb/framework2"
)

func Costr() framework2.ControllerHandler {
	return func(c *framework2.Context) error {
		start := time.Now()
		c.Next()
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri:%v,cost: %v",c.GetRequest().RequestURI,cost.Seconds())
		return nil
	}
}
