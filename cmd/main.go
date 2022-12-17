package main

import (
	"github.com/mattfan00/nycvbtracker/internal/scraper"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	checkErr(err)

	scraper, err := scraper.Default()
	checkErr(err)

	scraper.Scrape()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
