package notifier

import (
	"fmt"
	"time"

	vb "github.com/mattfan00/mangovb"
	"github.com/mattfan00/mangovb/internal/bot"
	"github.com/mattfan00/mangovb/internal/store"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"

	mapset "github.com/deckarep/golang-set/v2"
)

type Notifier struct {
	bot             *bot.Bot
	eventStore      *store.EventStore
	eventNotifStore *store.EventNotifStore
	logger          *logrus.Entry
}

func New(
	bot *bot.Bot,
	eventStore *store.EventStore,
	eventNotifStore *store.EventNotifStore,
	logger *logrus.Entry,
) *Notifier {
	return &Notifier{
		bot:             bot,
		eventStore:      eventStore,
		eventNotifStore: eventNotifStore,
		logger:          logger,
	}
}

func (n *Notifier) Notify() {
	events, err := n.eventStore.GetLatest()
	if err != nil {
		n.logger.Error(err)
		return
	}

	ids := make([]string, len(events))
	for i, event := range events {
		ids[i] = event.Id
	}

	notifMap, err := n.eventNotifStore.GetByEventIds(ids)
	if err != nil {
		n.logger.Error(err)
		return
	}

	notifs := createNotifs(events, notifMap)

	if len(notifs) > 0 {
		err = n.eventNotifStore.InsertMultiple(notifs)
		if err != nil {
			n.logger.Error(err)
			return
		}

		message := generateNotifMessage(notifs)
		err = n.bot.SendMessageToAllChannels(message)
		for _, err := range multierr.Errors(err) {
			n.logger.Warn(err)
		}
	}
}

func createNotifs(events []vb.Event, notifMap map[string][]vb.EventNotif) []vb.EventNotif {
	notifs := []vb.EventNotif{}
	notifCreatedOn := time.Now()

	for _, e := range events {
		prevNotifs := notifMap[e.Id]
		prevNotifIdSet := mapset.NewSet[vb.EventNotifType]()
		for _, prevNotif := range prevNotifs {
			prevNotifIdSet.Add(prevNotif.TypeId)
		}

		if !prevNotifIdSet.Contains(vb.NewEvent) {
			notifs = append(notifs, vb.EventNotif{
				TypeId:    vb.NewEvent,
				EventId:   e.Id,
				Event:     e,
				CreatedOn: notifCreatedOn,
			})
		}

		if !prevNotifIdSet.Contains(vb.LimitedSpots) &&
			e.IsAvailable &&
			e.SpotsLeft > 0 &&
			e.SpotsLeft < 5 {
			notifs = append(notifs, vb.EventNotif{
				TypeId:    vb.LimitedSpots,
				EventId:   e.Id,
				Event:     e,
				CreatedOn: notifCreatedOn,
			})
		}
	}

	return notifs
}

func generateNotifMessage(notifs []vb.EventNotif) string {
	m := ""
	for _, notif := range notifs {
		switch notif.TypeId {
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
