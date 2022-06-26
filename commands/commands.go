package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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
		Usage: `You: !ping
Bot: pong`,
		Command: Ping,
	},
	"help": {
		Name:        "help",
		Description: "Will provide you with a link to the list of commands",
		Usage: `You: !help
Bot: [link to commands]`,
		Command: Help,
	},
	"image": {
		Name:        "image",
		Description: "Will post a random image based on your search term",
		Usage: `You: !image [SEARCH TERM]
Bot: [LINK TO RANDOM IMAGE WHICH WILL GET EMBEDDED]`,
		Command: RandomImage,
	},
}

func GenerateReadme() error {
	fmt.Println("generating readme file")
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	readmePath := fmt.Sprintf("%s/README.md", wd)
	fmt.Printf("readme file at %s\n", readmePath)

	data, err := ioutil.ReadFile(readmePath)

	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	startLine := -1
	for i, line := range lines {
		if strings.Contains(line, "## commands") {
			startLine = i + 1
		}
	}

	if startLine == -1 {
		return errors.New("Could not find '## commands'")
	}

	lines = lines[0 : startLine+1]

	for name, command := range Commands {
		lines = append(lines, "<details>")
		lines = append(lines, fmt.Sprintf("<summary>%s</summary>", name))
		lines = append(lines, "<br>")
		lines = append(lines, command.Description)
		if command.Usage != "" {
			lines = append(lines, "\n### usage\n")
			lines = append(lines, "```")
			lines = append(lines, strings.Trim(command.Usage, "\n"))
			lines = append(lines, "```")
		}
		lines = append(lines, "</details>")
	}

	ioutil.WriteFile(readmePath, []byte(strings.Join(lines, "\n")), 0777)
	fmt.Println("Generated README file")
	return nil
}
