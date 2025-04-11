// internal/source/local.go
package source

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolveLocalPath checks that a path exists and returns its absolute form.
func ResolveLocalPath(p string) string {
	if p == "" {
		fmt.Println("❌ Please specify a local path using --path or use --git to scan a repo.")
		os.Exit(1)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to resolve path: %v\n", err)
		os.Exit(1)
	}
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "❌ Invalid path: %s\n", abs)
		os.Exit(1)
	}
	return abs
}
