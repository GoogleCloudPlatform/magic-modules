package test

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func parseHCLBytes(src []byte, filePath string) (map[string]map[string]struct{}, error) {
	parser := hclparse.NewParser()
	hclFile, diags := parser.ParseHCL(src, filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse HCL: %w", diags)
	}

	if hclFile == nil {
		return nil, fmt.Errorf("parsed HCL file %s is nil cannot proceed", filePath)
	}

	parsed := make(map[string]map[string]struct{})

	for _, block := range hclFile.Body.(*hclsyntax.Body).Blocks {
		if block.Type == "resource" {
			if len(block.Labels) != 2 {
				log.Printf("Skipping address block with unexpected number of labels: %v", block.Labels)
				continue
			}

			resType := block.Labels[0]
			resName := block.Labels[1]
			addr := fmt.Sprintf("%s.%s", resType, resName)
			attrs, procDiags := parseHCLBody(block.Body)

			if procDiags.HasErrors() {
				log.Printf("Diagnostics while processing address %s.%s body in %s:", resType, resName, filePath)
				for _, diag := range procDiags {
					log.Printf("  - %s (Severity)", diag.Error())
				}
			}

			flattenedAttrs := make(map[string]struct{})
			flatten(attrs, "", flattenedAttrs)
			parsed[addr] = flattenedAttrs
		}
	}
	return parsed, nil
}

// parseHCLBody recursively parses attributes and nested blocks from an HCL body.
func parseHCLBody(body hcl.Body) (
	attributes map[string]any,
	diags hcl.Diagnostics,
) {
	attributes = make(map[string]any)
	var allDiags hcl.Diagnostics

	if syntaxBody, ok := body.(*hclsyntax.Body); ok {
		for _, attr := range syntaxBody.Attributes {
			insert(struct{}{}, attr.Name, attributes)
		}

		for _, block := range syntaxBody.Blocks {
			nestedAttr, diags := parseHCLBody(block.Body)
			if diags.HasErrors() {
				allDiags = append(allDiags, diags...)
			}

			insert(nestedAttr, block.Type, attributes)
		}
	} else {
		allDiags = append(allDiags, &hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Body type assertion to *hclsyntax.Body failed",
			Detail:   fmt.Sprintf("Cannot directly parse attributes for body of type %T. Attribute parsing may be incomplete.", body),
		})
	}

	return attributes, allDiags
}

func insert(data any, key string, parent map[string]any) {
	if existing, ok := parent[key]; ok {
		if existingSlice, ok := existing.([]any); ok {
			existingSlice = append(existingSlice, data)
		} else {
			// Until we see a second instance of a repeated block or attribute, it will look non-repeated.
			parent[key] = []any{existing, data}
		}
	} else {
		parent[key] = data
	}
}

func flatten(data interface{}, prefix string, result map[string]struct{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = prefix + "." + key
			}
			flatten(value, newPrefix, result)
		}
	case []interface{}:
		if len(v) == 0 && prefix != "" {
			result[prefix] = struct{}{}
		}
		for i, value := range v {
			newPrefix := fmt.Sprintf("%s.%d", prefix, i)
			flatten(value, newPrefix, result)
		}
	default:
		if prefix != "" {
			result[prefix] = struct{}{}
		}
	}
}
