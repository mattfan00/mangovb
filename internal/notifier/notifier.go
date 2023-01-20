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
		for _, notif := range notifs {
			n.logger.WithFields(logrus.Fields{
				"notif_type": notif.TypeId,
				"event_id":   notif.EventId,
			}).Info("Created notif")
		}

		err = n.eventNotifStore.InsertMultiple(notifs)
		if err != nil {
			n.logger.Error(err)
			return
		}

		messages := generateNotifMessages(notifs)
		err = n.bot.SendMessagesToAllChannels(messages)
		for _, err := range multierr.Errors(err) {
			n.logger.Warn(err)
		}
	} else {
		n.logger.Info("No notifs")
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

func generateNotifMessages(notifs []vb.EventNotif) []string {
	messages := []string{}

	newEventCount := map[vb.EventSource]int{}
	limitedSpotsMessages := []string{}
	for _, notif := range notifs {
		switch notif.TypeId {
		case vb.LimitedSpots:
			msg := fmt.Sprintf(
				"**%s**: %d spot(s) left | %s :calendar_spiral: %s :round_pushpin: %s",
				vb.EventSourceMap[notif.Event.SourceId],
				notif.Event.SpotsLeft,
				notif.Event.Name,
				notif.Event.StartTime.Format("Mon Jan 02 3:04 PM"),
				notif.Event.Location,
			)
			limitedSpotsMessages = append(limitedSpotsMessages, msg)
		case vb.NewEvent:
			newEventCount[notif.Event.SourceId] += 1
		}
	}

	for sourceId, count := range newEventCount {
		msg := fmt.Sprintf("**%s**: %d new event(s)", vb.EventSourceMap[sourceId], count)
		messages = append(messages, msg)
	}

	messages = append(messages, limitedSpotsMessages...)

	return messages
}
