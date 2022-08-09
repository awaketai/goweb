package framework

import (
	"context"
	"log"
	"time"
)

func Timeout(f ControllerHandler, d time.Duration) ControllerHandler {
	return func(c *Context) error {
		finishCh := make(chan struct{}, 1)
		panicCh := make(chan any, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.request.WithContext(durationCtx)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					panicCh <- e
				}
			}()
			// specific business
			f(c)
			finishCh <- struct{}{}
		}()
		select {
		case p := <-panicCh:
			log.Println(p)
		case <-finishCh:
			c.responseWriter.WriteHeader(500)
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.responseWriter.Write([]byte("time out"))
		}

		return nil
	}
}
