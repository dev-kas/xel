package cmds

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
						if err.Error() != "cannot find manifest" {
							return err
						}
					}

					// Lockfile integration
					lockfilePath := ""
					if mainManifest != nil {
						lockfilePath = filepath.Join(filepath.Dir(mainManifestPath), "xel.lock")
					}
					lockfileData := []byte("{}")
					if lockfilePath != "" {
						data, err := os.ReadFile(lockfilePath)
						if err != nil {
							if !os.IsNotExist(err) {
								return err
							} else {
								lockfileData = []byte("{}")
							}
						} else {
							lockfileData = data
						}
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
								if mainManifest != nil && mainManifest.Deps != nil {
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

								if mainManifest != nil {
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
								}
								return nil
							}

							// download the package
							manifestPath, manifest, err := helpers.DownloadModule(name, version, &lockfile)
							if err != nil {
								return err
							}

							if mainManifest != nil {
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
							}
							for depName, depVersion := range *manifest.Deps {
								exePath, err := os.Executable()
								if err != nil {
									return err
								}
								xelPath, err := filepath.EvalSymlinks(exePath)
								if err != nil {
									return err
								}

								cmd := exec.Command(xelPath, "pkg", "add", fmt.Sprintf("%s@%s", depName, depVersion))
								cmd.Stderr = os.Stderr
								cmd.Stdout = os.Stdout
								cmd.Dir = filepath.Dir(manifestPath)
								if err := cmd.Run(); err != nil {
									return err
								}
							}

							shared.ColorPalette.Info.Printf("Package `%s@%s` installed\n", name, manifest.Version)
						}
						return nil
					}()
					if err != nil {
						return err
					}

					// update lockfile
					if lockfilePath != "" {
						lockfileData, err = json.MarshalIndent(lockfile, "", "  ")
						if err != nil {
							return err
						}
						if err := os.WriteFile(lockfilePath, lockfileData, 0644); err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"uninstall"},
				Usage:   "Remove a package",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "global",
						Aliases: []string{"g"},
						Usage:   "Removes the package from global modules",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "local",
						Aliases: []string{"l"},
						Usage:   "Removes the package from local modules",
						Value:   false,
					},
				},
				Action: func(c *cli.Context) error {
					removalMode := "local"

					// Get current working dir
					cwd, err := os.Getwd()
					if err != nil {
						return err
					}

					// Get closest manifest
					mainManifest, mainManifestPath, err := helpers.FetchManifest(cwd, cwd)
					if err != nil {
						return err
					}
					if mainManifest == nil {
						removalMode = "global"
					}

					for i := 0; i < c.Args().Len(); i++ {
						// Validate package name
						if len(strings.TrimSpace(c.Args().Get(i))) == 0 {
							return fmt.Errorf("package name is required")
						}

						packageName := strings.TrimSpace(c.Args().Get(i))
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
						default:
							name = strings.Join(split[:len(split)-1], "@")
							version = split[len(split)-1]
						}

						var ok bool
						if mainManifest != nil {
							_, ok = (*mainManifest.Deps)[name]
						}
						if !ok {
							removalMode = "global"
						}

						manifestPath, _, err := helpers.ResolveModuleLocal(name, version)
						if err != nil && removalMode != "local" {
							removalMode = "local"
						}

						if mainManifest.Deps == nil {
							removalMode = "global"
						} else {
							if _, ok := (*mainManifest.Deps)[name]; !ok {
								removalMode = "global"
							} else {
								removalMode = "local"
							}
						}

						globalFlag := c.Bool("global")
						localFlag := c.Bool("local")
						if globalFlag && localFlag {
							return fmt.Errorf("cannot specify both global and local")
						}
						if globalFlag {
							removalMode = "global"
						} else if localFlag {
							removalMode = "local"
						}

						if removalMode == "local" {
							delete(*mainManifest.Deps, name)
							manifestData, err := json.MarshalIndent(mainManifest, "", "  ")
							if err != nil {
								return err
							}
							if err := os.WriteFile(mainManifestPath, manifestData, 0644); err != nil {
								return err
							}

							lfPath := filepath.Join(filepath.Dir(mainManifestPath), "xel.lock")
							lfData, err := os.ReadFile(lfPath)
							if err != nil {
								if !os.IsNotExist(err) {
									return err
								}
							}
							var lockfile shared.Lockfile
							if err := json.Unmarshal(lfData, &lockfile); err != nil {
								return err
							}
							delete(lockfile, name)
							lfData, err = json.MarshalIndent(lockfile, "", "  ")
							if err != nil {
								return err
							}
							if err := os.WriteFile(lfPath, lfData, 0644); err != nil {
								return err
							}
						} else {
							if manifestPath == "" {
								return fmt.Errorf("package `%s` not found", name)
							}
							if err := os.RemoveAll(filepath.Dir(manifestPath)); err != nil {
								return err
							}
						}
					}

					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List installed packages",
				Action: func(c *cli.Context) error {
					cwd, err := os.Getwd()
					if err != nil {
						return err
					}
					mainManifest, _, err := helpers.FetchManifest(cwd, cwd)
					if err != nil {
						return err
					}
					if mainManifest == nil {
						return fmt.Errorf("no manifest found")
					}
					shared.ColorPalette.Info.Printf("Installed Packages in %s@%s:\n", mainManifest.Name, mainManifest.Version)
					shared.ColorPalette.Info.Println(strings.Repeat("-", 20))
					for name, version := range *mainManifest.Deps {
						shared.ColorPalette.Info.Printf("- %s@%s\n", name, version)
					}
					if len(*mainManifest.Deps) == 0 {
						shared.ColorPalette.Info.Println("No packages installed")
					}
					return nil
				},
			},
		},
	}
}
