package worker

import (
	"fmt"
	"strings"
	"sync"

	"github.com/stonicio/distonic/module"
)

type Stage struct {
	name string
	jobs []*Job
}

func (s *Stage) Run() (*module.Result, error) {
	var wg sync.WaitGroup
	result := &module.Result{}
	errorJobs := []string{}
	failedJobs := []string{}

	for _, job := range s.jobs {
		wg.Add(1)
		go func() {
			jobResult, err := job.Run()
			if err != nil {
				errorJobs = append(errorJobs, job.name)
			} else if !jobResult.Success {
				failedJobs = append(failedJobs, job.name)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(errorJobs) > 0 {
		return result, fmt.Errorf(
			"Error jobs: %s", strings.Join(errorJobs, ", "))
	}

	if len(failedJobs) > 0 {
		result.Description = fmt.Sprintf(
			"Failed jobs: %s", strings.Join(failedJobs, ", "))
		return result, nil
	}

	result.Success = true

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
