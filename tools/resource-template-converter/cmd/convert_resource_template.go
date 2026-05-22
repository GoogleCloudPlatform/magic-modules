package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/GoogleCloudPlatform/magic-modules/tools/resource-template-converter/copy"
	"github.com/GoogleCloudPlatform/magic-modules/tools/resource-template-converter/github"
	"github.com/GoogleCloudPlatform/magic-modules/tools/resource-template-converter/migrate"
	"github.com/spf13/cobra"
)

var resourceFileRegex = regexp.MustCompile(`(?:mmv1/)?products/([^/]+)/([^/]+\.yaml)`)

var filePath string
var skipOpenPR bool

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
	var productsPath, examplesSourceDir, samplesDestDir string

	if _, err := os.Stat(filepath.Join(basePath, "mmv1")); err == nil {
		// Public magic-modules repository
		productsPath = filepath.Join(basePath, "mmv1", "products")
		templatesPath := filepath.Join(basePath, "mmv1", "templates", "terraform")
		examplesSourceDir = filepath.Join(templatesPath, "examples")
		samplesDestDir = filepath.Join(templatesPath, "samples", "services")
	} else if _, err := os.Stat(filepath.Join(basePath, "products")); err == nil {
		// EAP private overrides repository
		productsPath = filepath.Join(basePath, "products")
		examplesSourceDir = basePath
		samplesDestDir = filepath.Join(basePath, "templates", "terraform", "samples", "services")
	} else {
		log.Fatalf("Neither 'mmv1' nor 'products' directory structure found. Please ensure this tool is run from a magic-modules or magic-modules-private-overrides directory.")
	}

	var touchedFiles map[string][]int
	if skipOpenPR {
		fmt.Println("Fetching open PRs updated in the last 2 months from GitHub...")
		var err error
		touchedFiles, err = github.GetFilesTouchedByOpenPRs()
		if err != nil {
			return fmt.Errorf("failed to get touched files from open PRs: %w", err)
		}
		fmt.Printf("Successfully fetched open PRs. Found %d modified files in open PRs.\n", len(touchedFiles))
	}

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

		relPath, err := filepath.Rel(basePath, resolvedPath)
		if err == nil && skipOpenPR {
			normPath := github.NormalizePath(relPath)
			if prs, touched := touchedFiles[normPath]; touched {
				fmt.Printf("Skipping single target file %s: modified in active open PR(s) %v\n", targetFile, prs)
				return nil
			}
		}

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

			relPath, err := filepath.Rel(basePath, path)
			if err == nil && skipOpenPR {
				normPath := github.NormalizePath(relPath)
				if prs, touched := touchedFiles[normPath]; touched {
					log.Printf("Skipping file %s: modified in active open PR(s) %v\n", path, prs)
					return nil
				}
			}

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
	convertResourceTemplateCmd.Flags().BoolVar(&skipOpenPR, "skip-open-pr", false, "Skip files modified by active open PRs updated in the last 2 months")
	rootCmd.AddCommand(convertResourceTemplateCmd)
}
