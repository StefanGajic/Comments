package database

import (
	"github.com/Comments/internal/comment"
	"github.com/jinzhu/gorm"
)

// MigrateDB migrates database and creates comment table
func MigrateDB(db *gorm.DB) error {
	if result := db.AutoMigrate(&comment.Comment{}); result.Error != nil {
		return result.Error
	}
	return nil
}
