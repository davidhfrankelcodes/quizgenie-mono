// backend/cmd/api/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/auth"
    "github.com/davidhfrankelcodes/quizgenie-backend/internal/bucket"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
)

func main() {
	// Load .env first, so that os.Getenv can pick up JWT_SECRET
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	// Now initialize the JWT module (this will fatal if JWT_SECRET is missing)
	auth.InitJWT()

	// Initialize GORM + Postgres
	db.InitDB()

	// Auto-migrate only the User model for now
	if err := db.DB.AutoMigrate(&auth.User{}); err != nil {
		log.Fatal("AutoMigrate User failed:", err)
	}

	mux := http.NewServeMux()

	// Public routes:
	mux.HandleFunc("/signup", auth.SignupHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.Handle("/buckets", auth.AuthMiddleware(http.HandlerFunc(bucket.CreateBucketHandler)))
	mux.Handle("/buckets/", auth.AuthMiddleware(http.HandlerFunc(bucket.ListBucketsHandler)))


	// Protected example route:
	mux.Handle("/ping", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})))

	// Read port from environment, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Starting API server on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
