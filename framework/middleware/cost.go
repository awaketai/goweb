package middleware

import (
	"log"
	"time"

	"github.com/awaketai/goweb/framework/gin"
)

func Cost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		log.Printf("api uri start:%v", ctx.Request.RequestURI)
		ctx.Next()
		log.Printf("api uri end:%v,cost:%v", ctx.Request.RequestURI, time.Since(start))
	}
}
