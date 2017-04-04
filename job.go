package distonic

import (
	"fmt"

	"github.com/stonicio/distonic/module"
	"github.com/stonicio/distonic/registry"
)

type Result struct {
	success     bool
	description string
}

type Job struct {
	name   string
	module module.Callable
}

func (j *Job) Run() (*Result, error) {
	return &Result{}, nil
}

func (j *Job) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var d map[string]interface{}

	if err := unmarshal(&d); err != nil {
		return err
	}

	for k, v := range d {
		switch k {
		case "name":
			j.name = v.(string)
		default:
			if j.module != nil {
				return fmt.Errorf("Extra module `%s` defined in job", k)
			}
			m, err := registry.Get(k)
			if err != nil {
				return err
			}
			params := map[string]interface{}{}
			if v != nil {
				for vk, vv := range v.(map[interface{}]interface{}) {
					params[vk.(string)] = vv
				}
			}
			b, err := m.Bind(params)
			if err != nil {
				return err
			}
			j.module = b
		}
	}

	return nil
}
