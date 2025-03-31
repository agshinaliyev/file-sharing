package db

import (
	"file-sharing/config"
	"file-sharing/model/entity"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Baku",
		config.Conf.DBHost, config.Conf.DBPort, config.Conf.DBUser, config.Conf.DBPass, config.Conf.DBName,
	)

	log.Println("Starting database connection: ", config.Conf.DBHost, config.Conf.Port)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal("ActionLog.DbInit.error ", err)
		return err
	}

	DB = db

	err = Migrate()

	if err != nil {
		return err
	}

	return nil
}

func GetDb() *gorm.DB { return DB }

func Migrate() error {
	log.Println("Migration database started: ", config.Conf.DBHost, config.Conf.Port)

	err := GetDb().AutoMigrate(entity.User{})

	if err != nil {
		return err
	}
	return nil
}
