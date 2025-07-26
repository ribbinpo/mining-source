package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Dsn string
}

func NewDatabase(Dsn string) *Database {
	return &Database{Dsn}
}

func (d *Database) Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(d.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
