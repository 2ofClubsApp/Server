package config

import (
	"github.com/2-of-clubs/2ofclubs-server/app/handler"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
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
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		}
}

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	}
}

func GetAdminConfig() *model.User {
	credentials := model.NewCredentials()
	credentials.Username = os.Getenv("ADMIN_USERNAME")
	credentials.Email = os.Getenv("ADMIN_EMAIL")
	hashedPass, err := handler.Hash(os.Getenv("ADMIN_PASSWORD"))
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
