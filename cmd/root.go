package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

// Flag variables
var (
    watchMode  bool
    pathToScan string
    enableTUI  bool
)

// rootCmd is the main command for the CLI tool.
var rootCmd = &cobra.Command{
    Use:   "scantrix",
    Short: "Scantrix is a security scanning tool for your PHP code.",
    Long: `Scantrix is a command-line tool that can scan your PHP code 
for potential vulnerabilities, using an AST approach (not just naive regex).`,
    // The RunE is executed if no subcommand is provided.
    RunE: func(cmd *cobra.Command, args []string) error {
        return runRoot()
    },
}

// Execute is called by main.go to run the CLI.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}

func init() {
    // Define flags (shared by the root command and any subcommands)
    rootCmd.PersistentFlags().StringVarP(&pathToScan, "path", "p", ".", "Path to scan for vulnerabilities")
    rootCmd.PersistentFlags().BoolVarP(&watchMode, "watch", "w", false, "Enable watch mode (rescan on file changes)")
    rootCmd.PersistentFlags().BoolVarP(&enableTUI, "tui", "t", false, "Enable TUI (text-based user interface)")
}

// runRoot is our default action if the user runs "scantrix" without subcommands
func runRoot() error {
    // For now, just demonstrate that flags are being parsed:
    fmt.Println("You must specify a subcommand, e.g.:")
    fmt.Println("  scantrix scan --path ./somefolder")
    fmt.Println("  scantrix tui")
    fmt.Println("Available subcommands:")
    fmt.Println("  scan - Perform a security scan")
    fmt.Println("  tui - Start the text-based user interface")

    // In later steps, we will:
    // 1. Initialize the parser/analyzer.
    // 2. Start the TUI if --tui is enabled.
    // 3. Or do a one-time scan, possibly in a subcommand, etc.

    return nil
}
