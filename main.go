package main

import (
	"./app"
	"./config"
)

func main() {
	dbConfig := config.GetDBConfig()
	redisConfig := config.GetRedisConfig()
	adminConfig := config.GetAdminConfig()
	if adminConfig != nil {
		api := app.App{}
		api.Initialize(dbConfig, redisConfig, adminConfig)
		api.Run(":8080")
	}
}
