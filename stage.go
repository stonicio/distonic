package distonic

import (
	"fmt"
	"strings"
)

type Stage struct {
	name string
	jobs []*Job
}

func (s *Stage) Run() (*Result, error) {
	result := &Result{}
	errorJobs := []string{}
	failedJobs := []string{}

	for _, job := range s.jobs {
		go func() {
			jobResult, err := job.Run()
			if err != nil {
				errorJobs = append(errorJobs, job.name)
			} else if !jobResult.success {
				failedJobs = append(failedJobs, job.name)
			}
		}()
	}

	if len(errorJobs) > 0 {
		return result, fmt.Errorf(
			"Error jobs: %s", strings.Join(errorJobs, ", "))
	}

	if len(failedJobs) > 0 {
		result.description = fmt.Sprintf(
			"Failed jobs: %s", strings.Join(failedJobs, ", "))
		return result, nil
	}

	result.success = true

	return result, nil
}

func (s *Stage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s.jobs = []*Job{}
	var d []Job

	if err := unmarshal(&d); err != nil {
		return err
	}

	for _, job := range d {
		s.jobs = append(s.jobs, &job)
	}
	return nil
}
