package main

import (
	"fmt"
	"toussaint/backend/app"
)

func main() {
	srv := app.SetupRestApi("localhost", 9999)

	updater, err := app.NewUpdater("24h", "02:00")
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
