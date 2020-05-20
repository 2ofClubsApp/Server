package config

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

func GetConfig() *DBConfig {
	return &DBConfig{
		Host:     "localhost",
		Port:     5432,
		Name:     "cdb",
		User:     "postgres",
		Password: "postgres",
	}
}
