package demo

import (
	"github.com/awaketai/goweb/app/provider/demo"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/gin"
)

type DemoApi struct {
	service *Service
}

func Register(r *gin.Engine) error {
	api := NewDemoApi()
	r.Bind(&demo.DemoServiceProvider{})
	r.GET("/demo/demo", api.Demo)
	r.GET("/demo/demo2", api.Demo2)
	r.GET("/demo/pwd", api.Pwd)
	r.POST("/demo/demoPost", api.DemoPost)

	return nil
}

func NewDemoApi() *DemoApi {
	service := NewService()
	return &DemoApi{service: service}
}

func (api *DemoApi) Demo(c *gin.Context) {
	log := c.MustMake(contract.LogKey).(contract.Log)
	m := make(map[string]any)
	m["f1"] = "f1val"
	m["f2"] = 23
	log.Info(c, "demo/demo", m)
	users := api.service.GetUsers()
	userDTO := UserModelsToUserDTOs(users)
	c.JSON(200, userDTO)
}

func (api *DemoApi) Pwd(c *gin.Context) {
	service := c.MustMake(contract.ConfigKey).(contract.Config)
	password := service.GetString("database.mysql.password")
	c.JSON(200, password)
}

func (apo *DemoApi) Demo2(c *gin.Context) {
	provider := c.MustMake(demo.DemoKey).(demo.IService)
	students := provider.GetAllStudent()
	usersDTO := StudentToUserDTOs(students)
	c.JSON(200, usersDTO)
}

func (api *DemoApi) DemoPost(c *gin.Context) {
	type Foo struct {
		Name string
	}

	foo := &Foo{}
	err := c.BindJSON(&foo)
	if err != nil {
		c.AbortWithError(500, err)
	}
	c.JSON(200, nil)
}
