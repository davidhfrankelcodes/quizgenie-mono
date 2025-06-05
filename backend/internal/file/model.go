// internal/file/model.go
package file

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID          uint           `gorm:"primaryKey"`
	BucketID    uint           `gorm:"index;not null"`
	Filename    string         `gorm:"size:255;not null"`
	StoragePath string         `gorm:"size:500;not null"`
	Status      string         `gorm:"size:20;not null"`        // "pending", "processing", "completed", "failed"
	ErrorMsg    *string        `gorm:"type:text"`               // nullable if no error
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
