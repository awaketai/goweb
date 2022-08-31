package demo

import (
	"fmt"

	"github.com/awaketai/goweb/framework"
)

type DemoService struct {
	// 实现接口
	Service
	// 参数
	container framework.Container
}

func NewDemoService(params ...any) (any, error) {
	container := params[0].(framework.Container)
	fmt.Println("new demo service")
	return &DemoService{container: container}, nil

}

func (service *DemoService) GetFoo() Foo {
	return Foo{
		Name: "i am foo",
	}
}
