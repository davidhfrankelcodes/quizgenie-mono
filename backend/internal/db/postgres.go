// backend/internal/db/postgres.go
package db

import (
  "fmt"
  "log"
  "os"

  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

var DB *gorm.DB

// InitDB opens a connection to Postgres via GORM and sets the global DB var.
func InitDB() {
  // Read connection info from environment variables:
  host := os.Getenv("DB_HOST")
  port := os.Getenv("DB_PORT")
  user := os.Getenv("DB_USER")
  pass := os.Getenv("DB_PASSWORD")
  name := os.Getenv("DB_NAME")

  dsn := fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
    host, user, pass, name, port,
  )

  var err error
  DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    log.Fatalf("failed to connect to database: %v", err)
  }
  log.Println("✅ Connected to Postgres via GORM")
}
