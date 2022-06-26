package main

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type MqttModel struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Topic     string
	Payload   string
	Qos       uint8
	Retained  bool
	MessageID uint16
	CreatedAt time.Time
	TName     string `gorm:"-"`
}

func InitDataBase(config ConfigT) (*gorm.DB, error) {
	switch config.DbType {
	case "memory":
		return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	case "sqlite":
		return gorm.Open(sqlite.Open(config.DbPath), &gorm.Config{})
	case "mysql":
		dsn := config.DbUsername + ":" + config.DbPassword + "@tcp(" + config.DbHost + ":" + strconv.Itoa(config.DbPort) + ")/" + config.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		dsn := "host=" + config.DbHost + " user=" + config.DbUsername + " password=" + config.DbPassword + " dbname=" + config.DbName + " port=" + strconv.Itoa(config.DbPort) + " sslmode=disable TimeZone=Asia/Shanghai"
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, errors.New("'DbType' error in the config file")
	}
}

func CheckDbTable(db *gorm.DB, name string) error {
	if db.Migrator().HasTable(name) == false {
		if db.Migrator().HasTable(&MqttModel{}) == false {
			log.Println("Create Table: ", name)
			err := db.Migrator().CreateTable(&MqttModel{})
			if err != nil {
				log.Println(err.Error())
				return err
			}
		}
		err := db.Migrator().RenameTable(&MqttModel{}, name)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}
