package distonic

import (
	"log"
	"path"
	"path/filepath"
	"time"

	git "github.com/libgit2/git2go"
	"github.com/spf13/viper"
)

type Watcher struct {
	name        string
	dir         string
	url         string
	repo        *git.Repository
	branchSpecs []string
	branchRefs  map[string]*git.Reference
}

func NewWatcher(name, url string, branchSpecs []string) (*Watcher, error) {
	dataDir := viper.GetString("data_dir")
	dir := path.Join(dataDir, "watcher", name)

	repo, err := git.InitRepository(dir, true)
	if err != nil {
		log.Printf("Cannot init repo `%s`: %s", name, err)
		return nil, err
	}

	config, err := repo.Config()
	if err != nil {
		log.Printf("Cannot open config for repo `%s`: %s", name, err)
		return nil, err
	}

	if err := config.SetBool("remote.origin.mirror", true); err != nil {
		log.Printf("Cannot set mirror mode for repo `%s`: %s", name, err)
		return nil, err
	}

	config.Free()

	w := &Watcher{
		name:        name,
		dir:         dir,
		url:         url,
		repo:        repo,
		branchSpecs: branchSpecs,
		branchRefs:  map[string]*git.Reference{}}
	return w, nil
}

func (w *Watcher) Run(orders chan<- *Order) error {
	interval := viper.GetDuration("repos." + w.name + ".interval")

	for {
		branchRefs, err := w.getBranchRefs()
		if err != nil {
			log.Printf("Cannot update refs for repo `%s`: %s", w.name, err)
			return err
		}

		for branchName, ref := range branchRefs {
			oldRef, ok := w.branchRefs[branchName]
			if ok && oldRef.Cmp(ref) == 0 {
				log.Printf(
					"Branch `%s` in repo `%s` is up to date",
					branchName, w.name)
				continue
			}
			w.branchRefs[branchName] = ref
			log.Printf("Updated `%s` branch for repo `%s`", branchName, w.name)

			objectCommit, err := ref.Peel(git.ObjectCommit)
			if err != nil {
				log.Printf(
					"Cannot read repo `%s` commit object for ref: %s",
					w.name, ref, err)
				return err
			}
			commit, err := objectCommit.AsCommit()
			if err != nil {
				log.Printf(
					"Cannot read repo `%s` commit for object: %s",
					w.name, objectCommit, err)
				return err
			}
			order := &Order{
				repoName:   w.name,
				repo:       w.repo,
				branchName: branchName,
				commit:     commit}
			orders <- order
			log.Printf(
				"New order for `%s:%s`", order.repoName, order.branchName)
		}

		time.Sleep(interval)
	}
}

func (w *Watcher) getBranchRefs() (map[string]*git.Reference, error) {
	remote, err := w.repo.Remotes.Lookup("origin")
	if err != nil {
		remote, err = w.repo.Remotes.CreateWithFetchspec(
			"origin", w.url, "+refs/*:refs/*")
		if err != nil {
			log.Printf("Cannot create remote for repo `%s`: %s", w.name, err)
			return nil, err
		}
	}

	err = remote.Fetch(
		[]string{},
		&git.FetchOptions{Prune: git.FetchPruneOn, UpdateFetchhead: true},
		"")
	if err != nil {
		log.Printf("Cannot fetch repo `%s`: %s", w.name, err)
		return nil, err
	}

	branchRefs := map[string]*git.Reference{}

	for _, branchSpec := range w.branchSpecs {
		branchMatches, err := filepath.Glob(
			path.Join(w.dir, "refs/heads", branchSpec))
		if err != nil {
			log.Printf("Cannot read refs for repo `%s`: %s", w.name, err)
			return nil, err
		}

		for _, branchMatch := range branchMatches {
			branchName, err := filepath.Rel(
				path.Join(w.dir, "refs/heads"), branchMatch)
			if err != nil {
				log.Printf(
					"Impossible ref path `%s` for branch in repo `%s`: %s",
					branchMatch, w.name, err)
				return nil, err
			}

			branch, err := w.repo.LookupBranch(branchName, git.BranchLocal)
			if err != nil {
				log.Printf(
					"Cannot find branch `%s` for repo `%s`: %s",
					branchName, w.name, err)
				return nil, err
			}
			branchRefs[branchName] = branch.Reference
		}
	}

	return branchRefs, nil
}
