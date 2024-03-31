package summary

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

type OpenAISummarizer struct {
	client  *openai.Client
	promt   string
	enabled bool
	mu      sync.Mutex
}

func NewOpenAISummarizer(apiKey string, promt string) *OpenAISummarizer {
	s := &OpenAISummarizer{
		client: openai.NewClient(apiKey),
		promt:  promt,
	}

	log.Printf("[INFO] OpenAI summarizer enabled: %v", apiKey != "")

	if apiKey != "" {
		s.enabled = true
	}

	return s
}

func (s *OpenAISummarizer) Summarize(ctx context.Context, text string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		return "", nil
	}

	request := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("%s%s", text, s.promt),
			},
		},
		MaxTokens:   256,
		Temperature: 0.7,
		TopP:        1,
	}

	resp, err := s.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	rawSumary := strings.TrimSpace(resp.Choices[0].Message.Content)
	if strings.HasSuffix(rawSumary, ".") {
		return rawSumary, nil
	}

	sentences := strings.Split(rawSumary, ".")

	return strings.Join(sentences[:len(sentences)-1], ".") + ".", nil
}
