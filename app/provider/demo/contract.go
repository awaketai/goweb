package demo

// 关键字凭证
const Key = "web:demo"

type Service interface {
	GetFoo() Foo
}

type Foo struct {
	Name string
}
