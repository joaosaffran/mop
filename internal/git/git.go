package git

import (
	"fmt"
	"strings"

	"github.com/joaosaffran/mob/internal/shell"
)

// Run executes a git command with stdout and stderr connected to the terminal
func Run(args ...string) error {
	return shell.Run("git", args...)
}

// Output executes a git command and returns the output
func Output(args ...string) (string, error) {
	output, err := shell.Output("git", args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Checkout switches to the specified branch
func Checkout(branch string) error {
	return Run("checkout", branch)
}

// CheckoutNewBranch creates and switches to a new branch
func CheckoutNewBranch(branch string) error {
	return Run("checkout", "-b", branch)
}

// Pull fetches and merges the latest changes
func Pull() error {
	return Run("pull")
}

// CurrentBranch returns the current branch name
func CurrentBranch() (string, error) {
	return Output("rev-parse", "--abbrev-ref", "HEAD")
}

// BranchExists checks if a branch exists
func BranchExists(branch string) bool {
	_, err := Output("rev-parse", "--verify", branch)
	return err == nil
}

// GetCommitHash returns the commit hash for a ref
func GetCommitHash(ref string) (string, error) {
	return Output("rev-parse", ref)
}

// GetCommitsBetween returns commit hashes between two refs (exclusive base, inclusive head)
func GetCommitsBetween(base, head string) ([]string, error) {
	output, err := Output("log", "--format=%H", fmt.Sprintf("%s..%s", base, head))
	if err != nil {
		return nil, err
	}
	if output == "" {
		return []string{}, nil
	}
	return strings.Split(output, "\n"), nil
}

// GetCommitMessage returns the commit message for a ref
func GetCommitMessage(ref string) (string, error) {
	return Output("log", "-1", "--format=%B", ref)
}

// CherryPick cherry-picks a commit
func CherryPick(commit string) error {
	return Run("cherry-pick", commit)
}

// CommitSquash creates a squash commit with a message
func CommitSquash(message string) error {
	return Run("commit", "-m", message)
}

// Reset resets to a commit
func Reset(commit string, mode string) error {
	return Run("reset", mode, commit)
}

// StashPush stashes current changes
func StashPush() error {
	return Run("stash", "push")
}

// StashPop pops the latest stash
func StashPop() error {
	return Run("stash", "pop")
}

// Merge merges a branch with squash
func MergeSquash(branch string) error {
	return Run("merge", "--squash", "-X", "theirs", branch)
}

// AbortMerge aborts an in-progress merge
func AbortMerge() error {
	return Run("merge", "--abort")
}

// ResetHard resets the current branch to a commit, discarding all changes
func ResetHard(commit string) error {
	return Run("reset", "--hard", commit)
}

// DeleteBranch deletes a local branch
func DeleteBranch(branch string) error {
	return Run("branch", "-D", branch)
}

// Push pushes the current branch to the remote
func Push() error {
	return Run("push")
}

// PushSetUpstream pushes the current branch and sets the upstream
func PushSetUpstream(remote, branch string) error {
	return Run("push", "-u", remote, branch)
}

// Diff returns the diff between two refs
func Diff(base, head string) (string, error) {
	return Output("diff", base, head)
}

// DiffFiles returns list of changed files between two refs
func DiffFiles(base, head string) ([]string, error) {
	output, err := Output("diff", "--name-only", base, head)
	if err != nil {
		return nil, err
	}
	if output == "" {
		return []string{}, nil
	}
	return strings.Split(output, "\n"), nil
}

// DiffStat returns diff statistics between two refs
func DiffStat(base, head string) (string, error) {
	return Output("diff", "--stat", base, head)
}
