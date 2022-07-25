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

type Configs struct {
	Appconfig   AppConfig
	Dbconfig    DbConfig
	Jwtconfig   JwtConfig
	Redisconfig RedisConfig
}

var lock = &sync.Mutex{}
var configs *Configs

func GetInstance() *Configs {
	if configs == nil {
		lock.Lock()
		JwtExpired, _ := strconv.Atoi(os.Getenv("JWT_EXPIRED_SECOND"))

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
		}
		lock.Unlock()
	}
	// fmt.Println("[CONFIG] : ", &configs)
	return configs
}
