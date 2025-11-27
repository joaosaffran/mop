package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configDir = ".mob"
const checklistFile = "checklist.yaml"

// Checklist represents the review checklist configuration
type Checklist struct {
	Items []ChecklistItem `yaml:"items"`
}

// ChecklistItem represents a single checklist item
type ChecklistItem struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
}

// getChecklistPath returns the path to the checklist file
func getChecklistPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, configDir, checklistFile), nil
}

// LoadChecklist loads the checklist from the yaml file
func LoadChecklist() (*Checklist, error) {
	path, err := getChecklistPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// Return default checklist if file doesn't exist
		return &Checklist{
			Items: []ChecklistItem{
				{ID: "tests", Description: "All tests pass"},
				{ID: "review", Description: "Code has been self-reviewed"},
				{ID: "docs", Description: "Documentation updated if needed"},
				{ID: "no-debug", Description: "No debug code left behind"},
			},
		}, nil
	}
	if err != nil {
		return nil, err
	}

	var checklist Checklist
	if err := yaml.Unmarshal(data, &checklist); err != nil {
		return nil, err
	}

	return &checklist, nil
}

// SaveChecklist saves the checklist to the yaml file
func SaveChecklist(checklist *Checklist) error {
	path, err := getChecklistPath()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(checklist)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// CreateDefaultChecklist creates a default checklist file
func CreateDefaultChecklist() error {
	checklist := &Checklist{
		Items: []ChecklistItem{
			{ID: "tests", Description: "All tests pass"},
			{ID: "review", Description: "Code has been self-reviewed"},
			{ID: "docs", Description: "Documentation updated if needed"},
			{ID: "no-debug", Description: "No debug code left behind"},
			{ID: "lint", Description: "No lint errors"},
		},
	}
	return SaveChecklist(checklist)
}
