package helpers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"xel/shared"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type VersionData struct {
	manifestPath string
	manifest     *shared.ProjectManifest
}

// ResolveModuleLocal takes a module name and attempts to find it in the configured module paths
// It returns the path to the module's manifest, the manifest itself and an error if the module is not found
func ResolveModuleLocal(moduleName string, constraint string) (string, *shared.ProjectManifest, error) {
	// Convert constraint to semver constraint
	versionConstraint, err := semver.NewConstraint(constraint)
	if err != nil {
		return "", nil, err
	}

	// Collect versions
	collectedVersions := GetVersions(moduleName)

	// filter out versions that don't match the constraint
	filteredVersions := map[*semver.Version]VersionData{}
	for version, manifest := range collectedVersions {
		if versionConstraint.Check(version) {
			filteredVersions[version] = VersionData{
				manifestPath: manifest.manifestPath,
				manifest:     manifest.manifest,
			}
		}
	}

	// Let's see if we got anything collected in our filteredVersions slice
	if len(filteredVersions) != 0 {
		// Yep, we definitely got some

		// Pick the highest version
		var highestVersion *semver.Version
		for version := range filteredVersions {
			if highestVersion == nil || version.GreaterThan(highestVersion) {
				highestVersion = version
			}
		}

		return filteredVersions[highestVersion].manifestPath, filteredVersions[highestVersion].manifest, nil
	}

	// If we get here, the module wasn't found in any paths
	return "", nil, fmt.Errorf("cannot resolve any module named `%s` with version constraint `%s`", moduleName, constraint)
}

var versionsCache = map[string]map[*semver.Version]VersionData{}

// GetVersions takes a module name and returns a map of versions to their manifests
func GetVersions(moduleName string) map[*semver.Version]VersionData {
	// Check if we have a cached version
	if _, ok := versionsCache[moduleName]; ok {
		// Return the cached version
		return versionsCache[moduleName]
	}

	// Not cached, let's collect versions
	collectedVersions := map[*semver.Version]VersionData{}
	// Get module paths from shared config
	config := shared.XelConfig

	// Let's warn the user if no module paths are configured
	if len(config.ModulePaths) == 0 {
		shared.ColorPalette.Warning.Println("No module paths configured")
	}

	// Loop through each module path
	for _, modulePath := range config.ModulePaths {
		// Read all directories in the current module path
		dirs, err := os.ReadDir(modulePath)
		if err != nil {
			continue // Skip if we can't read this path
		}

		// Loop through each directory
		for _, dir := range dirs {
			if dir.IsDir() {
				dirPath := filepath.Join(modulePath, dir.Name())
				subdirs, err := os.ReadDir(dirPath)
				if err != nil {
					continue // Skip if we can't read this path
				}

				for _, subdir := range subdirs {
					if subdir.IsDir() {
						subdirPath := filepath.Join(dirPath, subdir.Name())
						subdirVersion, err := semver.NewVersion(subdir.Name())
						if err != nil {
							continue // Skip if we can't parse the version
						}

						// Read its manifest
						manifestPath := filepath.Join(subdirPath, "xel.json")
						manifestContent, err := os.ReadFile(manifestPath)
						if err != nil {
							continue // Skip if we can't read the manifest
						}

						manifest := shared.ProjectManifest{}
						if err := json.Unmarshal(manifestContent, &manifest); err != nil {
							continue // Skip if we can't parse the manifest
						}

						if manifest.Name != moduleName {
							continue // Skip if the module name doesn't match
						}

						// Double check module version from the manifest to ensure it matches the constraint
						if manifest.Version != subdirVersion.String() {
							continue // Skip if the module version doesn't match
						}

						collectedVersions[subdirVersion] = struct {
							manifestPath string
							manifest     *shared.ProjectManifest
						}{
							manifestPath: manifestPath,
							manifest:     &manifest,
						}
					}
				}
			}
		}
	}

	// Let's cache the collected versions for future use
	versionsCache[moduleName] = collectedVersions

	return collectedVersions
}

