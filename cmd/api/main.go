package main

import (
	"flag"

	"github.com/mattfan00/mangovb/internal/api"
	"github.com/mattfan00/mangovb/internal/logger"
	"github.com/mattfan00/mangovb/internal/store"
	"github.com/mattfan00/mangovb/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configPath := flag.String("c", "./config.yaml", "path to config file")
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	conf, err := config.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	log := logger.New(conf.Env)
	log.Infof("read config: %+v", conf)

	apiLogger := logger.SetSource(log, "api")

	db, err := sqlx.Connect("sqlite3", conf.DbConn)
	if err != nil {
		log.Panic(err)
	}
	log.Infof("Connected to DB: %s", conf.DbConn)

	eventStore := store.NewEventStore(db)

	a := api.New(apiLogger, eventStore, conf)
	a.Start(*port)
}
