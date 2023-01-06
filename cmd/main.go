package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mattfan00/mangovb/internal/bot"
	"github.com/mattfan00/mangovb/internal/logger"
	"github.com/mattfan00/mangovb/internal/notifier"
	"github.com/mattfan00/mangovb/internal/scraper"
	"github.com/mattfan00/mangovb/internal/store"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const SCRAPER_NAME = "scraper"
const NOTIFIER_NAME = "notifier"

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	checkErr(err)

	log := logger.New()
	scraperLogger := logger.SetSource(log, SCRAPER_NAME)
	notifierLogger := logger.SetSource(log, NOTIFIER_NAME)

	dbConn := viper.GetString("db_conn")
	db, err := sqlx.Connect("sqlite3", dbConn)
	checkErr(err)
	log.Infof("Connected to DB: %s", dbConn)

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
		scraperLogger.Infof("Started %s", SCRAPER_NAME)
		defer logExecTime(SCRAPER_NAME, scraperLogger)()
		scraper.Scrape()
	})

	c.AddFunc(viper.GetString("cron_notify"), func() {
		notifierLogger.Infof("Started %s", NOTIFIER_NAME)
		defer logExecTime(NOTIFIER_NAME, notifierLogger)()
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

func logExecTime(name string, logger *logrus.Entry) func() {
	start := time.Now()
	return func() {
		logger.Infof("%s execution time: %v\n", name, time.Since(start))
	}
}
