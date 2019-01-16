package main

import (
	"fmt"
	"os"
	"toussaint/backend/app"
	"toussaint/backend/cli"
)

func main() {

	params, err := cli.ParseCLI()
	if err != nil {
		os.Exit(0)
	}

	srv := app.SetupRestApi(params.Host, params.Post)

	updater, err := app.NewUpdater(params.TimeBetweenUpdatesStr, params.UpdateTimeStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting Updater")

	go updater.Start()
	defer updater.Stop()

	fmt.Println("Starting server...")

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
