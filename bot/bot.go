package bot

import (
	"DiscordBot/commands"
	"DiscordBot/config"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BotId string
)

func Start() *discordgo.Session {
	goBot, err := discordgo.New("Bot " + config.Config.Token)

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

	if len(msg) == 0 || msg[0:1] != config.Config.BotPrefix {
		return
	}

	msg = strings.Trim(msg[1:], " ")
	if msg == "" {
		s.ChannelMessageSend(m.ChannelID, "No command provided, you ho")
		return
	}

	messageSplit := strings.Split(msg, " ")
	for i := 0; i < len(messageSplit); i++ {
		messageSplit[i] = strings.Trim(messageSplit[i], " ")
	}

	commandType := messageSplit[0]
	args := messageSplit[1:]

	if command, ok := commands.Commands[commandType]; ok {
		command.Command(args, s, m)
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The '%s' command doesn't exist yet, maybe bother Erhan about it and leave me alone", commandType))
	}
}
