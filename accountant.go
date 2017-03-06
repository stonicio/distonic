package distonic

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Accountant struct {
	repos map[string]*Watcher
}

func NewAccountant() (*Accountant, error) {
	var err error

	a := Accountant{repos: map[string]*Watcher{}}
	reposConfig := viper.Sub("repos")

	for repoName := range viper.GetStringMap("repos") {
		repoSettings := reposConfig.Sub(repoName)
		a.repos[repoName], err = NewWatcher(
			repoName,
			repoSettings.GetString("url"),
			repoSettings.GetStringSlice("branches"))
		if err != nil {
			log.Printf("Cannot create watcher for repo `%s`: %s", repoName, err)
			return nil, err
		}
		log.Printf("Created watcher for repo: %s", repoName)
	}
	return &a, nil
}

func (a *Accountant) Run() error {
	var wg sync.WaitGroup
	var errorCount int

	for name, watcher := range a.repos {
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
