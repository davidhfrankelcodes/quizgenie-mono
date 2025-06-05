// backend/internal/auth/model.go
package auth

import (
  "time"

  "gorm.io/gorm"
)

type User struct {
  ID           uint           `gorm:"primaryKey"`
  Username     string         `gorm:"uniqueIndex;size:50"`
  PasswordHash string         `gorm:"size:255"`
  Email        string         `gorm:"size:100"`
  CreatedAt    time.Time
  UpdatedAt    time.Time
  DeletedAt    gorm.DeletedAt `gorm:"index"`
}
