package module

type Registerable interface {
	Register() (string, Bindable)
}

type Bindable interface {
	Bind(params map[string]interface{}) (Callable, error)
}

type Callable interface {
	Call(context *Context) error
}
