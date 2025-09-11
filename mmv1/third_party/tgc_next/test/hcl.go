package test

import (
	"fmt"
	"log"
	"sort"
	"strings"

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
			parent[key] = append(existingSlice, data)
		} else {
			// Until we see a second instance of a repeated block or attribute, it will look non-repeated.
			parent[key] = []any{existing, data}
		}
	} else {
		parent[key] = data
	}
}

func flatten(data any, prefix string, result map[string]struct{}) {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = prefix + "." + key
			}
			flatten(value, newPrefix, result)
		}
	case []any:
		flattenSlice(prefix, v, result)
	default:
		if prefix != "" {
			result[prefix] = struct{}{}
		}
	}
}

func flattenSlice(prefix string, v []any, result map[string]struct{}) {
	if len(v) == 0 && prefix != "" {
		result[prefix] = struct{}{}
		return
	}

	type sortableElement struct {
		flatKeys  string
		flattened map[string]struct{}
	}

	sortable := make([]sortableElement, len(v))
	for i, value := range v {
		flattened := make(map[string]struct{})
		flatten(value, "", flattened)
		keys := make([]string, 0, len(flattened))
		for k := range flattened {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sortable[i] = sortableElement{
			flatKeys:  strings.Join(keys, ";"),
			flattened: flattened,
		}
	}

	sort.Slice(sortable, func(i, j int) bool {
		return sortable[i].flatKeys < sortable[j].flatKeys
	})

	for i, element := range sortable {
		newPrefix := fmt.Sprintf("%s.%d", prefix, i)
		if len(element.flattened) == 0 {
			if newPrefix != "" {
				result[newPrefix] = struct{}{}
			}
		} else {
			for k := range element.flattened {
				newKey := newPrefix
				if k != "" {
					newKey = newPrefix + "." + k
				}
				result[newKey] = struct{}{}
			}
		}
	}
}
