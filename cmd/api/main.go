package main

import (
	"flag"

	"github.com/mattfan00/mangovb/internal/api"
	"github.com/mattfan00/mangovb/internal/logger"
	"github.com/mattfan00/mangovb/internal/store"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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

	dbConn := viper.GetString("db_conn")
	db, err := sqlx.Connect("sqlite3", dbConn)
	if err != nil {
		log.Panic(err)
	}
	log.Infof("Connected to DB: %s", dbConn)

	eventStore := store.NewEventStore(db)

	a := api.New(apiLogger, eventStore)
	a.Start()
}
