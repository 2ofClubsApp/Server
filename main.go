package main

import (
	"github.com/2-of-clubs/2ofclubs-server/app"
	"github.com/2-of-clubs/2ofclubs-server/config"
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
