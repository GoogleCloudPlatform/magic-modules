package detector

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	if currentFunc != nil {
		functions = append(functions, currentFunc)
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
	entries, err := os.ReadDir(dir)
	if err != nil {
		return true, fmt.Errorf("error reading directory %s: %w", dir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		if !strings.HasPrefix(entry.Name(), resourcePrefix) {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return true, err
		}
		// Updated to check for the actual ImportStateKind configuration
		if strings.Contains(string(content), "ImportStateKind: resource.ImportBlockWithResourceIdentity,") {
			return false, nil
		}
	}
	return true, nil
}
