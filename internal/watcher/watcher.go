package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// OnChangeFunc is called with a file path whenever that file changes.
type OnChangeFunc func(filePath string)

// Example slice if you only want to watch these file extensions:
var WatchedExtensions = []string{".php"}

// WatchDir watches a directory (and all subdirectories) for file changes
// and calls onChange for each relevant event. This version:
// 1) Recursively adds watchers for subdirectories.
// 2) If a new directory is created, adds a watcher for it too.
// 3) Optionally filters by file extension.
func WatchDir(root string, onChange OnChangeFunc) error {
	// Create a new watcher
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// We can’t defer w.Close() here because we run an infinite loop.
	// We'll close it when we exit or user kills the process.
	// (In a CLI tool, that might be fine, or you might need graceful shutdown.)

	// Add watchers recursively
	if err := addAllSubdirs(w, root); err != nil {
		w.Close()
		return fmt.Errorf("failed to watch subdirectories: %w", err)
	}

	// Start event loop
	go func() {
		defer w.Close()

		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}

				if event.Op&(fsnotify.Create|fsnotify.Write) != 0 {
					// If it's a directory, we want to watch it too
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() {
						_ = addAllSubdirs(w, event.Name)
					} else {
						// If it's a file, check extension
						if shouldWatchFile(event.Name) {
							onChange(event.Name)
						}
					}
				}
				// If you want to handle remove or rename, do so here
				// e.g., if event.Op & fsnotify.Remove != 0 ...

			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v\n", err)
			}
		}
	}()

	return nil
}

// addAllSubdirs walks the given path. If path is a directory, it watches it
// and recursively watches all subdirectories as well.
func addAllSubdirs(w *fsnotify.Watcher, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		// If it's a file, just watch the file’s directory
		// (though fsnotify tends to watch directories).
		// Or you can skip if not a directory.
		return nil
	}

	// Watch this directory
	if err := w.Add(path); err != nil {
		return fmt.Errorf("failed to watch directory %s: %w", path, err)
	}

	// Recursively watch children
	return filepath.Walk(path, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			// skip that path
			return nil
		}
		if info.IsDir() {
			// watch the subdirectory
			if e := w.Add(fullPath); e != nil {
				log.Printf("failed to watch %s: %v", fullPath, e)
			}
		}
		return nil
	})
}

// shouldWatchFile checks if the file’s extension is in WatchedExtensions.
// If you only want to watch certain file types, add them to WatchedExtensions.
func shouldWatchFile(path string) bool {
	if len(WatchedExtensions) == 0 {
		// If the user hasn't set any extension filters, watch everything
		return true
	}
	ext := filepath.Ext(path)
	for _, wext := range WatchedExtensions {
		if strings.EqualFold(ext, wext) {
			return true
		}
	}
	return false
}
