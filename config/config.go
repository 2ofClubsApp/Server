package config

import (
	"github.com/2-of-clubs/2ofclubs-server/app/handler"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
)

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func GetDBConfig() *DBConfig {
	return &DBConfig{
		Host:     "localhost",
		Port:     5432,
		Name:     "cdb",
		User:     "postgres",
		Password: "postgres",
	}
}

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     "localhost:2345",
		Password: "",
		DB:       0,
	}
}

func GetAdminConfig() *model.User {
	credentials := model.NewCredentials()
	credentials.Username = "admin"
	credentials.Email = "admin@utmsu.ca"
	hashedPass, err := handler.Hash("password")
	if err != nil {
		return nil
	}
	credentials.Password = hashedPass
	user := model.NewUser()
	user.Credentials = credentials
	user.IsAdmin = true
	user.IsApproved = true
	return user
}
