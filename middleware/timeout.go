package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/awaketai/goweb/framework2"
)

func TimeoutHandler(fn framework2.ControllerHandler, d time.Duration) framework2.ControllerHandler {
	return func(c *framework2.Context) error {
		finishCh := make(chan struct{}, 1)
		panicCh := make(chan any, 1)
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.GetRequest().WithContext(durationCtx)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicCh <- p
				}
			}()
			// execution specific business logic
			c.Next()
			finishCh <- struct{}{}
		}()
		select {
		case p := <-panicCh:
			log.Println(p)
		case <-finishCh:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.IsTimeout()
			c.GetResponse().Write([]byte("time out"))
		}

		return nil
	}
}
