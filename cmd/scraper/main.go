package main

import (
	"fmt"
	"log"

	"github.com/mattfan00/nycvbtracker/internal/nyurban"
	"github.com/mattfan00/nycvbtracker/pkg/query"
)

func main() {
	client := query.DefaultClient()

	nyurbanService := nyurban.NewService(client)
	events, err := nyurbanService.Scrape()
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range events {
		fmt.Printf("%+v\n", event)
	}
}