// DownloadModuleGit takes a module git url and attempts to fetch and download a version satisfying the constraint
func DownloadModuleGit(url string, constraint string) (string, *shared.ProjectManifest, error) {
	// First we create a temp dir to clone the repo
	tmpDir, err := os.MkdirTemp("", "xel-module-*")
	if err != nil {
		return "", nil, err
	}

	defer os.RemoveAll(tmpDir)

	// Clone the repo shallowly
	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:        url,
		Depth:      1,
		Tags:       git.AllTags,
		NoCheckout: true,
		Auth:       nil,
		Progress:   nil,
	})
	if err != nil {
		return "", nil, err
	}

	// list the all tags
	tags, err := repo.Tags()
	if err != nil {
		return "", nil, err
	}

	// collect valid tags
	// a tag is valid if it follows this format:
	// vX.Y.Z
	validTags := make([]*semver.Version, 0)
	tagMap := make(map[string]*plumbing.Reference)

	err = tags.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().String()

		name = strings.TrimPrefix(name, "refs/tags/")

		if strings.HasPrefix(name, "v") {
			version, err := semver.NewVersion(name[1:])
			if err != nil {
				return err
			}
			validTags = append(validTags, version)
			tagMap[version.String()] = ref
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}

	if len(validTags) == 0 {
		// Oops, no valid tags found
		return "", nil, fmt.Errorf("no valid semver tags found")
	}

	// sort tags descending (newer first)
	sort.Sort(sort.Reverse(semver.Collection(validTags)))

	// now we check each version and see if that thingy satisfy the given constraint.
	// actually, lets first parse the constraint
	versionConstraint, err := semver.NewConstraint(constraint)
	if err != nil {
		return "", nil, err
	}

	// finally lets check it...
	for _, version := range validTags {
		if versionConstraint.Check(version) {
			// we found a tag name that satisfies the constraint
			tagRef := tagMap[version.String()]

			// prepare destination for cloning the tag
			dest := filepath.Join(shared.XelConfig.ModulePaths[0],
				// the format of this is mod-[hash of url]/[version]
				// this kinda keeps uniqueness but still, its not the best option i could opt for...
				fmt.Sprintf("mod-%x", sha256.Sum256([]byte(url))), version.String())

			// create the dirs
			if err := os.MkdirAll(dest, 0755); err != nil {
				return "", nil, err
			}

			// clone ONLY the this specific tag
			_, err := git.PlainClone(dest, false, &git.CloneOptions{
				URL:           url,
				Depth:         1,
				SingleBranch:  true,
				ReferenceName: tagRef.Name(),
				Tags:          git.NoTags,
			})

			if err != nil {
				// we skip if we can't clone
				continue
			}

			// lets actually validate it has a real manifest
			manifestPath := filepath.Join(dest, "xel.json")
			manifestContent, err := os.ReadFile(manifestPath)
			if err != nil {
				// manifest might not exist, so we remove this
				os.RemoveAll(dest)
				continue
			}

			manifest := shared.ProjectManifest{}
			if err := json.Unmarshal(manifestContent, &manifest); err != nil {
				// manifest might be corrupted, so we remove this
				os.RemoveAll(dest)
				continue
			}

			// while we're here, lets also double verify the version
			// we dont waste memory by using the semver package here because we already verified the version
			// string earlier, so we can simply check if the manifest has same version as the tag specifies
			// this also means that the tag should have exactly the same version as listed in the package
			if manifest.Version != version.String() {
				// we skip if the version doesn't match
				shared.ColorPalette.Warning.Printf("Version spoof detected, tag name (`%s`) is not matching the version in the manifest (`%s`). Skipping.\n", version.String(), manifest.Version)
				os.RemoveAll(dest)
				continue
			}

			// we also need to verify that this package version is compatible with the current version of xel
			if shared.RuntimeVersion != "" { // Xel is not in development mode
				constraint, err := semver.NewConstraint(*manifest.Xel)
				if err != nil {
					shared.ColorPalette.Warning.Printf("Invalid Xel version constraint in manifest: %v\n", err)
					os.RemoveAll(dest)
					continue
				}

				runtimeVersion, err := semver.NewVersion(shared.RuntimeVersion)
				if err != nil {
					shared.ColorPalette.Warning.Printf("Invalid runtime version format: %v\n", err)
					os.RemoveAll(dest)
					continue
				}

				if !constraint.Check(runtimeVersion) {
					shared.ColorPalette.Warning.Printf("Xel version %s does not satisfy required version %s from xel.json, please upgrade your runtime\n", shared.RuntimeVersion, *manifest.Xel)
					os.RemoveAll(dest)
					continue
				}
			} else { // Xel is in development mode
				// we skip the version check
				shared.ColorPalette.Warning.Println("Xel is in development mode, skipping version check")
			}

			// we also need to verify that this package version is compatible with the current version of engine
			if shared.EngineVersion != "" { // Engine is not in development mode
				constraint, err := semver.NewConstraint(*manifest.Engine)
				if err != nil {
					shared.ColorPalette.Warning.Printf("Invalid engine version constraint in manifest: %v\n", err)
					os.RemoveAll(dest)
					continue
				}

				engineVersion, err := semver.NewVersion(shared.EngineVersion)
				if err != nil {
					shared.ColorPalette.Warning.Printf("Invalid engine version format: %v\n", err)
					os.RemoveAll(dest)
					continue
				}

				if !constraint.Check(engineVersion) {
					shared.ColorPalette.Warning.Printf("Engine version %s does not satisfy required version %s from xel.json, please upgrade your engine\n", shared.EngineVersion, *manifest.Engine)
					os.RemoveAll(dest)
					continue
				}
			} else { // Engine is in development mode
				// we skip the version check
				shared.ColorPalette.Warning.Println("Engine is in development mode, skipping version check")
			}

			// finally, we return the manifest
			return manifestPath, &manifest, nil
		}
	}

	return "", nil, fmt.Errorf("no valid version found that satisfies the constraint")
}

