package cli

import (
	"fmt"
	"strings"

	"github.com/joaosaffran/mob/internal/config"
	"github.com/joaosaffran/mob/internal/git"
	"github.com/joaosaffran/mob/internal/tracking"
	"github.com/joaosaffran/mob/internal/ui"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review changes before updating PR",
	Long: `Shows a diff of all changes and a checklist to verify before updating.
All checklist items must be checked before update is allowed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get current branch
		currentBranch, err := git.CurrentBranch()
		if err != nil {
			return fmt.Errorf("error getting current branch: %w", err)
		}

		// Verify we're on a wip branch
		if !strings.HasPrefix(currentBranch, "wip/") {
			return fmt.Errorf("not on a wip branch. Please checkout a wip/<issue> branch first")
		}

		// Extract issue from branch name
		issue := strings.TrimPrefix(currentBranch, "wip/")
		wipBranch := currentBranch

		// Load tracking data
		trackingData, err := tracking.Load()
		if err != nil {
			return fmt.Errorf("error loading tracking data: %w", err)
		}

		// Get fork point from tracking
		forkPoint := trackingData.GetForkPoint(issue)
		if forkPoint == "" {
			return fmt.Errorf("no fork point found. Was this branch created with 'mob init'?")
		}

		// Get diff
		diff, err := git.Diff(forkPoint, wipBranch)
		if err != nil {
			return fmt.Errorf("error getting diff: %w", err)
		}

		if diff == "" {
			fmt.Println("No changes to review")
			return nil
		}

		// Get diff stats
		diffStat, err := git.DiffStat(forkPoint, wipBranch)
		if err != nil {
			diffStat = "Unable to get diff stats"
		}

		// Load checklist
		checklist, err := config.LoadChecklist()
		if err != nil {
			return fmt.Errorf("error loading checklist: %w", err)
		}

		// Convert config.ChecklistItem to ui.ChecklistItem
		uiItems := make([]ui.ChecklistItem, len(checklist.Items))
		for i, item := range checklist.Items {
			uiItems[i] = ui.ChecklistItem{Description: item.Description}
		}

		// Run review UI
		completed, err := ui.RunReview(diff, diffStat, issue, uiItems)
		if err != nil {
			return fmt.Errorf("error running review UI: %w", err)
		}

		// Check if review is complete
		if completed {
			fmt.Println("\n✓ Review complete! You can now run 'mob update' to push changes.")
		} else {
			fmt.Println("\n✗ Review incomplete. Please check all items before updating.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}
