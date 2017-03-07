package distonic

import (
	"log"
	"time"
)

type Worker struct {
}

func NewWorker() (*Worker, error) {
	return &Worker{}, nil
}

func (w *Worker) Run(jobs <-chan *Job) error {
	for job := range jobs {
		log.Printf("Received job: %s", job)
		time.Sleep(time.Second * 10)
	}
	return nil
}
