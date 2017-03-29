package registry

import (
	"github.com/stonicio/distonic/module"
	"github.com/stonicio/distonic/module/docker_build"
)

var registry map[string]module.Bindable

func init() {
	registry = map[string]module.Bindable{}
	Register(docker_build.NewDockerBuild())
}
