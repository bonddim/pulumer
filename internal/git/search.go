package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

// SearchFiles searches for files matching the pattern in the working directory
func (g Client) SearchFiles(cwd, pattern string) ([]string, error) {
	var includePattern = regexp.MustCompile(pattern)
	var paths []string

	ignorePatterns, err := g.getIgnorePatterns()
	if err != nil {
		return nil, fmt.Errorf("failed to get gitignore patterns: %v", err)
	}

	// Walk the working directory to find files matching the pattern and not ignored
	err = filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}

		if !d.IsDir() {
			if ignorePatterns.Match(strings.Split(relPath, string(os.PathSeparator)), false) {
				return nil
			}

			if includePattern.MatchString(d.Name()) {
				paths = append(paths, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed search files: %v", err)
	}

	return paths, nil
}

// getIgnorePatterns searches for ".gitignore" file from working to parent directory.
// Returns the matcher with founded patterns.
func (g Client) getIgnorePatterns() (gitignore.Matcher, error) {
	repoRoot, err := g.getRoot()
	if err != nil {
		return nil, err
	}

	// Create a filesystem abstraction from the root directory
	fs := osfs.New(repoRoot)

	// Read gitignore patterns
	patterns, err := gitignore.ReadPatterns(fs, nil)
	if err != nil {
		return nil, err
	}

	// Create a matcher with the patterns
	return gitignore.NewMatcher(patterns), err
}
