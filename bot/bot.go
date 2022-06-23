package bot

import (
	"DiscordBot/config"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BotId string
)

func Start() *discordgo.Session {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// Making our bot a user
	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fmt.Println("Bot is running!")

	return goBot
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// so the bot never responds to its self
	if m.Author.ID == BotId {
		return
	}

	msg := strings.Trim(m.Content, " ")

	if msg[0:1] != config.BotPrefix {
		return
	}

	msg = strings.Trim(msg[1:], " ")

	if msg == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}
