package docker_build

import (
	"fmt"
	"log"

	"github.com/stonicio/distonic/module"
)

const (
	moduleId = "docker_build"
)

type dockerBuild struct {
	dockerfile string
}

func NewDockerBuild() *dockerBuild {
	return &dockerBuild{dockerfile: "Dockerfile"}
}

func (m *dockerBuild) Register() (string, module.Bindable) {
	return moduleId, m
}

func (m *dockerBuild) Bind(params map[string]interface{}) (module.Callable, error) {
	bound := boundDockerBuild(*NewDockerBuild())

	for key, value := range params {
		switch key {
		case "dockerfile":
			bound.dockerfile = value.(string)
		default:
			err := fmt.Errorf("Unknown `%s` field: %s", moduleId, key)
			log.Print(err)
			return nil, err
		}
	}

	return &bound, nil
}

type boundDockerBuild dockerBuild

func (m *boundDockerBuild) Call(context *module.Context) error {
	log.Printf("Module `%s` called with context: %s", moduleId, context)
	return nil
}
