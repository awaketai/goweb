package goweb

import (
	"context"
	"time"
)

func TimeoutHandler(fn ControllerHandler, d time.Duration) ControllerHandler {
	//return func(c *Context) error {
	//	finishCh := make(chan struct{}, 1)
	//	panicCh := make(chan any, 1)
	//	durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
	//	defer cancel()
	//	c.request.WithContext(durationCtx)
	//	go func() {
	//		defer func() {
	//			if p := recover(); p != nil {
	//				panicCh <- p
	//			}
	//		}()
	//		// execution
	//		fn(c)
	//		finishCh <- struct{}{}
	//	}()
	//	select {
	//	case p := <-panicCh:
	//		log.Println(p)
	//	case <-durationCtx.Done():
	//		c.IsTimeout()
	//		c.responseWriter.Write([]byte("time out"))
	//
	//	}
	//	return nil
	//}
	return func(c context.Context) error {
		return nil
	}
}
