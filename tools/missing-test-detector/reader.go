package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Resource map[string]any // config of one resource in a test

type Resources map[string]Resource // map of resource names to resource configs

type Step map[string]Resources // map of resource types to resources of that type

type Test struct {
	Name  string
	Steps []Step
}

func (t *Test) String() string {
	return fmt.Sprintf("%s: %#v", t.Name, t.Steps)
}

// Return a slice of tests as well as a map of file or test names to errors encountered.
func readAllTests(servicesDir string) ([]*Test, map[string]error) {
	dirs, err := os.ReadDir(servicesDir)
	if err != nil {
		return nil, map[string]error{servicesDir: err}
	}
	allTests := make([]*Test, 0)
	allErrs := make(map[string]error)
	for _, dir := range dirs {
		servicePath := filepath.Join(servicesDir, dir.Name())
		files, err := os.ReadDir(servicePath)
		if err != nil {
			return nil, map[string]error{servicePath: err}
		}
		var testFileNames []string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), "_test.go") {
				testFileNames = append(testFileNames, filepath.Join(servicePath, file.Name()))
			}
		}
		serviceTests, serviceErrs := readTestFiles(testFileNames)
		for fileName, err := range serviceErrs {
			allErrs[fileName] = err
		}
		allTests = append(allTests, serviceTests...)
	}
	if len(allErrs) > 0 {
		return allTests, allErrs
	}
	return allTests, nil
}

// Read all the test files in a service directory together to capture cross-file function usage.
func readTestFiles(filenames []string) ([]*Test, map[string]error) {
	funcDecls := make(map[string]*ast.FuncDecl) // map of function names to function declarations
	varDecls := make(map[string]*ast.BasicLit)  // map of variable names to value expressions
	errs := make(map[string]error)              // map of file or test names to errors encountered parsing
	fset := token.NewFileSet()
	for _, filename := range filenames {
		f, err := parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			errs[filename] = err
			continue
		}
		for _, decl := range f.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				// This is a function declaration.
				funcDecls[funcDecl.Name.Name] = funcDecl
			} else if genDecl, ok := decl.(*ast.GenDecl); ok {
				// This is an import, constant, type, or variable declaration
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						if len(valueSpec.Values) > 0 {
							if basicLit, ok := valueSpec.Values[0].(*ast.BasicLit); ok {
								varDecls[valueSpec.Names[0].Name] = basicLit
							}
						}
					}
				}
			}
		}
	}
	tests := make([]*Test, 0)
	for name, funcDecl := range funcDecls {
		if strings.HasPrefix(name, "TestAcc") {
			funcTests, err := readTestFunc(funcDecl, funcDecls, varDecls)
			if err != nil {
				errs[name] = err
			}
			tests = append(tests, funcTests...)
		}
	}
	if len(errs) > 0 {
		return tests, errs
	}
	return tests, nil
}

func readTestFunc(testFunc *ast.FuncDecl, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) ([]*Test, error) {
	// This is an exported test function.
	var tests []*Test
	var errs []error
	vars := make(map[string]*ast.CompositeLit, len(testFunc.Body.List)) // map of variable names to composite literal values in function body
	for _, stmt := range testFunc.Body.List {
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				// This is a call expression.
				ident, isIdent := callExpr.Fun.(*ast.Ident)
				selExpr, isSelExpr := callExpr.Fun.(*ast.SelectorExpr)
				if isIdent && ident.Name == "VcrTest" || isSelExpr && selExpr.Sel.Name == "VcrTest" {
					test, err := readVcrTestCall(callExpr, funcDecls, varDecls)
					if err != nil {
						errs = append(errs, err)
					}
					test.Name = testFunc.Name.Name
					tests = append(tests, test)
				}
			}
		} else if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
			if len(assignStmt.Lhs) == 1 && len(assignStmt.Rhs) == 1 {
				// For now, only allow single assignment variables for serial test maps.
				// e.g. testCases := map[string]func(t *testing.T) {...
				if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
					if rhsCompLit, ok := assignStmt.Rhs[0].(*ast.CompositeLit); ok {
						vars[ident.Name] = rhsCompLit
					}
				}
			}
		} else if rangeStmt, ok := stmt.(*ast.RangeStmt); ok {
			if ident, ok := rangeStmt.X.(*ast.Ident); ok {
				if varCompLit, ok := vars[ident.Name]; ok {
					serialTests, serialErrs := readSerialTestCompLit(varCompLit, funcDecls, varDecls)
					errs = append(errs, serialErrs...)
					tests = append(tests, serialTests...)
				}
			}
		}
	}
	if len(errs) > 0 {
		return tests, fmt.Errorf("errors reading test func %s: %v", testFunc.Name.Name, errs)
	}
	return tests, nil
}

