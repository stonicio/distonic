package registry

import (
	"github.com/stonicio/distonic/module"
	"github.com/stonicio/distonic/module/docker_build"
	"github.com/stonicio/distonic/module/docker_registry"
	"github.com/stonicio/distonic/module/docker_run"
)

var registry map[string]module.Module

func init() {
	registry = map[string]module.Module{}
	Register(docker_build.ID, docker_build.New())
	Register(docker_run.ID, docker_run.New())
	Register(docker_registry.ID, docker_registry.New())
}
