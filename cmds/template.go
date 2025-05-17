package cmds

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	"xel/helpers"
	"xel/shared"
)

// TemplateCommand returns the cli.Command for the template command
func TemplateCommand() *cli.Command {
	return &cli.Command{
		Name:  "template",
		Usage: "Manage project templates",
		Subcommands: []*cli.Command{
			{
				Name:  "pull",
				Usage: "Download a template from a URL",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("please provide a template URL")
					}

					templateURL := c.Args().First()
					templateName := filepath.Base(templateURL)

					if err := helpers.DownloadTemplate(templateURL, templateName); err != nil {
						return fmt.Errorf("failed to download template: %v", err)
					}

					shared.ColorPalette.Info.Printf("Template '%s' downloaded successfully.\n", templateName)
					return nil
				},
			},
			{
				Name:  "use",
				Usage: "Use an existing template",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("please provide a template name")
					}

					templateName := c.Args().First()
					templatesPath := filepath.Join(shared.XelDir(), "templates")
					templatePath := filepath.Join(templatesPath, templateName)

					// Check if template exists
					if _, err := os.Stat(templatePath); os.IsNotExist(err) {
						return fmt.Errorf("template '%s' not found", templateName)
					}

					// Check if current directory is empty
					cwd, err := os.Getwd()
					if err != nil {
						return fmt.Errorf("could not get current working directory: %v", err)
					}

					files, err := os.ReadDir(cwd)
					if err != nil {
						return fmt.Errorf("could not read current directory: %v", err)
					}

					if len(files) > 0 {
						shared.ColorPalette.Warning.Println("Warning: The current directory is not empty.")
						shared.ColorPalette.Warning.Println("Using this template here may overwrite existing files.")

						reader := bufio.NewReader(os.Stdin)
						shared.ColorPalette.Prompt.Print("Do you want to continue? (y/N): ")
						response, _ := reader.ReadString('\n')
						response = strings.TrimSpace(strings.ToLower(response))

						if response != "y" && response != "yes" {
							return fmt.Errorf("template application cancelled")
						}
					}

					// Copy template files to current directory
					if err := helpers.CopyDir(templatePath, cwd); err != nil {
						return fmt.Errorf("failed to copy template files: %v", err)
					}

					shared.ColorPalette.Info.Printf("Template '%s' applied successfully.\n", templateName)
					return nil
				},
			},
			{
				Name:  "remove",
				Usage: "Remove an existing template",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("please provide a template name")
					}

					templateName := c.Args().First()
					templatesPath := filepath.Join(shared.XelDir(), "templates")
					templatePath := filepath.Join(templatesPath, templateName)

					// Check if template exists
					if _, err := os.Stat(templatePath); os.IsNotExist(err) {
						return fmt.Errorf("template '%s' not found", templateName)
					}

					// Ask for confirmation
					shared.ColorPalette.Warning.Printf("Are you sure you want to delete template '%s'? This action cannot be undone.\n", templateName)
					reader := bufio.NewReader(os.Stdin)
					shared.ColorPalette.Prompt.Print("Type 'DELETE' to confirm: ")
					response, _ := reader.ReadString('\n')
					response = strings.TrimSpace(strings.ToUpper(response))

					if response != "DELETE" {
						return fmt.Errorf("template deletion cancelled")
					}

					// Remove the template directory
					if err := os.RemoveAll(templatePath); err != nil {
						return fmt.Errorf("failed to remove template: %v", err)
					}

					shared.ColorPalette.Info.Printf("Template '%s' removed successfully.\n", templateName)
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List all available templates",
				Action: func(c *cli.Context) error {
					templatesPath := filepath.Join(shared.XelDir(), "templates")

					// Check if templates directory exists
					if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
						shared.ColorPalette.Info.Println("No templates installed yet.")
						return nil
					}

					// Read all directories in templates folder
					entries, err := os.ReadDir(templatesPath)
					if err != nil {
						return fmt.Errorf("failed to read templates directory: %v", err)
					}

					// Print header
					shared.ColorPalette.Info.Println("Installed Templates:")
					shared.ColorPalette.Info.Println(strings.Repeat("-", 20))

					// Print each template
					for _, entry := range entries {
						if !entry.IsDir() {
							continue
						}
						shared.ColorPalette.Info.Printf("- %s\n", entry.Name())
					}

					if len(entries) == 0 {
						shared.ColorPalette.Info.Println("No templates installed yet.")
					}

					return nil
				},
			},
		},
	}
}
