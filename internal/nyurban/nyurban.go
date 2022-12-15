package nyurban

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	vb "github.com/mattfan00/nycvbtracker"
	"github.com/mattfan00/nycvbtracker/pkg/query"
)

type NyurbanService struct {
	query *query.Query
}

func NewService(client *http.Client) *NyurbanService {
	q := query.New(client)
	q.UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"

	return &NyurbanService{
		query: q,
	}
}

func (n *NyurbanService) parse(doc *goquery.Document) ([]vb.Event, error) {
	events := []vb.Event{}

	tableDiv := doc.FindMatcher(goquery.Single("div.time_schedule_table"))
	location := strings.TrimSuffix(tableDiv.Find("h3 > span").Text(), ":")

	tableDiv.Find("table tr").Each(func(rowNum int, row *goquery.Selection) {
		if rowNum == 0 {
			return
		}

		e := vb.Event{}
		e.Location = location
		e.Name = row.Find("td:nth-child(3)").Text()

		price, err := strconv.ParseFloat(row.Find("td:nth-child(5)").Text(), 64)
		if err != nil {
			log.Println(err)
		} else {
			e.Price = price
		}

		parsedStartDate, err := time.Parse("Mon 01/02", strings.TrimSpace(row.Find("td:nth-child(2)").Text()))
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

			e.StartDate = time.Date(
				eventYear,
				parsedStartDate.Month(),
				parsedStartDate.Day(),
				0, 0, 0, 0,
				parsedStartDate.Location(),
			)
		}

		rawTimesSplit := strings.Split(row.Find("td:nth-child(4)").Text(), "-")
		rawStartTime := strings.Trim(rawTimesSplit[0], " ")
		rawEndTime := strings.Trim(rawTimesSplit[1], " ")

		parsedStartTime, err := time.Parse("3:04 pm", rawStartTime)
		if err != nil {
			log.Println(err)
		} else {
			e.StartTime = parsedStartTime.Format("15:04")
		}

		parsedEndTime, err := time.Parse("3:04 pm", rawEndTime)
		if err != nil {
			log.Println(err)
		} else {
			e.EndTime = parsedEndTime.Format("15:04")
		}

		rawAvail := strings.TrimSpace(row.Find("td:nth-child(6)").Text())
		if rawAvail == "Yes" {
			e.IsAvailable = true
		} else if rawAvail == "Sold Out" {
			e.IsAvailable = false
		} else if strings.Contains(strings.ToLower(rawAvail), "space") {
			re := regexp.MustCompile("[0-9]+")
			rawSpotsLeft := re.FindString(rawAvail)

			e.IsAvailable = true
			e.SpotsLeft, err = strconv.Atoi(rawSpotsLeft)
			if err != nil {
				log.Println(err)
			}
		}

		events = append(events, e)
	})

	return events, nil
}

func (n *NyurbanService) Scrape() ([]vb.Event, error) {
	docs, err := n.query.VisitMulitple([]string{
		"https://www.nyurban.com/open-play-registration-vb/?id=35&gametypeid=1&filter_id=1",
	})
	if err != nil {
		return []vb.Event{}, err
	}

	allEvents := []vb.Event{}

	for _, doc := range docs {
		events, err := n.parse(doc)
		if err != nil {
			return []vb.Event{}, err
		}

		allEvents = append(allEvents, events...)
	}

	return allEvents, nil
}
