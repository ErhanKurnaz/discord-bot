package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Name        string
	Description string
	Usage       string
	Command     func(args []string, s *discordgo.Session, m *discordgo.MessageCreate)
}

var Commands = map[string]Command{
	"ping": {
		Name:        "ping",
		Description: "Will return the text 'pong'. this is used to test if the bot is still running and working.",
		Usage: `
You: !ping
Bot: pong
		`,
		Command: Ping,
	},
}
