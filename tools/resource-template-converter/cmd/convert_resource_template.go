package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/resource-template-converter/copy"
	"github.com/GoogleCloudPlatform/magic-modules/tools/resource-template-converter/github"
	"github.com/GoogleCloudPlatform/magic-modules/tools/resource-template-converter/migrate"
	"github.com/spf13/cobra"
)

var resourceFileRegex = regexp.MustCompile(`(?:mmv1/)?products/([^/]+)/([^/]+\.yaml)`)

var filePath string
var targetProduct string
var skipFilesFlag string
var skipProductsFlag string
var onlyMigration bool
var onlyFormat bool
var skipOpenPR bool
var skipOpenPRDays int

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
		if cmd.Flags().Changed("skip-open-pr-days") {
			skipOpenPR = true
		}
		return exeCconvertResourceTemplate(root, filePath)
	},
}

func exeCconvertResourceTemplate(basePath string, targetFile string) error {
	if targetFile != "" && targetProduct != "" {
		return fmt.Errorf("cannot specify both --file and --product")
	}
	if onlyMigration && onlyFormat {
		return fmt.Errorf("cannot specify both --only-migration and --only-format")
	}

	skipProductsMap := make(map[string]bool)
	if skipProductsFlag != "" {
		for _, p := range strings.Split(skipProductsFlag, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				skipProductsMap[p] = true
			}
		}
	}

	skipFilesMap := make(map[string]bool)
	if skipFilesFlag != "" {
		for _, f := range strings.Split(skipFilesFlag, ",") {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			resolved := f
			if !filepath.IsAbs(resolved) {
				resolved = filepath.Clean(filepath.Join(basePath, f))
			}
			skipFilesMap[resolved] = true
		}
	}

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

	var pathsToWalk []string
	if targetProduct != "" {
		products := strings.Split(targetProduct, ",")
		for _, p := range products {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			productDir := filepath.Join(productsPath, p)
			if _, err := os.Stat(productDir); os.IsNotExist(err) {
				return fmt.Errorf("product directory does not exist: %s", productDir)
			}
			pathsToWalk = append(pathsToWalk, productDir)
		}
	} else {
		pathsToWalk = []string{productsPath}
	}

	var touchedFiles map[string][]int
	if skipOpenPR {
		fmt.Printf("Fetching open PRs updated in the last %d days from GitHub...\n", skipOpenPRDays)
		var err error
		touchedFiles, err = github.GetFilesTouchedByOpenPRs(skipOpenPRDays)
		if err != nil {
			return fmt.Errorf("failed to get touched files from open PRs: %w", err)
		}
		fmt.Printf("Successfully fetched open PRs. Found %d modified files in open PRs.\n", len(touchedFiles))
	}

	if targetFile != "" {
		type targetYAML struct {
			resolvedPath string
			serviceName  string
			relPath      string
		}
		var targets []targetYAML

		files := strings.Split(targetFile, ",")
		for _, f := range files {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			resolvedPath := f
			if !filepath.IsAbs(resolvedPath) {
				resolvedPath = filepath.Clean(filepath.Join(basePath, f))
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
			if err != nil {
				return fmt.Errorf("failed to resolve relative path for %s: %w", resolvedPath, err)
			}

			targets = append(targets, targetYAML{
				resolvedPath: resolvedPath,
				serviceName:  serviceName,
				relPath:      relPath,
			})
		}

		for _, t := range targets {
			if skipOpenPR {
				normPath := github.NormalizePath(t.relPath)
				if prs, touched := touchedFiles[normPath]; touched {
					fmt.Printf("Skipping target file %s: modified in active open PR(s) %v\n", t.resolvedPath, prs)
					continue
				}
			}

			if _, skip := skipFilesMap[t.resolvedPath]; skip {
				fmt.Printf("Skipping target file %s: matched --skip-file filter\n", t.resolvedPath)
				continue
			}

			if _, skip := skipProductsMap[t.serviceName]; skip {
				fmt.Printf("Skipping target file %s: product %s matched --skip-product filter\n", t.resolvedPath, t.serviceName)
				continue
			}

			fmt.Printf("Processing product YAML file: %s (service: %s)\n", t.resolvedPath, t.serviceName)

			if !onlyFormat {
				if err := copy.ProcessResourceFile(t.resolvedPath, t.serviceName, examplesSourceDir, samplesDestDir); err != nil {
					return fmt.Errorf("error copying templates for %s: %w", t.resolvedPath, err)
				}
			}
			if err := migrate.MigrateFile(t.resolvedPath, t.serviceName, onlyMigration, onlyFormat); err != nil {
				return fmt.Errorf("failed to migrate file %s: %w", t.resolvedPath, err)
			}
		}

		fmt.Println("Processing complete.")
		return nil
	}

	for _, pPath := range pathsToWalk {
		fmt.Printf("Starting processing of product YAML files in: %s\n", pPath)

		err := filepath.Walk(pPath, func(path string, info os.FileInfo, err error) error {
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

				if _, skip := skipProductsMap[serviceName]; skip {
					return nil
				}

				resolvedPath := path
				if !filepath.IsAbs(resolvedPath) {
					resolvedPath = filepath.Clean(filepath.Join(basePath, path))
				}
				if _, skip := skipFilesMap[resolvedPath]; skip {
					log.Printf("Skipping file %s: matched --skip-file filter\n", path)
					return nil
				}

				relPath, err := filepath.Rel(basePath, path)
				if err == nil && skipOpenPR {
					normPath := github.NormalizePath(relPath)
					if prs, touched := touchedFiles[normPath]; touched {
						log.Printf("Skipping file %s: modified in active open PR(s) %v\n", path, prs)
						return nil
					}
				}

				if !onlyFormat {
					if err := copy.ProcessResourceFile(path, serviceName, examplesSourceDir, samplesDestDir); err != nil {
						log.Printf("Error copying templates registered in file %s: %v\n", path, err)
						// Continue processing other files even if one fails.
					}
				}
				if err := migrate.MigrateFile(path, serviceName, onlyMigration, onlyFormat); err != nil {
					log.Printf("Failed to migrate file %s: %v\n", path, err)
					// Continue migrating other files even if one fails.
				}
			}
			return nil
		})

		if err != nil {
			log.Fatalf("Error walking the products path %q: %v\n", pPath, err)
		}
	}

	fmt.Println("Processing complete.")
	return nil
}

