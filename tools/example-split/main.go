package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Resource struct {
	Examples []Example `yaml:"examples"`
}

type Example struct {
	Name string `yaml:"name"`
}

func main() {
	basePath := "../.."

	productsPath := filepath.Join(basePath, "mmv1", "products")
	templatesPath := filepath.Join(basePath, "mmv1", "templates", "terraform")
	examplesSourceDir := filepath.Join(templatesPath, "examples")
	samplesDestDir := filepath.Join(templatesPath, "samples", "services")

	fmt.Printf("Starting processing of product YAML files in: %s\n", productsPath)

	// Walk through the products directory to find all resource.yaml files.
	err := filepath.Walk(productsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// We are only interested in files, not directories.
		if info.IsDir() {
			return nil
		}

		// Check if the file is a YAML file.
		if filepath.Ext(path) == ".yaml" {
			err := processResourceFile(path, examplesSourceDir, samplesDestDir)
			if err != nil {
				log.Printf("Error processing file %s: %v\n", path, err)
				// Continue processing other files even if one fails.
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the products path %q: %v\n", productsPath, err)
	}

	fmt.Println("Processing complete.")
}

// processResourceFile reads a resource YAML file, extracts example names,
// and copies the corresponding template files.
func processResourceFile(filePath, examplesSourceDir, samplesDestDir string) error {
	// Extract the service name from the file path.
	// The path is like .../mmv1/products/<service_name>/<resource_name>.yaml
	parts := strings.Split(filepath.Dir(filePath), string(os.PathSeparator))
	if len(parts) < 2 {
		return fmt.Errorf("could not determine service name from path: %s", filePath)
	}
	serviceName := parts[len(parts)-1]

	// Read the YAML file content.
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	var resource Resource
	// Unmarshal the YAML content into our Resource struct.
	err = yaml.Unmarshal(yamlFile, &resource)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// If there are no examples, there's nothing to do.
	if len(resource.Examples) == 0 {
		return nil
	}

	fmt.Printf("Found %d examples in %s for service '%s'\n", len(resource.Examples), filepath.Base(filePath), serviceName)

	// Process each example found in the file.
	for _, example := range resource.Examples {
		if example.Name == "" {
			log.Printf("Skipping example with empty name in file: %s\n", filePath)
			continue
		}

		// Construct the source and destination paths for the template file.
		sourceFileName := fmt.Sprintf("%s.tf.tmpl", example.Name)
		sourcePath := filepath.Join(examplesSourceDir, sourceFileName)

		// Before trying to copy, check if the source template file actually exists.
		// It's possible for an example to be defined without a corresponding template.
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			// Silently skip if the template doesn't exist.
			continue
		}

		destDir := filepath.Join(samplesDestDir, serviceName)
		destPath := filepath.Join(destDir, sourceFileName)

		// Copy the file from source to destination.
		err := copyFile(sourcePath, destPath)
		if err != nil {
			// Log the error but don't stop processing other examples.
			log.Printf("Failed to copy '%s' to '%s': %v\n", sourcePath, destPath, err)
		} else {
			fmt.Printf("  - Copied '%s'\n", sourceFileName)
		}
	}

	return nil
}

// copyFile handles the copying of a single file from a source to a destination,
// replacing "$.Vars" with "$.PrefixedVars" in the process.
// It creates the destination directory if it doesn't exist.
func copyFile(src, dst string) error {
	// Read the content of the source file.
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("could not read source file %s: %w", src, err)
	}

	// Perform the string replacement.
	output := strings.ReplaceAll(string(input), "$.Vars", "$.PrefixedVars")

	// Create the destination directory if it doesn't already exist.
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("could not create destination directory for %s: %w", dst, err)
	}

	// Get the source file's permissions to apply to the destination file.
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("could not stat source file %s: %w", src, err)
	}

	// Write the modified content to the destination file.
	err = ioutil.WriteFile(dst, []byte(output), info.Mode())
	if err != nil {
		return fmt.Errorf("could not write to destination file %s: %w", dst, err)
	}

	return nil
}