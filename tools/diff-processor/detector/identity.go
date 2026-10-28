package detector

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// MissingIdentityInfo holds the results for a single resource.
type MissingIdentityInfo struct {
	// MissingCRUD lists which CRUD functions are missing SetResourceIdentityAttributes.
	MissingCRUD []string `json:"missingCRUD,omitempty"`
	// MissingImportTest is true if no step uses ImportStateKind: ImportBlockWithResourceIdentity.
	MissingImportTest bool `json:"missingImportTest"`
}

// DetectMissingIdentityCoverage scans the provided resource files.
// servicesDir must be the direct parent of the resource_*.go files.
func DetectMissingIdentityCoverage(servicesDir string, changedResources []string) (map[string]*MissingIdentityInfo, error) {
	results := make(map[string]*MissingIdentityInfo)
	for _, resourceName := range changedResources {
		path := filepath.Join(servicesDir, "resource_"+resourceName+".go")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		hasIdentity, err := fileContainsIdentityBlock(path)
		if err != nil {
			return nil, err
		}
		if !hasIdentity {
			continue
		}

		missingCRUD, err := checkCRUDCoverage(path)
		if err != nil {
			return nil, err
		}

		missingImportTest, err := checkImportTest(path)
		if err != nil {
			return nil, err
		}

		if len(missingCRUD) > 0 || missingImportTest {
			results[resourceName] = &MissingIdentityInfo{
				MissingCRUD:       missingCRUD,
				MissingImportTest: missingImportTest,
			}
		}
	}
	return results, nil
}

func fileContainsIdentityBlock(path string) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(content), "ResourceIdentity"), nil
}

func checkCRUDCoverage(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	type funcInfo struct {
		name        string
		crudType    string
		startLine   int
		hasIdentity bool
	}

	funcPattern := regexp.MustCompile(`^func\s+\w+(Create|Read|Update)\w*\(`)
	setIdentityPattern := "SetResourceIdentityAttributes"

	var functions []*funcInfo
	var currentFunc *funcInfo

	scanner := bufio.NewScanner(file)
	lineNum := 0
	braceDepth := 0
	inFunc := false

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if matches := funcPattern.FindStringSubmatch(line); matches != nil {
			if currentFunc != nil {
				functions = append(functions, currentFunc)
			}
			currentFunc = &funcInfo{
				name:      line,
				crudType:  matches[1],
				startLine: lineNum,
			}
			inFunc = true
			braceDepth = 0
		}

		if inFunc && currentFunc != nil {
			braceDepth += strings.Count(line, "{") - strings.Count(line, "}")

			if strings.Contains(line, setIdentityPattern) {
				currentFunc.hasIdentity = true
			}

			if braceDepth <= 0 && lineNum > currentFunc.startLine {
				functions = append(functions, currentFunc)
				currentFunc = nil
				inFunc = false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	crudFound := map[string]bool{"Create": false, "Read": false, "Update": false}
	crudCovered := map[string]bool{"Create": false, "Read": false, "Update": false}

	for _, f := range functions {
		crudFound[f.crudType] = true
		if f.hasIdentity {
			crudCovered[f.crudType] = true
		}
	}

	var missing []string
	for _, crud := range []string{"Create", "Read", "Update"} {
		if crudFound[crud] && !crudCovered[crud] {
			missing = append(missing, crud)
		}
	}

	return missing, nil
}

func checkImportTest(resourcePath string) (bool, error) {
	dir := filepath.Dir(resourcePath)
	baseName := filepath.Base(resourcePath)
	resourcePrefix := strings.TrimSuffix(baseName, ".go")
	allowLegacyRefreshState := isLegacyIdentityImportOptedOut(resourcePath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return true, fmt.Errorf("error reading directory %s: %w", dir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") || !strings.HasPrefix(entry.Name(), resourcePrefix) {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return true, err
		}
		if hasValidIdentityImportStep(content, allowLegacyRefreshState) {
			return false, nil
		}
	}
	return true, nil
}

var mmv1ConfigURLPattern = regexp.MustCompile(`Configuration:\s+\S+/mmv1/(products/[^/\s]+/[^/\s]+\.yaml)`)

type legacyIdentityOptOutYAML struct {
	LegacyIdentityImportTest bool `yaml:"legacy_identity_import_test"`
}

// isLegacyIdentityImportOptedOut checks whether the MMv1 resource or its parent product
// has `legacy_identity_import_test: true` set in YAML.
func isLegacyIdentityImportOptedOut(resourcePath string) bool {
	content, err := os.ReadFile(resourcePath)
	if err != nil {
		return false
	}
	matches := mmv1ConfigURLPattern.FindSubmatch(content)
	if matches == nil {
		return false
	}
	relResourceYAML := filepath.FromSlash(string(matches[1]))
	relProductYAML := filepath.Join(filepath.Dir(relResourceYAML), "product.yaml")

	mmv1Root := findMMv1Root(filepath.Dir(resourcePath))
	if mmv1Root == "" {
		return false
	}

	if yamlHasLegacyIdentityOptOut(filepath.Join(mmv1Root, relResourceYAML)) {
		return true
	}
	if yamlHasLegacyIdentityOptOut(filepath.Join(mmv1Root, relProductYAML)) {
		return true
	}
	return false
}

func findMMv1Root(startDir string) string {
	candidates := []string{startDir}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, cwd)
	}
	for _, base := range candidates {
		dir := base
		for {
			mmv1Dir := filepath.Join(dir, "mmv1")
			if info, err := os.Stat(filepath.Join(mmv1Dir, "products")); err == nil && info.IsDir() {
				return mmv1Dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return ""
}

func yamlHasLegacyIdentityOptOut(yamlPath string) bool {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return false
	}
	var cfg legacyIdentityOptOutYAML
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return false
	}
	return cfg.LegacyIdentityImportTest
}

func hasValidIdentityImportStep(content []byte, allowLegacyRefreshState bool) bool {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, 0)
	if err != nil {
		return false
	}

	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}
		compLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		hasIdentityKind := false
		hasImportState := false
		hasRefreshState := false

		for _, elt := range compLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			keyIdent, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			switch keyIdent.Name {
			case "ImportStateKind":
				switch v := kv.Value.(type) {
				case *ast.SelectorExpr:
					if v.Sel.Name == "ImportBlockWithResourceIdentity" {
						hasIdentityKind = true
					}
				case *ast.Ident:
					if v.Name == "ImportBlockWithResourceIdentity" {
						hasIdentityKind = true
					}
				}
			case "ImportState":
				if ident, ok := kv.Value.(*ast.Ident); ok && ident.Name == "true" {
					hasImportState = true
				}
			case "RefreshState":
				if ident, ok := kv.Value.(*ast.Ident); ok && ident.Name == "true" {
					hasRefreshState = true
				}
			}
		}

		if hasIdentityKind {
			if (hasImportState && !hasRefreshState) || (allowLegacyRefreshState && hasRefreshState) {
				found = true
				return false
			}
		}
		return true
	})

	return found
}
