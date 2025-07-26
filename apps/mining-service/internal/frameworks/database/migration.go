package database

import (
	model "github.com/ribbinpo/mining-service/internal/frameworks/database/model"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.StoneModel{},
		&model.PathModel{},
	)
}
