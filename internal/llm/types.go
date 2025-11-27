package llm

// Message represents a chat message for LLM APIs
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RoleSystem is the system role for chat messages
const RoleSystem = "system"

// RoleUser is the user role for chat messages
const RoleUser = "user"

// RoleAssistant is the assistant role for chat messages
const RoleAssistant = "assistant"
