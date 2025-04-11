package tui

import (
	"os"
	"path/filepath"
	"strings"
)

// Example: we only watch these file extensions.
var WatchedExtensions = []string{".php", ".twig", ".js", ".blade", ".vue", ".ts", ".html"}

// BuildTreeFromPath returns a slice of top-level TreeNodes, each representing
// either a single file (if path is a file) or a directory node with children.
// Directories that contain no watched files (directly or indirectly) are
// pruned so empty directories don't appear.
func BuildTreeFromPath(path string) ([]*TreeNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// If it's a file, just create one node if the scanner wants it:
	if !info.IsDir() {
		if shouldWatchFile(path) {
			node := &TreeNode{
				Name:     filepath.Base(path),
				IsDir:    false,
				Expanded: false,
				Parent:   nil,
				// Optional: Severity or other fields
			}
			return []*TreeNode{node}, nil
		}
		// If not watched, return empty
		return []*TreeNode{}, nil
	}

	// Build a root directory node
	root := &TreeNode{
		Name:     filepath.Base(path),
		IsDir:    true,
		Expanded: true, // expand top-level by default
		Parent:   nil,
	}

	// Recursively walk the directory
	err = filepath.Walk(path, func(fullPath string, info os.FileInfo, werr error) error {
		if werr != nil {
			return werr
		}
		if fullPath == path {
			// skip the root path itself (already created)
			return nil
		}

		rel, _ := filepath.Rel(path, fullPath)
		parts := strings.Split(rel, string(filepath.Separator))

		if info.IsDir() {
			createDirectoryNode(root, parts)
		} else {
			if shouldWatchFile(fullPath) {
				createFileNode(root, parts)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Now remove any empty directories
	pruneEmptyDirs(root)

	// If the root itself became empty, return an empty slice
	if root.IsDir && len(root.Children) == 0 {
		return []*TreeNode{}, nil
	}
	// Otherwise, return a slice with just our single root
	return []*TreeNode{root}, nil
}

func shouldWatchFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, wext := range WatchedExtensions {
		if ext == wext {
			return true
		}
	}
	return false
}

func createDirectoryNode(parent *TreeNode, parts []string) {
	current := parent
	for i, p := range parts {
		if i == len(parts)-1 {
			_ = findOrCreateChildDir(current, p)
		} else {
			current = findOrCreateChildDir(current, p)
		}
	}
}

func createFileNode(parent *TreeNode, parts []string) {
	current := parent
	for i, p := range parts {
		if i == len(parts)-1 {
			f := &TreeNode{
				Name:     p,
				IsDir:    false,
				Expanded: false,
				Parent:   current,
			}
			current.Children = append(current.Children, f)
		} else {
			current = findOrCreateChildDir(current, p)
		}
	}
}

func findOrCreateChildDir(parent *TreeNode, dirName string) *TreeNode {
	for _, c := range parent.Children {
		if c.IsDir && c.Name == dirName {
			return c
		}
	}
	d := &TreeNode{
		Name:     dirName,
		IsDir:    true,
		Expanded: false,
		Parent:   parent,
	}
	parent.Children = append(parent.Children, d)
	return d
}

// pruneEmptyDirs removes any directory node that has zero children
// after pruning. This is a post-order traversal approach:
// 1) prune children
// 2) if this directory is empty, remove it from parent's child list.
func pruneEmptyDirs(node *TreeNode) {
	if !node.IsDir {
		return
	}

	// prune children first
	newChildren := []*TreeNode{}
	for _, c := range node.Children {
		if c.IsDir {
			pruneEmptyDirs(c)
		}
		// If c is a directory with no children, skip it
		if c.IsDir && len(c.Children) == 0 {
			continue // don't add it
		}
		newChildren = append(newChildren, c)
	}
	node.Children = newChildren
}
