package main

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

func TestTemplatesStillNeedToBeTemplates(t *testing.T) {
	// Get the directory where this test file is located
	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current test file path")
	}
	testDir := filepath.Dir(testFilePath)

	// Define the third_party directory relative to the test file
	thirdPartyDir := filepath.Join(testDir, "third_party", "terraform")

	// Regular expression to match Go template syntax
	templateSyntaxRegex := regexp.MustCompile(`\{\{.*?\}\}`)

	// Track files that no longer need to be templates
	unnecessaryTemplates := []string{}

	// Walk through the third_party directory
	err := filepath.Walk(thirdPartyDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Handle case where third_party directory doesn't exist
			if os.IsNotExist(err) && path == thirdPartyDir {
				t.Logf("Warning: third_party directory not found at %s", thirdPartyDir)
				return nil
			}
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only check .tmpl files
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			t.Logf("Error reading file %s: %v", path, err)
			return nil
		}

		// Check if file contains any Go template syntax
		hasTemplateSyntax := templateSyntaxRegex.Match(content)

		// If no template syntax found, add to the list
		if !hasTemplateSyntax {
			// Get relative path for cleaner output
			relPath, _ := filepath.Rel(testDir, path)
			unnecessaryTemplates = append(unnecessaryTemplates, relPath)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Error walking directory: %v", err)
	}

	// Output results at the end
	if len(unnecessaryTemplates) > 0 {
		t.Errorf("\nThe following %d .tmpl files in third_party directory don't contain any template syntax "+
			"and no longer need to be templates:\n", len(unnecessaryTemplates))

		for _, file := range unnecessaryTemplates {
			t.Errorf("  - %s", file)
		}

		t.Errorf("\nConsider removing the .tmpl extension from these files.")
	} else {
		t.Logf("All .tmpl files in third_party directory properly contain template syntax.")
	}
}
