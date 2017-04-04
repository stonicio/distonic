package docker_registry

import (
	"fmt"
	"log"

	"github.com/stonicio/distonic/module"
)

const (
	moduleId = "docker_registry"
)

type dockerRegistry struct {
	repo string
	tags []string
}

func NewDockerRegistry() *dockerRegistry {
	return &dockerRegistry{tags: []string{}}
}

func (m *dockerRegistry) Register() (string, module.Bindable) {
	return moduleId, m
}

func (m *dockerRegistry) Bind(params map[string]interface{}) (module.Callable, error) {
	bound := boundDockerRegistry(*NewDockerRegistry())

	for key, value := range params {
		switch key {
		case "repo":
			bound.repo = value.(string)
		case "tags":
			for _, v := range value.([]interface{}) {
				bound.tags = append(bound.tags, v.(string))
			}
		default:
			err := fmt.Errorf("Unknown `%s` field: %s", moduleId, key)
			log.Print(err)
			return nil, err
		}
	}

	return &bound, nil
}

type boundDockerRegistry dockerRegistry

func (m *boundDockerRegistry) Call(context *module.Context) error {
	log.Printf("Module `%s` called with context: %s", moduleId, context)
	return nil
}