func init() {
	convertResourceTemplateCmd.Flags().StringVarP(&filePath, "file", "f", "", "Comma-separated list of resource yaml files to convert (e.g. mmv1/products/vertexai/Dataset.yaml)")
	convertResourceTemplateCmd.Flags().StringVarP(&targetProduct, "product", "p", "", "Comma-separated list of product directories to convert (e.g. vertexai,pubsublite)")
	convertResourceTemplateCmd.Flags().StringVarP(&skipFilesFlag, "skip-file", "F", "", "Comma-separated list of resource yaml files to skip from migration")
	convertResourceTemplateCmd.Flags().StringVarP(&skipProductsFlag, "skip-product", "P", "", "Comma-separated list of product directories to skip from migration")
	convertResourceTemplateCmd.Flags().BoolVar(&onlyMigration, "only-migration", false, "Only run migration steps (examples -> samples, copy templates), skip formatting")
	convertResourceTemplateCmd.Flags().BoolVar(&onlyFormat, "only-format", false, "Only run formatting steps (sort keys, strip quotes), skip migration")
	convertResourceTemplateCmd.Flags().BoolVar(&skipOpenPR, "skip-open-pr", false, "Skip files modified by active open PRs updated in the last N days (configured by --skip-open-pr-days)")
	convertResourceTemplateCmd.Flags().IntVar(&skipOpenPRDays, "skip-open-pr-days", 60, "Number of days of open PR history to verify when checking open PRs")
	rootCmd.AddCommand(convertResourceTemplateCmd)
}
