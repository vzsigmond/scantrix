// internal/source/source.go
package source

import (
	"fmt"
	"os"
	"os/exec"
)

// CloneGitRepo clones a git repo and returns the temp dir path or an error.
func CloneGitRepo(repoURL string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "scantrix-git-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	cmd := exec.Command("git", "clone", "--depth=1", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return tmpDir, nil
}
