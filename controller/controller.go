package controller

import (
	"context"
	"fmt"
	"goweb/framework"
	"net/http"
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
		ctx.SetStatus(http.StatusOK).Json("ok")
		finish <- struct{}{}
	}()

	select {
	case p := <-panicCh:
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.SetStatus(http.StatusInternalServerError).Json(p)

	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.SetStatus(http.StatusInternalServerError).Json("time out")
		ctx.SetHasTimeout()
	}

	ctx.SetStatus(http.StatusOK).Json(map[string]any{
		"code": 0,
	})
	return nil
}
