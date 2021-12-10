//package db decouples db logic from the rest of our applications.
package db

import (
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
		var port *Port
		rows, err := db.Model(&Port{}).Rows()
		defer rows.Close()
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			db.ScanRows(rows, &port)
			ret <- port
		}
		close(ret)
	}()
	return ret
}
