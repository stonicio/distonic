package distonic

import (
	"log"
	"os"
	"path"

	git "github.com/libgit2/git2go"
	"github.com/spf13/viper"
)

type Worker struct {
}

func NewWorker() (*Worker, error) {
	return &Worker{}, nil
}

func (w *Worker) Run(orders <-chan *Order) {
	for order := range orders {
		log.Printf("Received order: %s", order)
		err := w.processOrder(order)
		if err != nil {
			log.Printf("Error processing order `%s`: %s", order, err)
		}
	}
}

func (w *Worker) processOrder(order *Order) error {
	workdir, err := w.prepareWorkdir(order)
	if err != nil {
		log.Printf("Error preparing workdir for order `%s`: %s", order, err)
		return err
	}

	pipeline, err := w.readPipeline(workdir)
	if err != nil {
		log.Printf("Could not read pipeline for order `%s`: %s", order, err)
		return err
	}
	log.Printf("Read pipeline %s", pipeline)
	return nil
}

func (w *Worker) prepareWorkdir(order *Order) (string, error) {
	var err error
	var repo *git.Repository

	dataDir := viper.GetString("data_dir")
	workDir := path.Join(
		dataDir,
		"worker",
		order.repoName,
		order.branchName,
		order.commit.Object.Id().String())

	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		repo, err = git.Clone(
			order.repo.Path(),
			workDir,
			&git.CloneOptions{
				Bare:           false,
				CheckoutBranch: order.branchName,
				CheckoutOpts:   &git.CheckoutOpts{Strategy: git.CheckoutForce}})
		if err != nil {
			log.Printf(
				"Cannot make working clone for repo `%s`: %s",
				order.repoName, err)
			return "", err
		}
	} else {
		repo, err = git.OpenRepository(workDir)
		if err != nil {
			log.Printf(
				"Cannot open working clone for repo `%s`: %s",
				order.repoName, err)
			return "", err
		}
	}

	err = repo.SetHeadDetached(order.commit.Object.Id())
	if err != nil {
		log.Printf("Cannot set head on repo `%s`: %s", order.repoName, err)
		return "", err
	}

	err = repo.CheckoutHead(&git.CheckoutOpts{Strategy: git.CheckoutForce})
	if err != nil {
		log.Printf(
			"Cannot checkout workdir for repo `%s`: %s",
			order.repoName, err)
		return "", err
	}

	log.Printf("Working dir `%s` is ready", workDir)
	return workDir, nil
}

func (w *Worker) readPipeline(dir string) (*Pipeline, error) {
	return &Pipeline{}, nil
}
