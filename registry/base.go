package registry

import (
	"fmt"

	"github.com/stonicio/distonic/module"
)

func Register(id string, m module.Module) error {
	if _, ok := registry[id]; ok {
		return fmt.Errorf("Module `%s` is already registered", id)
	}

	registry[id] = m
	return nil
}

func Get(id string) (module.Module, error) {
	m, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("Could not find `%s` module", id)
	}

	return m, nil
}
