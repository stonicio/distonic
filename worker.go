package distonic

import (
	"log"
)

type Worker struct {
}

func NewWorker() (*Worker, error) {
	return &Worker{}, nil
}

func (w *Worker) Run() error {
	log.Println(w)
	return nil
}
