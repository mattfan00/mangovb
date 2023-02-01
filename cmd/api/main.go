package main

import (
	"flag"

	"github.com/mattfan00/mangovb/internal/api"
	"github.com/mattfan00/mangovb/internal/logger"

	"github.com/spf13/viper"
)

func main() {
	configPath := flag.String("c", "./config.yaml", "path to config file")
	flag.Parse()

	viper.SetConfigFile(*configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	log := logger.New(viper.GetString("env"))
	apiLogger := logger.SetSource(log, "api")

	api.Start(apiLogger)
}