type RegistryPackageResp struct {
	ID int `json:"id"`
	GID int `json:"gid"`
	Name string `json:"name"`
	Latest int `json:"latest"`
	Description string `json:"description"`
	Author string `json:"author"`
	RepoName string `json:"repo_name"`
	RepoURL string `json:"url"`
	MirrorURL string `json:"mirror"`
	Tags []string `json:"tags"`
	IsDeprecated bool `json:"isDeprecated"`
	DeprecationReason string `json:"deprecatedReason"`
}

type RegistryPackageVersionRespMetadata struct {
	Semver struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Patch int `json:"patch"`
	} `json:"semver"`
	Version string `json:"version"`
	Package int `json:"package"`
	Downloads int `json:"downloads"`
	License string `json:"license"`
	DistMode string `json:"dist_mode"`
	Xel string `json:"xel"`
	Engine string `json:"engine"`
	ID int `json:"id"`
	GID int `json:"gid"`
}

type RegistryPackageVersionResp struct {
	Versions []RegistryPackageVersionRespMetadata
	Total int `json:"total"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

type RegistryTarballMetadataResp struct {
	Tarballs []struct {
		Integrity struct {
			Algorithm string `json:"algorithm"`
			Hash      string `json:"hash"`
		} `json:"integrity"`
		Package int `json:"package"`
		Version int `json:"version"`
		URL     string `json:"url"`
		SizeBytes int `json:"size_bytes"`
		ID        int `json:"id"`
		GID       int `json:"gid"`
		Downloads int `json:"downloads"`
	} `json:"tarballs"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}


func DownloadModuleOnline(moduleName string, constraint string) (string, *shared.ProjectManifest, string, string, string, error) {
	// Fetch the package details from the registry
	resp, err := http.Get(fmt.Sprintf("%s%s%s", shared.XelConfig.PackageRegistryURI, "packages/name/", moduleName))
	if err != nil {
		return "", nil, "", "", "", err
	}
	defer resp.Body.Close()
	
	// Parse the response body
	var packageDetails RegistryPackageResp
	if err := json.NewDecoder(resp.Body).Decode(&packageDetails); err != nil {
		return "", nil, "", "", "", err
	}

	// fetch the package versions from the registry
	resp, err = http.Get(fmt.Sprintf("%s%s%d?limit=100&offset=0", shared.XelConfig.PackageRegistryURI, "versions/pkg/", packageDetails.ID))
	if err != nil {
		return "", nil, "", "", "", err
	}
	defer resp.Body.Close()
	
	// Parse the response body
	var versions RegistryPackageVersionResp
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return "", nil, "", "", "", err
	}
	
	// Sort it in descending order (new -> old)
	sort.Slice(versions.Versions, func(i, j int) bool {
		a := versions.Versions[i].Semver
		b := versions.Versions[j].Semver

		if a.Major != b.Major {
			return a.Major > b.Major
		}
		if a.Minor != b.Minor {
			return a.Minor > b.Minor
		}
		return a.Patch > b.Patch
	})

	// find the version that satisfies the constraint
	// targetVersion := versions.Versions[0]
	targetVersion := RegistryPackageVersionRespMetadata{
		ID: -1,
	}
	for _, v := range versions.Versions {
		constraint, err := semver.NewConstraint(constraint)
		if err != nil {
			return "", nil, "", "", "", err
		}

		version, err := semver.NewVersion(v.Version)
		if err != nil {
			return "", nil, "", "", "", err
		}

		if constraint.Check(version) {
			targetVersion = v
			break
		}
	}

	if targetVersion.ID == -1 {
		return "", nil, "", "", "", fmt.Errorf("no version found that satisfies the constraint %s for package %s", constraint, moduleName)
	}

	// download the tarball metadata specific to that version constraint
	resp, err = http.Get(fmt.Sprintf("%s%s%d?limit=1&offset=0", shared.XelConfig.PackageRegistryURI, "tarballs/ver/", targetVersion.ID))
	if err != nil {
		return "", nil, "", "", "", err
	}
	defer resp.Body.Close()
	
	// Parse the response body
	var tarballMetadata RegistryTarballMetadataResp
	if err := json.NewDecoder(resp.Body).Decode(&tarballMetadata); err != nil {
		return "", nil, "", "", "", err
	}

	if len(tarballMetadata.Tarballs) == 0 {
		return "", nil, "", "", "", fmt.Errorf("no tarballs found for version %s", constraint)
	}

	manifestPath, manifest, err := DownloadFromTarball(tarballMetadata.Tarballs[0].URL, tarballMetadata.Tarballs[0].Integrity.Algorithm, tarballMetadata.Tarballs[0].Integrity.Hash, moduleName, targetVersion.Version)
	if err != nil {
		return "", nil, "", "", "", err
	}
	
	return manifestPath, manifest, tarballMetadata.Tarballs[0].Integrity.Algorithm, tarballMetadata.Tarballs[0].Integrity.Hash, tarballMetadata.Tarballs[0].URL, nil
}

