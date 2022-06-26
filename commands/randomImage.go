package commands

import (
	"DiscordBot/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	query := "?content_filter=high"
	if len(args) != 0 {
		query = fmt.Sprintf("%s&query=%s", query, strings.Join(args, " "))
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/photos/random%s", unsplash_url, url.QueryEscape(query)), nil)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not get images")
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", config.Config.UnsplashAccessKey))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not get images")
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Couldn't map body")
		return
	}

	var response unsplashPhotoResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Couldn't map body")
		return
	}

	s.ChannelMessageSend(m.ChannelID, response.Urls.Regular)
}
