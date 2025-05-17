package helpers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"xel/shared"
)

// DownloadTemplate downloads a template from a given URL to the templates directory
// Supports:
// - Local file paths
// - HTTP/HTTPS URLs for zip files
// - Git repository URLs
func DownloadTemplate(templateURL, templateName string) error {
	// Validate template name
	if templateName == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	// Prepare the destination path
	templatesPath := filepath.Join(shared.XelDir(), "templates")
	templatePath := filepath.Join(templatesPath, templateName)

	// Check if template already exists
	if _, err := os.Stat(templatePath); err == nil {
		return fmt.Errorf("template '%s' already exists", templateName)
	}

	// Parse the URL
	parsedURL, err := url.Parse(templateURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Determine download method based on URL or file path
	switch {
	case parsedURL.Scheme == "" && isLocalPath(templateURL):
		// Local file or directory
		return copyLocalTemplate(templateURL, templatePath)
	case parsedURL.Scheme == "http" || parsedURL.Scheme == "https":
		// HTTP/HTTPS URL
		// Check if it's a GitHub repository URL
		if strings.Contains(templateURL, "github.com") {
			return downloadGitHubTemplate(templateURL, templatePath)
		}
		// Regular HTTP/HTTPS download
		return downloadHTTPTemplate(templateURL, templatePath)
	case isGitURL(templateURL):
		// Git repository
		return cloneGitTemplate(templateURL, templatePath)
	default:
		return fmt.Errorf("unsupported template source: %s", templateURL)
	}
}

// isLocalPath checks if the path is a local file or directory
func isLocalPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// copyLocalTemplate copies a local file or directory to the templates directory
func copyLocalTemplate(src, dest string) error {
	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	// Check if source is a directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return CopyDir(src, dest)
	}

	// If it's a file (including zip), copy it
	return copyFile(src, dest)
}

// downloadHTTPTemplate downloads a zip file from HTTP/HTTPS URL
func downloadHTTPTemplate(templateURL, templatePath string) error {
	// Download the file
	resp, err := http.Get(templateURL)
	if err != nil {
		return fmt.Errorf("failed to download template: %v", err)
	}
	defer resp.Body.Close()

	// Check if it's a zip file
	if !isZipFile(templateURL) {
		return fmt.Errorf("only zip files are supported for HTTP downloads")
	}

	// Create the destination file
	out, err := os.Create(templatePath + ".zip")
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Unzip the file
	return unzipFile(templatePath+".zip", templatePath)
}

// downloadGitHubTemplate downloads a GitHub repository as a zip
func downloadGitHubTemplate(templateURL, templatePath string) error {
	// Parse the URL to extract query parameters
	parsedURL, err := url.Parse(templateURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Extract repository owner and name from the original URL
	repoOwner := "dev-kas"
	repoName := "XelTemplates"
	branch := "master" // default branch

	// Get the specific template from query parameter
	specificTemplate := parsedURL.Query().Get("path")

	zipURL := fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.zip", repoOwner, repoName, branch)

	// Download the zip file
	resp, err := http.Get(zipURL)
	if err != nil {
		return fmt.Errorf("failed to download GitHub repository: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: HTTP status %d", resp.StatusCode)
	}

	// Create a temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "xel-template-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the destination zip file
	zipFilePath := filepath.Join(tempDir, "template.zip")
	out, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer out.Close()

	// Write the body to file
	writtenBytes, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write zip file: %v", err)
	}

	// Verify file size
	if writtenBytes == 0 {
		return fmt.Errorf("downloaded zip file is empty")
	}

	// Close the file before unzipping
	out.Close()

	// Unzip the file
	err = unzipFile(zipFilePath, tempDir)
	if err != nil {
		return fmt.Errorf("failed to unzip file: %v", err)
	}

	// Find the extracted repository directory
	extractedDirs, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read extracted directory: %v", err)
	}

	var extractedRepoDir string
	for _, entry := range extractedDirs {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), fmt.Sprintf("%s-", repoName)) {
			extractedRepoDir = filepath.Join(tempDir, entry.Name())
			break
		}
	}

	if extractedRepoDir == "" {
		return fmt.Errorf("no repository directory found in extracted files")
	}

	// If a specific template is requested, look for it within the extracted repository
	if specificTemplate != "" {
		// Look for the specific template within the extracted repository
		var specificTemplateDir string
		found := false

		// First, check if the template is a direct subdirectory of the extracted repo
		err = filepath.Walk(extractedRepoDir, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			// Check if the current directory matches the template name
			if info.IsDir() && info.Name() == specificTemplate {
				specificTemplateDir = path
				found = true
				return filepath.SkipDir // Stop walking once found
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("error searching for template: %v", err)
		}

		// If not found directly, try a more exhaustive search
		if !found {
			return fmt.Errorf("specific template '%s' not found in repository", specificTemplate)
		}

		// Update the extracted repo directory to the specific template directory
		extractedRepoDir = specificTemplateDir
	}

	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(templatePath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create template parent directory: %v", err)
	}

	// Copy the extracted files to the final template path
	err = CopyDir(extractedRepoDir, templatePath)
	if err != nil {
		return fmt.Errorf("failed to copy template files: %v", err)
	}

	return nil
}

