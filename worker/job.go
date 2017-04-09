package worker

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/stonicio/distonic/module"
	"github.com/stonicio/distonic/registry"
)

type Job struct {
	name   string
	module module.Module
}

func (j *Job) Run() (*module.Result, error) {
	r, err := j.module.Call(&module.Context{})
	return r, err
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

			if err := mapstructure.Decode(v, m); err != nil {
				return err
			}

			j.module = m
		}
	}

	return nil
}
