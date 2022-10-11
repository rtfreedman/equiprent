package db

import (
	"equiprent/internal/util/config"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type db struct {
	connection *gorm.DB
}

var database db

func Connect() (err error) {
	sslmode := "disable"
	if config.Conf.DB.SSL {
		sslmode = "enable"
	}
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=",
		config.Conf.DB.Host,
		os.Getenv(config.Conf.DB.UserENV),
		os.Getenv(config.Conf.DB.PassENV),
		config.Conf.DB.DB,
		config.Conf.DB.Port,
		sslmode,
	)
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	database.connection = db
	return
}
