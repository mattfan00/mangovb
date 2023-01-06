package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mattfan00/mangovb/internal/bot"
	"github.com/mattfan00/mangovb/internal/logger"
	"github.com/mattfan00/mangovb/internal/notifier"
	"github.com/mattfan00/mangovb/internal/scraper"
	"github.com/mattfan00/mangovb/internal/store"

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

	log := logger.New()
	scraperLogger := logger.SetSource(log, "scraper")
	notifierLogger := logger.SetSource(log, "notifier")

	dbConn := viper.GetString("db_conn")
	db, err := sqlx.Connect("sqlite3", dbConn)
	checkErr(err)
	log.Info("Connected to %s", dbConn)

	eventStore := store.NewEventStore(db)
	eventNotifStore := store.NewEventNotifStore(db)

	bot, err := bot.New()
	checkErr(err)

	err = bot.Start()
	checkErr(err)
	defer func() {
		bot.Stop()
	}()

	scraper := scraper.New(eventStore, scraperLogger)
	notifier := notifier.New(bot, eventStore, eventNotifStore, notifierLogger)

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
