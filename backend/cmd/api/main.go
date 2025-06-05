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

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/auth"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/bucket"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/file"
)

func main() {
	// Load .env first, so that os.Getenv can pick up JWT_SECRET, REDIS_ADDR, etc.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	// Initialize JWT (will log.Fatal if JWT_SECRET missing)
	auth.InitJWT()

	// Initialize GORM + Postgres
	db.InitDB()

	// Auto-migrate User, Bucket, and File models
	if err := db.DB.AutoMigrate(&auth.User{}, &bucket.Bucket{}, &file.File{}); err != nil {
		log.Fatal("AutoMigrate models failed:", err)
	}

	mux := http.NewServeMux()

	// Public routes:
	mux.HandleFunc("/signup", auth.SignupHandler)
	mux.HandleFunc("/login", auth.LoginHandler)

	// All /buckets/* routes require authentication
	mux.Handle("/buckets/", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// POST   /buckets               → create a new bucket
		// GET    /buckets               → list all buckets
		// POST   /buckets/{id}/files    → upload a file into bucket {id}
		// GET    /buckets/{id}/files    → (later) list files in bucket {id}

		path := r.URL.Path
		method := r.Method

		// If exactly "/buckets" or "/buckets/" with POST → CreateBucketHandler
		if (path == "/buckets" || path == "/buckets/") && method == http.MethodPost {
			bucket.CreateBucketHandler(w, r)
			return
		}

		// If exactly "/buckets" or "/buckets/" with GET → ListBucketsHandler
		if (path == "/buckets" || path == "/buckets/") && method == http.MethodGet {
			bucket.ListBucketsHandler(w, r)
			return
		}

		// If path matches "/buckets/{id}/files" with POST → UploadFileHandler
		if strings.HasPrefix(path, "/buckets/") && strings.HasSuffix(path, "/files") && method == http.MethodPost {
			file.UploadFileHandler(w, r)
			return
		}

		// (Placeholder) If path matches "/buckets/{id}/files" with GET → ListFilesHandler (not implemented yet)
		// if strings.HasPrefix(path, "/buckets/") && strings.HasSuffix(path, "/files") && method == http.MethodGet {
		//     file.ListFilesHandler(w, r)
		//     return
		// }

		// If none matched, return 404
		http.NotFound(w, r)
	})))

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
