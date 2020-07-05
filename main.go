package main

import (
	"./app"
	"./config"
)

func main() {
	dbConfig := config.GetDBConfig()
	redisConfig := config.GetRedisConfig()
	api := app.App{}
	api.Initialize(dbConfig, redisConfig)
	api.Run(":8080")
}
