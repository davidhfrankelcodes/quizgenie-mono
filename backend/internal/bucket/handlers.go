// backend/internal/bucket/handlers.go
package bucket

import (
  "encoding/json"
  "net/http"

  "github.com/davidhfrankelcodes/quizgenie-backend/internal/auth"
  "github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
)

type createBucketRequest struct {
  // we’ll add fields later (e.g. initial file). For now, empty.
}

type bucketResponse struct {
  ID   uint   `json:"id"`
  Name string `json:"name"`
}

// POST /buckets
func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
  // pull userID from context
  claims, ok := auth.FromContext(r.Context())
  if !ok {
    http.Error(w, "unauthorized", http.StatusUnauthorized)
    return
  }

  // create bucket placeholder
  b := Bucket{
    UserID: claims.UserID,
    Name:   "(processing…)",
  }
  if err := db.DB.Create(&b).Error; err != nil {
    http.Error(w, "could not create bucket", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(bucketResponse{ID: b.ID, Name: b.Name})
}

// GET /buckets
func ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
  claims, ok := auth.FromContext(r.Context())
  if !ok {
    http.Error(w, "unauthorized", http.StatusUnauthorized)
    return
  }

  var buckets []Bucket
  if err := db.DB.Where("user_id = ?", claims.UserID).Find(&buckets).Error; err != nil {
    http.Error(w, "could not fetch buckets", http.StatusInternalServerError)
    return
  }

  var resp []bucketResponse
  for _, b := range buckets {
    resp = append(resp, bucketResponse{ID: b.ID, Name: b.Name})
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(resp)
}
