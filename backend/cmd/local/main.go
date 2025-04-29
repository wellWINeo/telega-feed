package main

import (
	"TelegaFeed/cmd/local/local_app"
	"os"
	"strconv"
)

const defaultPort int64 = 8080

func main() {
	app := localapp.NewLocalApp()

	if err := app.Build(); err != nil {
		panic(err)
	}

	if err := app.Setup(); err != nil {
		panic(err)
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
	if err != nil {
		port = defaultPort
	}

	_ = app.Start(int32(port))
}
