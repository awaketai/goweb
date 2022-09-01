package demo

import (
	"fmt"

	"github.com/awaketai/goweb/framework"
)

type DemoServiceProvider struct {
}

func (dsp *DemoServiceProvider) Name() string {
	return Key
}

func (dsp *DemoServiceProvider) Register(container framework.Container) framework.NewInstance {
	return NewDemoService
}

func (dsp *DemoServiceProvider) IsDefer() bool {
	return true

}

func (dsp *DemoServiceProvider) Params(container framework.Container) []any {
	return []any{container}
}

func (dsp *DemoServiceProvider) Boot(container framework.Container) error {
	fmt.Println("demo service boot execute")
	return nil
}
