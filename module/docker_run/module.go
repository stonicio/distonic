package docker_run

import (
	"fmt"

	"github.com/stonicio/distonic/module"
)

const (
	ID = "docker_run"
)

func New() *dockerRun {
	return &dockerRun{Cmd: ""}
}

type dockerRun struct {
	Cmd string
}

func (m *dockerRun) Call(context *module.Context) (*module.Result, error) {
	description := fmt.Sprintf("Module `%s` called with context: %s", ID, context)
	return &module.Result{Success: false, Description: description}, nil
}
