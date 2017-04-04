package module

type Module interface {
	Call(context *Context) error
}
