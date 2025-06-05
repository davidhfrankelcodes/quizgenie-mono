// backend/cmd/api/main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/auth"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
)

func main() {
	// Load .env (if it exists) so that os.Getenv can read DB_HOST, DB_USER, etc.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	// Initialize GORM + Postgres
	db.InitDB()

	// Auto‐migrate the User model
	if err := db.DB.AutoMigrate(&auth.User{}); err != nil {
		log.Fatal("AutoMigrate User failed:", err)
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"pong"}`))
	})

	// Read port from environment (loaded from .env or actual env), default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Starting API server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
