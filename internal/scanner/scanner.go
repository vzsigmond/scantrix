package scanner

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"scantrix/internal/types"
)

func ScanDirectory(root string, rules []types.Rule, exclude *regexp.Regexp, severityFilter string) ([]Finding, error) {
	var findings []Finding

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if exclude != nil && exclude.MatchString(path) {
			return nil
		}

		for _, rule := range rules {
			if severityFilter != "" && rule.Severity != severityFilter {
				continue
			}

			if !matchesExtension(path, rule.FileTypes) {
				continue
			}

			file, readErr := os.Open(path)
			if readErr != nil {
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNumber := 1

			for scanner.Scan() {
				line := scanner.Text()
				if rule.Pattern.MatchString(line) {
					findings = append(findings, Finding{
						File:     path,
						Line:     lineNumber,
						RuleID:   rule.ID,
						Severity: rule.Severity,
						Title:    rule.Title,
						Advice:   rule.Advice,
					})
				}
				lineNumber++
			}
		}
		return nil
	})

	return findings, err
}

func matchesExtension(filename string, extensions []string) bool {
	for _, ext := range extensions {
		if filepath.Ext(filename) == ext {
			return true
		}
	}
	return false
}

// Finding is the result of a scan
type Finding struct {
	File     string
	Line     int
	RuleID   string
	Severity string
	Title    string
	Advice   string
}
