package main

import "toussaint/backend/app"

func main() {
	srv := app.SetupRestApi("localhost", 9999)
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
