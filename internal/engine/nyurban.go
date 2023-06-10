package engine

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	vb "github.com/mattfan00/mangovb"
	"github.com/mattfan00/mangovb/internal/hash"
	"github.com/mattfan00/mangovb/pkg/query"

	"go.uber.org/multierr"
)

type NyUrbanEngine struct {
	query    *query.Query
	sourceId vb.EventSource
}

func NewNyUrbanEngine(client *http.Client) *NyUrbanEngine {
	q := query.New(client)
	q.UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"

	return &NyUrbanEngine{
		query:    q,
		sourceId: vb.NyUrban,
	}
}

func (n *NyUrbanEngine) parse(doc *goquery.Document, sourceUrl string) ([]vb.Event, error) {
	events := []vb.Event{}

	tableDiv := doc.FindMatcher(goquery.Single("div.time_schedule_table"))
	location := strings.TrimSuffix(tableDiv.Find("h3 > span").Text(), ":")

	var err error
	tableDiv.Find("table tr").Each(func(rowNum int, row *goquery.Selection) {
		if rowNum == 0 {
			return
		}

		// skip row if it doesn't have expected number of columns
		// thus, it skips rows that indicate there are no spots available
		if row.Find("td").Length() != 6 {
			return
		}

		e := vb.Event{}
		e.SourceId = n.sourceId
		e.Url = sourceUrl
		e.Location = location

		parseErr := n.parseRow(row, &e)

		e.Id = hash.Hash(e.Id)

		if parseErr != nil {
			err = multierr.Append(err, ParseEventErr{e, parseErr})
		} else {
			events = append(events, e)
		}
	})

	return events, err
}

func (n *NyUrbanEngine) parseRow(row *goquery.Selection, event *vb.Event) error {
	if id, exists := row.Find("td:nth-child(1) input").Attr("value"); exists {
		event.Id = id
	}

	rawLevel := row.Find("td:nth-child(3)").Text()
	event.Name = rawLevel
	event.SkillLevel = getSkillLevel(rawLevel)

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

	newYorkTimezone, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}

	startDate := time.Date(
		eventYear,
		parsedStartDate.Month(),
		parsedStartDate.Day(),
		0, 0, 0, 0,
		newYorkTimezone,
	)

	rawTimesSplit := strings.Split(row.Find("td:nth-child(4)").Text(), "-")
	rawStartTime := strings.Trim(rawTimesSplit[0], " ")
	rawEndTime := strings.Trim(rawTimesSplit[1], " ")

	parsedStartTime, err := time.Parse("3:04 pm", rawStartTime)
	if err != nil {
		return err
	}
	event.StartTime = combineDateAndTime(startDate, parsedStartTime)

	parsedEndTime, err := time.Parse("3:04 pm", rawEndTime)
	if err != nil {
		return err
	}
	event.EndTime = combineDateAndTime(startDate, parsedEndTime)

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

func getSkillLevel(rawLevel string) vb.EventSkillLevel {
	if strings.Contains(rawLevel, "Beginner") {
		return vb.Beginner
	} else if strings.Contains(rawLevel, "Intermediate") && strings.Contains(rawLevel, "Advanced") {
		return vb.AdvancedIntermediate
	} else if strings.Contains(rawLevel, "Intermediate") {
		return vb.Intermediate
	} else if strings.Contains(rawLevel, "Advanced") {
		return vb.Advanced
	}

	return vb.None
}

func combineDateAndTime(d time.Time, t time.Time) time.Time {
	return time.Date(
		d.Year(),
		d.Month(),
		d.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		d.Location(),
	)
}

func (n *NyUrbanEngine) Run() ([]vb.Event, error) {
	var err error

	urls := []string{
		"https://www.nyurban.com/open-play-volleyball-no/?id=35&gametypeid=1&filter_id=1",
		"https://www.nyurban.com/open-play-volleyball-no/?id=18&gametypeid=1&filter_id=1",
		"https://www.nyurban.com/open-play-volleyball-no/?id=34&gametypeid=1&filter_id=1",
		"https://www.nyurban.com/open-play-volleyball-no/?id=32&gametypeid=11&filter_id=1",
	}

	allEvents := []vb.Event{}
	for _, url := range urls {
		doc, queryErr := n.query.Document(url)
		if queryErr != nil {
			err = multierr.Append(err, QueryErr{url, queryErr})
		} else {
			events, parseErrs := n.parse(doc, url)

			err = multierr.Append(err, parseErrs)
			allEvents = append(allEvents, events...)
		}

	}

	return allEvents, err
}
