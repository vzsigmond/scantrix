package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"scantrix/internal/logger"
	"scantrix/internal/rules"
	"scantrix/internal/scanner"
	"scantrix/internal/ui"
)

func main() {
	logger.Log("🚨 Scantrix started")

	// CLI flags
	excludeFlag := flag.String("exclude", "", "Regex to exclude files/folders (e.g. 'node_modules|tests')")
	severityFlag := flag.String("severity", "", "Filter by severity: critical, warning, or info")
	watchFlag := flag.Bool("watch", false, "Enable real-time scanning (rescan on file change)")
	gitRepo := flag.String("git", "", "Scan a GitHub repo by URL (e.g. --git https://github.com/user/repo)")
	pathFlag := flag.String("path", "", "Path to local directory to scan")
	debugFlag := flag.Bool("debug", false, "Enable debug logging to logs/debug.log")
	helpFlag := flag.Bool("help", false, "Show help message")


	// Custom help
	flag.Usage = func() {
		fmt.Println(`🛡️ Scantrix - Code Security Scanner

Usage:
  scantrix --path ./myapp [--watch] [--severity=...] [--exclude=...]
  scantrix --git https://github.com/user/repo [--severity=...]

Flags:`)
		flag.PrintDefaults()
		fmt.Println(`
Examples:
  scantrix --path ./my-project --exclude="vendor|tests --watch"
  scantrix --git https://github.com/drupal/drupal --severity=critical`)
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *debugFlag {
		logger.EnableDebug()
	}

	var targetPath string

	if *gitRepo != "" {
		tmpDir, err := os.MkdirTemp("", "scantrix-git-")
		exitOnError(err, "Failed to create temp dir")

		logger.Log("⬇️ Cloning repo %s into %s", *gitRepo, tmpDir)
		cmd := exec.Command("git", "clone", "--depth=1", *gitRepo, tmpDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		exitOnError(cmd.Run(), "Git clone failed")

		targetPath = tmpDir
	} else {
		if *pathFlag == "" {
			fmt.Println("❌ Please specify a local path using --path or use --git to scan a repo.")
			flag.Usage()
			os.Exit(1)
		}
		targetPath = *pathFlag
	}

	// Compile exclude regex
	var excludeRegex *regexp.Regexp
	if *excludeFlag != "" {
		var err error
		excludeRegex, err = regexp.Compile(*excludeFlag)
		exitOnError(err, "Invalid exclude regex")
	}

	// Load rules
	allRules := rules.LoadAll()

	// Scan
	fmt.Println("🔍 Scanning", targetPath)
	findings, err := scanner.ScanDirectory(targetPath, allRules, excludeRegex, *severityFlag)
	exitOnError(err, "Scan error")

	if len(findings) == 0 {
		fmt.Println("✅ No vulnerabilities found!")
		return
	}

	// Run TUI
	if *watchFlag {
		err = ui.RunRealtime(targetPath, allRules, excludeRegex, *severityFlag)
	} else {
		err = ui.Run(findings)
	}
	exitOnError(err, "TUI error")
}

func exitOnError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
		os.Exit(1)
	}
}
