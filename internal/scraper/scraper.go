package scraper

import (
	"log"

	vb "github.com/mattfan00/nycvbtracker"
	"github.com/mattfan00/nycvbtracker/internal/engine"
	"github.com/mattfan00/nycvbtracker/internal/store"
	"github.com/mattfan00/nycvbtracker/pkg/query"
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
	scraper.RegisterEngine(engine.NewNyurbanEngine(client))

	return scraper
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

	/*
		ids := make([]string, len(parsedEvents))
		for i, event := range parsedEvents {
			ids[i] = event.Id
		}

		notifMap, err := s.eventNotifStore.GetByEventIds(ids)
		if err != nil {
			log.Fatal(err)
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
	*/
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
