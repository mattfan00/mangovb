package scraper

import (
	"log"

	vb "github.com/mattfan00/mangovb"
	"github.com/mattfan00/mangovb/internal/engine"
	"github.com/mattfan00/mangovb/internal/store"
	"github.com/mattfan00/mangovb/pkg/query"
	"go.uber.org/multierr"
)

type Scraper struct {
	eventStore *store.EventStore
	engines    []engine.Engine
}

func New(eventStore *store.EventStore) *Scraper {
	client := query.DefaultClient()

	scraper := &Scraper{
		eventStore: eventStore,
		engines:    []engine.Engine{},
	}
	scraper.RegisterEngine(engine.NewNyUrbanEngine(client))

	return scraper
}

func (s *Scraper) RegisterEngine(engines ...engine.Engine) {
	s.engines = append(s.engines, engines...)
}

func (s *Scraper) Scrape() {
	parsedEvents := []vb.Event{}

	for _, engine := range s.engines {
		events, err := engine.Run()
		for _, err := range multierr.Errors(err) {
			log.Println(err)
		}

		parsedEvents = append(parsedEvents, events...)
	}

	err := s.eventStore.UpsertMultiple(parsedEvents)
	if err != nil {
		log.Fatal(err)
	}
}
