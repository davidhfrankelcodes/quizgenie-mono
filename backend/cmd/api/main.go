// cmd/api/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/ai"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/auth"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/bucket"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/file"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/quiz"
)

func main() {
	// 1) Load .env so JWT_SECRET, REDIS_ADDR, OPENAI_API_KEY, etc. are available
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	// 2) Initialize JWT (will log.Fatal if JWT_SECRET is missing)
	auth.InitJWT()

	// 3) Initialize OpenAI client
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatal("OPENAI_API_KEY must be set")
	}
	ai.InitOpenAI(openaiKey)

	// 4) Initialize GORM + Postgres
	db.InitDB()

	// 5) Auto-migrate all models:
	//    - User (auth)
	//    - Bucket (bucket)
	//    - File and FileChunk (file)
	//    - Quiz, Question, Answer, Attempt, AttemptAnswer (quiz)
	if err := db.DB.AutoMigrate(
		&auth.User{},
		&bucket.Bucket{},
		&file.File{},
		&file.FileChunk{},
		&quiz.Quiz{},
		&quiz.Question{},
		&quiz.Answer{},
		&quiz.Attempt{},
		&quiz.AttemptAnswer{},
	); err != nil {
		log.Fatal("AutoMigrate models failed:", err)
	}

	mux := http.NewServeMux()

	// Public routes:
	mux.HandleFunc("/signup", auth.SignupHandler)
	mux.HandleFunc("/login", auth.LoginHandler)

	// All /buckets and nested routes require authentication.
	mux.Handle("/buckets", auth.AuthMiddleware(http.HandlerFunc(handleBucketsRoot)))
	mux.Handle("/buckets/", auth.AuthMiddleware(http.HandlerFunc(handleBucketsRoot)))

	// Quiz‐specific routes:
	mux.Handle("/quizzes/", auth.AuthMiddleware(http.HandlerFunc(handleQuizzesRoot)))

	// Attempt‐detail route:
	mux.Handle("/attempts/", auth.AuthMiddleware(http.HandlerFunc(quiz.GetAttemptDetailsHandler)))

	// Protected ping (example)
	mux.Handle("/ping", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})))

   // Public healthcheck endpoint (no auth)
   mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
       w.WriteHeader(http.StatusOK)
       w.Write([]byte(`{"status":"ok"}`))
   })

	// Read PORT from environment, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Starting API server on", addr)

	log.Fatal(http.ListenAndServe(addr, mux))
}

// handleBucketsRoot dispatches:
//   - POST   /buckets                → CreateBucketHandler
//   - GET    /buckets                → ListBucketsHandler
//   - POST   /buckets/{id}/files     → UploadFileHandler
//   - GET    /buckets/{id}/files     → ListFilesHandler
//   - POST   /buckets/{id}/quizzes   → CreateQuizHandler
//   - GET    /buckets/{id}/attempts  → ListAttemptsHandler
func handleBucketsRoot(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	// 1) POST   /buckets
	if (path == "/buckets" || path == "/buckets/") && method == http.MethodPost {
		bucket.CreateBucketHandler(w, r)
		return
	}

	// 2) GET    /buckets
	if (path == "/buckets" || path == "/buckets/") && method == http.MethodGet {
		bucket.ListBucketsHandler(w, r)
		return
	}

	// 3) POST   /buckets/{id}/files
	if strings.HasPrefix(path, "/buckets/") && strings.HasSuffix(path, "/files") && method == http.MethodPost {
		file.UploadFileHandler(w, r)
		return
	}

	// 4) GET    /buckets/{id}/files
	if strings.HasPrefix(path, "/buckets/") && strings.HasSuffix(path, "/files") && method == http.MethodGet {
		file.ListFilesHandler(w, r)
		return
	}

	// 5) POST   /buckets/{id}/quizzes
	if strings.HasPrefix(path, "/buckets/") && strings.HasSuffix(path, "/quizzes") && method == http.MethodPost {
		quiz.CreateQuizHandler(w, r)
		return
	}

	// 6) GET    /buckets/{id}/attempts
	if strings.HasPrefix(path, "/buckets/") && strings.HasSuffix(path, "/attempts") && method == http.MethodGet {
		quiz.ListAttemptsHandler(w, r)
		return
	}

	http.NotFound(w, r)
}

// handleQuizzesRoot dispatches:
//   - GET  /quizzes/{quizId}           → GetQuizStatusHandler
//   - GET  /quizzes/{quizId}/questions → GetQuizQuestionsHandler
//   - POST /quizzes/{quizId}/attempts  → SubmitQuizHandler
func handleQuizzesRoot(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	// GET /quizzes/{quizId}
	if segments := strings.Split(path, "/"); len(segments) == 3 && segments[1] == "quizzes" && method == http.MethodGet {
		quiz.GetQuizStatusHandler(w, r)
		return
	}

	// GET /quizzes/{quizId}/questions
	if strings.HasPrefix(path, "/quizzes/") && strings.HasSuffix(path, "/questions") && method == http.MethodGet {
		quiz.GetQuizQuestionsHandler(w, r)
		return
	}

	// POST /quizzes/{quizId}/attempts
	if strings.HasPrefix(path, "/quizzes/") && strings.HasSuffix(path, "/attempts") && method == http.MethodPost {
		quiz.SubmitQuizHandler(w, r)
		return
	}

	http.NotFound(w, r)
}
