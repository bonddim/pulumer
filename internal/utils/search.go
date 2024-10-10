package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SearchFiles searches for files matching the pattern in the root directory and subdirectories with a maximum depth.
func SearchFiles(root string, maxDepth int, pattern *regexp.Regexp) (paths []string) {
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && withinDepth(root, maxDepth, path) && pattern.MatchString(d.Name()) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return paths
}

// withinDepth checks if the path is within the maximum depth relative to the root.
func withinDepth(root string, maxDepth int, path string) bool {
	// For some reason, filepath.SplitList() doesn't work as expected
	separator := string(filepath.Separator)
	return len(strings.Split(path, separator))-len(strings.Split(root, separator)) <= maxDepth
}
