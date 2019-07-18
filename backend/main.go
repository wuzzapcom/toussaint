package main

import (
	"log"
	"net/http"
	"os"
	"toussaint/backend/app"
	"toussaint/backend/cli"
)

func startUpdater(params *cli.CLIParams) *app.Updater {
	updater, err := app.NewUpdater(params.TimeBetweenUpdatesStr, params.UpdateTimeStr)
	if err != nil {
		panic(err)
	}

	go updater.Start()
	log.Println("started updater")

	return updater
}

func startHTTP(params *cli.CLIParams) *http.Server {
	srv := app.SetupRestAPI(params.Host, params.Post, params.Debug)

	log.Println("start HTTP server")
	var starter = func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("http server has been stopped: %+v", err)
		}
	}

	starter() // run as go starter() to achieve asynchrony

	return srv
}

func main() {

	params, err := cli.ParseCLI()
	if err != nil {
		os.Exit(0)
	}

	startUpdater(params)

	startHTTP(params)

}
