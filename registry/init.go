package registry

import (
	"github.com/stonicio/distonic/module"
	"github.com/stonicio/distonic/module/docker_build"
	"github.com/stonicio/distonic/module/docker_registry"
	"github.com/stonicio/distonic/module/docker_run"
)

var registry map[string]module.Bindable

func init() {
	registry = map[string]module.Bindable{}
	Register(docker_build.NewDockerBuild())
	Register(docker_run.NewDockerRun())
	Register(docker_registry.NewDockerRegistry())
}
