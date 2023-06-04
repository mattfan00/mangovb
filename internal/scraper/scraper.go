package scraper

import (
	vb "github.com/mattfan00/mangovb"
	"github.com/mattfan00/mangovb/internal/engine"
	"github.com/mattfan00/mangovb/internal/store"
	"github.com/mattfan00/mangovb/pkg/query"
	"github.com/sirupsen/logrus"

	"go.uber.org/multierr"
)

type Scraper struct {
	eventStore *store.EventStore
	engines    []engine.Engine
	logger     *logrus.Entry
}

func New(eventStore *store.EventStore, logger *logrus.Entry) *Scraper {
	client := query.DefaultClient()

	scraper := &Scraper{
		eventStore: eventStore,
		engines:    []engine.Engine{},
		logger:     logger,
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
		for _, event := range events {
			s.logger.WithFields(logrus.Fields{
				"event": logrus.Fields{
					"id":           event.Id,
					"is_available": event.IsAvailable,
					"spots_left":   event.SpotsLeft,
				},
			}).Info("Scraped event")
		}

		for _, err := range multierr.Errors(err) {
			s.logger.Error(err)
		}

		parsedEvents = append(parsedEvents, events...)
	}

	err := s.eventStore.UpsertMultiple(parsedEvents)
	if err != nil {
		s.logger.Error(err)
		return
	}
}