// CopyDir recursively copies a directory
func CopyDir(src, dest string) error {
	// Create destination directory
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := CopyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy files
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

// unzipFile extracts a zip file to a destination directory
func unzipFile(zipPath, destPath string) error {
	// Open the zip file
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	// Create the destination directory
	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Extract files
	var extractedFiles []string
	for _, f := range r.File {
		fpath := filepath.Join(destPath, f.Name)

		// Check for zip slip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create parent directories if they don't exist
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create parent directory: %v", err)
		}

		// Create the file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open zip entry: %v", err)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("failed to write file %s: %v", f.Name, err)
		}

		extractedFiles = append(extractedFiles, f.Name)
	}

	// Log extracted files for debugging
	if len(extractedFiles) == 0 {
		return fmt.Errorf("no files extracted from zip")
	}

	fmt.Printf("Extracted %d files from zip\n", len(extractedFiles))

	return nil
}

// cloneGitTemplate clones a git repository
func cloneGitTemplate(repoURL, templatePath string) error {
	// Normalize the URL to ensure it's a valid git clone URL
	repoURL = normalizeGitURL(repoURL)

	// Clone the repository
	_, err := git.PlainClone(templatePath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: nil, // You can replace with a progress writer if needed
	})

	return err
}

// isZipFile checks if the URL points to a zip file
func isZipFile(urlStr string) bool {
	return strings.HasSuffix(strings.ToLower(urlStr), ".zip")
}

// isGitURL checks if the URL is a git repository URL
func isGitURL(urlStr string) bool {
	gitSuffixes := []string{".git"}
	gitPrefixes := []string{"git://", "git@", "https://github.com/", "https://gitlab.com/"}

	for _, prefix := range gitPrefixes {
		if strings.HasPrefix(urlStr, prefix) {
			return true
		}
	}

	for _, suffix := range gitSuffixes {
		if strings.HasSuffix(urlStr, suffix) {
			return true
		}
	}

	return false
}

// normalizeGitURL converts GitHub URLs to a git clone-compatible format
func normalizeGitURL(urlStr string) string {
	// If it's already a git URL, return as-is
	if strings.HasSuffix(urlStr, ".git") {
		return urlStr
	}

	// Convert GitHub web URLs to git clone URLs
	if strings.Contains(urlStr, "github.com") {
		// Remove any specific folder paths
		parts := strings.Split(urlStr, "/")
		if len(parts) >= 5 {
			repoURL := fmt.Sprintf("https://github.com/%s/%s.git", parts[3], parts[4])
			return repoURL
		}
	}

	return urlStr
}
