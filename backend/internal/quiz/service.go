// internal/quiz/service.go
package quiz

import (
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
)

// GenerateQuiz is the “service” version of what used to live inline in cmd/worker/main.go.
// It:
//   1) marks quiz.status = "generating"
//   2) (in this dummy example) sleeps, then inserts one question + two answers
//   3) marks quiz.status = "ready" (or "failed")
// If you plan to replace this stub with a real AI‐driven generator, swap out the “sleep + dummy insert” accordingly.
func GenerateQuiz(quizID uint) error {
	// 1) Mark quiz.status = "generating"
	if err := db.DB.Model(&Quiz{}).
		Where("id = ?", quizID).
		Update("status", "generating").Error; err != nil {
		log.Printf("[quiz.GenerateQuiz] failed to set generating status: %v\n", err)
		// continue anyway so we don’t get stuck
	}

	// 2) Simulate “quiz generation” work
	time.Sleep(1 * time.Second)

	// 3) Insert a dummy question + two dummy answers inside a transaction
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// a) Create one question
		q := Question{
			QuizID:      quizID,
			Text:        "What is 2 + 2?",
			Explanation: "Basic arithmetic: 2 + 2 = 4.",
		}
		if err := tx.Create(&q).Error; err != nil {
			return err
		}

		// b) Create two answers
		answers := []Answer{
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
		log.Printf("[quiz.GenerateQuiz] failed to insert dummy question/answers: %v\n", err)
		// mark quiz as "failed"
		_ = db.DB.Model(&Quiz{}).
			Where("id = ?", quizID).
			Update("status", "failed").Error
		return err
	}

	// 4) Finally, mark quiz.status = "ready"
	if err := db.DB.Model(&Quiz{}).
		Where("id = ?", quizID).
		Update("status", "ready").Error; err != nil {
		log.Printf("[quiz.GenerateQuiz] failed to set ready status: %v\n", err)
		return err
	}

	log.Printf("[quiz.GenerateQuiz] successfully generated quiz_id=%d\n", quizID)
	return nil
}