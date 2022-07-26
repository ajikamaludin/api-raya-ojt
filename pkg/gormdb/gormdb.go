package gormdb

import (
	"fmt"
	"sync"

	"github.com/ajikamaludin/api-raya-ojt/app/configs"
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/helper"

	// "github.com/ajikamaludin/api-raya-ojt/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var lock = &sync.Mutex{}

type GormDB struct {
	db *gorm.DB
}

func New() *GormDB {
	return &GormDB{}
}

func (gdb *GormDB) GetInstance() (*gorm.DB, error) {
	if gdb.db == nil {
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
		gdb.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		lock.Unlock()
		if err != nil {
			return nil, err
		}

		if configs.Dbconfig.DbIsMigrate {
			// Migrate Here
			gdb.db.AutoMigrate(
				&models.User{},
				&models.Account{},
				&models.Bank{},
				&models.BankAccount{},
				&models.BankAccountFavorite{},
				&models.BankTransaction{},
			)
			helper.Seed(gdb.db)
		}
		return gdb.db, nil
	}
	// fmt.Println("[DATABASE] : ", gdb.db)
	return gdb.db, nil
}
