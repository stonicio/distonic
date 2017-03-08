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

func (w *Worker) Run(jobs <-chan *Job) {
	for job := range jobs {
		log.Printf("Received job: %s", job)
		err := w.processJob(job)
		if err != nil {
			log.Printf("Error processing job `%s`: %s", job, err)
		}
	}
}

func (w *Worker) processJob(job *Job) error {
	if err := w.prepareWorkdir(job); err != nil {
		log.Printf("Error preparing workdir for job `%s`: %s", job, err)
		return err
	}
	return nil
}

func (w *Worker) prepareWorkdir(job *Job) error {
	var err error
	var repo *git.Repository

	dataDir := viper.GetString("data_dir")
	workDir := path.Join(
		dataDir,
		"worker",
		job.repoName,
		job.branchName,
		job.commit.Object.Id().String())

	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		repo, err = git.Clone(
			job.repo.Path(),
			workDir,
			&git.CloneOptions{
				Bare:           false,
				CheckoutBranch: job.branchName,
				CheckoutOpts:   &git.CheckoutOpts{Strategy: git.CheckoutForce}})
		if err != nil {
			log.Printf(
				"Cannot make working clone for repo `%s`: %s",
				job.repoName, err)
			return err
		}
	} else {
		repo, err = git.OpenRepository(workDir)
		if err != nil {
			log.Printf(
				"Cannot open working clone for repo `%s`: %s",
				job.repoName, err)
			return err
		}
	}

	err = repo.SetHeadDetached(job.commit.Object.Id())
	if err != nil {
		log.Printf("Cannot set head on repo `%s`: %s", job.repoName, err)
		return err
	}

	err = repo.CheckoutHead(&git.CheckoutOpts{Strategy: git.CheckoutForce})
	if err != nil {
		log.Printf(
			"Cannot checkout workdir for repo `%s`: %s",
			job.repoName, err)
		return err
	}

	log.Printf("Working dir `%s` is ready", workDir)
	return nil
}
