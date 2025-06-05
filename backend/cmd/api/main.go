// backend/cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/auth"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
)

func main() {
	// Initialize GORM + Postgres
	db.InitDB()

	// Auto‐migrate the User model (creates users table if it doesn’t exist yet)
	if err := db.DB.AutoMigrate(&auth.User{}); err != nil {
		log.Fatal("AutoMigrate User failed:", err)
	}

	// Simple ping endpoint
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"pong"}`))
	})

	// Read port from environment, default to 8080 if not set
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
