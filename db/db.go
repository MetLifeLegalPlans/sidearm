package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/MetLifeLegalPlans/sidearm/config"
	"github.com/MetLifeLegalPlans/sidearm/db/models"
)

var Conn *gorm.DB
var setup bool

func Setup(conf *config.Config) {
	if setup {
		return
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(sqlite.Open(conf.DbPath), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(models.All...)

	Conn = db
}
