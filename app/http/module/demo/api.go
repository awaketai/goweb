package demo

import (
	"github.com/awaketai/goweb/app/provider/demo"
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
	r.POST("/demo/demoPost", api.DemoPost)
	return nil
}

func NewDemoApi() *DemoApi {
	service := NewService()
	return &DemoApi{service: service}
}

func (api *DemoApi) Demo(c *gin.Context) {
	users := api.service.GetUsers()
	userDTO := UserModelsToUserDTOs(users)
	c.JSON(200, userDTO)
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