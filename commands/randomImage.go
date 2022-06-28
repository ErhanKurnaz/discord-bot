package commands

import (
	"DiscordBot/config"
	"DiscordBot/util/rest"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const unsplash_url = "https://api.unsplash.com"

type unsplashPhotoResponse struct {
	Id   string `json: "id"`
	Urls struct {
		Raw     string `json: "raw"`
		Full    string `json: "full"`
		Regular string `json: "regular"`
		Small   string `json: "small"`
		Thumb   string `json: "thumb"`
	} `json: "urls"`
}

func RandomImage(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Client-ID %s", config.Config.UnsplashAccessKey),
	}

	queries := map[string]string{}

	queries["content_filter"] = "high"
	if len(args) != 0 {
		queries["query"] = strings.Join(args, " ")
	}

	response, err := rest.Get[unsplashPhotoResponse](fmt.Sprintf("%s/photos/random", unsplash_url), queries, headers)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not get images")
		return
	}

	s.ChannelMessageSend(m.ChannelID, response.Urls.Regular)
}
