package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/GoogleCloudPlatform/magic-modules/tools/resource-template-converter/copy"
	"github.com/GoogleCloudPlatform/magic-modules/tools/resource-template-converter/migrate"
	"github.com/spf13/cobra"
)

var resourceFileRegex = regexp.MustCompile(`mmv1/products/([^/]+)/([^/]+\.yaml)`)

var filePath string

var convertResourceTemplateCmd = &cobra.Command{
	Use:   "convert-resource-template",
	Short: "convert resource template from using examples to samples",
	Long: `This command convert resource yaml template to use new version samples within existing legacy examples.


	The command expects the following argument(s):
	1. Root directory path

	It then performs the following operations:
	1. Updates existing example config to new vars and then copy them to the new samples/services/<service> dir
	2. Updates resource yaml teplate to use new samples blocks from existing legacy examples block`,

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := args[0]
		return exeCconvertResourceTemplate(root, filePath)
	},
}

func exeCconvertResourceTemplate(basePath string, targetFile string) error {
	if _, err := os.Stat(filepath.Join(basePath, "mmv1")); os.IsNotExist(err) {
		log.Fatalf("magic-modules directory structure not found. Please ensure this tool is run from 'magic-modules/tools/example-split'.")
	}

	productsPath := filepath.Join(basePath, "mmv1", "products")
	templatesPath := filepath.Join(basePath, "mmv1", "templates", "terraform")
	examplesSourceDir := filepath.Join(templatesPath, "examples")
	samplesDestDir := filepath.Join(templatesPath, "samples", "services")

	if targetFile != "" {
		resolvedPath := targetFile
		if !filepath.IsAbs(resolvedPath) {
			resolvedPath = filepath.Clean(filepath.Join(basePath, targetFile))
		}

		if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
			return fmt.Errorf("target file does not exist: %s", resolvedPath)
		}

		matches := resourceFileRegex.FindStringSubmatch(resolvedPath)
		if matches == nil {
			return fmt.Errorf("file path %s does not match expected pattern mmv1/products/<service>/<resource>.yaml", resolvedPath)
		}
		serviceName := matches[1]

		fmt.Printf("Processing single product YAML file: %s (service: %s)\n", resolvedPath, serviceName)

		if err := copy.ProcessResourceFile(resolvedPath, serviceName, examplesSourceDir, samplesDestDir); err != nil {
			return fmt.Errorf("error copying templates: %w", err)
		}
		if err := migrate.MigrateFile(resolvedPath, serviceName); err != nil {
			return fmt.Errorf("failed to migrate file: %w", err)
		}

		fmt.Println("Processing complete.")
		return nil
	}

	fmt.Printf("Starting processing of product YAML files in: %s\n", productsPath)

	err := filepath.Walk(productsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			matches := resourceFileRegex.FindStringSubmatch(path)
			if matches == nil {
				log.Printf("Skipping non-resource file: %s\n", path)
				return nil
			}
			serviceName := matches[1]

			if err := copy.ProcessResourceFile(path, serviceName, examplesSourceDir, samplesDestDir); err != nil {
				log.Printf("Error copying templates registered in file %s: %v\n", path, err)
				// Continue processing other files even if one fails.
			}
			if err := migrate.MigrateFile(path, serviceName); err != nil {
				log.Printf("Failed to migrate file %s: %v\n", path, err)
				// Continue migrating other files even if one fails.
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the products path %q: %v\n", productsPath, err)
	}

	fmt.Println("Processing complete.")
	return nil
}

func init() {
	convertResourceTemplateCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to a single resource yaml file to convert")
	rootCmd.AddCommand(convertResourceTemplateCmd)
}
