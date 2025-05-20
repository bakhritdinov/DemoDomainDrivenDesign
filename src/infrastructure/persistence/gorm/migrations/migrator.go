package migrations

import (
	"DDD/src/domain"
	"gorm.io/gorm"
)

type GormMigrator struct {
	db *gorm.DB
}

func NewGormMigrator(db *gorm.DB) *GormMigrator {
	return &GormMigrator{db: db}
}

func (m *GormMigrator) Run() error {
	return m.db.AutoMigrate(
		&domain.Post{},
		&domain.PostComment{},
	)
}
