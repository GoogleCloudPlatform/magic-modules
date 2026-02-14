package copy

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
	Name       string `yaml:"name"`
	ConfigPath string `yaml:"config_path,omitempty"`
}

func ProcessResourceFile(filePath, serviceName, examplesSourceDir, samplesDestDir string) error {
	originalBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	var resource Resource
	err = yaml.Unmarshal(originalBytes, &resource)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if len(resource.Examples) == 0 {
		return nil
	}

	// fmt.Printf("Found %d examples in %s for service '%s'\n", len(resource.Examples), filepath.Base(filePath), serviceName)

	for _, example := range resource.Examples {
		if example.Name == "" {
			log.Printf("Skipping example with empty name in file: %s\n", filePath)
			continue
		}

		// Determine the source file name. Use config_path if it exists,
		// otherwise fall back to the example's name.
		var sourceFileName string
		if example.ConfigPath != "" {
			sourceFileName = filepath.Base(example.ConfigPath)
		} else {
			sourceFileName = fmt.Sprintf("%s.tf.tmpl", example.Name)
		}

		sourcePath := filepath.Join(examplesSourceDir, sourceFileName)

		// Check if the source template file actually exists.
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			// Skip if the template doesn't exist.
			continue
		}

		destDir := filepath.Join(samplesDestDir, serviceName)
		destPath := filepath.Join(destDir, sourceFileName)

		err := copyFile(sourcePath, destPath)
		if err != nil {
			// Log the error but don't stop processing other examples.
			log.Printf("Failed to copy '%s' to '%s': %v\n", sourcePath, destPath, err)
		} else {
			// fmt.Printf("  - Copied '%s'\n", sourceFileName)
			fmt.Print()
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("could not read source file %s: %w", src, err)
	}

	// Switch existing Vars to PrefixedVars
	output := strings.ReplaceAll(string(input), "$.Vars", "$.PrefixedVars")

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
