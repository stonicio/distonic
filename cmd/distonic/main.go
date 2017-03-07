package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/stonicio/distonic"
)

func main() {
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Fatal error in config file: %s", err)
	}
	log.Println("Config read successfully")

	supervisor, err := distonic.NewSupervisor()
	if err != nil {
		log.Panicf("Cannot create supervisor: %s", err)
	}
	log.Println("Created supervisor")

	if err := supervisor.Run(); err != nil {
		log.Panicf("Fatal error in supervisor service: %s", err)
	}

	log.Println("Shutting down")
}
