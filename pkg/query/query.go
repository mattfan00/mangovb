package query

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Query struct {
	UserAgent string
	client    *http.Client
}

func New(client *http.Client) *Query {
	return &Query{
		client: client,
	}
}

func DefaultClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func (q *Query) Visit(url string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", q.UserAgent)

	res, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d, not 200", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
