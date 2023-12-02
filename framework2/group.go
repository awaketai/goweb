package framework2

type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)

	// Group nested group implement
	Group(string) IGroup
	// Use nested middleware
	Use(middlewares ...ControllerHandler)
}

type Group struct {
	core   *Core
	prefix string
	// group general prefix
	parent *Group
	middlewares []ControllerHandler
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

func (g *Group) Get(uri string, handler ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allhandlers := append(g.getMiddlewares(),handler...)
	g.core.Get(uri, allhandlers...)
}

func (g *Group) Post(uri string, handler ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allhandlers := append(g.getMiddlewares(),handler...)
	g.core.Post(uri, allhandlers...)
}

func (g *Group) Put(uri string, handler ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allhandlers := append(g.getMiddlewares(),handler...)
	g.core.Put(uri, allhandlers...)
}

func (g *Group) Delete(uri string, handler ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allhandlers := append(g.getMiddlewares(),handler...)
	g.core.Delete(uri, allhandlers...)
}

func (g *Group) Group(uri string) IGroup {
	cgroup := NewGroup(g.core, uri)
	cgroup.parent = g
	return cgroup
}

func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}

	return g.parent.getAbsolutePrefix() + g.prefix
}

// getMiddlewares get the middleware for a group
// except set for Get/Post/Put/Delete method
func (g *Group) getMiddlewares() []ControllerHandler{
	if g.parent == nil {
		return g.middlewares
	}

	return append(g.parent.getMiddlewares(),g.middlewares...)
}

func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

