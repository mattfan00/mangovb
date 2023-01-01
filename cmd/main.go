package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattfan00/nycvbtracker/internal/bot"
	"github.com/mattfan00/nycvbtracker/internal/notifier"
	"github.com/mattfan00/nycvbtracker/internal/scraper"
	"github.com/mattfan00/nycvbtracker/internal/store"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	checkErr(err)

	dbConn := viper.GetString("db_conn")
	db, err := sqlx.Connect("sqlite3", dbConn)
	checkErr(err)
	log.Printf("Connected to %s", dbConn)

	eventStore := store.NewEventStore(db)
	eventNotifStore := store.NewEventNotifStore(db)

	bot, err := bot.New()
	checkErr(err)

	err = bot.Start()
	checkErr(err)
	defer func() {
		fmt.Println("end")
		bot.Stop()
	}()

	scraper := scraper.New(eventStore)
	notifier := notifier.New(bot, eventStore, eventNotifStore)

	c := cron.New()
	c.AddFunc(viper.GetString("cron_scrape"), func() {
		log.Println("Started scraper")
		scraper.Scrape()
	})

	c.AddFunc(viper.GetString("cron_notify"), func() {
		log.Println("Started notifer")
		notifier.Notify()
	})

	c.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
