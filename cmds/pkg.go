package cmds

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"xel/helpers"
	"xel/shared"

	"github.com/urfave/cli/v2"
)

// PackageCommands returns the cli.Command for package management
func PackageCommands() *cli.Command {
	return &cli.Command{
		Name:  "pkg",
		Usage: "Manage packages",
		Subcommands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"install"},
				Usage:   "Add a package",
				Action: func(c *cli.Context) error {
					// Get current working directory
					cwd, err := os.Getwd()
					if err != nil {
						return err
					}

					// Get closest manifest
					mainManifest, mainManifestPath, err := helpers.FetchManifest(cwd, cwd)
					if err != nil {
						return err
					}

					// Lockfile integration
					lockfilePath := filepath.Join(filepath.Dir(mainManifestPath), "xel.lock")
					lockfileData, err := os.ReadFile(lockfilePath)
					if err != nil {
						if !os.IsNotExist(err) {
							return err
						}
						lockfileData = []byte("{}")
					}
					var lockfile shared.Lockfile
					if err := json.Unmarshal(lockfileData, &lockfile); err != nil {
						return err
					}

					err = func() error {
						for i := 0; i < c.Args().Len(); i++ {
							// Validate package name
							packageName := c.Args().Get(i)
							if packageName == "" {
								return fmt.Errorf("package name is required")
							}

							// Split package name and version if provided
							split := strings.Split(packageName, "@")
							if len(split) > 2 {
								return fmt.Errorf("invalid package format: %s", packageName)
							}

							var name, version string
							switch len(split) {
							case 0:
								return fmt.Errorf("invalid package format: %s", packageName)
							case 1:
								name = strings.Join(split, "@")
								version = "*"
								if mainManifest.Deps != nil {
									if ver, ok := (*mainManifest.Deps)[name]; ok {
										version = ver
									}
								}

							default:
								name = strings.Join(split[:len(split)-1], "@")
								version = split[len(split)-1]
							}

							// Validate name
							if len(strings.TrimSpace(name)) == 0 {
								return fmt.Errorf("package name is required")
							}

							// Check if an installed package satisfies the requirements
							versions := helpers.GetVersions(name)

							if len(versions) != 0 && len(version) != 0 && version != "latest" {
								_, manifest, err := helpers.ResolveModuleLocal(name, version)
								if err != nil {
									return err
								}

								shared.ColorPalette.Warning.Printf("Package `%s@%s` already installed\n", name, manifest.Version)

								if mainManifest.Deps == nil {
									mainManifest.Deps = &map[string]string{}
								}
								(*mainManifest.Deps)[manifest.Name] = manifest.Version
								manifestData, err := json.MarshalIndent(mainManifest, "", "  ")
								if err != nil {
									return err
								}
								if err := os.WriteFile(mainManifestPath, manifestData, 0644); err != nil {
									return err
								}
								return nil
							}

							// download the package
							manifest, err := helpers.DownloadModule(name, version, &lockfile)
							if err != nil {
								return err
							}
							
							if mainManifest.Deps == nil {
								mainManifest.Deps = &map[string]string{}
							}
							(*mainManifest.Deps)[manifest.Name] = manifest.Version
							manifestData, err := json.MarshalIndent(mainManifest, "", "  ")
							if err != nil {
								return err
							}
							if err := os.WriteFile(mainManifestPath, manifestData, 0644); err != nil {
								return err
							}
							
							shared.ColorPalette.Info.Printf("Package `%s@%s` installed\n", name, manifest.Version)
						}
						return nil
					}()
					if err != nil {
						return err
					}

					// update lockfile
					lockfileData, err = json.MarshalIndent(lockfile, "", "  ")
					if err != nil {
						return err
					}
					if err := os.WriteFile(lockfilePath, lockfileData, 0644); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"uninstall"},
				Usage:   "Remove a package",
				Action: func(c *cli.Context) error {
					// TODO: Implement package removal logic
					// - Validate package name
					// - Check if package is installed
					// - Remove package files
					// - Update package manifest
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List installed packages",
				Action: func(c *cli.Context) error {
					// TODO: Implement package listing logic
					// - Read package manifest
					// - Format and display installed packages
					return nil
				},
			},
			{
				Name:  "update",
				Usage: "Update installed packages",
				Action: func(c *cli.Context) error {
					// TODO: Implement package update logic
					// - Check for package updates
					// - Download and install updates
					// - Update package manifest
					return nil
				},
			},
		},
	}
}
