package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/joaosaffran/mob/internal/errors"
	"github.com/joaosaffran/mob/internal/git"
	"github.com/joaosaffran/mob/internal/github"
	"github.com/joaosaffran/mob/internal/tracking"
	"github.com/joaosaffran/mob/internal/ui"
	"github.com/spf13/cobra"
)

// buildOptions converts issues to huh options
func buildOptions(issues []github.Issue) []huh.Option[int] {
	options := make([]huh.Option[int], len(issues))
	for i, issue := range issues {
		label := fmt.Sprintf("#%d - %s", issue.Number, issue.Title)
		options[i] = huh.NewOption(label, issue.Number)
	}
	return options
}

// selectIssue displays a menu to select an issue and returns the selected issue number
func selectIssue(issues []github.Issue) (int, error) {
	if len(issues) == 0 {
		return 0, fmt.Errorf("no issues assigned to you")
	}

	options := buildOptions(issues)
	return ui.ShowForm(options, "Select an issue to work on")
}

var initCmd = &cobra.Command{
	Use:   "init [work]",
	Short: "Create and checkout a new wip branch",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var work string

		if len(args) == 0 {
			// Fetch and display issues
			issues, err := github.GetAssignedIssues()
			errors.ExitOnError(err, "Error fetching issues")

			issueNumber, err := selectIssue(issues)
			errors.ExitOnError(err, "Error selecting issue")

			work = fmt.Sprintf("%d", issueNumber)
		} else {
			work = args[0]
		}

		// Sanitize work name for branch
		work = strings.ReplaceAll(work, " ", "-")
		branchName := fmt.Sprintf("wip/%s", work)
		baseBranch, _ := cmd.Flags().GetString("base-branch")

		// If base branch is specified, checkout and pull latest
		if baseBranch != "" {
			errors.ExitOnErrorf(git.Checkout(baseBranch), "Error checking out base branch '%s'", baseBranch)
			errors.ExitOnError(git.Pull(), "Error pulling latest updates")
		}

		// Get current commit hash as fork point before creating branch
		forkPoint, err := git.GetCommitHash("HEAD")
		errors.ExitOnError(err, "Error getting current commit")

		// Create and checkout the wip branch
		errors.ExitOnErrorf(git.CheckoutNewBranch(branchName), "Error creating branch '%s'", branchName)

		// Save fork point to tracking
		trackingData, err := tracking.Load()
		errors.ExitOnError(err, "Error loading tracking data")
		trackingData.SetForkPoint(work, forkPoint)
		errors.ExitOnError(trackingData.Save(), "Error saving tracking data")

		fmt.Printf("Created and switched to branch '%s'\n", branchName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("base-branch", "b", "", "Base branch to create the wip branch from")
}
