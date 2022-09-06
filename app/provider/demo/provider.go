package demo

import (
	"fmt"

	"github.com/awaketai/goweb/framework"
)

type DemoServiceProvider struct {
	framework.ServiceProvider
	c framework.Container
}

func (dsp *DemoServiceProvider) Name() string {
	return DemoKey
}

func (dsp *DemoServiceProvider) Register(container framework.Container) framework.NewInstance {
	return NewDemoService
}

func (dsp *DemoServiceProvider) IsDefer() bool {
	return true

}

func (dsp *DemoServiceProvider) Params(container framework.Container) []any {
	return []any{dsp.c}
}

func (dsp *DemoServiceProvider) Boot(container framework.Container) error {
	dsp.c = container
	fmt.Println("demo service boot execute")
	return nil
}
