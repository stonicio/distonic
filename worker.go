package distonic

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	git "github.com/libgit2/git2go"
	"github.com/spf13/viper"
	"github.com/stonicio/distonic/module"
	"gopkg.in/yaml.v2"
)

type Worker struct {
}

func NewWorker() (*Worker, error) {
	return &Worker{}, nil
}

func (w *Worker) Run(orders <-chan *Order) {
	for order := range orders {
		log.Printf(
			"Received order for `%s:%s`", order.repoName, order.branchName)
		err := w.processOrder(order)
		if err != nil {
			log.Printf(
				"Error processing `%s:%s`: %s",
				order.repoName, order.branchName, err)
		}
	}
}

func (w *Worker) processOrder(order *Order) error {
	workdir, err := w.prepareWorkdir(order)
	if err != nil {
		log.Printf(
			"Error preparing workdir for `%s:%s`: %s",
			order.repoName, order.branchName, err)
		return err
	}

	context := &module.Context{
		Workdir:      workdir,
		Branch:       order.branchName,
		BranchDashed: strings.Replace(order.branchName, "/", "-", -1),
		Commit:       order.commit.Object.Id().String()}

	pipeline, err := w.readPipeline(context)
	if err != nil {
		log.Printf(
			"Could not read pipeline for `%s:%s`: %s",
			order.repoName, order.branchName, err)
		return err
	}

	result, err := pipeline.Run()
	log.Printf("### %+v", result)
	log.Printf("$$$ %s", err)

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

func (w *Worker) readPipeline(context *module.Context) (*Pipeline, error) {
	configFilename := path.Join(context.Workdir, "distonic.yml")

	t, err := template.ParseFiles(configFilename)
	if err != nil {
		log.Printf("Could not load distonic pipeline template: %s", err)
		return nil, err
	}

	config, err := os.Create(configFilename)
	if err != nil {
		log.Printf("Could not open distonic pipeline for writing: %s", err)
		return nil, err
	}

	if err := t.Execute(config, context); err != nil {
		log.Printf("Could not execute distonic pipeline template: %s", err)
		return nil, err
	}

	pData, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Printf("Could not read distonic pipeline config: %s", err)
		return nil, err
	}

	var p Pipeline
	err = yaml.Unmarshal(pData, &p)
	if err != nil {
		log.Printf("Could not initialize pipeline: %s", err)
		return nil, err
	}

	return &p, nil
}
