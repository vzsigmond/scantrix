package phpast

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// getParserScriptPath locates the parser.sh script.
// It assumes that the repository root is one directory above the binary’s directory,
// and that parser.sh is located in the repository root.
func getParserScriptPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	baseDir := filepath.Dir(exePath)
	// Assume repo root is one directory up from the binary’s directory.
	repoRoot := filepath.Join(baseDir, "..")
	parserScript := filepath.Join(repoRoot, "parser.sh")
	if _, err := os.Stat(parserScript); os.IsNotExist(err) {
		// Fallback: use the current working directory.
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", cwdErr)
		}
		parserScript = filepath.Join(cwd, "parser.sh")
		if _, err := os.Stat(parserScript); os.IsNotExist(err) {
			return "", fmt.Errorf("parser script not found in repository root (%s) or current working directory (%s)", filepath.Join(repoRoot, "parser.sh"), parserScript)
		}
	}
	return parserScript, nil
}


// Parse executes the parser.sh script with the provided arguments.
// For example, to run a PHP AST parser on a given PHP file,
// you could call Parse("path/to/test.php").
// The output of the script (usually JSON representing the AST) is returned.
func Parse(args ...string) (string, error) {
	parserScript, err := getParserScriptPath()
	if err != nil {
		return "", err
	}
	cmd := exec.Command(parserScript, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("failed to run parser script: %w; output: %s", err, out.String())
	}
	return out.String(), nil
}
