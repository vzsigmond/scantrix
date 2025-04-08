package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"scantrix/internal/logger"
	"scantrix/internal/rules"
	"scantrix/internal/scanner"
	"scantrix/internal/ui"
)

var version = "dev" // will be overridden by -ldflags during build

func main() {
	debugFlag := flag.Bool("debug", false, "Enable debug logging to logs/debug.log")
	excludeFlag := flag.String("exclude", "", "Regex to exclude files/folders (e.g. 'node_modules|tests')")
	severityFlag := flag.String("severity", "", "Filter by severity: critical, warning, or info")
	watchFlag := flag.Bool("watch", false, "Enable real-time scanning (rescan on file change)")
	gitRepo := flag.String("git", "", "Scan a GitHub repo by URL (e.g. --git https://github.com/user/repo)")
	pathFlag := flag.String("path", "", "Path to local directory to scan")
	helpFlag := flag.Bool("help", false, "Show help message")
	upgradeFlag := flag.Bool("self-upgrade", false, "Upgrade Scantrix to the latest version")

	flag.Usage = func() {
		fmt.Println(`🛡️ Scantrix - Code Security Scanner

Usage:
  scantrix --path ./myapp [--watch] [--severity=...] [--exclude=...]
  scantrix --git https://github.com/user/repo [--severity=...]

Flags:`)
		flag.PrintDefaults()
		fmt.Println(`
Examples:
  scantrix --path ./my-project --exclude="vendor|tests"
  scantrix --git https://github.com/drupal/drupal --severity=critical --watch`)
	}

	flag.Parse()

	if *debugFlag {
		logger.EnableDebug()
	}

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *upgradeFlag {
		fmt.Println("🔄 Upgrading Scantrix...")
		cmd := exec.Command("bash", "-c", "curl -sSL https://raw.githubusercontent.com/vzsigmond/scantrix/main/scripts/install.sh | bash")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Upgrade failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	checkForNewVersion()

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

	var excludeRegex *regexp.Regexp
	if *excludeFlag != "" {
		var err error
		excludeRegex, err = regexp.Compile(*excludeFlag)
		exitOnError(err, "Invalid exclude regex")
	}

	allRules := rules.LoadAll()
	fmt.Println("🔍 Scanning", targetPath)
	findings, err := scanner.ScanDirectory(targetPath, allRules, excludeRegex, *severityFlag)
	exitOnError(err, "Scan error")

	if len(findings) == 0 {
		fmt.Println("✅ No vulnerabilities found!")
		return
	}

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

func checkForNewVersion() {
	resp, err := http.Get("https://api.github.com/repos/vzsigmond/scantrix/releases/latest")
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	var result struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	if result.TagName != "" && result.TagName != version {
		fmt.Printf("\n🔔 New version available: %s (you are on %s)\n", result.TagName, version)
		fmt.Println("👉 Run `scantrix --self-upgrade` to upgrade.")
	}
}
