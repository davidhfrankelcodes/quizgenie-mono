// internal/utils/pdf.go
package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtractTextFromPDF opens a PDF file at `path` and returns its full text.
// If the file isn’t a valid PDF, it returns an error.
func ExtractTextFromPDF(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("could not open PDF: %w", err)
	}
	defer f.Close()

	// ledongthuc/pdf requires reading the entire file into memory
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, f); err != nil {
		return "", fmt.Errorf("could not read PDF into buffer: %w", err)
	}

	// Use pdf.NewReader on the buffer
	reader, err := pdf.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return "", fmt.Errorf("pdf.NewReader error: %w", err)
	}

	var output strings.Builder
	// Iterate over each page
	numPages := reader.NumPage()
	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("could not extract text from page %d: %w", i, err)
		}
		output.WriteString(text)
		output.WriteString("\n\n")
	}

	return output.String(), nil
}
