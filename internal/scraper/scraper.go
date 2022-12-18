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
	allEvents := []vb.Event{}

	for _, engine := range s.engines {
		events, err := engine.Run()
		if err != nil {
			log.Println(err)
		} else {
			allEvents = append(allEvents, events...)
		}
	}

	for _, event := range allEvents {
		fmt.Printf("%+v\n", event)
	}

	err := s.eventStore.InsertMultiple(allEvents)
	if err != nil {
		log.Fatal(err)
	}
}
