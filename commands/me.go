package commands

import (
	"DiscordBot/util/rest"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getCacheDir(nestedPath string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	path := cacheDir + "/DiscordBot"

	if nestedPath != "" {
		path += "/" + nestedPath
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if e := os.MkdirAll(path, os.ModePerm); e != nil {
			return "", e
		}
	}

	return path, nil
}

func Me(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Author.Avatar)
	avatarUrl := m.Author.AvatarURL("")
	var extension string
	{
		splitUrl := strings.Split(avatarUrl, ".")
		fmt.Println(avatarUrl)
		fmt.Println(splitUrl)
		extension = splitUrl[len(splitUrl)-1]
	}

	cacheDir, err := getCacheDir("avatars")
	if err != nil {
		fmt.Printf("Error getCacheDir() %s\n", err.Error())
		s.ChannelMessageSend(m.ChannelID, "Could not get image")
		return
	}

	path := cacheDir + m.Author.Avatar + "." + extension
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Downloading file %s into file %s\n", avatarUrl, path)
		if e := rest.DownloadFile(avatarUrl, path); e != nil {
			fmt.Printf("Error Download file %s\n", e.Error())
			s.ChannelMessageSend(m.ChannelID, "Could not fetch profile image from discord servers")
			return
		}
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s\n", err.Error())
		s.ChannelMessageSend(m.ChannelID, "Could not open cached profile image")
		return
	}
	defer f.Close()

	message := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://" + path,
			},
		},
		Files: []*discordgo.File{
			{
				Name:   path,
				Reader: f,
			},
		},
	}

	s.ChannelMessageSendComplex(m.ChannelID, message)
	// s.ChannelMessageSend(m.ChannelID, m.Author.AvatarURL(""))
}
