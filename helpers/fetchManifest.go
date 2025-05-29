package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"xel/shared"
)

// FetchManifest implements a recursive search for the project manifest file (xel.json).
// The search begins in the specified directory and proceeds upward through the directory
// tree until either the manifest is found or the filesystem root is reached.
//
// This implementation is inspired by Git's repository discovery mechanism, which
// uses a similar approach to locate the .git directory.
//
// The function employs a caching mechanism to optimize performance by storing
// previously resolved manifest locations. The cache uses the current working directory
// as the key and stores the manifest data, its filesystem path, and a boolean
// indicating whether the manifest was found.

type manifestCacheData struct {
	manifest *shared.ProjectManifest
	path     string
	found    bool
}

var manifestCache = map[string]*manifestCacheData{}

// FetchManifest locates and parses the xel.json manifest file by searching upward from the specified directory.
//
// Parameters:
//   - pwd: The current working directory to start searching from
//   - initialPwd: The original directory where the search was initiated (for error reporting)
//
// Returns:
//   - *shared.ProjectManifest: The parsed manifest object if found
//   - string: The absolute path to the manifest file if found
//   - error: An error if the manifest cannot be found or parsed
func FetchManifest(pwd string, initialPwd string) (*shared.ProjectManifest, string, error) {
	// First, check if we have a cached result for this directory
	// This optimization prevents repeated filesystem operations for the same directory
	if cached, exists := manifestCache[pwd]; exists {
		if cached.found {
			// Return cached manifest data if available
			return cached.manifest, cached.path, nil
		}
	}

	// Base case: Stop recursion if we've reached the root directory
	// without finding a manifest file
	if pwd != "/" {
		// Construct the full path to the potential manifest file
		manifestPath := filepath.Join(pwd, "xel.json")

		// Attempt to read the manifest file
		manifestContent, err := os.ReadFile(manifestPath)
		if err != nil {
			// If the manifest doesn't exist in the current directory,
			// recursively search in the parent directory
			return FetchManifest(filepath.Dir(pwd), initialPwd)
		}

		// Initialize a new ProjectManifest structure to hold the parsed data
		manifest := shared.ProjectManifest{}

		// Parse the JSON content into the manifest structure
		if err := json.Unmarshal(manifestContent, &manifest); err != nil {
			return nil, "", fmt.Errorf("failed to parse manifest at %s: %v", manifestPath, err)
		}

		// Cache the successful result to optimize future lookups
		// This includes both the parsed manifest and its filesystem location
		manifestCache[pwd] = &manifestCacheData{
			manifest: &manifest,
			path:     manifestPath,
			found:    true,
		}

		// Check if we are not in the initial directory
		if pwd != initialPwd {
			// Also cache the manifest for the initial lookup directory
			manifestCache[initialPwd] = &manifestCacheData{
				manifest: &manifest,
				path:     manifestPath,
				found:    true,
			}
		}

		return &manifest, manifestPath, nil
	}

	// Cache negative result to avoid redundant searches in the same directory
	manifestCache[initialPwd] = &manifestCacheData{
		found: false,
	}

	// Manifest not found after reaching root
	return nil, "", errors.New("cannot find manifest")
}
