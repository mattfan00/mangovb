package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/multierr"
)

type Bot struct {
	Session  *discordgo.Session
	Channels []*Channel
}

type Channel struct {
	Id      string
	GuildId string
}

func New(botToken string) (*Bot, error) {
	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	s.Identify.Intents = discordgo.IntentsGuilds

	bot := &Bot{
		Session:  s,
		Channels: []*Channel{},
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

	for _, channel := range guild.Channels {
		if channel.Name == "volleyball-events" {
			b.Channels = append(b.Channels, &Channel{
				Id:      channel.ID,
				GuildId: guild.ID,
			})
			return
		}
	}

	// create volleyball-events channel if it doesn't exist
	newChannel, err := s.GuildChannelCreate(guild.ID, "volleyball-events", 0)
	if err != nil {
		log.Println(err)
	}
	b.Channels = append(b.Channels, &Channel{
		Id:      newChannel.ID,
		GuildId: guild.ID,
	})

	s.ChannelMessageSend(newChannel.ID, "This channel is for volleyball events")
}

func (b *Bot) SendMessageToAllChannels(message string) error {
	var err error
	for _, channel := range b.Channels {
		_, sendErr := b.Session.ChannelMessageSend(channel.Id, message)
		if sendErr != nil {
			err = multierr.Append(
				err,
				fmt.Errorf("cannot send message to channel %s: %w", channel.Id, sendErr),
			)
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
