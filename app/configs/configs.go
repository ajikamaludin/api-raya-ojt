package configs

import (
	"os"
	"strconv"
	"sync"
)

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DbConfig struct {
	Host        string
	Port        string
	Dbname      string
	Username    string
	Password    string
	DbTimeZone  string
	DbIsMigrate bool
}

type JwtConfig struct {
	Secret  string
	Expired int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type GooglePubSubConfig struct {
	ProjectName string
	Credentials string
}

type Configs struct {
	Appconfig          AppConfig
	Dbconfig           DbConfig
	Jwtconfig          JwtConfig
	Redisconfig        RedisConfig
	GooglePubSubConfig GooglePubSubConfig
}

var lock = &sync.Mutex{}
var configs *Configs

func GetInstance() *Configs {
	if configs == nil {
		lock.Lock()
		JwtExpired, _ := strconv.Atoi(os.Getenv("JWT_EXPIRED_SECOND"))

		GOOGLE_CREDENTIALS := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

		if _, err := os.Stat(GOOGLE_CREDENTIALS); err != nil {
			panic("GOOGLE APP CREDENTIALS FILE NOT EXISTS")
		}

		info, _ := os.Stat(GOOGLE_CREDENTIALS)
		if info.IsDir() {
			panic("GOOGLE APP CREDENTIALS MUST BE FILE, NOT DIR")
		}

		configs = &Configs{
			Appconfig: AppConfig{
				Name: os.Getenv("APP_NAME"),
				Env:  os.Getenv("APP_ENV"),
				Port: os.Getenv("APP_PORT"),
			},
			Dbconfig: DbConfig{
				Host:        os.Getenv("DB_HOST"),
				Port:        os.Getenv("DB_PORT"),
				Dbname:      os.Getenv("DB_NAME"),
				Username:    os.Getenv("DB_USER"),
				Password:    os.Getenv("DB_PASS"),
				DbTimeZone:  os.Getenv("DB_TIMEZONE"),
				DbIsMigrate: os.Getenv("DB_ISMIGRATE") == "true",
			},
			Jwtconfig: JwtConfig{
				Secret:  os.Getenv("JWT_SECRET"),
				Expired: int(JwtExpired),
			},
			Redisconfig: RedisConfig{
				Host:     os.Getenv("REDIS_HOST"),
				Port:     os.Getenv("REDIS_PORT"),
				Password: os.Getenv("REDIS_PASSWORD"),
			},
			GooglePubSubConfig: GooglePubSubConfig{
				ProjectName: os.Getenv("GOOGLE_PROJECT_NAME"),
				Credentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
			},
		}
		lock.Unlock()
	}
	// fmt.Println("[CONFIG] : ", &configs)
	return configs
}
