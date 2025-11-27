package llm

import (
	"encoding/json"
	"fmt"
)

// Recommendation represents a code improvement suggestion
type Recommendation struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // "high", "medium", "low"
}

// GetRecommendations fetches code review recommendations from an LLM
func GetRecommendations(diff string) ([]Recommendation, error) {
	if !IsAPIKeySet() {
		return getDefaultRecommendations(), nil
	}

	// Truncate diff if too long
	maxDiffLen := 8000
	if len(diff) > maxDiffLen {
		diff = diff[:maxDiffLen] + "\n... (truncated)"
	}

	// Load prompts from templates
	systemPrompt, err := GetSystemPrompt()
	if err != nil {
		return nil, fmt.Errorf("failed to load system prompt: %w", err)
	}

	userPrompt, err := GetUserPrompt(diff)
	if err != nil {
		return nil, fmt.Errorf("failed to load user prompt: %w", err)
	}

	// Create ChatGPT client and make request
	client := NewChatGPTClient(GetAPIKey())
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	content, err := client.Complete(messages, 0.3, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	// Parse JSON response
	var recommendations []Recommendation
	if err := json.Unmarshal([]byte(content), &recommendations); err != nil {
		// If parsing fails, return default recommendations
		return getDefaultRecommendations(), nil
	}

	return recommendations, nil
}

// getDefaultRecommendations returns default recommendations when LLM is unavailable
func getDefaultRecommendations() []Recommendation {
	return []Recommendation{
		{
			Title:       "Set OPENAI_API_KEY for AI recommendations",
			Description: "Export OPENAI_API_KEY environment variable to enable AI-powered code review suggestions.",
			Severity:    "low",
		},
	}
}

// SeverityColor returns the color for a severity level
func SeverityColor(severity string) string {
	switch severity {
	case "high":
		return "196" // Red
	case "medium":
		return "208" // Orange
	case "low":
		return "81" // Cyan
	default:
		return "250" // Gray
	}
}
