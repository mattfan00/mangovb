package query

import (
	"encoding/json"
	"fmt"
	"io"
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

func (q *Query) Document(url string) (*goquery.Document, error) {
	res, err := q.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (q *Query) Json(method string, url string, body io.Reader, data any) error {
	res, err := q.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(data)
	return err
}

func (q *Query) doRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if q.UserAgent != "" {
		req.Header.Set("User-Agent", q.UserAgent)
	}

	res, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		res.Body.Close()
		return nil, fmt.Errorf("status code is %d, not 200", res.StatusCode)
	}

	return res, nil
}
