package main

import (
	"fmt"
	"os"
	"toussaint/clients/telegram/app"
	"toussaint/clients/telegram/app/cache"
	"toussaint/clients/telegram/app/cli"
	"toussaint/clients/telegram/app/srv"
)

func main() {

	cache.Init()

	params, err := cli.ParseCLI()
	if err != nil {
		os.Exit(0)
	}

	srv.APIEndpoint = fmt.Sprintf("http://%s/v1", params.Backend)

	tg, err := app.NewTelegram(params.TelegramBotToken, params.Debug)
	if err != nil {
		fmt.Printf("[ERR] Failed initialization of telegram bot: %+v", err)
		os.Exit(-1)
	}

	tg.Start()
}
