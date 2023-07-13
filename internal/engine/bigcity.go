package engine

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	vb "github.com/mattfan00/mangovb"
	"github.com/mattfan00/mangovb/internal/hash"
	"github.com/mattfan00/mangovb/pkg/query"
	"go.uber.org/multierr"
)

type BigCityEngine struct {
	query    *query.Query
	sourceId vb.EventSource
}

func NewBigCityEngine(client *http.Client) *BigCityEngine {
	q := query.New(client)

	return &BigCityEngine{
		query:    q,
		sourceId: vb.BigCity,
	}
}

func (b *BigCityEngine) parseEvent(event BigCityEvent) vb.Event {
	e := vb.Event{}

	e.Id = hash.HashInt(event.Id)
	e.Name = event.Title
	e.Location = event.Place.Title
	e.StartTime = event.StartTime
	e.EndTime = event.EndTime
	e.SkillLevel = b.getSkillLevel(event.Data.Level.Title)

	ticket := event.TicketsSummary[0]
	e.Price = float64(ticket.Price)
	e.SpotsLeft = ticket.QuantityTotal - ticket.QuantitySold

	e.IsAvailable = e.SpotsLeft != 0
	e.Url = fmt.Sprintf("https://opensports.net/posts/%s", event.AliasId)
	e.SourceId = b.sourceId

	return e
}

func (b *BigCityEngine) getSkillLevel(rawLevel string) vb.EventSkillLevel {
	if strings.Contains(rawLevel, "C") {
		return vb.Beginner
	} else if strings.Contains(rawLevel, "B") {
		return vb.Intermediate
	} else if strings.Contains(rawLevel, "A") {
		return vb.Advanced
	}

	return vb.None
}

func (b *BigCityEngine) Run() ([]vb.Event, error) {
	var err error

	// only fetches first 48 records
	url := "https://osapi.opensports.ca/app/posts/listFiltered?limit=48&limitedFields=true&groupID=1962&rootTags[]=Open%20Play"

	allEvents := []vb.Event{}

	res := &BigCityResponse{}
	queryErr := b.query.Json(http.MethodGet, url, nil, &res)
	if queryErr != nil {
		err = multierr.Append(err, QueryErr{url, queryErr})
		return []vb.Event{}, err
	}

	for _, event := range res.Result {
		allEvents = append(allEvents, b.parseEvent(event))
	}

	return allEvents, err
}

type BigCityResponse struct {
	Response int            `json:"response"`
	Result   []BigCityEvent `json:"result"`
}

type BigCityEvent struct {
	Id             int                          `json:"id"`
	Title          string                       `json:"title"`
	StartTime      time.Time                    `json:"start"`
	EndTime        time.Time                    `json:"end"`
	AliasId        string                       `json:"aliasId"`
	Place          BigCityEventPlace            `json:"place"`
	TicketsSummary []BigCityEventTicketsSummary `json:"ticketsSummary"`
	Data           BigCityEventData             `json:"data"`
}

type BigCityEventPlace struct {
	Title string `json:"title"`
}

type BigCityEventTicketsSummary struct {
	Price         int `json:"price"`
	QuantitySold  int `json:"quantitySold"`
	QuantityTotal int `json:"quantityTotal"`
}

type BigCityEventData struct {
	Level BigCityEventDataLevel `json:"level"`
}

type BigCityEventDataLevel struct {
	Title string `json:"title"`
}
