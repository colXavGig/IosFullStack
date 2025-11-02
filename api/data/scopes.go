package data

import (
	"log"

	"gorm.io/gorm"
)

func GenericPaginate(page uint, pageSize uint) func(statement *gorm.Statement) {
	return func(statement *gorm.Statement) {
		if page == 0 {
			log.Panicln("page must be greater than 0")
		}
		offset := (page - 1) * pageSize

		statement.SQL.WriteString("OFFSET ? LIMIT ?")
		statement.Vars = append(statement.Vars, offset, pageSize)
	}
}
func Paginate(page uint, pageSize uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			log.Panicln("page must be greater than 0")
		}

		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}
