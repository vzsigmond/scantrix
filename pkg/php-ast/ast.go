package phpast

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Parse runs the AST parser on a given PHP file.
// It simply looks for parser.sh in the source root and calls it.
// Usage: Parse("path/to/phpfile.php")
func Parse(args ...string) (string, error) {
	// Validate input.
	if len(args) < 1 {
		return "", fmt.Errorf("Usage: Parse(path/to/phpfile.php)")
	}
	phpFile := args[0]
	info, err := os.Stat(phpFile)
	if err != nil || info.IsDir() {
		return "", fmt.Errorf("Error: file '%s' does not exist or is not a regular file", phpFile)
	}

	// Use the current working directory as the source (app) root.
	repoRoot, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Build the path to the parser.sh script.
	parserScript := filepath.Join(repoRoot, "pkg", "php-ast", "bin", "parser.sh")
	if _, err := os.Stat(parserScript); os.IsNotExist(err) {
		return "", fmt.Errorf("parser script not found at %s", parserScript)
	}

	// Execute the parser.sh script with the PHP file as an argument.
	cmd := exec.Command(parserScript, phpFile)
	// Set the working directory to repoRoot so that relative paths in parser.sh resolve correctly.
	cmd.Dir = repoRoot

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("failed to run parser script: %w; output: %s", err, out.String())
	}
	return out.String(), nil
}
