package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
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

func readAllTests(providerDir string) ([]*Test, error) {
	files, err := os.ReadDir(providerDir)
	if err != nil {
		return nil, err
	}
	allTests := make([]*Test, 0)
	errs := make([]error, 0)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_test.go") {
			tests, err := readTestFile(filepath.Join(providerDir, file.Name()))
			if err != nil {
				errs = append(errs, err)
			}
			allTests = append(allTests, tests...)
		}
	}
	if len(errs) > 0 {
		return allTests, fmt.Errorf("errors reading tests: %v", errs)
	}
	return allTests, nil
}

func readTestFile(filename string) ([]*Test, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}
	funcDecls := make(map[string]*ast.FuncDecl) // map of function names to function declarations
	varDecls := make(map[string]*ast.BasicLit)  // map of variable names to value expressions
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
	tests := make([]*Test, 0)
	errs := make([]error, 0)
	for name, funcDecl := range funcDecls {
		if strings.HasPrefix(name, "TestAcc") {
			test, err := readTestFunc(funcDecl, funcDecls, varDecls)
			if err != nil {
				errs = append(errs, err)
			}
			if test != nil {
				test.Name = name
				tests = append(tests, test)
			}
		}
	}
	if len(errs) > 0 {
		return tests, fmt.Errorf("errors encountered parsing test file: %v", errs)
	}
	return tests, nil
}

func readTestFunc(testFunc *ast.FuncDecl, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (*Test, error) {
	// This is an exported test function.
	for _, stmt := range testFunc.Body.List {
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				// This is a call expression.
				if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "VcrTest" {
					return readVcrTestCall(callExpr, funcDecls, varDecls)
				}
			}
		}
	}
	return nil, nil
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
						if configCallExpr, ok := keyValueExpr.Value.(*ast.CallExpr); ok {
							step, err := readConfigCallExpr(configCallExpr, funcDecls, varDecls)
							if err != nil {
								errs = append(errs, err)
							}
							test.Steps = append(test.Steps, step)
						} else if ident, ok := keyValueExpr.Value.(*ast.Ident); ok {
							if configVar, ok := varDecls[ident.Name]; ok {
								step, err := readConfigBasicLit(configVar)
								if err != nil {
									errs = append(errs, err)
								}
								test.Steps = append(test.Steps, step)
							}
						}
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

func readConfigCallExpr(configCallExpr *ast.CallExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (Step, error) {
	if ident, ok := configCallExpr.Fun.(*ast.Ident); ok {
		if configFunc, ok := funcDecls[ident.Name]; ok {
			return readConfigFunc(configFunc)
		}
		return nil, fmt.Errorf("failed to find function declaration %s", ident.Name)
	}
	return nil, fmt.Errorf("failed to get ident for %v", configCallExpr.Fun)
}

func readConfigFunc(configFunc *ast.FuncDecl) (Step, error) {
	for _, stmt := range configFunc.Body.List {
		if returnStmt, ok := stmt.(*ast.ReturnStmt); ok {
			for _, result := range returnStmt.Results {
				if callExpr, ok := result.(*ast.CallExpr); ok {
					if len(callExpr.Args) == 0 {
						return nil, fmt.Errorf("no arguments found for call expression %v in %v", callExpr, result)
					}
					if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
						return readConfigBasicLit(basicLit)
					}
					return nil, fmt.Errorf("no string literal found in arguments to call expression %v", callExpr)
				}
			}
			return nil, fmt.Errorf("failed to find a call expression in results %v", returnStmt.Results)
		}
	}
	return nil, fmt.Errorf("failed to find a return statement in %v", configFunc.Body.List)
}

func readConfigBasicLit(configBasicLit *ast.BasicLit) (Step, error) {
	if configStr, err := strconv.Unquote(configBasicLit.Value); err != nil {
		return nil, err
	} else {
		// Remove template variables because they interfere with hcl parsing.
		configStr = strings.ReplaceAll(configStr, "%", "")
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
