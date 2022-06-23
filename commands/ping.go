package commands

import "github.com/bwmarrin/discordgo"

func Ping(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
}
