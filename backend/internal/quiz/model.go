// internal/quiz/model.go
package quiz

import (
  "time"

  "gorm.io/gorm"
)

type Quiz struct {
  ID           uint           `gorm:"primaryKey"`
  BucketID     uint           `gorm:"index;not null"`
  Status       string         `gorm:"size:20;not null"` // 'pending','generating','ready','failed'
  TimedMode    bool           `gorm:"not null"`
  PracticeMode bool           `gorm:"not null"`
  ErrorMsg     *string        `gorm:"type:text"`
  CreatedAt    time.Time
  UpdatedAt    time.Time
  DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Question struct {
  ID        uint           `gorm:"primaryKey"`
  QuizID    uint           `gorm:"index;not null"`
  Text      string         `gorm:"type:text;not null"`
  Explanation string       `gorm:"type:text"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Answer struct {
  ID         uint           `gorm:"primaryKey"`
  QuestionID uint           `gorm:"index;not null"`
  Text       string         `gorm:"type:text;not null"`
  IsCorrect  bool           `gorm:"not null"`
  Explanation string        `gorm:"type:text"`
  CreatedAt  time.Time
  UpdatedAt  time.Time
  DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Attempt struct {
  ID        uint           `gorm:"primaryKey"`
  QuizID    uint           `gorm:"index;not null"`
  UserID    uint           `gorm:"index;not null"`
  Score     float64        `gorm:"not null"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}

type AttemptAnswer struct {
  ID          uint           `gorm:"primaryKey"`
  AttemptID   uint           `gorm:"index;not null"`
  QuestionID  uint           `gorm:"index;not null"`
  AnswerID    uint           `gorm:"index;not null"`
  IsCorrect   bool           `gorm:"not null"`
  CreatedAt   time.Time
  UpdatedAt   time.Time
  DeletedAt   gorm.DeletedAt `gorm:"index"`
}

