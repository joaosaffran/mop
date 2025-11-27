package llm

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed prompts/*.txt
var promptsFS embed.FS

// PromptTemplate represents a loaded prompt template
type PromptTemplate struct {
	tmpl *template.Template
}

// LoadPrompt loads a prompt template from the prompts folder
func LoadPrompt(name string) (*PromptTemplate, error) {
	content, err := promptsFS.ReadFile(fmt.Sprintf("prompts/%s.txt", name))
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt %s: %w", name, err)
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse prompt template %s: %w", name, err)
	}

	return &PromptTemplate{tmpl: tmpl}, nil
}

// Execute renders the prompt template with the given data
func (p *PromptTemplate) Execute(data any) (string, error) {
	var buf bytes.Buffer
	if err := p.tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}
	return buf.String(), nil
}

// MustLoadPrompt loads a prompt template and panics on error
func MustLoadPrompt(name string) *PromptTemplate {
	p, err := LoadPrompt(name)
	if err != nil {
		panic(err)
	}
	return p
}

// GetSystemPrompt returns the system prompt for code review
func GetSystemPrompt() (string, error) {
	tmpl, err := LoadPrompt("code_review_system")
	if err != nil {
		return "", err
	}
	return tmpl.Execute(nil)
}

// GetUserPrompt returns the user prompt with the diff
func GetUserPrompt(diff string) (string, error) {
	tmpl, err := LoadPrompt("code_review_user")
	if err != nil {
		return "", err
	}
	return tmpl.Execute(map[string]string{"Diff": diff})
}
