package cli_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestScantrixCLI(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	fixturesPath := filepath.Join("tests", "fixtures") // relative to project root

	cmd := exec.Command("go", "run", "./cmd/scantrix", fixturesPath, "--severity=critical")
	cmd.Dir = projectRoot // set working directory for go run
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Scantrix failed to run: %v\nOutput:\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "❌") && !strings.Contains(output, "CRITICAL") {
		t.Errorf("Expected critical finding in output, got:\n%s", output)
	}
}
