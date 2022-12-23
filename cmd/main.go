package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattfan00/nycvbtracker/internal/bot"
	"github.com/mattfan00/nycvbtracker/internal/scraper"

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

	err = bot.Session.Open()
	checkErr(err)

	scraper, err := scraper.New(bot)
	checkErr(err)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	scraper.Scrape()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	bot.Session.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
