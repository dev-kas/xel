package cmds

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"

	"xel/helpers"
	"xel/shared"
)

// InitCommand returns the cli.Command for the init command
func InitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize a new project.",
		Action: func(c *cli.Context) error {
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
				shared.ColorPalette.Warning.Println("Initializing a project here may overwrite existing files.")

				reader := bufio.NewReader(os.Stdin)
				shared.ColorPalette.Prompt.Print("Do you want to continue? (y/N): ")
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))

				if response != "y" && response != "yes" {
					return fmt.Errorf("project initialization cancelled")
				}
			}

			// Check if a template name is provided
			var templateName string
			if c.NArg() == 0 {
				// Use default template from config if no template is specified
				templateName = shared.XelConfig.DefaultTemplate
				shared.ColorPalette.Info.Printf("No template specified. Using default template: %s\n", templateName)
			} else {
				// Get the template name from the first argument
				templateName = c.Args().First()
			}

			// Check if the template exists in the templates directory
			templatesPath := filepath.Join(shared.XelDir(), "templates")
			templatePath := filepath.Join(templatesPath, templateName)

			// Check if the template directory exists
			if _, err := os.Stat(templatePath); os.IsNotExist(err) {
				// Check if template exists in online repository
				shared.ColorPalette.Warning.Printf("Template '%s' not found locally.\n", templateName)

				reader := bufio.NewReader(os.Stdin)
				shared.ColorPalette.Prompt.Print("Would you like to download this template? (y/N): ")
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))

				if response == "y" || response == "yes" {
					// Construct the full GitHub repository URL for cloning
					onlineTemplateCloneURL := fmt.Sprintf("https://github.com/dev-kas/XelTemplates/archive/refs/heads/master.zip?path=%s", templateName)

					// Print the URL for debugging
					fmt.Printf("Attempting to download from URL: %s\n", onlineTemplateCloneURL)

					// Use the DownloadTemplate helper to fetch the template
					err := helpers.DownloadTemplate(onlineTemplateCloneURL, templateName)
					if err != nil {
						shared.ColorPalette.Error.Printf("Failed to download template: %v\n", err)
						return fmt.Errorf("template download failed: %v", err)
					}

					shared.ColorPalette.Info.Printf("Template '%s' downloaded successfully.\n", templateName)
				} else {
					return fmt.Errorf("template '%s' not found", templateName)
				}
			}

			// Define questionnaire for project manifest
			manifest := shared.ProjectManifest{
				Name: templateName, // Default to template name
				Deps: &map[string]string{},
			}

			type extras struct {
				git string
			}

			extra := extras{}

			// Questionnaire questions with their corresponding manifest fields
			questions := []struct {
				prompt       string
				field        *string
				defaultValue string
			}{
				{
					prompt:       fmt.Sprintf("What would you like to name this project? (%s): ", templateName),
					field:        &manifest.Name,
					defaultValue: templateName,
				},
				{
					prompt: "Provide a short description for your project: ",
					field:  &manifest.Description,
				},
				{
					prompt:       "What version would you like to start with? (0.1.0): ",
					field:        &manifest.Version,
					defaultValue: "0.1.0",
				},
				{
					prompt: "Enter project author name: ",
					field:  &manifest.Author,
				},
				{
					prompt:       "Choose a license (MIT): ",
					field:        &manifest.License,
					defaultValue: "MIT",
				},
				{
					prompt:       "Would you like to initialize a git repository? (Y/n): ",
					field:        &extra.git,
					defaultValue: "Y",
				},
			}

			// Interactive questionnaire
			reader := bufio.NewReader(os.Stdin)
			for _, q := range questions {
				shared.ColorPalette.Prompt.Print(q.prompt)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(response)

				// Use default if no response provided
				if response == "" {
					response = q.defaultValue
				}

				*q.field = response
			}

			// Prepare and copy template
			for _, file := range files {
				if !strings.HasPrefix(file.Name(), ".") {
					if err := os.RemoveAll(filepath.Join(cwd, file.Name())); err != nil {
						return fmt.Errorf("could not remove file/directory %s: %v", file.Name(), err)
					}
				}
			}

			// Copy template recursively
			if err := helpers.CopyDir(templatePath, cwd); err != nil {
				return fmt.Errorf("failed to copy template: %v", err)
			}

			// Try to read existing xel.json if present
			manifestPath := filepath.Join(cwd, "xel.json")
			existingManifest := shared.ProjectManifest{}

			if manifestBytes, err := os.ReadFile(manifestPath); err == nil {
				err = json.Unmarshal(manifestBytes, &existingManifest)
				if err != nil {
					shared.ColorPalette.Warning.Printf("Could not parse existing xel.json: %v\n", err)
				}
			}

			// Merge questionnaire data with existing manifest
			// Prioritize questionnaire data over existing manifest
			if manifest.Name != "" {
				existingManifest.Name = manifest.Name
			}
			if manifest.Description != "" {
				existingManifest.Description = manifest.Description
			}
			if manifest.Version != "" {
				existingManifest.Version = manifest.Version
			}
			existingManifest.Deps = manifest.Deps
			if manifest.Author != "" {
				existingManifest.Author = manifest.Author
			}
			if manifest.License != "" {
				existingManifest.License = manifest.License
			}

			(*existingManifest.Xel) = fmt.Sprintf("^%s", shared.RuntimeVersion)
			(*existingManifest.Engine) = fmt.Sprintf("^%s", shared.EngineVersion)

			// Save merged manifest
			manifestJSON, err := json.MarshalIndent(existingManifest, "", "  ")
			if err != nil {
				shared.ColorPalette.Error.Printf("Failed to create manifest: %v\n", err)
				return err
			}

			err = os.WriteFile(manifestPath, manifestJSON, 0644)
			if err != nil {
				shared.ColorPalette.Error.Printf("Failed to write manifest file: %v\n", err)
				return err
			}

			// Initialize git repository if requested
			if extra.git == "Y" || extra.git == "y" {
				// Initialize git repository
				_, err = git.PlainInit(cwd, false)
				if err != nil {
					shared.ColorPalette.Error.Printf("Failed to initialize git repository: %v\n", err)
					return err
				}
			}

			shared.ColorPalette.Info.Printf("Project '%s' initialized successfully!\n", existingManifest.Name)
			return nil
		},
	}
}
