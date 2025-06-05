package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pgvector/pgvector-go"

	"github.com/davidhfrankelcodes/quizgenie-backend/internal/ai"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/bucket"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/db"
	"github.com/davidhfrankelcodes/quizgenie-backend/internal/utils"
)

// ProcessFile does the full pipeline for a given fileID:
//  1) mark status="processing"
//  2) extract text (PDF or plain text)
//  3) chunk it (~2000 chars each)
//  4) store each chunk in file_chunks and immediately compute + store its embedding
//  5) once at least one chunk is stored, call AI to generate a bucket name
//  6) mark file.status="completed" (or "failed" on error)
func ProcessFile(fileID uint) {
	// 1) Fetch file record
	var frec File
	if err := db.DB.First(&frec, fileID).Error; err != nil {
		log.Printf("[ProcessFile] could not find File ID=%d: %v\n", fileID, err)
		return
	}

	// 2) Update status = "processing"
	if err := db.DB.Model(&frec).Update("status", "processing").Error; err != nil {
		log.Printf("[ProcessFile] failed to set processing status: %v\n", err)
		// continue anyway
	}

	// 3) Extract text from disk
	rawText := ""
	ext := strings.ToLower(filepath.Ext(frec.StoragePath))
	if ext == ".pdf" {
		content, err := utils.ExtractTextFromPDF(frec.StoragePath)
		if err != nil {
			failFile(fileID, fmt.Errorf("PDF extract error: %w", err))
			return
		}
		rawText = content
	} else {
		// treat everything else as plain text
		b, err := os.ReadFile(frec.StoragePath)
		if err != nil {
			failFile(fileID, fmt.Errorf("text file read error: %w", err))
			return
		}
		rawText = string(b)
	}

	// 4) Chunk the text (2000 chars per chunk)
	chunks := utils.ChunkText(rawText, 2000)
	if len(chunks) == 0 {
		failFile(fileID, fmt.Errorf("no text chunks produced"))
		return
	}

	// 5) For each chunk: insert FileChunk row, then compute embedding
	for idx, txt := range chunks {
		ch := FileChunk{
			FileID:     fileID,
			ChunkIndex: idx,
			Content:    txt,
			// Embedding will be populated below
		}
		if err := db.DB.Create(&ch).Error; err != nil {
			log.Printf("[ProcessFile] could not create chunk row: %v\n", err)
			// skip this chunk
			continue
		}

		// 5.a) Compute embedding
		embedVec, err := ai.GetEmbedding(txt)
		if err != nil {
			log.Printf("[ProcessFile] GetEmbedding error (chunk %d): %v\n", idx, err)
			// but keep going for other chunks
			continue
		}

		// 5.b) Convert []float32 → pgvector.Vector and save
		vec := pgvector.NewVector(embedVec)
		if err := db.DB.Model(&FileChunk{}).
			Where("id = ?", ch.ID).
			Update("embedding", vec).Error; err != nil {
			log.Printf("[ProcessFile] could not save embedding for chunk %d: %v\n", idx, err)
		}
	}

	// 6) Once at least one chunk is inserted, generate a bucket name
	//    We’ll fetch the first 3 chunks’ content, concatenate them, and call AI.
	var firstChunks []FileChunk
	if err := db.DB.Where("file_id = ?", fileID).
		Order("chunk_index ASC").
		Limit(3).
		Find(&firstChunks).Error; err != nil {
		log.Printf("[ProcessFile] could not load first few chunks: %v\n", err)
	} else if len(firstChunks) > 0 {
		var combined []string
		for _, c := range firstChunks {
			combined = append(combined, c.Content)
		}
		preview := strings.Join(combined, "\n\n")
		newName, err := ai.GenerateBucketName(preview)
		if err != nil {
			log.Printf("[ProcessFile] GenerateBucketName error: %v\n", err)
		} else {
			// Update bucket name
			if err := db.DB.Model(&bucket.Bucket{}).
				Where("id = ?", frec.BucketID).
				Update("name", newName).Error; err != nil {
				log.Printf("[ProcessFile] could not update bucket name: %v\n", err)
			} else {
				log.Printf("[ProcessFile] bucket %d renamed → %s\n", frec.BucketID, newName)
			}
		}
	}

	// 7) Finally, mark file as "completed"
	if err := db.DB.Model(&frec).Update("status", "completed").Error; err != nil {
		log.Printf("[ProcessFile] failed to set completed status: %v\n", err)
	}
}

// failFile updates file.status="failed" and records the error message.
func failFile(fileID uint, procErr error) {
	errMsg := procErr.Error()
	_ = db.DB.Model(&File{}).
		Where("id = ?", fileID).
		Updates(map[string]interface{}{
			"status":    "failed",
			"error_msg": &errMsg,
		}).Error
	log.Printf("[ProcessFile] file %d failed: %v\n", fileID, procErr)
}
