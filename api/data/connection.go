package data

import (
	"errors"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dataSourceName string = "host=postgresql user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
)

var (
	db     *gorm.DB = nil
	models []any
)

func MightInitDB() {
	var err error

	if db != nil {
		return
	}

	db, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		panic("db do not want to talk with you")
	}

	log.Println("migrating...")
	err = db.AutoMigrate(models...)
	if err != nil {
		panic("cant migrate")
	}
}

func GetGormDB() (*gorm.DB, error) {
	if db == nil {
		return db, errors.New("No db instanced found")
	}

	return db, nil
}

func AddModelsToGormSetUp(newModels ...any) {
	models = append(models, newModels...)
}

func CloseGormDB() {
	if db != nil {
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatalln(err)
		}
		err = sqlDb.Close()
		if err != nil {
			return
		}
	}
}
