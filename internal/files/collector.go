// internal/files/collector.go
package files

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Collector discovers PHP files under a root directory, excluding logs and optional patterns.
type Collector struct {
	Root     string
	Exclude  *regexp.Regexp
	Exts     map[string]bool
}

func NewCollector(root string, excludePattern string) *Collector {
	var re *regexp.Regexp
	if excludePattern != "" {
		var err error
		re, err = regexp.Compile(excludePattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Invalid exclude regex: %v\n", err)
			os.Exit(1)
		}
	}

	return &Collector{
		Root:    root,
		Exclude: re,
		Exts: map[string]bool{
			".php":    true,
			".module": true,
			".inc":    true,
		},
	}
}

// Collect walks the root directory and returns matching PHP files.
func (c *Collector) Collect() []string {
	var files []string
	_ = filepath.WalkDir(c.Root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".log") {
			return nil
		}
		ext := filepath.Ext(path)
		if c.Exts[ext] && (c.Exclude == nil || !c.Exclude.MatchString(path)) {
			abs, _ := filepath.Abs(path)
			files = append(files, abs)
		}
		return nil
	})
	return files
}
