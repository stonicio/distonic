package registry

import (
	"fmt"

	"github.com/stonicio/distonic/module"
)

func Register(reg module.Registerable) error {
	id, module := reg.Register()
	if _, ok := registry[id]; ok {
		return fmt.Errorf("Module `%s` is already registered", id)
	}

	registry[id] = module
	return nil
}

func Get(id string) (module.Bindable, error) {
	module, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("Could not find `%s` module", id)
	}

	return module, nil
}
