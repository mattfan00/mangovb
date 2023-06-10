package engine

import (
	"fmt"
	"net/http"
	"time"

	vb "github.com/mattfan00/mangovb"
	"github.com/mattfan00/mangovb/pkg/query"
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

func (b *BigCityEngine) Run() ([]vb.Event, error) {
	url := "https://osapi.opensports.ca/app/posts/listFiltered?limit=48&limitedFields=true&groupID=1962&rootTags[]=Open%20Play"

	res := &BigCityResponse{}
	b.query.Json(http.MethodGet, url, nil, &res)

    for _, event := range res.Result {
        fmt.Println(event.Id)
    }

	return []vb.Event{}, nil
}

type BigCityResponse struct {
	Response int            `json:"response"`
	Result   []BigCityEvent `json:"result"`
}

type BigCityEvent struct {
	Id             int                        `json:"id"`
	Title          string                     `json:"title"`
	StartTime      time.Time                  `json:"start"`
	EndTime        time.Time                  `json:"end"`
	AliasId        string                     `json:"aliasId"`
	Place          BigCityEventPlace          `json:"place"`
	TicketsSummary BigCityEventTicketsSummary `json:"ticketsSummary"`
	Data           BigCityEventData           `json:"data"`
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

func (b *BigCityEvent) Url() string {
	return fmt.Sprintf("https://opensports.net/posts/%s", b.AliasId)
}

func (b *BigCityEvent) SpotsLeft() int {
	return b.TicketsSummary.QuantityTotal - b.TicketsSummary.QuantitySold
}

func (b *BigCityEvent) IsAvailable() bool {
	return b.SpotsLeft() != 0
}
