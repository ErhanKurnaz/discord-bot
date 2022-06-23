package main

import (
	"DiscordBot/bot"
	"DiscordBot/commands"
	"DiscordBot/config"
	"fmt"
	"os"
	"strings"
)

func checkForQuit(quitChannel chan struct{}) {
	for {
		var text string
		fmt.Scan(&text)

		if text == "q" || text == "quit" {
			fmt.Println("Shutting down")
			break
		}
	}

	quitChannel <- struct{}{}
}

func startBot() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	goBot := bot.Start()

	if goBot == nil {
		return
	}

	defer goBot.Close()

	quitChannel := make(chan struct{})
	go checkForQuit(quitChannel)
	<-quitChannel
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		startBot()
		return
	}

	if len(args) == 2 && args[0] == "gen" && args[1] == "readme" {
		if err := commands.GenerateReadme(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		return
	}

	fmt.Printf("Unknown command '%s'\n", strings.Join(args, " "))
	fmt.Println("enter no arguments to start the bot or enter 'gen readme' to generate the readme file")
}
