package engine

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	vb "github.com/mattfan00/nycvbtracker"
	"github.com/mattfan00/nycvbtracker/internal/hash"
	"github.com/mattfan00/nycvbtracker/pkg/query"
)

type NyurbanEngine struct {
	query  *query.Query
	source string
}

func NewNyurbanEngine(client *http.Client) *NyurbanEngine {
	q := query.New(client)
	q.UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"

	return &NyurbanEngine{
		query:  q,
		source: "nyurban",
	}
}

func (n *NyurbanEngine) parse(doc *goquery.Document) []vb.Event {
	events := []vb.Event{}
	scrapedOn := time.Now()

	tableDiv := doc.FindMatcher(goquery.Single("div.time_schedule_table"))
	location := strings.TrimSuffix(tableDiv.Find("h3 > span").Text(), ":")

	tableDiv.Find("table tr").Each(func(rowNum int, row *goquery.Selection) {
		if rowNum == 0 {
			return
		}

		e := vb.Event{}
		e.Source = n.source
		e.Location = location

		err := n.parseRow(row, &e)
		if err != nil {
			log.Printf("Error parsing %+v: %s", e, err.Error())
			return
		}

		hashedEvent, err := hash.HashEvent(e)
		if err != nil {
			log.Printf("Error hashing %+v: %s", e, err.Error())
			return
		}
		e.Id = hashedEvent
		e.ScrapedOn = scrapedOn

		events = append(events, e)
	})

	return events
}

func (n *NyurbanEngine) parseRow(row *goquery.Selection, event *vb.Event) error {
	event.Name = row.Find("td:nth-child(3)").Text()

	price, err := strconv.ParseFloat(row.Find("td:nth-child(5)").Text(), 64)
	if err != nil {
		return err
	} else {
		event.Price = price
	}

	parsedStartDate, err := time.Parse("Mon 01/02", strings.TrimSpace(row.Find("td:nth-child(2)").Text()))
	if err != nil {
		return err
	}
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

	rawTimesSplit := strings.Split(row.Find("td:nth-child(4)").Text(), "-")
	rawStartTime := strings.Trim(rawTimesSplit[0], " ")
	rawEndTime := strings.Trim(rawTimesSplit[1], " ")

	parsedStartTime, err := time.Parse("3:04 pm", rawStartTime)
	if err != nil {
		return err
	}
	event.StartTime = parsedStartTime.Format("15:04")

	parsedEndTime, err := time.Parse("3:04 pm", rawEndTime)
	if err != nil {
		return err
	}
	event.EndTime = parsedEndTime.Format("15:04")

	rawAvail := strings.TrimSpace(row.Find("td:nth-child(6)").Text())
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
			return err
		}
	}

	return nil
}

func (n *NyurbanEngine) Run() ([]vb.Event, error) {
	docs, err := n.query.VisitMulitple([]string{
		"https://www.nyurban.com/open-play-registration-vb/?id=35&gametypeid=1&filter_id=1",
	})
	if err != nil {
		return []vb.Event{}, err
	}

	allEvents := []vb.Event{}

	for _, doc := range docs {
		events := n.parse(doc)

		allEvents = append(allEvents, events...)
	}

	return allEvents, nil
}
