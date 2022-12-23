package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func Default() (*discordgo.Session, error) {
	bot, err := discordgo.New("Bot " + viper.GetString("bot_token"))
	if err != nil {
		return nil, err
	}

	bot.Identify.Intents = discordgo.IntentsGuilds

	bot.AddHandler(ready)
	bot.AddHandler(guildCreate)

	return bot, nil
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "volleyball")
}

func guildCreate(s *discordgo.Session, guild *discordgo.GuildCreate) {
	fmt.Println(guild.ID)
	fmt.Println(guild.Name)
	fmt.Println(guild.Permissions)

	for _, channel := range guild.Channels {
		fmt.Printf("channel ID: %s\t", channel.ID)
		fmt.Printf("channel name: %s\t", channel.Name)
		fmt.Printf("channel type: %d\n", channel.Type)
	}

	categoryChannel, err := s.GuildChannelCreate(guild.ID, "volleyball-events", 4)
	if err != nil {
		panic(err)
	}

	advancedChannel, err := s.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
		Name:     "advanced",
		Type:     0,
		ParentID: categoryChannel.ID,
	})
	if err != nil {
		panic(err)
	}

	s.ChannelMessageSend(advancedChannel.ID, "hello")
}
