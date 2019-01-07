package main

import "toussaint/backend/app"

func main() {
	srv := app.SetupRestApi("localhost", 9999)
	srv.ListenAndServe()
}
