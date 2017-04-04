package module

type Module interface {
	Call(context *Context) (*Result, error)
}

type Context struct {
	Workdir      string
	Branch       string
	BranchDashed string
	Commit       string
}

type Result struct {
	Success     bool
	Description string
}
