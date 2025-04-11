package cmd

import (
    "fmt"

    "github.com/spf13/cobra"

    "scantrix/internal/scanner"
    "scantrix/internal/watcher"
)

var scanCmd = &cobra.Command{
    Use:   "scan",
    Short: "Scans the specified path for vulnerabilities",
    Long: `Scans the specified path for PHP code vulnerabilities
using an AST-based scanning engine (placeholder for now).`,
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Running the scan subcommand...")

        // Create the scanner
        s := scanner.NewScanner()

        if watchMode {
            // Start the watcher
            fmt.Printf("Starting watch mode on path: %s\n", pathToScan)
            err := watcher.WatchDir(pathToScan, func(changedFile string) {
                // Re-scan the changed file
                _ = s.ScanFile(changedFile)
            })
            if err != nil {
                return err
            }
            // The WatchDir is a blocking call, so we probably won't return here
            // until the user kills the process.
            return nil
        }

        // No watch mode => just do a one-time scan
        fmt.Printf("Scanning path: %s\n", pathToScan)
        if err := s.ScanPath(pathToScan); err != nil {
            return err
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(scanCmd)
}
