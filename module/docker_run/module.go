package docker_run

import (
	"fmt"
	"log"

	"github.com/stonicio/distonic/module"
)

const (
	moduleId = "docker_run"
)

type dockerRun struct {
	cmd string
}

func NewDockerRun() *dockerRun {
	return &dockerRun{cmd: ""}
}

func (m *dockerRun) Register() (string, module.Bindable) {
	return moduleId, m
}

func (m *dockerRun) Bind(params map[string]interface{}) (module.Callable, error) {
	bound := boundDockerRun(*NewDockerRun())

	for key, value := range params {
		switch key {
		case "cmd":
			bound.cmd = value.(string)
		default:
			err := fmt.Errorf("Unknown `%s` field: %s", moduleId, key)
			log.Print(err)
			return nil, err
		}
	}

	return &bound, nil
}

type boundDockerRun dockerRun

func (m *boundDockerRun) Call(context *module.Context) error {
	log.Printf("Module `%s` called with context: %s", moduleId, context)
	return nil
}
