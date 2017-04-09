package artefactory

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var artefactory Artefactory

func init() {
	var err error
	artefactory, err = NewArtefactory()
	if err != nil {
		log.Fatalf("Cannot init artefactory: %s", err)
	}
}

func NewArtefactory() (Artefactory, error) {
	backend := viper.GetString("artefactory.storage.backend")
	switch backend {
	case "mongodb":
		return NewMongoDBArtefactory(
			viper.GetString("artefactory.storage.mongodb.url"))
	default:
		return nil, fmt.Errorf(
			"Unknown artefactory storage backend: `%s`", backend)
	}
}

type Artefactory interface {
}
