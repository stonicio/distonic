package docker_build

import (
	"fmt"

	"github.com/stonicio/distonic/module"
)

const (
	ID = "docker_build"
)

func New() *dockerBuild {
	return &dockerBuild{Dockerfile: "Dockerfile"}
}

type dockerBuild struct {
	Dockerfile string
}

func (m *dockerBuild) Call(context *module.Context) (*module.Result, error) {
	description := fmt.Sprintf("Module `%s` called with context: %s", ID, context)
	return &module.Result{Success: true, Description: description}, nil
}
