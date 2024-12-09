package goredis

import "context"

type Context struct {
	context.Context
	command Value
}

func (ctx *Context) Ctx() context.Context {
	return ctx.Context
}

func (ctx *Context) Command() Value {
	return ctx.command
}

func (ctx *Context) Args() []Value {
	return ctx.command.array[1:]
}
