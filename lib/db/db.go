//package db decouples db logic from the rest of our applications.
package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

type Port struct {
	Code string `gorm:"primaryKey"`
	Data string
}

func Init() error {
	var err error
	db, err = gorm.Open(sqlite.Open("data/techtest.db"), &gorm.Config{})

	if err != nil {
		return err
	}
	//check db sanity
	err = db.Exec("select 1+1").Error
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Port{})
	if err != nil {
		return err
	}

	return nil
}

func PutPort(p *Port) error {
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(p).Error
}

func GetPorts() chan *Port {
	ret := make(chan *Port)
	go func() {
		close(ret)
		var port *Port
		rows, err := db.Model(&Port{}).Rows()
		if err != nil {
			log.Printf("Error getting ports: %s", err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			db.ScanRows(rows, &port)
			ret <- port
		}

	}()
	return ret
}

func Close() {
	rdb, err := db.DB()
	if err != nil {
		log.Printf("Error closing db: %s", err.Error())
		return
	}
	err = rdb.Close()
	if err != nil {
		log.Printf("Error closing db: %s", err.Error())
	}
}
