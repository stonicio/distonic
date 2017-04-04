package docker_run

import (
	"log"

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

func (m *dockerRun) Call(context *module.Context) error {
	log.Printf("Module `%s` called with context: %s", ID, context)
	return nil
}
