package scraper

import (
	"fmt"
	"log"
	"time"

	vb "github.com/mattfan00/nycvbtracker"
	"github.com/mattfan00/nycvbtracker/internal/bot"
	"github.com/mattfan00/nycvbtracker/internal/engine"
	"github.com/mattfan00/nycvbtracker/internal/store"
	"github.com/mattfan00/nycvbtracker/pkg/query"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Scraper struct {
	db              *sqlx.DB
	bot             *bot.Bot
	eventStore      *store.EventStore
	eventNotifStore *store.EventNotifStore
	engines         []engine.Engine
}

func New(bot *bot.Bot) (*Scraper, error) {
	db, err := sqlx.Connect("sqlite3", viper.GetString("db_conn"))
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to %s", viper.GetString("db_conn"))

	eventStore := store.NewEventStore(db)
	eventNotifStore := store.NewEventNotifStore(db)

	client := query.DefaultClient()

	scraper := &Scraper{
		db:              db,
		bot:             bot,
		eventStore:      eventStore,
		eventNotifStore: eventNotifStore,
		engines:         []engine.Engine{},
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

	err := s.eventStore.UpsertMultiple(parsedEvents)
	if err != nil {
		log.Fatal(err)
	}

	ids := make([]string, len(parsedEvents))
	for i, event := range parsedEvents {
		ids[i] = event.Id
	}

	notifMap, err := s.eventNotifStore.GetByEventIds(ids)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range notifMap {
		fmt.Printf("%s: ", key)
		for _, notif := range value {
			fmt.Printf("[%d] ", notif.Type)
		}
		fmt.Print("\n")
	}

	notifs := []vb.EventNotif{}
	notifCreatedOn := time.Now()
	for i := range parsedEvents {
		parsedEvent := parsedEvents[i]
		if notif, created := createNotification(parsedEvent, notifMap); created {
			notif.CreatedOn = notifCreatedOn
			notifs = append(notifs, notif)
		}
	}

	err = s.eventNotifStore.InsertMultiple(notifs)
	if err != nil {
		log.Fatal(err)
	}

	s.bot.NotifyAllChannels(notifs)
}

func createNotification(e vb.Event, notifMap map[string][]vb.EventNotif) (vb.EventNotif, bool) {
	if prevNotifs, found := notifMap[e.Id]; found {
		if e.IsAvailable && e.SpotsLeft > 0 && e.SpotsLeft < 5 {
			hasNotifiedLimitedSpots := false
			for _, prevNotif := range prevNotifs {
				if prevNotif.Type == vb.LimitedSpots {
					hasNotifiedLimitedSpots = true
				}
			}

			// only notify if haven't notified limited spots in the past
			if !hasNotifiedLimitedSpots {
				return vb.EventNotif{
					Type:    vb.LimitedSpots,
					EventId: e.Id,
					Event:   e,
				}, true
			}
		}
	} else {
		return vb.EventNotif{
			Type:    vb.NewEvent,
			EventId: e.Id,
			Event:   e,
		}, true
	}

	return vb.EventNotif{}, false
}
