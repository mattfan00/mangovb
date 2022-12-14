package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Event struct {
	Name        string
	Location    string
	StartDate   time.Time
	StartTime   string
	EndTime     string
	SkillLevel  string
	Price       float64
	IsAvailable bool
	SpotsLeft   int
	Url         string
	Source      string
}

func main() {
	c := colly.NewCollector()

	const source = "NY Urban"
	const url = "https://www.nyurban.com/open-play-registration-vb/?id=35&gametypeid=1&filter_id=1"

	c.OnHTML("div.time_schedule_table", func(tableDiv *colly.HTMLElement) {
		events := []Event{}

		location := strings.TrimSuffix(tableDiv.ChildText("h3 > span"), ":")
		fmt.Println(location)

		tableDiv.ForEach("table tr", func(rowNum int, row *colly.HTMLElement) {
			if rowNum == 0 {
				return
			}

			event := Event{}

			event.Name = row.ChildText("td:nth-child(3)")
			event.Location = location
			event.Url = url
			event.Source = source

			price, err := strconv.ParseFloat(row.ChildText("td:nth-child(5)"), 64)
			if err != nil {
				log.Println(err)
			} else {
				event.Price = price
			}

			parsedStartDate, err := time.Parse("Mon 01/02", row.ChildText("td:nth-child(2)"))
			if err != nil {
				log.Println(err)
			} else {
				currentTime := time.Now()
				eventYear := currentTime.Year()

				// if the parsed month is less than the current month, that means the
				// event takes place in the next year
				if parsedStartDate.Month() < currentTime.Month() {
					eventYear = currentTime.Year() + 1
				}

				event.StartDate = time.Date(
					eventYear,
					parsedStartDate.Month(),
					parsedStartDate.Day(),
					0, 0, 0, 0,
					parsedStartDate.Location(),
				)
			}

			rawTimesSplit := strings.Split(row.ChildText("td:nth-child(4)"), "-")
			rawStartTime := strings.Trim(rawTimesSplit[0], " ")
			rawEndTime := strings.Trim(rawTimesSplit[1], " ")

			parsedStartTime, err := time.Parse("3:04 pm", rawStartTime)
			if err != nil {
				log.Println(err)
			} else {
				event.StartTime = parsedStartTime.Format("15:04")
			}

			parsedEndTime, err := time.Parse("3:04 pm", rawEndTime)
			if err != nil {
				log.Println(err)
			} else {
				event.EndTime = parsedEndTime.Format("15:04")
			}

			rawAvail := row.ChildText("td:nth-child(6)")
			if rawAvail == "Yes" {
				event.IsAvailable = true
			} else if rawAvail == "Sold Out" {
				event.IsAvailable = false
			} else if strings.Contains(strings.ToLower(rawAvail), "space") {
				re := regexp.MustCompile("[0-9]+")
				rawSpotsLeft := re.FindString(rawAvail)

				event.IsAvailable = true
				event.SpotsLeft, err = strconv.Atoi(rawSpotsLeft)
				if err != nil {
					log.Println(err)
				}
			}

			events = append(events, event)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"
	c.Visit(url)
}
