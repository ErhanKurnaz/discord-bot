package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Config *configStruct
)

type configStruct struct {
	Token             string `json: "Token"`
	BotPrefix         string `json: "BotPrefix"`
	UnsplashAccessKey string `json: "UnsplashAccessKey"`
	UnsplashSecretKey string `json: "UnsplashSecretKey"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))
	err = json.Unmarshal(file, &Config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
