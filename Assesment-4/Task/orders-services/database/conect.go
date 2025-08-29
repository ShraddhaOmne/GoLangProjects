package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	RETRY_COUNT    = 10
	RETRY_DURATION = 5 // int seconds
)

func GetConnection(dsn string) (*gorm.DB, error) {
	count := 1
retry:
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		if count > RETRY_COUNT {
			return nil, err
		}
		count++
		log.Println("Trying to connect to the database ...for the number of times..", count)
		time.Sleep(time.Second * RETRY_DURATION)
		goto retry
		//return nil, err
	} else {
		return db, nil
	}
}
