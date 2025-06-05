// internal/ai/openai.go
package ai

import (
	"context"
	"fmt"
	"strings"

	// Make sure you have run:
	//     go get github.com/sashabaranov/go-openai
	goopenai "github.com/sashabaranov/go-openai"
)

var OpenAIClient *goopenai.Client

// InitOpenAI initializes the global OpenAIClient using the provided API key.
// Call this (e.g., in cmd/api/main.go and cmd/worker/main.go) before making any embeddings or chat calls.
func InitOpenAI(apiKey string) {
	OpenAIClient = goopenai.NewClient(apiKey)
}

// GetEmbedding sends `text` to OpenAI’s AdaEmbeddingV2 model and returns a 1536-dimensional embedding.
func GetEmbedding(text string) ([]float32, error) {
	if OpenAIClient == nil {
		return nil, fmt.Errorf("OpenAI client not initialized")
	}
	ctx := context.Background()
	req := goopenai.EmbeddingRequest{
		Model: goopenai.AdaEmbeddingV2,
		Input: []string{text},
	}
	resp, err := OpenAIClient.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateEmbeddings error: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return resp.Data[0].Embedding, nil
}

// GenerateBucketName takes up to the first few chunks (concatenated)
// and returns a short, descriptive bucket name.
func GenerateBucketName(chunks string) (string, error) {
	if OpenAIClient == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}
	ctx := context.Background()
	prompt := "Based on the following text snippets, provide a concise bucket name:\n\n" + chunks
	req := goopenai.ChatCompletionRequest{
		Model: goopenai.GPT4,
		Messages: []goopenai.ChatCompletionMessage{
			{Role: "system", Content: "You are an AI that generates short, descriptive names."},
			{Role: "user", Content: prompt},
		},
		MaxTokens:   16,
		Temperature: 0.7,
	}
	resp, err := OpenAIClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("GenerateBucketName ChatCompletion error: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choice returned from GenerateBucketName")
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

// GenerateQuestions takes a long context string, desired question count and choice count,
// and returns a JSON string that the caller should parse into question objects.
func GenerateQuestions(contextText string, questionCount int, choiceCount int, difficulty string) (string, error) {
	if OpenAIClient == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}
	ctx := context.Background()
	prompt := fmt.Sprintf(
		"Generate %d multiple-choice questions (each with %d answer choices, one correct) from the following context. "+
			"Include an explanation for each question and each choice. Output as valid JSON array of objects. "+
			"Context:\n\n%s",
		questionCount, choiceCount, contextText,
	)

	req := goopenai.ChatCompletionRequest{
		Model: goopenai.GPT4,
		Messages: []goopenai.ChatCompletionMessage{
			{Role: "system", Content: "You are an AI that creates detailed multiple-choice quizzes."},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.7,
	}
	resp, err := OpenAIClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("GenerateQuestions ChatCompletion error: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choice returned from GenerateQuestions")
	}
	return resp.Choices[0].Message.Content, nil
}
