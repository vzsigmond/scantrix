package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	phpast "scantrix/pkg/php-ast"
)

// Scanner struct could hold configuration/options if needed.
type Scanner struct {
	// Future config fields, e.g. exclude patterns, etc.
}

// NewScanner creates a new instance of the scanning engine.
func NewScanner() *Scanner {
	return &Scanner{}
}

// ScanPath checks if it's a file or directory. If it's a directory, recursively scan .php files.
func (s *Scanner) ScanPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat path %q: %w", path, err)
	}

	if !info.IsDir() {
		// Single file scan
		return s.ScanFile(path)
	}

	// If it's a directory, walk it
	return filepath.Walk(path, func(fullPath string, fileInfo fs.FileInfo, walkErr error) error {
		if walkErr != nil {
			// Skip any path that caused an error
			return nil
		}
		if fileInfo.IsDir() {
			// Don’t scan directories themselves
			return nil
		}
		if strings.HasSuffix(strings.ToLower(fileInfo.Name()), ".php") {
			// Instead of returning error immediately, log the error and continue scanning the rest.
			if err := s.ScanFile(fullPath); err != nil {
				fmt.Printf("Error scanning file %s: %v\n", fullPath, err)
			}
		}
		return nil
	})
}

// ScanFile is a stubbed method that would parse & analyze the file in a real scenario.
// For now, we print a placeholder message and attempt to parse the file using the AST parser.
func (s *Scanner) ScanFile(filePath string) error {
	fmt.Printf("Scanning file: %s\n", filePath)

	// Attempt to parse the file using the AST parser.
	// In a future step, you'll deserialize the AST output and run vulnerability checks.
	astOutput, err := phpast.Parse(filePath)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", filePath, err)
		return err
	}
	fmt.Printf("AST output for %s:\n%s\n", filePath, astOutput)
	return nil
}
