package cmd

import (
    "fmt"

    "github.com/spf13/cobra"

    "scantrix/internal/tui" // your TUI package
)


var tuiCmd = &cobra.Command{
    Use:   "tui",
    Short: "Launch the TUI for directory navigation",
    RunE: func(cmd *cobra.Command, args []string) error {
        // 1. Build the tree from the provided path
        roots, err := tui.BuildTreeFromPath(pathToScan)
        if err != nil {
            return fmt.Errorf("failed to build tree: %w", err)
        }

        // If there's nothing to show, exit
        if len(roots) == 0 {
            fmt.Println("No watched files or directories found.")
            return nil
        }

        // 2. Create the TUI model
        m := tui.NewModel(roots)

        // 3. Run it
        return tui.StartProgram(m)
    },
}

func init() {
    rootCmd.AddCommand(tuiCmd)
    tuiCmd.Flags().StringVarP(&pathToScan, "path", "p", ".", "Path to scan")
}
