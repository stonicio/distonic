package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stonicio/distonic"
)

func main() {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error in config file: %s \n", err))
	}

	accountant := distonic.NewAccountant()
	err = accountant.Run()
	if err != nil {
		panic(fmt.Errorf("Fatal error in accountant service file: %s \n", err))
	}
}
