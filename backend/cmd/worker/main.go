// cmd/worker/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/file"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/quiz"
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

	// 5) Build a ServeMux and register handlers
	mux := asynq.NewServeMux()

	// ─── ProcessFile ────────────────────────────────────────────────────────────
	mux.HandleFunc("ProcessFile", func(ctx context.Context, t *asynq.Task) error {
		// payload is JSON: {"file_id":123}
		var payload struct {
			FileID uint `json:"file_id"`
		}
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

		// 5.b) (stub) simulate some work
		time.Sleep(500 * time.Millisecond)

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

	// ─── GenerateQuiz ───────────────────────────────────────────────────────────
	mux.HandleFunc("GenerateQuiz", func(ctx context.Context, t *asynq.Task) error {
		// **CRITICAL:** payload is JSON: {"quiz_id":123}
		//      so we need the `json:"quiz_id"` tag here.
		var payload struct {
			QuizID uint `json:"quiz_id"`
		}
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		quizID := payload.QuizID
		log.Printf("Worker: starting GenerateQuiz for quiz_id=%d\n", quizID)

		// 1) Mark quiz.status = "generating"
		if err := db.DB.Model(&quiz.Quiz{}).
			Where("id = ?", quizID).
			Update("status", "generating").Error; err != nil {
			log.Printf("Worker: failed to set generating status: %v\n", err)
			// continue anyway so we don’t get stuck
		}

		// 2) Simulate “quiz generation” work
		time.Sleep(1 * time.Second)

		// 3) Insert a dummy question + two dummy answers
		err := db.DB.Transaction(func(tx *gorm.DB) error {
			// a) Create one question
			q := quiz.Question{
				QuizID:      quizID,
				Text:        "What is 2 + 2?",
				Explanation: "Basic arithmetic: 2 + 2 = 4.",
			}
			if err := tx.Create(&q).Error; err != nil {
				return err
			}

			// b) Create two answers
			answers := []quiz.Answer{
				{QuestionID: q.ID, Text: "3", IsCorrect: false, Explanation: "No, 2 + 2 is not 3."},
				{QuestionID: q.ID, Text: "4", IsCorrect: true, Explanation: "Yes, 2 + 2 = 4."},
			}
			for _, a := range answers {
				if err := tx.Create(&a).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("Worker: failed to insert dummy question/answers: %v\n", err)
			// mark quiz as "failed"
			_ = db.DB.Model(&quiz.Quiz{}).
				Where("id = ?", quizID).
				Update("status", "failed").Error
			return err
		}

		// 4) Finally, mark quiz.status = "ready"
		if err := db.DB.Model(&quiz.Quiz{}).
			Where("id = ?", quizID).
			Update("status", "ready").Error; err != nil {
			log.Printf("Worker: failed to set ready status: %v\n", err)
			return err
		}

		log.Printf("Worker: finished GenerateQuiz for quiz_id=%d\n", quizID)
		return nil
	})

	// 6) Run the Asynq server
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Asynq server failed: %v", err)
	}
}