// Reads a composite literal which is either a slice or a map of serialized test functions.
func readSerialTestCompLit(varCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) ([]*Test, []error) {
	var tests []*Test
	var errs []error
	for _, elt := range varCompLit.Elts {
		if eltKeyValueExpr, ok := elt.(*ast.KeyValueExpr); ok {
			eltTests, err := readSerialTestEltKeyValueExpr(eltKeyValueExpr, funcDecls, varDecls)
			if err != nil {
				errs = append(errs, err)
			}
			tests = append(tests, eltTests...)
		}
	}
	return tests, errs
}

func readSerialTestEltKeyValueExpr(eltKeyValueExpr *ast.KeyValueExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) ([]*Test, error) {
	if ident, ok := eltKeyValueExpr.Value.(*ast.Ident); ok {
		if testFunc, ok := funcDecls[ident.Name]; ok {
			return readTestFunc(testFunc, funcDecls, varDecls)
		}
		return nil, fmt.Errorf("failed to find function with name %s", ident.Name)
	}
	return nil, fmt.Errorf("element key value expression with key %+v had non-ident value %+v", eltKeyValueExpr.Key, eltKeyValueExpr.Value)
}

func readVcrTestCall(vcrTestCall *ast.CallExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (*Test, error) {
	for _, arg := range vcrTestCall.Args {
		if vcrTestArgCompLit, ok := arg.(*ast.CompositeLit); ok {
			if selExpr, ok := vcrTestArgCompLit.Type.(*ast.SelectorExpr); ok {
				if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "resource" && selExpr.Sel.Name == "TestCase" {
					return readTestCaseCompLit(vcrTestArgCompLit, funcDecls, varDecls)
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to find TestCase in %v", vcrTestCall.Args)
}

func readTestCaseCompLit(testCaseCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (*Test, error) {
	for _, elt := range testCaseCompLit.Elts {
		if keyValueExpr, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := keyValueExpr.Key.(*ast.Ident); ok && ident.Name == "Steps" {
				if stepsCompLit, ok := keyValueExpr.Value.(*ast.CompositeLit); ok {
					return readStepsCompLit(stepsCompLit, funcDecls, varDecls)
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to find Steps in %v", testCaseCompLit.Elts)
}

func readStepsCompLit(stepsCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (*Test, error) {
	test := &Test{}
	errs := make([]error, 0)
	for _, elt := range stepsCompLit.Elts {
		if eltCompLit, ok := elt.(*ast.CompositeLit); ok {
			for _, eltCompLitElt := range eltCompLit.Elts {
				if keyValueExpr, ok := eltCompLitElt.(*ast.KeyValueExpr); ok {
					if ident, ok := keyValueExpr.Key.(*ast.Ident); ok && ident.Name == "Config" {
						var configStr string
						var err error
						if configCallExpr, ok := keyValueExpr.Value.(*ast.CallExpr); ok {
							configStr, err = readConfigCallExpr(configCallExpr, funcDecls, varDecls)
						} else if ident, ok := keyValueExpr.Value.(*ast.Ident); ok {
							if configVar, ok := varDecls[ident.Name]; ok {
								configStr, err = strconv.Unquote(configVar.Value)
							}
						}
						if err != nil {
							errs = append(errs, err)
						}
						step, err := readConfigStr(configStr)
						if err != nil {
							errs = append(errs, err)
						}
						test.Steps = append(test.Steps, step)
					}
				}
			}
		}
	}
	if len(errs) > 0 {
		return test, fmt.Errorf("errors reading test steps: %v", errs)
	}
	return test, nil
}

// Read the call expression in the public test function that returns the config.
func readConfigCallExpr(configCallExpr *ast.CallExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (string, error) {
	if ident, ok := configCallExpr.Fun.(*ast.Ident); ok {
		if configFunc, ok := funcDecls[ident.Name]; ok {
			return readConfigFunc(configFunc, funcDecls, varDecls)
		}
		return "", fmt.Errorf("failed to find function declaration %s", ident.Name)
	}
	return "", fmt.Errorf("failed to get ident for %v", configCallExpr.Fun)
}

func readConfigFunc(configFunc *ast.FuncDecl, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (string, error) {
	for _, stmt := range configFunc.Body.List {
		if returnStmt, ok := stmt.(*ast.ReturnStmt); ok {
			if len(returnStmt.Results) > 0 {
				return readConfigFuncResult(returnStmt.Results[0], funcDecls, varDecls)
			}
			return "", fmt.Errorf("failed to find a config string in results %v", returnStmt.Results)
		}
	}
	return "", fmt.Errorf("failed to find a return statement in %v", configFunc.Body.List)
}

// Read the return result of a config func and return the config string.
func readConfigFuncResult(result ast.Expr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (string, error) {
	if basicLit, ok := result.(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
		return strconv.Unquote(basicLit.Value)
	} else if callExpr, ok := result.(*ast.CallExpr); ok {
		return readConfigFuncCallExpr(callExpr, funcDecls, varDecls)
	} else if binaryExpr, ok := result.(*ast.BinaryExpr); ok {
		xConfigStr, err := readConfigFuncResult(binaryExpr.X, funcDecls, varDecls)
		if err != nil {
			return "", err
		}
		yConfigStr, err := readConfigFuncResult(binaryExpr.Y, funcDecls, varDecls)
		if err != nil {
			return "", err
		}
		return xConfigStr + yConfigStr, nil
	}
	return "", fmt.Errorf("unknown config func result %v (%T)", result, result)
}

// Read the call expression in the config function that returns the config string.
// The call expression can contain a nested call expression.
// Return the config string.
func readConfigFuncCallExpr(configFuncCallExpr *ast.CallExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (string, error) {
	if len(configFuncCallExpr.Args) > 0 {
		if basicLit, ok := configFuncCallExpr.Args[0].(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
			return strconv.Unquote(basicLit.Value)
		} else if nestedCallExpr, ok := configFuncCallExpr.Args[0].(*ast.CallExpr); ok {
			return readConfigFuncCallExpr(nestedCallExpr, funcDecls, varDecls)
		}
	}
	// Config string not readable from args, attempt to read call expression as a helper function.
	return readConfigCallExpr(configFuncCallExpr, funcDecls, varDecls)
}

// Read the config string and return a test step.
func readConfigStr(configStr string) (Step, error) {
	// Remove template variables because they interfere with hcl parsing.
	pattern := regexp.MustCompile("%({[^{}]*}|[vdts])")
	// Replace with a value that can be parsed outside quotation marks.
	configStr = pattern.ReplaceAllString(configStr, "true")
	parser := hclparse.NewParser()
	file, diagnostics := parser.ParseHCL([]byte(configStr), "config.hcl")
	if diagnostics.HasErrors() {
		return nil, fmt.Errorf("errors parsing hcl: %v", diagnostics.Errs())
	}
	content, diagnostics := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
			},
			{
				Type:       "data",
				LabelNames: []string{"type", "name"},
			},
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
			{
				Type: "locals",
			},
		},
	})
	if diagnostics.HasErrors() {
		return nil, fmt.Errorf("errors getting hcl body content: %v", diagnostics.Errs())
	}
	m := make(map[string]Resources)
	errs := make([]error, 0)
	for _, block := range content.Blocks {
		if len(block.Labels) != 2 {
			continue
		}
		if _, ok := m[block.Labels[0]]; !ok {
			// Create an empty map for this resource type.
			m[block.Labels[0]] = make(Resources)
		}
		// Use the resource name as a key.
		resourceConfig, err := readHCLBlockBody(block.Body, file.Bytes)
		if err != nil {
			errs = append(errs, err)
		}
		m[block.Labels[0]][block.Labels[1]] = resourceConfig
	}
	if len(errs) > 0 {
		return m, fmt.Errorf("errors reading hcl blocks: %v", errs)
	}
	return m, nil
}

func readHCLBlockBody(body hcl.Body, fileBytes []byte) (Resource, error) {
	var m Resource
	gohcl.DecodeBody(body, nil, &m)
	for k, v := range m {
		if attr, ok := v.(*hcl.Attribute); ok {
			m[k] = string(attr.Expr.Range().SliceBytes(fileBytes))
		}
	}
	syntaxBody, ok := body.(*hclsyntax.Body)
	if !ok {
		return m, fmt.Errorf("couldn't get hclsyntax body from %v", body)
	}
	errs := make([]error, 0)
	for _, block := range syntaxBody.Blocks {
		blockConfig, err := readHCLBlockBody(block.Body, fileBytes)
		if err != nil {
			errs = append(errs, err)
		}
		if existing, ok := m[block.Type]; ok {
			// Merge the fields from the current block into the existing resource config.
			if existingResource, ok := existing.(Resource); ok {
				mergeResources(existingResource, blockConfig)
			}
		} else {
			m[block.Type] = blockConfig
		}
	}
	if len(errs) > 0 {
		return m, fmt.Errorf("errors reading hcl blocks: %v", errs)
	}
	return m, nil
}

// Perform a recursive one-way merge of b into a.
func mergeResources(a, b Resource) {
	for k, bv := range b {
		if av, ok := a[k]; ok {
			if avr, ok := av.(Resource); ok {
				if bvr, ok := bv.(Resource); ok {
					mergeResources(avr, bvr)
				}
			}
		} else {
			a[k] = bv
		}
	}
}
