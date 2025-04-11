// internal/cli/flags.go
package cli

import (
	"flag"
	"fmt"
)

// Options represents parsed command-line options.
type Options struct {
	Debug    bool
	Exclude  string
	Severity string
	Watch    bool
	GitRepo  string
	Path     string
	Help     bool
	Upgrade  bool
}

// ParseFlags parses all CLI flags and returns an Options struct.
func ParseFlags() Options {
	opts := Options{}
	flag.BoolVar(&opts.Debug, "debug", false, "Enable debug logging to logs/debug.log")
	flag.StringVar(&opts.Exclude, "exclude", "", "Regex to exclude files/folders")
	flag.StringVar(&opts.Severity, "severity", "", "Filter by severity: critical, warning, or info")
	flag.BoolVar(&opts.Watch, "watch", false, "Enable real-time scanning")
	flag.StringVar(&opts.GitRepo, "git", "", "Scan a GitHub repo by URL")
	flag.StringVar(&opts.Path, "path", "", "Path to local directory to scan")
	flag.BoolVar(&opts.Help, "help", false, "Show help message")
	flag.BoolVar(&opts.Upgrade, "self-upgrade", false, "Upgrade Scantrix to the latest version")
	flag.Parse()
	return opts
}

// PrintUsage displays the usage message for the Scantrix CLI.
func PrintUsage() {
	fmt.Println("🛡️ Scantrix - Code Security Scanner")
	fmt.Println("\nUsage:")
	fmt.Println("  scantrix --path ./myapp [--watch] [--severity=...] [--exclude=...]")
	fmt.Println("  scantrix --git https://github.com/user/repo [--severity=...]")
	fmt.Println("\nFlags:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  scantrix --path ./project --exclude=vendor|tests")
	fmt.Println("  scantrix --git https://github.com/drupal/drupal --watch")
}