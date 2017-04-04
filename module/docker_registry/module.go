package docker_registry

import (
	"fmt"

	"github.com/stonicio/distonic/module"
)

const (
	ID = "docker_registry"
)

func New() *dockerRegistry {
	return &dockerRegistry{Tags: []string{}}
}

type dockerRegistry struct {
	Repo string
	Tags []string
}

func (m *dockerRegistry) Call(context *module.Context) (*module.Result, error) {
	description := fmt.Sprintf("Module `%s` called with context: %s", ID, context)
	return &module.Result{Success: true, Description: description}, nil
}
