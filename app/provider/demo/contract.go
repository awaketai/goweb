package demo

// 关键字凭证
const DemoKey = "web:demo"

type IService interface {
	GetAllStudent() []Student
}
type Student struct {
	ID   int
	Name string
}
