package distonic

import (
	"container/list"
	"fmt"
	"log"
	"sync"

	git "github.com/libgit2/git2go"
	"github.com/spf13/viper"
)

type Supervisor struct {
	repos   map[string]*Watcher
	workers []*Worker
	queue   *list.List
}

type Job struct {
	repo       *git.Repository
	ref        *git.Reference
	branchName string
}

func NewSupervisor() (*Supervisor, error) {
	var err error
	reposConfig := viper.Sub("repos")
	numWorkers := viper.GetInt("num_workers")

	s := &Supervisor{
		repos:   map[string]*Watcher{},
		workers: []*Worker{},
		queue:   list.New()}

	for repoName := range viper.GetStringMap("repos") {
		repoSettings := reposConfig.Sub(repoName)
		s.repos[repoName], err = NewWatcher(
			repoName,
			repoSettings.GetString("url"),
			repoSettings.GetStringSlice("branches"))
		if err != nil {
			log.Printf("Cannot create watcher for repo `%s`: %s", repoName, err)
			return nil, err
		}
		log.Printf("Created watcher for repo: %s", repoName)
	}

	for n := 0; n < numWorkers; n++ {
		w, err := NewWorker()
		if err != nil {
			log.Printf("Cannot create worker #%s: %s", n, err)
			return nil, err
		}
		s.workers = append(s.workers, w)
	}

	return s, nil
}

func (s *Supervisor) Run() error {
	var wg sync.WaitGroup
	var errWorkers error
	var errWatchers error

	wg.Add(1)
	go func() {
		errWorkers = s.runWorkers()
		if errWorkers != nil {
			log.Print(errWorkers)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errWatchers = s.runWatchers()
		if errWatchers != nil {
			log.Print(errWatchers)
		}
		wg.Done()
	}()

	wg.Wait()
	if errWorkers != nil || errWatchers != nil {
		return fmt.Errorf("Supervisor exited with errors")
	}
	return nil
}

func (s *Supervisor) runWorkers() error {
	var wg sync.WaitGroup
	var errorCount int

	for n, worker := range s.workers {
		wg.Add(1)
		go func() {
			err := worker.Run()
			if err != nil {
				log.Printf("Error in worker #%s: %s", n, err)
				errorCount += 1
			}
			wg.Done()
		}()
	}

	wg.Wait()
	if errorCount > 0 {
		return fmt.Errorf("There was %s errors running workers", errorCount)
	}
	return nil
}

func (s *Supervisor) runWatchers() error {
	var wg sync.WaitGroup
	var errorCount int

	for name, watcher := range s.repos {
		wg.Add(1)
		go func() {
			err := watcher.Run()
			if err != nil {
				log.Printf("Error in watcher for repo %s: %s", name, err)
				errorCount += 1
			}
			wg.Done()
		}()
	}

	wg.Wait()
	if errorCount > 0 {
		return fmt.Errorf("There was %s errors running watchers", errorCount)
	}
	return nil
}

func (s *Supervisor) schedule(job *Job) error {
	return nil
}
