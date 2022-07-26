package repository

import (
	"github.com/ajikamaludin/api-raya-ojt/pkg/gormdb"
	"github.com/ajikamaludin/api-raya-ojt/pkg/redisclient"
)

// Any Kind Of DataNeeds from database or redis as data holder is access from repository
type Repository struct {
	Gormdb      *gormdb.GormDB
	RedisClient *redisclient.RedisClient
}

func New(db *gormdb.GormDB, rd *redisclient.RedisClient) *Repository {
	return &Repository{
		Gormdb:      db,
		RedisClient: rd,
	}
}
