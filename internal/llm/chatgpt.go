package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	httpClient "github.com/joaosaffran/mob/internal/http"
)

const (
	OpenAiApiUrl = "https://api.openai.com/v1/chat/completions"
	defaultModel = "gpt-4o-mini"
)

// ChatGPTClient handles communication with the OpenAI ChatGPT API
type ChatGPTClient struct {
	client *httpClient.Client
	model  string
}

// ChatGPTOption is a function that configures a ChatGPTClient
type ChatGPTOption func(*ChatGPTClient)

// WithModel sets the model to use
func WithModel(model string) ChatGPTOption {
	return func(c *ChatGPTClient) {
		c.model = model
	}
}

// NewChatGPTClient creates a new ChatGPT client
func NewChatGPTClient(apiKey string, opts ...ChatGPTOption) *ChatGPTClient {
	client := &ChatGPTClient{
		client: httpClient.NewClient(
			httpClient.WithBearerToken(apiKey),
			httpClient.WithTimeout(60*time.Second),
		),
		model: defaultModel,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Complete sends a chat completion request to the API
func (c *ChatGPTClient) Complete(messages []Message, temperature float64, maxTokens int) (string, error) {
	reqBody := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(OpenAiApiUrl, jsonBody)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(resp.Body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}

// GetAPIKey retrieves the OpenAI API key from environment
func GetAPIKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

// IsAPIKeySet checks if the OpenAI API key is configured
func IsAPIKeySet() bool {
	return GetAPIKey() != ""
}
