package framework

import "net/http"

type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO

}
