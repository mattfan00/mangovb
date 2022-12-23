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

	_, err = scraper.Default()
	checkErr(err)

	// scraper.Scrape()

	bot, err := bot.Default()
	checkErr(err)

	err = bot.Open()
	checkErr(err)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
