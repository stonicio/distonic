package distonic

import (
	"fmt"
)

var modulesRegistry map[string]BindableModule

func init() {
	modulesRegistry = map[string]BindableModule{}
}

func RegisterModule(id string, module BindableModule) error {
	if _, ok := modulesRegistry[id]; ok {
		return fmt.Errorf("Module `%s` is already registered", id)
	}

	modulesRegistry[id] = module

	return nil
}

func getModule(id string) (BindableModule, error) {
	module, ok := modulesRegistry[id]
	if !ok {
		return nil, fmt.Errorf("Could not find `%s` module", id)
	}

	return module, nil
}
