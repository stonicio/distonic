package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/stonicio/distonic/supervisor"
)

func main() {
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Fatal error in config file: %s", err)
	}
	log.Println("Config read successfully")

	supervisor, err := supervisor.NewSupervisor()
	if err != nil {
		log.Panicf("Cannot create supervisor: %s", err)
	}
	log.Println("Created supervisor")

	supervisor.Run()

	log.Println("Shutting down")
}
