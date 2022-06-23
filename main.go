package main

import (
	"DiscordBot/bot"
	"DiscordBot/config"
	"fmt"
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

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
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
