package docker_build

import (
	"log"

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

func (m *dockerBuild) Call(context *module.Context) error {
	log.Printf("Module `%s` called with context: %s", ID, context)
	return nil
}
