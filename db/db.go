package db

import (
	"database/sql"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*sql.DB
}

var db *gorm.DB

func Init() {
	var err error
	db, err = ConnectDB(os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{ SkipDefaultTransaction: true })
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDB() *gorm.DB {
	return db
}
