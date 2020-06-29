package config

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
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

func GetRedisConfig() * DBConfig {
	return &DBConfig{

	}
}
