// cmd/worker/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/ai"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/file"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/quiz"
)

func main() {
	// 1) Load .env so REDIS_ADDR, OPENAI_API_KEY, and DB_* are available
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}

	// 2) Initialize OpenAI client
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatal("OPENAI_API_KEY must be set")
	}
	ai.InitOpenAI(openaiKey)

	// 3) Connect to Postgres
	db.InitDB()

	// 4) Get Redis address from env
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR must be set")
	}

	// 5) Create an Asynq server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 5,
		},
	)

	// 6) Build a ServeMux and register handlers
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

		// Delegate to the file service
		file.ProcessFile(payload.FileID)

		log.Printf("Worker: finished ProcessFile (service) for file_id=%d\n", payload.FileID)
		return nil
	})

	// ─── GenerateQuiz ───────────────────────────────────────────────────────────
	mux.HandleFunc("GenerateQuiz", func(ctx context.Context, t *asynq.Task) error {
		// payload is JSON: {"quiz_id":123}
		var payload struct {
			QuizID uint `json:"quiz_id"`
		}
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		quizID := payload.QuizID
		log.Printf("Worker: starting GenerateQuiz for quiz_id=%d\n", quizID)

		// Delegate to the quiz service
		if err := quiz.GenerateQuiz(quizID); err != nil {
			log.Printf("Worker: GenerateQuiz service error for quiz_id=%d: %v\n", quizID, err)
			return err
		}

		log.Printf("Worker: finished GenerateQuiz for quiz_id=%d\n", quizID)
		return nil
	})

	// 7) Run the Asynq server
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Asynq server failed: %v", err)
	}
}
