package middlewares

import (
	"go-koa/contexts"
)

type Next func() error
type Middleware func(ctx *contexts.Context, next Next) error
