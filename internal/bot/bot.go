package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

type Bot struct {
	Session  *discordgo.Session
	Channels []*Channel
}

type Channel struct {
	Id      string
	GuildId string
}

func New() (*Bot, error) {
	s, err := discordgo.New("Bot " + viper.GetString("bot_token"))
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
	for _, channel := range b.Channels {
		_, err := b.Session.ChannelMessageSend(channel.Id, message)
		if err != nil {
			log.Println(err)
		}
	}

	// TODO: return multiple errors
	return nil
}

func (b *Bot) Start() error {
	err := b.Session.Open()

	return err
}

func (b *Bot) Stop() {
	b.Session.Close()
}
