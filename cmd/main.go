package main

import (
	"fmt"
	"log"

	"github.com/mattfan00/nycvbtracker/internal/scraper"
	"github.com/mattfan00/nycvbtracker/internal/scraper/engine"
	"github.com/mattfan00/nycvbtracker/internal/store"
	"github.com/mattfan00/nycvbtracker/pkg/query"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Connect("sqlite3", viper.GetString("db_conn"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to %s", viper.GetString("db_conn"))

	_ = store.NewEventStore(db)

	client := query.DefaultClient()

	nyurbanEngine := engine.NewNyurbanEngine(client)

	scraper := scraper.New()
	scraper.Register(nyurbanEngine)

	events := scraper.Scrape()

	for _, event := range events {
		fmt.Printf("%+v\n", event)
	}
}
