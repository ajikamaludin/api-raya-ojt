package gormdb

import (
	"fmt"
	"sync"

	"github.com/ajikamaludin/api-raya-ojt/app/configs"
	// "github.com/ajikamaludin/api-raya-ojt/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var lock = &sync.Mutex{}
var db *gorm.DB

func GetInstance() (*gorm.DB, error) {
	if db == nil {
		configs := configs.GetInstance()

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			configs.Dbconfig.Host,
			configs.Dbconfig.Username,
			configs.Dbconfig.Password,
			configs.Dbconfig.Dbname,
			configs.Dbconfig.Port,
		)
		lock.Lock()
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		lock.Unlock()
		if err != nil {
			return nil, err
		}

		if configs.Dbconfig.DbIsMigrate {
			// Migrate Here
			// db.AutoMigrate(&models.Note{})
		}
		return db, nil
	}
	// fmt.Println("[DATABASE] : ", db)
	return db, nil
}
