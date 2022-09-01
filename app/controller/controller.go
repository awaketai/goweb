package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/awaketai/goweb/framework/gin"
)

func FooControllerHandler(ctx *gin.Context) error {
	durationCtx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Second))
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
		ctx.ISetStatus(http.StatusOK).IJson("ok")
		finish <- struct{}{}
	}()

	select {
	case p := <-panicCh:
		ctx.ISetStatus(http.StatusInternalServerError).IJson(p)

	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		ctx.ISetStatus(http.StatusInternalServerError).IJson("time out")
	}

	ctx.ISetStatus(http.StatusOK).IJson(map[string]any{
		"code": 0,
	})
	return nil
}
