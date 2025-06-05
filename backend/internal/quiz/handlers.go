// internal/quiz/handlers.go
package quiz

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/auth"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/bucket"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
	"github.com/hibiken/asynq"
)

var queueClient *asynq.Client

func init() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr != "" {
		queueClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	}
}

type createQuizRequest struct {
	TimedMode    bool `json:"timedMode"`
	PracticeMode bool `json:"practiceMode"`
}

type createQuizResponse struct {
	QuizID uint   `json:"quizId"`
	Status string `json:"status"`
}

// POST /buckets/{bucketId}/quizzes
func CreateQuizHandler(w http.ResponseWriter, r *http.Request) {
	// 1) Auth check
	claims, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2) Parse bucketId from URL
	parts := strings.Split(r.URL.Path, "/")
	bucketID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid bucket ID", http.StatusBadRequest)
		return
	}

	// 3) Verify bucket belongs to user
	var b bucket.Bucket
	if err := db.DB.Where("id = ? AND user_id = ?", bucketID, claims.UserID).First(&b).Error; err != nil {
		http.Error(w, "bucket not found", http.StatusNotFound)
		return
	}

	// 4) Decode request body
	var req createQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// 5) Create initial Quiz record with status='pending'
	q := Quiz{
		BucketID:     uint(bucketID),
		Status:       "pending",
		TimedMode:    req.TimedMode,
		PracticeMode: req.PracticeMode,
	}
	if err := db.DB.Create(&q).Error; err != nil {
		http.Error(w, "could not create quiz", http.StatusInternalServerError)
		return
	}

	// 6) Enqueue GenerateQuizTask
	payload, _ := json.Marshal(map[string]interface{}{"quiz_id": q.ID})
	task := asynq.NewTask("GenerateQuiz", payload)
	if queueClient != nil {
		if _, err := queueClient.Enqueue(task); err != nil {
			log.Printf("failed to enqueue quiz task: %v", err)
		}
	}

	// 7) Return 202 Accepted with quizId
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(createQuizResponse{QuizID: q.ID, Status: q.Status})
}

// GET /quizzes/{quizId}
func GetQuizStatusHandler(w http.ResponseWriter, r *http.Request) {
	// We don’t actually use the claims here, so no need to call auth.FromContext
	parts := strings.Split(r.URL.Path, "/")
	quizID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid quiz ID", http.StatusBadRequest)
		return
	}

	var qrec Quiz
	if err := db.DB.First(&qrec, quizID).Error; err != nil {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": qrec.Status})
}

