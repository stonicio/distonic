package main

import (
	"github.com/spf13/viper"
	"github.com/stonicio/distonic"
	"log"
)

func main() {
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Fatal error in config file: %s", err)
	}
	log.Println("Config read successfully")

	accountant, err := distonic.NewAccountant()
	if err != nil {
		log.Panicf("Cannot create accountant: %s", err)
	}
	log.Println("Created accountant")

	if err := accountant.Run(); err != nil {
		log.Panicf("Fatal error in accountant service file: %s", err)
	}
	log.Println("Shutting down")
}
