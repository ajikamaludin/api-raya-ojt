package gormdb

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/ajikamaludin/api-raya-ojt/app/configs"
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/helper"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var lock = &sync.Mutex{}

type GormDB struct {
	Db   *gorm.DB
	Conn *sql.DB
}

func New() *GormDB {
	return &GormDB{}
}

func (gdb *GormDB) GetInstance() (*gorm.DB, error) {
	if gdb.Db == nil {
		configs := configs.GetInstance()

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
			configs.Dbconfig.Host,
			configs.Dbconfig.Username,
			configs.Dbconfig.Password,
			configs.Dbconfig.Dbname,
			configs.Dbconfig.Port,
			configs.Dbconfig.DbTimeZone,
		)
		lock.Lock()
		var err error
		gdb.Conn, err = sql.Open("pgx", dsn)
		gdb.Db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: gdb.Conn,
		}), &gorm.Config{})
		lock.Unlock()
		if err != nil {
			return nil, err
		}

		if configs.Dbconfig.DbIsMigrate {
			// Migrate Here
			gdb.Db.AutoMigrate(
				&models.User{},
				&models.BankTransaction{},
				&models.Bank{},
				&models.Account{},
				&models.BankAccount{},
				&models.BankAccountFavorite{},
			)
			helper.Seed(gdb.Db)
		}
		return gdb.Db, nil
	}
	// fmt.Println("[DATABASE] : ", gdb.db)
	return gdb.Db, nil
}
