package test

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func parseHCLBytes(src []byte, filePath string) (map[string]map[string]any, error) {
	parser := hclparse.NewParser()
	hclFile, diags := parser.ParseHCL(src, filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse HCL: %w", diags)
	}

	if hclFile == nil {
		return nil, fmt.Errorf("parsed HCL file %s is nil cannot proceed", filePath)
	}

	parsed := make(map[string]map[string]any)

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

			flattenedAttrs := make(map[string]any)
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
			insert(getValue(attr.Expr), attr.Name, attributes)
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

func flatten(data any, prefix string, result map[string]any) {
	switch v := data.(type) {
	case map[string]any:
		if len(v) == 0 && prefix != "" {
			result[prefix] = v
		} else {
			for key, value := range v {
				newPrefix := key
				if prefix != "" {
					newPrefix = prefix + "." + key
				}
				flatten(value, newPrefix, result)
			}
		}
	case []any:
		flattenSlice(prefix, v, result)
	default:
		if prefix != "" {
			result[prefix] = v
		}
	}
}

func flattenSlice(prefix string, v []any, result map[string]any) {
	if len(v) == 0 && prefix != "" {
		result[prefix] = struct{}{}
		return
	}

	type sortableElement struct {
		flatKeys  string
		flattened map[string]any
	}

	sortable := make([]sortableElement, len(v))
	for i, value := range v {
		flattened := make(map[string]any)
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
			for k, v := range element.flattened {
				newKey := newPrefix
				if k != "" {
					newKey = newPrefix + "." + k
				}
				result[newKey] = v
			}
		}
	}
}

// Gets the value of the expression of an attribute
func getValue(expr hcl.Expression) any {
	switch expr := expr.(type) {
	case *hclsyntax.ScopeTraversalExpr:
		// Example: id = google_instance.web.id
		return getTraveralExprVal(expr.Traversal)
	case *hclsyntax.LiteralValueExpr:
		// Example: region = "us-west1"
		return convertValue(expr.Val)
	case *hclsyntax.TemplateExpr:
		// Example: ip_address = "IP: ${var.ip}"
		vStr := ""
		parts := expr.Parts
		for _, part := range parts {
			vStr += getValue(part).(string)
		}
		return vStr
	case *hclsyntax.TupleConsExpr:
		// Example: methods = ["GET", "POST", "DELETE"]
		exprV := make([]string, 0, len(expr.Exprs))
		for _, elem := range expr.Exprs {
			exprV = append(exprV, fmt.Sprint(getValue(elem)))
		}
		return strings.Join(exprV, ",")
	case *hclsyntax.ObjectConsKeyExpr:
		// Example: labels = {(local.key_name) = "value"}
		return getValue(expr.Wrapped).(string)
	case *hclsyntax.ObjectConsExpr:
		// Example: tags = { Env = "dev", Owner = var.user }
		return map[string]any{}
	case *hclsyntax.ParenthesesExpr:
		// Example: (local.key_name)
		return getValue(expr.Expression)
	default:
		log.Printf("Unsupported expression type: %T", expr)
		return nil
	}
}

// Converts a literal cty.Value to a standard Go value
func convertValue(val cty.Value) any {
	switch val.Type() {
	case cty.Number:
		f, _ := val.AsBigFloat().Float64()
		return f
	case cty.String:
		return val.AsString()
	case cty.Bool:
		var v bool
		_ = gocty.FromCtyValue(val, &v)
		return v
	default:
		return ""
	}
}

// Gets the value for the traveral expression (e.g. "google_pubsub_topic.example.id")
func getTraveralExprVal(traversal hcl.Traversal) string {
	exprV := make([]string, 0, len(traversal))

	for _, v := range traversal {
		switch v := v.(type) {
		// The starting point of the traversal (e.g. "google_pubsub_topic)
		case hcl.TraverseRoot:
			exprV = append(exprV, v.Name)
		// A list of hclsyntax.Traversal objects, which represent each step in the path.
		// In "google_pubsub_topic.example.id", the steps are "example", "id"
		case hcl.TraverseAttr:
			exprV = append(exprV, v.Name)
		}
	}

	return strings.Join(exprV, ".")
}
