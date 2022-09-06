package demo

import (
	"fmt"

	"github.com/awaketai/goweb/framework"
)

type DemoService struct {
	// 实现接口
	IService
	// 参数
	container framework.Container
}

func NewDemoService(params ...any) (any, error) {
	container := params[0].(framework.Container)
	fmt.Println("new demo service")
	return &DemoService{container: container}, nil

}

func (service *DemoService) GetAllStudent() []Student {
	return []Student{
		{
			ID:   1,
			Name: "foo",
		},
		{
			ID:   2,
			Name: "bar",
		},
	}
}