// GET /quizzes/{quizId}/questions
func GetQuizQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	// Again, we don’t need claims for permission here—just check status.
	parts := strings.Split(r.URL.Path, "/")
	quizID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid quiz ID", http.StatusBadRequest)
		return
	}

	var qrec Quiz
	if err := db.DB.First(&qrec, quizID).Error; err != nil {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}

	if qrec.Status != "ready" {
		http.Error(w, "quiz not ready", http.StatusBadRequest)
		return
	}

	// Fetch questions + answers
	var questions []Question
	db.DB.Where("quiz_id = ?", quizID).Find(&questions)

	type answerResp struct {
		ID          uint    `json:"id"`
		Text        string  `json:"text"`
		IsCorrect   *bool   `json:"isCorrect,omitempty"`
		Explanation *string `json:"explanation,omitempty"`
	}
	type questionResp struct {
		ID          uint          `json:"questionId"`
		Text        string        `json:"text"`
		Explanation *string       `json:"explanation,omitempty"`
		Answers     []answerResp  `json:"answers"`
	}

	var out []questionResp
	for _, q := range questions {
		var ans []Answer
		db.DB.Where("question_id = ?", q.ID).Find(&ans)

		var aresp []answerResp
		for _, a := range ans {
			if qrec.PracticeMode {
				aresp = append(aresp, answerResp{
					ID:          a.ID,
					Text:        a.Text,
					IsCorrect:   &a.IsCorrect,
					Explanation: &a.Explanation,
				})
			} else {
				aresp = append(aresp, answerResp{
					ID:   a.ID,
					Text: a.Text,
				})
			}
		}

		qout := questionResp{
			ID:   q.ID,
			Text: q.Text,
		}
		if qrec.PracticeMode {
			qout.Explanation = &q.Explanation
		}
		qout.Answers = aresp
		out = append(out, qout)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// POST /quizzes/{quizId}/attempts
type submitAnswersReq struct {
	Answers []struct {
		QuestionID uint `json:"questionId"`
		AnswerID   uint `json:"answerId"`
	} `json:"answers"`
}

type submitAnswersResp struct {
	AttemptID uint    `json:"attemptId"`
	Score     float64 `json:"score"`
}

func SubmitQuizHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	quizID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid quiz ID", http.StatusBadRequest)
		return
	}

	var qrec Quiz
	if err := db.DB.First(&qrec, quizID).Error; err != nil {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}
	if qrec.Status != "ready" {
		http.Error(w, "quiz not ready", http.StatusBadRequest)
		return
	}

	var payload submitAnswersReq
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// 1) Create Attempt
	att := Attempt{
		QuizID: uint(quizID),
		UserID: claims.UserID,
		Score:  0, // compute below
	}
	if err := db.DB.Create(&att).Error; err != nil {
		http.Error(w, "could not create attempt", http.StatusInternalServerError)
		return
	}

	// 2) For each answer, check correctness and insert AttemptAnswer
	correctCount := 0
	for _, ans := range payload.Answers {
		var arec Answer
		if err := db.DB.First(&arec, ans.AnswerID).Error; err != nil {
			continue // skip invalid
		}
		isCorr := arec.IsCorrect
		if isCorr {
			correctCount++
		}
		aa := AttemptAnswer{
			AttemptID:  att.ID,
			QuestionID: ans.QuestionID,
			AnswerID:   ans.AnswerID,
			IsCorrect:  isCorr,
		}
		db.DB.Create(&aa)
	}

	// 3) Compute score
	total := len(payload.Answers)
	var score float64
	if total > 0 {
		score = (float64(correctCount) / float64(total)) * 100
	}
	db.DB.Model(&att).Update("score", score)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submitAnswersResp{AttemptID: att.ID, Score: score})
}

// GET /buckets/{bucketId}/attempts
func ListAttemptsHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	bucketID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid bucket ID", http.StatusBadRequest)
		return
	}

	// Join quizzes→attempts to filter for this bucket
	type row struct {
		AttemptID uint      `json:"attemptId"`
		QuizID    uint      `json:"quizId"`
		Score     float64   `json:"score"`
		CreatedAt time.Time `json:"createdAt"`
	}
	var results []row
	db.DB.Table("attempts").
		Select("attempts.id AS attempt_id, attempts.quiz_id, attempts.score, attempts.created_at").
		Joins("JOIN quizzes ON quizzes.id = attempts.quiz_id").
		Where("quizzes.bucket_id = ? AND attempts.user_id = ?", bucketID, claims.UserID).
		Scan(&results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GET /attempts/{attemptId}
func GetAttemptDetailsHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	attemptID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid attempt ID", http.StatusBadRequest)
		return
	}

	var att Attempt
	if err := db.DB.First(&att, attemptID).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if att.UserID != claims.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	type detailRow struct {
		QuestionText       string  `json:"questionText"`
		SelectedAnswerID   uint    `json:"selectedAnswerId"`
		SelectedAnswerText string  `json:"selectedAnswerText"`
		IsCorrect          bool    `json:"isCorrect"`
		CorrectAnswerText  string  `json:"correctAnswerText"`
		Explanation        string  `json:"explanation"`
	}
	var details []detailRow
	db.DB.Raw(`
		SELECT 
			q.text AS question_text,
			aa.answer_id AS selected_answer_id,
			a.text AS selected_answer_text,
			aa.is_correct AS is_correct,
			(SELECT text FROM answers WHERE question_id = q.id AND is_correct = true LIMIT 1) AS correct_answer_text,
			q.explanation AS explanation
		FROM attempt_answers aa
		JOIN questions q ON q.id = aa.question_id
		JOIN answers a ON a.id = aa.answer_id
		WHERE aa.attempt_id = ?
	`, attemptID).Scan(&details)

	resp := map[string]interface{}{
		"attemptId": att.ID,
		"quizId":    att.QuizID,
		"score":     att.Score,
		"details":   details,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
