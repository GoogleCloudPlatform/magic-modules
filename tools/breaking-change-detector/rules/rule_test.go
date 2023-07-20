package rules

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestUniqueRuleIdentifiers(t *testing.T) {
	identifiers := getArrayOfIdentifiers()
	// Create a map to track the identifiers
	identifierMap := make(map[string]bool)
	for _, id := range identifiers {
		if identifierMap[id] {
			t.Errorf("Duplicate identifier found: %s", id)
		}

		// Add the identifier to the map
		identifierMap[id] = true
	}
}

func TestMarkdownIdentifiers(t *testing.T) {
	// Define the Markdown file path relative to the importer
	mdFilePath := "../../../docs/content/reference/breaking-change-detector.md"

	// Read the Markdown file
	mdContent, err := ioutil.ReadFile(mdFilePath)
	if err != nil {
		t.Fatalf("Failed to read or find Markdown file: %v", err)
	}

	// Convert the Markdown content to a string
	mdString := string(mdContent)

	// Define the identifiers to check
	identifiers := getArrayOfIdentifiers()

	// Iterate over the identifiers and check if they have a corresponding <h4> tag
	for _, identifier := range identifiers {
		// Define the expected <h4> tag
		expectedTag := fmt.Sprintf("<h4 id=\"%s\">", identifier)

		// Check if the <h4> tag exists in the Markdown string
		if !strings.Contains(mdString, expectedTag) {
			t.Errorf("Identifier %s does not have a corresponding <h4> tag", identifier)
		}
	}
}

func getArrayOfIdentifiers() []string {
	var output []string
	ruleCats := GetRules()
	var rules []Rule
	for _, rc := range ruleCats.Categories {
		rules = append(rules, rc.Rules...)
	}

	for _, r := range rules {
		identifier := r.Identifier()
		output = append(output, identifier)
	}

	return output
}
