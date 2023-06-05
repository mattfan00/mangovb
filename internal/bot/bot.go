package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/multierr"
)

type Bot struct {
	Session  *discordgo.Session
	Channels []string
}

func New(botToken string) (*Bot, error) {
	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	s.Identify.Intents = discordgo.IntentsGuilds

	bot := &Bot{
		Session:  s,
		Channels: []string{},
	}

	bot.Session.AddHandler(ready)

	return bot, nil
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "volleyball")
}

func (b *Bot) AddChannel(channelId string) {
	if b.Channels == nil {
		b.Channels = []string{}
	}
	b.Channels = append(b.Channels, channelId)
}

func (b *Bot) SendMessagesToAllChannels(messages []string) error {
	var err error
	for _, channel := range b.Channels {
		for _, msg := range messages {
			_, sendErr := b.Session.ChannelMessageSend(channel, msg)
			if sendErr != nil {
				err = multierr.Append(
					err,
					fmt.Errorf("cannot send message to channel %s: %w", channel, sendErr),
				)
			}
		}
	}

	return err
}

func (b *Bot) Start() error {
	err := b.Session.Open()

	return err
}

func (b *Bot) Stop() {
	b.Session.Close()
}
