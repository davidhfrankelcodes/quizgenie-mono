// backend/internal/bucket/model.go
package bucket

import (
  "time"

  "gorm.io/gorm"
)

type Bucket struct {
  ID          uint           `gorm:"primaryKey"`
  UserID      uint           `gorm:"index;not null"`
  Name        string         `gorm:"size:255;not null"`
  CreatedAt   time.Time
  UpdatedAt   time.Time
  DeletedAt   gorm.DeletedAt `gorm:"index"`
}
