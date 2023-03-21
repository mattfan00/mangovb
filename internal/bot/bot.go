package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type Bot struct {
	Session  *discordgo.Session
	Channels map[string]*Channel
	logger   *logrus.Entry
}

type Channel struct {
	Id      string
	GuildId string
}

func New(botToken string, logger *logrus.Entry) (*Bot, error) {
	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	s.Identify.Intents = discordgo.IntentsGuilds

	bot := &Bot{
		Session:  s,
		Channels: map[string]*Channel{},
		logger:   logger,
	}

	bot.Session.AddHandler(bot.ready)
	bot.Session.AddHandler(bot.guildCreate)

	return bot, nil
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "volleyball")
}

func (b *Bot) guildCreate(s *discordgo.Session, guild *discordgo.GuildCreate) {
	if guild.Unavailable {
		return
	}

	b.logger.WithFields(logrus.Fields{
		"guild_id": guild.ID,
	}).Info("Connected to guild")

	foundChannel := false
	for _, channel := range guild.Channels {
		if channel.Name == "volleyball-events" {
			b.Channels[channel.ID] = &Channel{
				Id:      channel.ID,
				GuildId: guild.ID,
			}
			foundChannel = true

			b.logger.WithFields(logrus.Fields{
				"guild_id":   guild.ID,
				"channel_id": channel.ID,
			}).Info("Found channel")

			break
		}
	}

	// create volleyball-events channel if it doesn't exist
	if !foundChannel {
		newChannel, err := s.GuildChannelCreate(guild.ID, "volleyball-events", 0)
		if err != nil {
			b.logger.Error(err)
		}
		b.Channels[newChannel.ID] = &Channel{
			Id:      newChannel.ID,
			GuildId: guild.ID,
		}

		b.logger.WithFields(logrus.Fields{
			"guild_id":   guild.ID,
			"channel_id": newChannel.ID,
		}).Info("Created new channel")

		s.ChannelMessageSend(newChannel.ID, "This channel is for volleyball events")
	}
}

func (b *Bot) SendMessagesToAllChannels(messages []string) error {
	var err error
	for _, channel := range b.Channels {
		for _, msg := range messages {
			_, sendErr := b.Session.ChannelMessageSend(channel.Id, msg)
			if sendErr != nil {
				err = multierr.Append(
					err,
					fmt.Errorf("cannot send message to channel %s: %w", channel.Id, sendErr),
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
