package scraper

import (
	"log"

	vb "github.com/mattfan00/nycvbtracker"
	"github.com/mattfan00/nycvbtracker/internal/scraper/engine"
)

type Scraper struct {
	engines []engine.Engine
}

func New() *Scraper {
	return &Scraper{}
}

func (s *Scraper) Register(engines ...engine.Engine) {
	s.engines = append(s.engines, engines...)
}

func (s *Scraper) Scrape() []vb.Event {
	allEvents := []vb.Event{}

	for _, engine := range s.engines {
		events, err := engine.Run()
		if err != nil {
			log.Println(err)
		} else {
			allEvents = append(allEvents, events...)
		}
	}

	return allEvents
}
