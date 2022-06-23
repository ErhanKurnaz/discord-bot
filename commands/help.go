package commands

import "github.com/bwmarrin/discordgo"

func Help(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "Look here for a list of all available commands: https://github.com/ErhanKurnaz/discord-bot#commands")
}
