// internal/file/handlers.go
package file

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/auth"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/bucket"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"

	"github.com/hibiken/asynq"
)

var queueClient *asynq.Client

// initQueue initializes a global Asynq client (called once at package init).
func init() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr != "" {
		queueClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	}
}

type uploadFileResponse struct {
	FileID uint   `json:"fileId"`
	Status string `json:"status"`
}

// POST /buckets/{bucketId}/files
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 1) Extract userID from JWT‐injected context
	claims, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2) Parse bucketId from URL path: "/buckets/{id}/files"
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}
	bucketID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid bucket ID", http.StatusBadRequest)
		return
	}

	// 3) Verify that this bucket belongs to the current user
	var b bucket.Bucket
	if err := db.DB.
		Where("id = ? AND user_id = ?", bucketID, claims.UserID).
		First(&b).Error; err != nil {
		http.Error(w, "bucket not found", http.StatusNotFound)
		return
	}

	// 4) Parse multipart/form-data (single file field named "file")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "could not parse multipart form", http.StatusBadRequest)
		return
	}
	fileHeader, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "could not read uploaded file", http.StatusBadRequest)
		return
	}
	defer fileHeader.Close()

	// 5) Save the file to disk: FILE_STORAGE_PATH/user_<userID>/bucket_<bucketID>/
	baseDir := os.Getenv("FILE_STORAGE_PATH")
	userDir := filepath.Join(baseDir, "user_"+strconv.Itoa(int(claims.UserID)))
	bucketDir := filepath.Join(userDir, "bucket_"+strconv.FormatUint(bucketID, 10))
	if err := os.MkdirAll(bucketDir, 0o755); err != nil {
		http.Error(w, "could not create storage directory", http.StatusInternalServerError)
		return
	}

	//    Use the "filename" form field if provided, otherwise fallback to original filename
	origFilename := r.FormValue("filename")
	if origFilename == "" {
		origFilename = "uploaded_file"
	}
	dstPath := filepath.Join(bucketDir, filepath.Base(origFilename))
	outFile, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "could not save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, fileHeader); err != nil {
		http.Error(w, "error saving file", http.StatusInternalServerError)
		return
	}

	// 6) Create a new File record in the DB with status="pending"
	f := File{
		BucketID:    uint(bucketID),
		Filename:    filepath.Base(dstPath),
		StoragePath: dstPath,
		Status:      "pending",
	}
	if err := db.DB.Create(&f).Error; err != nil {
		http.Error(w, "could not insert file record", http.StatusInternalServerError)
		return
	}

	// 7) Enqueue background job to process this file (if Asynq client is available)
	if queueClient != nil {
		// JSON-marshal the payload into []byte
		payload, err := json.Marshal(map[string]interface{}{"file_id": f.ID})
		if err != nil {
			log.Printf("failed to marshal ProcessFile payload: %v", err)
		} else {
			task := asynq.NewTask("ProcessFile", payload)
			if _, err := queueClient.Enqueue(task); err != nil {
				// Log enqueue failure but do not block the request
				log.Printf("failed to enqueue ProcessFile task: %v", err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uploadFileResponse{
		FileID: f.ID,
		Status: f.Status,
	})
}

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}
	bucketID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid bucket ID", http.StatusBadRequest)
		return
	}

	var b bucket.Bucket
	if err := db.DB.
		Where("id = ? AND user_id = ?", bucketID, claims.UserID).
		First(&b).Error; err != nil {
		http.Error(w, "bucket not found", http.StatusNotFound)
		return
	}

	var files []File
	if err := db.DB.Where("bucket_id = ?", bucketID).Find(&files).Error; err != nil {
		http.Error(w, "could not fetch files", http.StatusInternalServerError)
		return
	}

	type fileResp struct {
		ID       uint   `json:"id"`
		Filename string `json:"filename"`
		Status   string `json:"status"`
	}
	var out []fileResp
	for _, f := range files {
		out = append(out, fileResp{
			ID:       f.ID,
			Filename: f.Filename,
			Status:   f.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}