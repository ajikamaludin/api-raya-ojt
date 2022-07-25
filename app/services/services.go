package services

import (
	"github.com/ajikamaludin/api-raya-ojt/app/configs"
	"github.com/ajikamaludin/api-raya-ojt/app/repository"
	"github.com/ajikamaludin/api-raya-ojt/pkg/gormdb"
	"github.com/ajikamaludin/api-raya-ojt/pkg/jwtmanager"
	"github.com/ajikamaludin/api-raya-ojt/pkg/redisclient"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/validator"
)

type Services struct {
	Configs     *configs.Configs
	Db          *gormdb.GormDB
	RedisClient *redisclient.RedisClient
	Validator   *validator.Validator
	JwtManager  *jwtmanager.JwtManager
	Repository  *repository.Repository
}

func New() *Services {
	configs := configs.GetInstance()
	validator := validator.New()
	gormdb := gormdb.New()
	redisclient := redisclient.New()
	jwtmanager := jwtmanager.New()

	return &Services{
		Configs:     configs,
		Db:          gormdb,
		RedisClient: redisclient,
		Validator:   validator,
		JwtManager:  jwtmanager,
	}
}
