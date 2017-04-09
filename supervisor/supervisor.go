package supervisor

import (
	"container/list"
	"log"
	"sync"

	"github.com/spf13/viper"
	"github.com/stonicio/distonic/watcher"
	"github.com/stonicio/distonic/worker"
)

type Supervisor struct {
	repos   map[string]*watcher.Watcher
	workers []*worker.Worker
	queue   *list.List
	bell    chan bool
}

func NewSupervisor() (*Supervisor, error) {
	var err error
	reposConfig := viper.Sub("repos")
	numWorkers := viper.GetInt("worker.concurrency")

	s := &Supervisor{
		repos:   map[string]*watcher.Watcher{},
		workers: []*worker.Worker{},
		queue:   list.New(),
		bell:    make(chan bool, 1)}

	for repoName := range viper.GetStringMap("repos") {
		repoSettings := reposConfig.Sub(repoName)
		s.repos[repoName], err = watcher.NewWatcher(
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
		w, err := worker.NewWorker()
		if err != nil {
			log.Printf("Cannot create worker #%s: %s", n, err)
			return nil, err
		}
		s.workers = append(s.workers, w)
	}

	return s, nil
}

func (s *Supervisor) Run() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		s.runWorkers()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.runWatchers()
		wg.Done()
	}()

	wg.Wait()
}

func (s *Supervisor) runWatchers() {
	orders := make(chan *watcher.Order, len(s.repos))

	for name, watcher := range s.repos {
		go func() {
			err := watcher.Run(orders)
			if err != nil {
				log.Printf("Error in watcher for repo %s: %s", name, err)
			}
		}()
	}

	for order := range orders {
		s.schedule(order)
	}
}

func (s *Supervisor) runWorkers() {
	orders := make(chan *watcher.Order, len(s.workers))

	for _, worker := range s.workers {
		go func() {
			worker.Run(orders)
		}()
	}

	for _ = range s.bell {
		log.Printf("Bell rings")
		for s.queue.Len() > 0 {
			order := s.queue.Remove(s.queue.Front())
			orders <- order.(*watcher.Order)
		}
	}
}

func (s *Supervisor) schedule(order *watcher.Order) error {
	s.queue.PushBack(order)

	select {
	case <-s.bell:
		log.Print("Silenced the bell")
	default:
		log.Print("Bell is idle")
	}
	s.bell <- true
	log.Print("Rang the bell")
	return nil
}
