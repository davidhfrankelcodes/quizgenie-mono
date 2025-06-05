package file

import (
	"time"

	"github.com/pgvector/pgvector-go"
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

// FileChunk represents one chunk of text plus its embedding.
type FileChunk struct {
	ID          uint             `gorm:"primaryKey"`
	FileID      uint             `gorm:"index;not null"`
	ChunkIndex  int              `gorm:"not null"`
	Content     string           `gorm:"type:text;not null"`

	// Make Embedding a *pgvector.Vector so that a nil pointer
	// becomes SQL NULL on INSERT.  We’ll fill it later.
	Embedding   *pgvector.Vector `gorm:"type:vector(1536)"`

	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt   `gorm:"index"`
}
