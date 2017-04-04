package docker_registry

import (
	"log"

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

func (m *dockerRegistry) Call(context *module.Context) error {
	log.Printf("Module `%s` called with context: %s", ID, context)
	return nil
}
