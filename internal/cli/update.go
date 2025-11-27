package cli

import (
	"fmt"
	"strings"

	"github.com/joaosaffran/mob/internal/git"
	"github.com/joaosaffran/mob/internal/tracking"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Squash wip commits and merge into pr branch",
	Long: `Squash all new commits from wip/<issue> and merge them into pr/<issue>.
Only commits that haven't been merged yet will be included.`,
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
		prBranch := fmt.Sprintf("pr/%s", issue)

		// Load tracking data
		trackingData, err := tracking.Load()
		if err != nil {
			return fmt.Errorf("error loading tracking data: %w", err)
		}

		// Get fork point from tracking
		forkPoint := trackingData.GetForkPoint(issue)
		if forkPoint == "" {
			return fmt.Errorf("no fork point found for this issue. Was this branch created with 'mob init'?")
		}

		// Get all commits in wip branch since fork point
		allCommits, err := git.GetCommitsBetween(forkPoint, wipBranch)
		if err != nil {
			return fmt.Errorf("error getting commits: %w", err)
		}

		if len(allCommits) == 0 {
			fmt.Println("No commits to merge")
			return nil
		}

		// Filter out already merged commits
		unmergedCommits := trackingData.GetUnmergedCommits(issue, allCommits)

		if len(unmergedCommits) == 0 {
			fmt.Println("No new commits to merge")
			return nil
		}

		fmt.Printf("Found %d new commit(s) to merge\n", len(unmergedCommits))

		// Track state for rollback
		prBranchExisted := git.BranchExists(prBranch)
		var prBranchOriginalCommit string
		if prBranchExisted {
			prBranchOriginalCommit, _ = git.GetCommitHash(prBranch)
		}

		// Rollback function to restore state on failure
		rollback := func(errMsg string) error {
			fmt.Println("Rolling back changes...")

			// Abort any in-progress merge
			git.AbortMerge()

			// Go back to wip branch
			git.Checkout(wipBranch)

			// Restore pr branch to original state if it existed
			if prBranchExisted && prBranchOriginalCommit != "" {
				git.Checkout(prBranch)
				git.ResetHard(prBranchOriginalCommit)
				git.Checkout(wipBranch)
			}

			return fmt.Errorf("%s (changes rolled back)", errMsg)
		}

		// Create or checkout pr branch
		if prBranchExisted {
			if err := git.Checkout(prBranch); err != nil {
				return rollback(fmt.Sprintf("error checking out pr branch: %v", err))
			}
		} else {
			// Create pr branch from fork point
			if err := git.Run("checkout", forkPoint); err != nil {
				return rollback(fmt.Sprintf("error checking out fork point: %v", err))
			}
			if err := git.CheckoutNewBranch(prBranch); err != nil {
				return rollback(fmt.Sprintf("error creating pr branch: %v", err))
			}
		}

		// Merge squash from wip branch
		if err := git.MergeSquash(wipBranch); err != nil {
			return rollback(fmt.Sprintf("error merging wip branch: %v", err))
		}

		// Get commit message
		message, _ := cmd.Flags().GetString("message")

		// Create the squash commit
		if err := git.CommitSquash(message); err != nil {
			return rollback(fmt.Sprintf("error creating squash commit: %v", err))
		}

		// Update tracking data
		latestCommit := allCommits[0] // Most recent commit
		trackingData.UpdateIssueTracking(issue, latestCommit, unmergedCommits)
		if err := trackingData.Save(); err != nil {
			return rollback(fmt.Sprintf("error saving tracking data: %v", err))
		}

		// Push to remote
		if err := git.PushSetUpstream("origin", prBranch); err != nil {
			return rollback(fmt.Sprintf("error pushing to remote: %v", err))
		}

		// Switch back to wip branch
		if err := git.Checkout(wipBranch); err != nil {
			// Don't rollback here, the commit was successful
			return fmt.Errorf("error switching back to wip branch: %w (but merge was successful)", err)
		}

		fmt.Printf("Successfully merged %d commit(s) into '%s' and pushed to remote\n", len(unmergedCommits), prBranch)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("message", "m", "", "Commit message for the squash commit")
	updateCmd.MarkFlagRequired("message")
}
