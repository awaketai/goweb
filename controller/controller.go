package controller

import (
	"context"
	"fmt"
	"goweb/framework"
	"time"
)

func FooControllerHandler(ctx *framework.Context) error {
	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), time.Duration(1*time.Second))
	defer cancel()
	finish := make(chan struct{}, 1)
	panicCh := make(chan any, 1)

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicCh <- p
			}
		}()
		// do sth
		time.Sleep(10 * time.Second)
		ctx.Json(200, "ok")
		finish <- struct{}{}
	}()

	select {
	case p := <-panicCh:
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.Json(500, p)
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.Json(500, "time out")
		ctx.SetHasTimeout()
	}

	return ctx.Json(200, map[string]any{
		"code": 0,
	})
}
