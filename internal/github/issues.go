package github

import (
	"encoding/json"
	"fmt"

	"github.com/joaosaffran/mob/internal/shell"
)

// Issue represents a GitHub issue
type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

// GetAssignedIssues fetches issues assigned to the current user using gh CLI
func GetAssignedIssues() ([]Issue, error) {
	output, err := shell.Output("gh", "issue", "list", "--assignee", "joaosaffran", "--json", "number,title")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}

	var issues []Issue
	if err := json.Unmarshal(output, &issues); err != nil {
		return nil, fmt.Errorf("failed to parse issues: %w", err)
	}

	return issues, nil
}
