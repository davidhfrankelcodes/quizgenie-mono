// internal/utils/chunk.go
package utils

// ChunkText splits a large string into slices of at most maxChars characters each.
func ChunkText(fullText string, maxChars int) []string {
	var chunks []string
	for start := 0; start < len(fullText); start += maxChars {
		end := start + maxChars
		if end > len(fullText) {
			end = len(fullText)
		}
		chunks = append(chunks, fullText[start:end])
	}
	return chunks
}
