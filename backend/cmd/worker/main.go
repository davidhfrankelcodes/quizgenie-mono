// cmd/worker/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/file"
)

func main() {
	// 1) Load .env so REDIS_ADDR and DB_* are available
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}

	// 2) Connect to Postgres
	db.InitDB()

	// 3) Get Redis address from env
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR must be set")
	}

	// 4) Create an Asynq server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 5,
		},
	)

	// 5) Build a ServeMux and register the ProcessFile handler
	mux := asynq.NewServeMux()
	mux.HandleFunc("ProcessFile", func(ctx context.Context, t *asynq.Task) error {
		// payload is JSON: {"file_id":123}
		var payload struct{ FileID uint }
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		log.Printf("Worker: starting ProcessFile for file_id=%d\n", payload.FileID)

		// 5.a) Mark status = "processing"
		if err := db.DB.Model(&file.File{}).
			Where("id = ?", payload.FileID).
			Update("status", "processing").Error; err != nil {
			log.Printf("Worker: failed to set processing status: %v\n", err)
			// even if this fails, we’ll try to mark completed below
		}

		// 5.b) (stub) sleep or do nothing
		//    Later: call your real file.ProcessFileService here.

		// 5.c) Finally, mark status = "completed"
		if err := db.DB.Model(&file.File{}).
			Where("id = ?", payload.FileID).
			Update("status", "completed").Error; err != nil {
			log.Printf("Worker: failed to set completed status: %v\n", err)
			return err
		}
		log.Printf("Worker: finished ProcessFile for file_id=%d\n", payload.FileID)
		return nil
	})

	// 6) Run the Asynq server
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Asynq server failed: %v", err)
	}
}
