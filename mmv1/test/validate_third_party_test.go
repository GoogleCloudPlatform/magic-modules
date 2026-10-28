package test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTemplatesStillNeedToBeTemplates(t *testing.T) {
	// Get the directory where this test file is located
	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current test file path")
	}
	thirdPartyDir := filepath.Join(filepath.Dir(filepath.Dir(testFilePath)), "third_party")

	// Define the third_party directory relative to the test file
	terraformDir := filepath.Join(thirdPartyDir, "terraform")

	// Regular expression to match Go template syntax
	templateSyntaxRegex := regexp.MustCompile(`\{\{.*?\}\}`)

	// Track files that no longer need to be templates
	unnecessaryTemplates := []string{}

	// Walk through the third_party directory
	err := filepath.Walk(terraformDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Handle case where third_party directory doesn't exist
			if os.IsNotExist(err) && path == terraformDir {
				t.Logf("Warning: third_party/terraform directory not found at %s", terraformDir)
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
			relPath, _ := filepath.Rel(thirdPartyDir, path)
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

func TestWebsiteDocsFrontmatter(t *testing.T) {
	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current test file path")
	}
	docsDir := filepath.Join(filepath.Dir(filepath.Dir(testFilePath)), "third_party", "terraform", "website", "docs")

	err := filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) && path == docsDir {
				t.Logf("Warning: third_party/terraform/website/docs directory not found at %s", docsDir)
				return nil
			}
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".markdown" {
			return nil
		}

		relPath, _ := filepath.Rel(docsDir, path)
		subDir := strings.Split(relPath, string(filepath.Separator))[0]

		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: failed to read file: %v", relPath, err)
			return nil
		}

		if bytes.Contains(content, []byte("\r\n")) {
			t.Errorf("%s: contains CRLF (\\r\\n) line endings; use LF (\\n) so frontmatter headers are inserted properly", relPath)
		}

		text := strings.ReplaceAll(string(content), "\r\n", "\n")
		if !strings.HasPrefix(text, "---\n") {
			t.Errorf("%s: missing opening YAML frontmatter delimiter '---'", relPath)
			return nil
		}

		rest := strings.TrimPrefix(text, "---\n")
		endIdx := strings.Index(rest, "\n---\n")
		if endIdx == -1 {
			t.Errorf("%s: missing closing YAML frontmatter delimiter '---'", relPath)
			return nil
		}

		frontmatterRaw := rest[:endIdx]
		var fm map[string]any
		if err := yaml.Unmarshal([]byte(frontmatterRaw), &fm); err != nil {
			t.Errorf("%s: invalid YAML frontmatter: %v", relPath, err)
			return nil
		}

		if fm == nil {
			t.Errorf("%s: empty YAML frontmatter", relPath)
			return nil
		}

		if strings.Contains(relPath, string(filepath.Separator)) && subDir != "guides" && subDir != "functions" {
			subcat, ok := fm["subcategory"].(string)
			if !ok || strings.TrimSpace(subcat) == "" {
				t.Errorf("%s: missing or empty 'subcategory' in YAML frontmatter", relPath)
			}
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Error walking website docs directory: %v", err)
	}
}
