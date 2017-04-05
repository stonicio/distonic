package distonic

import (
	"github.com/stonicio/distonic/module"
)

type Pipeline struct {
	stages []*Stage
}

func (p *Pipeline) Run() (*module.Result, error) {
	result := &module.Result{}

	for _, stage := range p.stages {
		stageResult, err := stage.Run()
		if err != nil {
			return result, err
		}
		*result = *stageResult
		if !result.Success {
			break
		}
	}

	return result, nil
}

func (p *Pipeline) UnmarshalYAML(unmarshal func(interface{}) error) error {
	p.stages = []*Stage{}
	var d map[string]Stage

	if err := unmarshal(&d); err != nil {
		return err
	}

	for name, s := range d {
		p.stages = append(p.stages, &Stage{name: name, jobs: s.jobs})
	}

	return nil
}
