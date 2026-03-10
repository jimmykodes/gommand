package gommand

import "io"

type ExecutionOption interface {
	Apply(ctx *Context)
}

type ExecutionOptionFunc func(ctx *Context)

func (f ExecutionOptionFunc) Apply(ctx *Context) {
	f(ctx)
}

func WithStdin(stdin io.Reader) ExecutionOptionFunc {
	return func(ctx *Context) {
		ctx.stdin = stdin
	}
}

func WithStdout(stdout io.Writer) ExecutionOptionFunc {
	return func(ctx *Context) {
		ctx.stdout = stdout
	}
}

func WithStderr(stderr io.Writer) ExecutionOptionFunc {
	return func(ctx *Context) {
		ctx.stderr = stderr
	}
}
