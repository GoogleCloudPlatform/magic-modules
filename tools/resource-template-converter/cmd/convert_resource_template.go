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

var resourceFileRegex = regexp.MustCompile(`/mmv1/products/([^/]+)/([^/]+\.yaml)`)

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
		return exeCconvertResourceTemplate(root)
	},
}

func exeCconvertResourceTemplate(basePath string) error {
	if _, err := os.Stat(filepath.Join(basePath, "mmv1")); os.IsNotExist(err) {
		log.Fatalf("magic-modules directory structure not found. Please ensure this tool is run from 'magic-modules/tools/example-split'.")
	}

	productsPath := filepath.Join(basePath, "mmv1", "products")
	templatesPath := filepath.Join(basePath, "mmv1", "templates", "terraform")
	examplesSourceDir := filepath.Join(templatesPath, "examples")
	samplesDestDir := filepath.Join(templatesPath, "samples", "services")
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
	rootCmd.AddCommand(convertResourceTemplateCmd)
}
