package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattfan00/nycvbtracker/internal/bot"
	"github.com/mattfan00/nycvbtracker/internal/scraper"
	"github.com/robfig/cron"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	checkErr(err)

	bot, err := bot.New()
	checkErr(err)

	err = bot.Start()
	checkErr(err)
	defer func() {
		fmt.Println("end")
		bot.Stop()
	}()

	scraper, err := scraper.New(bot)
	checkErr(err)

	c := cron.New()
	c.AddFunc("0 * * * * *", func() {
		log.Println("Started scraping")
		scraper.Scrape()
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
