package goweb

import (
	"context"
)

type ControllerHandler func(c *context.Context) error
