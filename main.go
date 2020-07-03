package main

import (
	"github.com/2-of-Clubs/2ofclubs-server/app"
	"github.com/2-of-Clubs/2ofclubs-server/config"
)

func main() {
	dbConfig := config.GetDBConfig()
	redisConfig := config.GetRedisConfig()
	api := app.App{}
	api.Initialize(dbConfig, redisConfig)
	api.Run(":8080")
}
