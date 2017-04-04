package distonic

type Pipeline struct {
	stages []*Stage
}

func (p *Pipeline) Run() (*Result, error) {
	result := &Result{}

	for _, stage := range p.stages {
		stageResult, err := stage.Run()
		if err != nil {
			return result, err
		}
		result = stageResult
	}

	return result, nil
}

func (p *Pipeline) UnmarshalYAML(unmarshal func(interface{}) error) error {
	p.stages = []*Stage{}
	var d map[string]Stage

	if err := unmarshal(&d); err != nil {
		return err
	}

	for name, stage := range d {
		stage.name = name
		p.stages = append(p.stages, &stage)
	}
	return nil
}
