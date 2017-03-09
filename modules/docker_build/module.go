package docker_build

import (
	"fmt"
	"log"

	"github.com/stonicio/distonic"
)

const (
	moduleId = "docker_build"
)

func init() {
	distonic.RegisterModule(moduleId, NewDockerBuildModule())
}

type DockerBuildModule struct {
	dockerfile string
}

func NewDockerBuildModule() *DockerBuildModule {
	return &DockerBuildModule{
		dockerfile: "Dockerfile"}
}

func (m *DockerBuildModule) Bind(params map[string]interface{}) (distonic.CallableModule, error) {
	bound := BoundDockerBuildModule(*NewDockerBuildModule())

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

type BoundDockerBuildModule DockerBuildModule

func (m *BoundDockerBuildModule) Call(context *distonic.Context) error {
	log.Printf("Module `%s` called with context: %s", moduleId, context)
	return nil
}
