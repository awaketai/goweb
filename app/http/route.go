package http

import (
	"fmt"

	"github.com/awaketai/goweb/app/http/module/demo"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/gin"
	ginSwagger "github.com/awaketai/goweb/framework/middleware/gin-swagger"
	swaggerfiles "github.com/swaggo/files"
)

func Routes(r *gin.Engine) {
	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	// if config swagger
	if configService.GetBool("app.swagger") {
		fmt.Println("swagger...")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	r.Static("/dist", "./dist")
	demo.Register(r)
}
