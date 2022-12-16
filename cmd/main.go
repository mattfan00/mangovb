package main

import (
	"fmt"

	"github.com/mattfan00/nycvbtracker/internal/scraper"
	"github.com/mattfan00/nycvbtracker/internal/scraper/engine"
	"github.com/mattfan00/nycvbtracker/pkg/query"
)

func main() {
	client := query.DefaultClient()

	nyurbanEngine := engine.NewNyurbanEngine(client)

	scraper := scraper.New()
	scraper.Register(nyurbanEngine)

	events := scraper.Scrape()

	for _, event := range events {
		fmt.Printf("%+v\n", event)
	}
}
