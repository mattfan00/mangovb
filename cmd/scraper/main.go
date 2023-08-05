package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mattfan00/mangovb/internal/bot"
	"github.com/mattfan00/mangovb/internal/logger"
	"github.com/mattfan00/mangovb/internal/notifier"
	"github.com/mattfan00/mangovb/internal/scraper"
	"github.com/mattfan00/mangovb/internal/store"
	"github.com/mattfan00/mangovb/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

const SCRAPER_NAME = "scraper"
const NOTIFIER_NAME = "notifier"

func main() {
	configPath := flag.String("c", "./config.yaml", "path to config file")
	flag.Parse()

	conf, err := config.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	log := logger.New(conf.Env)
	log.Infof("read config: %+v", conf)

	scraperLogger := logger.SetSource(log, SCRAPER_NAME)
	notifierLogger := logger.SetSource(log, NOTIFIER_NAME)

	db, err := sqlx.Connect("sqlite3", conf.DbConn)
	if err != nil {
		log.Panic(err)
	}
	log.Infof("Connected to DB: %s", conf.DbConn)

	eventStore := store.NewEventStore(db)
	eventNotifStore := store.NewEventNotifStore(db)

	bot, err := bot.New(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}

	for _, channel := range conf.Channels {
		bot.AddChannel(channel)
	}

	err = bot.Start()
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		bot.Stop()
	}()

	scraper := scraper.New(eventStore, scraperLogger)
	notifier := notifier.New(bot, eventStore, eventNotifStore, notifierLogger)

	c := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))
	c.AddFunc(conf.CronScrape, func() {
		scraperLogger.Infof("Started %s", SCRAPER_NAME)
		defer logExecTime(SCRAPER_NAME, scraperLogger)()
		scraper.Scrape()
	})

	c.AddFunc(conf.CronNotify, func() {
		notifierLogger.Infof("Started %s", NOTIFIER_NAME)
		defer logExecTime(NOTIFIER_NAME, notifierLogger)()
		notifier.Notify()
	})

	c.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func logExecTime(name string, logger *logrus.Entry) func() {
	start := time.Now()
	return func() {
		logger.Infof("%s execution time: %v", name, time.Since(start))
	}
}
