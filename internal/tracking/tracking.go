package tracking

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const trackingDir = ".mob"
const trackingFile = "tracking.json"

// TrackingData holds the tracking information for all issues
type TrackingData struct {
	Issues map[string]IssueTracking `json:"issues"`
}

// IssueTracking holds the tracking information for a single issue
type IssueTracking struct {
	ForkPoint        string   `json:"fork_point"`
	LastMergedCommit string   `json:"last_merged_commit"`
	MergedCommits    []string `json:"merged_commits"`
}

// getTrackingPath returns the path to the tracking file
func getTrackingPath() (string, error) {
	// Get git root directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, trackingDir, trackingFile), nil
}

// Load loads the tracking data from disk
func Load() (*TrackingData, error) {
	path, err := getTrackingPath()
	if err != nil {
		return nil, err
	}

	data := &TrackingData{
		Issues: make(map[string]IssueTracking),
	}

	file, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return data, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, data); err != nil {
		return nil, err
	}

	return data, nil
}

// Save saves the tracking data to disk
func (t *TrackingData) Save() error {
	path, err := getTrackingPath()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetIssueTracking returns the tracking data for an issue
func (t *TrackingData) GetIssueTracking(issue string) IssueTracking {
	if tracking, ok := t.Issues[issue]; ok {
		return tracking
	}
	return IssueTracking{
		MergedCommits: []string{},
	}
}

// UpdateIssueTracking updates the tracking data for an issue
func (t *TrackingData) UpdateIssueTracking(issue string, lastCommit string, commits []string) {
	tracking := t.GetIssueTracking(issue)
	tracking.LastMergedCommit = lastCommit
	tracking.MergedCommits = append(tracking.MergedCommits, commits...)
	t.Issues[issue] = tracking
}

// SetForkPoint sets the fork point for an issue (called when wip branch is created)
func (t *TrackingData) SetForkPoint(issue string, forkPoint string) {
	tracking := t.GetIssueTracking(issue)
	tracking.ForkPoint = forkPoint
	t.Issues[issue] = tracking
}

// GetForkPoint returns the fork point for an issue
func (t *TrackingData) GetForkPoint(issue string) string {
	return t.GetIssueTracking(issue).ForkPoint
}

// GetUnmergedCommits returns commits that haven't been merged yet
func (t *TrackingData) GetUnmergedCommits(issue string, allCommits []string) []string {
	tracking := t.GetIssueTracking(issue)
	mergedSet := make(map[string]bool)
	for _, c := range tracking.MergedCommits {
		mergedSet[c] = true
	}

	var unmerged []string
	for _, c := range allCommits {
		if !mergedSet[c] {
			unmerged = append(unmerged, c)
		}
	}
	return unmerged
}
