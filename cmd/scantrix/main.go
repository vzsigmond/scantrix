package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"os/exec"

	"scantrix/internal/rules"
	"scantrix/internal/scanner"
	"scantrix/internal/ui"
	"scantrix/internal/logger"
)

func main() {
	logger.Log("🚨 Scantrix started")
	// Define CLI flags
	excludeFlag := flag.String("exclude", "", "Regex to exclude files/folders (e.g. 'node_modules|tests')")
	severityFlag := flag.String("severity", "", "Filter by severity (e.g. critical, warning, info)")
	watchFlag := flag.Bool("watch", false, "Enable real-time scanning (rescan on file change)")
	gitRepo := flag.String("git", "", "Scan a GitHub repo by URL (e.g. --git https://github.com/user/repo)")

	flag.Parse()

	var targetPath string

	if *gitRepo != "" {
		tmpDir, err := os.MkdirTemp("", "scantrix-git-")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create temp dir: %v\n", err)
			os.Exit(1)
		}

		logger.Log("⬇️ Cloning repo %s into %s", *gitRepo, tmpDir)
		cmd := exec.Command("git", "clone", "--depth=1", *gitRepo, tmpDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Git clone failed: %v\n", err)
			os.Exit(1)
		}

		targetPath = tmpDir
	} else {

		args := flag.Args()
		if len(args) < 1 {
			fmt.Println("Usage: scantrix [--watch] [--exclude=...] [--git=url] /path/to/project")
			os.Exit(1)
		}
		targetPath = args[0]
	}

	// Compile exclude regex if provided
	var excludeRegex *regexp.Regexp
	if *excludeFlag != "" {
		var err error
		excludeRegex, err = regexp.Compile(*excludeFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid exclude regex: %v\n", err)
			os.Exit(1)
		}
	}

	// Load all rules
	allRules := rules.LoadAll()

	// Scan
	fmt.Println("🔍 Scanning", targetPath)
	findings, err := scanner.ScanDirectory(targetPath, allRules, excludeRegex, *severityFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Scan error: %v\n", err)
		os.Exit(1)
	}

	// If no results, show confirmation
	if len(findings) == 0 {
		fmt.Println("✅ No vulnerabilities found!")
		return
	}

	// Run TUI
	if *watchFlag {
		err = ui.RunRealtime(targetPath, rules.LoadAll(), excludeRegex, *severityFlag)
	} else {
		err = ui.Run(findings)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}

}