func DownloadFromTarball(url, algorithm, hash, name, version string) (string, *shared.ProjectManifest, error) {
	tempDir, err := os.MkdirTemp("", "xel-module-*")
	if err != nil {
		return "", nil, err
	}
	defer os.RemoveAll(tempDir)

	// download the tarball
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("failed to download tarball: %s", resp.Status)
	}

	out, err := os.Create(filepath.Join(tempDir, "tarball.tar.gz"))
	if err != nil {
		return "", nil, err
	}
	defer out.Close()
	
	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", nil, err
	}

	// verify integrity
	err = VerifyTarballIntegrity(filepath.Join(tempDir, "tarball.tar.gz"), algorithm, hash)
	if err != nil {
		return "", nil, err
	}

	// prepare destination for package download
	dest := filepath.Join(shared.XelConfig.ModulePaths[0],
		// the format of this is mod-[hash of name]/[version]
		// this kinda keeps uniqueness but still, its not the best option i could opt for...
		fmt.Sprintf("mod-%x", sha256.Sum256([]byte(name))), version)
	
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return "", nil, err
	}

	// extract the tarball
	if err := ExtractTarGz(filepath.Join(tempDir, "tarball.tar.gz"), dest); err != nil {
		return "", nil, err
	}

	// read the manifest
	manifestPath := filepath.Join(dest, "xel.json")
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", nil, err
	}

	manifest := &shared.ProjectManifest{}
	if err := json.Unmarshal(manifestBytes, manifest); err != nil {
		return "", nil, err
	}

	return manifestPath, manifest, nil
}

func DownloadModule(moduleName string, version string, lockfile *shared.Lockfile) (string, *shared.ProjectManifest, error) {
	ver := version
	ver = strings.TrimSpace(ver)
	ver = strings.ToLower(ver)

	versionAliasTable := map[string]string{
		"latest": ">= 0.0.0",
		"stable": ">= 1.0.0",
		"any": "*",
	}

	if _, ok := versionAliasTable[ver]; ok {
		ver = versionAliasTable[ver]
	}

	var lockedModule struct{Algorithm string "json:\"algorithm\""; Hash string "json:\"hash\""; URL string "json:\"url\""; Version string "json:\"version\""}
	var ok bool
	if lockfile != nil {
		lockedModule, ok = (*lockfile)[moduleName]
	}
	if ok {
		url := lockedModule.URL
		return DownloadFromTarball(url, lockedModule.Algorithm, lockedModule.Hash, moduleName, lockedModule.Version)
	}

	var manifestPath string
	var manifest *shared.ProjectManifest
	var iAlgo string
	var iHash string
	var url string
	var err error
	if strings.HasPrefix(moduleName, "git+") {
		manifestPath, manifest, err = DownloadModuleGit(moduleName[4:], ver)
	} else {
		manifestPath, manifest, iAlgo, iHash, url, err = DownloadModuleOnline(moduleName, ver)
		if err == nil {
			if lockfile != nil {
				(*lockfile)[moduleName] = struct{Algorithm string "json:\"algorithm\""; Hash string "json:\"hash\""; URL string "json:\"url\""; Version string "json:\"version\""}{
					Algorithm: iAlgo,
					Hash: iHash,
					URL: url,
					Version: manifest.Version,
				}
			}
		}
	}
	
	return manifestPath, manifest, err
}
