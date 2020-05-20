package main

import (
	"./app"
	"./config"
)

func main() {
	dbConfig := config.GetConfig()
	api := app.App{}
	api.Initialize(dbConfig)
}
