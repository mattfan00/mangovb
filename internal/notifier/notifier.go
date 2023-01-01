package notifier

import (
	"fmt"
	"log"
	"time"

	vb "github.com/mattfan00/nycvbtracker"
	"github.com/mattfan00/nycvbtracker/internal/bot"
	"github.com/mattfan00/nycvbtracker/internal/store"
)

type Notifier struct {
	bot             *bot.Bot
	eventStore      *store.EventStore
	eventNotifStore *store.EventNotifStore
}

func New(bot *bot.Bot, eventStore *store.EventStore, eventNotifStore *store.EventNotifStore) *Notifier {
	return &Notifier{
		bot:             bot,
		eventStore:      eventStore,
		eventNotifStore: eventNotifStore,
	}
}

func (n *Notifier) Notify() {
	events, err := n.eventStore.GetLatest()
	if err != nil {
		log.Println(err)
	}

	ids := make([]string, len(events))
	for i, event := range events {
		ids[i] = event.Id
	}

	notifMap, err := n.eventNotifStore.GetByEventIds(ids)
	if err != nil {
		log.Fatal(err)
	}

	notifs := []vb.EventNotif{}
	notifCreatedOn := time.Now()
	for i := range events {
		event := events[i]
		if notif, created := n.createNotif(event, notifMap); created {
			notif.CreatedOn = notifCreatedOn
			notifs = append(notifs, notif)
		}
	}

	if len(notifs) > 0 {
		err = n.eventNotifStore.InsertMultiple(notifs)
		if err != nil {
			log.Fatal(err)
		}

		message := n.generateNotifMessage(notifs)
		n.bot.SendMessageToAllChannels(message)
	}
}

func (n *Notifier) createNotif(e vb.Event, notifMap map[string][]vb.EventNotif) (vb.EventNotif, bool) {
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

func (n *Notifier) generateNotifMessage(notifs []vb.EventNotif) string {
	m := ""
	for _, notif := range notifs {
		switch notif.Type {
		case vb.LimitedSpots:
			m += "Limited spots"
		case vb.NewEvent:
			m += "New event"
		}
		m += " - "
		m += fmt.Sprintf("%s on %s\n", notif.Event.Name, notif.Event.StartDate)
	}

	return m
}
