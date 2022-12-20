package scraper

import (
	"fmt"
	"log"

	vb "github.com/mattfan00/nycvbtracker"
	"github.com/mattfan00/nycvbtracker/internal/engine"
	"github.com/mattfan00/nycvbtracker/internal/store"
	"github.com/mattfan00/nycvbtracker/pkg/query"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Scraper struct {
	db         *sqlx.DB
	eventStore *store.EventStore
	engines    []engine.Engine
}

func Default() (*Scraper, error) {
	db, err := sqlx.Connect("sqlite3", viper.GetString("db_conn"))
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to %s", viper.GetString("db_conn"))

	eventStore := store.NewEventStore(db)

	client := query.DefaultClient()

	scraper := &Scraper{
		db:         db,
		eventStore: eventStore,
	}
	scraper.RegisterEngine(engine.NewNyurbanEngine(client))

	return scraper, nil
}

func (s *Scraper) RegisterEngine(engines ...engine.Engine) {
	s.engines = append(s.engines, engines...)
}

func (s *Scraper) Scrape() {
	parsedEvents := []vb.Event{}

	for _, engine := range s.engines {
		events, err := engine.Run()
		if err != nil {
			log.Println(err)
		} else {
			parsedEvents = append(parsedEvents, events...)
		}
	}

	ids := []string{}
	for _, event := range parsedEvents {
		ids = append(ids, event.Id)
	}

	latestEvents, err := s.eventStore.GetLatestByIds(ids)
	if err != nil {
		fmt.Println("in here")
		log.Fatal(err)
	}

	notifs := []vb.Notification{}
	for i := range parsedEvents {
		parsedEvent := &parsedEvents[i]
		if notif, created := createNotification(parsedEvent, latestEvents); created {
			notifs = append(notifs, notif)
		}
	}

	for _, notif := range notifs {
		fmt.Printf("%d - %+v\n", notif.Type, *notif.Event)
	}

	err = s.eventStore.InsertMultiple(parsedEvents)
	if err != nil {
		log.Fatal(err)
	}
}

func createNotification(e *vb.Event, latestEvents map[string]vb.Event) (vb.Notification, bool) {
	if _, found := latestEvents[e.Id]; found {
		// TODO: revise LimitedSpots notification logic
		if e.IsAvailable && e.SpotsLeft < 10 {
			return vb.Notification{
				Type:  vb.LimitedSpots,
				Event: e,
			}, true
		}
	} else {
		return vb.Notification{
			Type:  vb.NewEvent,
			Event: e,
		}, true
	}

	return vb.Notification{}, false
}
