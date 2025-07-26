package database

import (
	"time"

	"github.com/google/uuid"
)

type StoneModel struct {
	ID        string      `gorm:"primaryKey"`
	Domain    string      `gorm:"not null"`
	Status    string      `gorm:"not null,default:pending"`
	Paths     []PathModel `gorm:"many2many:stone_paths"`
	Version   int         `gorm:"not null,default:1"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime"`
}

type PathModel struct {
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func (s *StoneModel) TableName() string {
	return "stones"
}

func (p *PathModel) TableName() string {
	return "paths"
}

func (s *StoneModel) CreateNew() StoneModel {
	return StoneModel{
		ID:      uuid.New().String(),
		Domain:  s.Domain,
		Status:  "pending",
		Paths:   s.Paths,
		Version: 1,
	}
}
