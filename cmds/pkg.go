package cmds

import (
	"fmt"
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
					// Validate package name
					packageName := c.Args().First()
					if packageName == "" {
						// TODO: Install packages listed in the cwd's manifest
						return fmt.Errorf("package name is required")
					}

					// Split package name and version if provided
					split := strings.Split(packageName, "@")
					if len(split) > 2 {
						return fmt.Errorf("invalid package format: %s", packageName)
					}

					var name, version string
					switch len(split) {
					case 1:
						name = split[0]
					case 2:
						name = split[0]
						version = split[1]
					}

					// Validate name
					if len(strings.TrimSpace(name)) == 0 {
						return fmt.Errorf("package name is required")
					}

					// Check if an installed package satisfies the requirements
					versions := helpers.GetVersions(name)

					if len(versions) != 0 && len(version) != 0 {
						// TODO: Link the package
						version, _, err := helpers.ResolveModule(name, version)
						if err != nil {
							return err
						}

						// TODO: Link the package
						shared.ColorPalette.Warning.Printf("Package `%s@%s` already installed", name, version)
						return nil
					} else if len(versions) != 0 && len(version) == 0 {
						// TODO: Link the package
						shared.ColorPalette.Warning.Printf("Package `%s` already installed", name)
						return nil
					}

					// download the package
					manifest, err := helpers.ResolveOnline(name, version)
					if err != nil {
						return err
					}

					// TODO: Link the package
					shared.ColorPalette.Info.Printf("Package `%s@%s` installed", name, manifest.Version)
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
