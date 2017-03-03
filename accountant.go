package distonic

import (
	"fmt"
	"github.com/spf13/viper"
)

type Accountant struct {
	repos map[string]*Watcher
}

func NewAccountant() *Accountant {
	a := Accountant{repos: map[string]*Watcher{}}
	reposConfig := viper.Sub("repos")
	for repoName := range viper.GetStringMap("repos") {
		repoSettings := reposConfig.GetStringMapString(repoName)
		a.repos[repoName] = NewWatcher(repoName, repoSettings["url"])
	}
	return &a
}

func (a *Accountant) Run() error {
	fmt.Println(a)
	return nil
}
