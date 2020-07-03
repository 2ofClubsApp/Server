package config

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
		Addr:     "localhost:5432",
		Password: "",
		DB:       0,
	}
}
