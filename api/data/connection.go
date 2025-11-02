package data

import (
	"errors"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dataSourceName string = "host=postgresql user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=America/Toronto"
)

var (
	db     *gorm.DB = nil
	models []any    = make([]any, 5)
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

	log.Fatalln(db.AutoMigrate(models...))
}

func GetGormDB() (*gorm.DB, error) {
	if db == nil {
		return db, errors.New("No db instanced found")
	}

	return db, nil
}

func AddModelsToGormSetUp(new_models ...any) {
	models = append(models, new_models...)
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
