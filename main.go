package main

import (
	"./app"
	"./config"
)

func main() {
	dbConfig := config.GetDBConfig()
	api := app.App{}
	api.Initialize(dbConfig)
	api.Run(":8080")
}
