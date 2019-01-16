package main

import (
	"fmt"
	"os"
	"toussaint/clients/telegram/app"
	"toussaint/clients/telegram/app/cache"
	"toussaint/clients/telegram/app/cli"
)

func main() {

	cache.Init()

	params, err := cli.ParseCLI()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	tg, err := app.NewTelegram(params.TelegramBotToken, params.Debug)
	if err != nil {
		fmt.Printf("Failed initialization of telegram bot: %+v", err)
		os.Exit(-1)
	}

	tg.Start()
}
