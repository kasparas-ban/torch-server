package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(dsn string) {
	var err error
	db, err = ConnectDB(dsn, &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB(connString string, options *gorm.Config) (db *gorm.DB, err error) {
	numOfRetries := 6
	timeoutSec := 10

	for {
		if numOfRetries <= 0 {
			panic(fmt.Sprintf("Could not connect to the database after %ds", timeoutSec*numOfRetries))
		}

		db, err = gorm.Open(mysql.Open(connString), options)
		fmt.Println("Trying to connect: ", connString)
		if err != nil {
			time.Sleep(time.Duration(timeoutSec) * time.Second)
			numOfRetries--
			continue
		}

		log.Println("Connected to the database")
		return db, nil
	}
}

func GetDB() *gorm.DB {
	return db
}
