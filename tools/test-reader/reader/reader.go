package reader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Block is the flattened set of attributes in one block body, keyed by
// attribute path, e.g. {"boot_disk.initialize_params.image": "..."}.
type Block map[string]any

// Blocks is a map of name labels to blocks. The name label is a block's second
// label, e.g. "foobar" in `resource "google_compute_instance" "foobar"`.
type Blocks map[string]Block

// TypeLabels is a map of type labels to the blocks carrying that label. The
// type label is a block's first label, e.g. "google_compute_instance" in
// `resource "google_compute_instance" "foobar"`.
type TypeLabels map[string]Blocks

// BlockType is the type of a top-level block, in HCL's sense of the term:
// "resource", "data", "list" or "ephemeral". Not to be confused with a block's
// type label, which names the resource type (see TypeLabels).
//
// It holds the literal keyword from the config rather than a closed set, so
// block types Terraform adds in the future are recorded under their own type
// rather than being dropped or attributed to an existing one.
type BlockType string

const (
	ResourceBlock  BlockType = "resource"
	DataBlock      BlockType = "data"
	ListBlock      BlockType = "list"
	EphemeralBlock BlockType = "ephemeral"
)

// Step is a map of block types to the blocks declared with that type.
// Block types are kept separate because a data, list or ephemeral block can
// share a type label with a resource block without sharing its schema, so an
// attribute set on one says nothing about the same-named attribute on another.
type Step map[BlockType]TypeLabels

type Test struct {
	Name          string
	Steps         []Step
	resolvedSteps []Step
}

// EffectiveSteps returns per-step HCL blocks with Nprintf and Sprintf argument
// bindings resolved across multi-step tests when available, falling back to Steps.
func (t *Test) EffectiveSteps() []Step {
	if len(t.resolvedSteps) == len(t.Steps) && len(t.resolvedSteps) > 0 {
		return t.resolvedSteps
	}
	return t.Steps
}

func (t *Test) String() string {
	return fmt.Sprintf("%s: %#v", t.Name, t.Steps)
}

// Return a slice of tests as well as a map of file or test names to errors encountered.
func ReadAllTests(servicesDir string) ([]*Test, map[string]error) {
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
		serviceTests, serviceErrs := ReadTestFiles(testFileNames)
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
func ReadTestFiles(filenames []string) ([]*Test, map[string]error) {
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
	scope := newEvalScope(varDecls)
	for _, stmt := range testFunc.Body.List {
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				// This is a call expression.
				ident, isIdent := callExpr.Fun.(*ast.Ident)
				selExpr, isSelExpr := callExpr.Fun.(*ast.SelectorExpr)
				if isIdent && ident.Name == "VcrTest" || isSelExpr && selExpr.Sel.Name == "VcrTest" {
					test, err := readVcrTestCallWithScope(callExpr, funcDecls, varDecls, scope)
					if err != nil {
						errs = append(errs, err)
					}
					if test != nil {
						test.Name = testFunc.Name.Name
						tests = append(tests, test)
					}
				}
			}
		} else if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
			recordAssignInScope(assignStmt, scope)
			if len(assignStmt.Lhs) == 1 && len(assignStmt.Rhs) == 1 {
				// For now, only allow single assignment variables for serial test maps.
				// e.g. testCases := map[string]func(t *testing.T) {...
				if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
					if rhsCompLit, ok := assignStmt.Rhs[0].(*ast.CompositeLit); ok && isSerialTestMap(rhsCompLit) {
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

// Report whether a composite literal is a serial test map: a map from test
// names to test functions, e.g. map[string]func(t *testing.T){...}. Ranging
// over a map is not enough to identify one, since table-driven tests range
// over a map of test cases too.
func isSerialTestMap(compLit *ast.CompositeLit) bool {
	mapType, ok := compLit.Type.(*ast.MapType)
	if !ok {
		return false
	}
	_, ok = mapType.Value.(*ast.FuncType)
	return ok
}

// Reads a composite literal which is a map of serialized test functions.
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
	return readVcrTestCallWithScope(vcrTestCall, funcDecls, varDecls, nil)
}

func readVcrTestCallWithScope(vcrTestCall *ast.CallExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, scope *evalScope) (*Test, error) {
	for _, arg := range vcrTestCall.Args {
		if vcrTestArgCompLit, ok := arg.(*ast.CompositeLit); ok {
			if selExpr, ok := vcrTestArgCompLit.Type.(*ast.SelectorExpr); ok {
				if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "resource" && selExpr.Sel.Name == "TestCase" {
					return readTestCaseCompLitWithScope(vcrTestArgCompLit, funcDecls, varDecls, scope)
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to find TestCase in %v", vcrTestCall.Args)
}

func readTestCaseCompLit(testCaseCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (*Test, error) {
	return readTestCaseCompLitWithScope(testCaseCompLit, funcDecls, varDecls, nil)
}

func readTestCaseCompLitWithScope(testCaseCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, scope *evalScope) (*Test, error) {
	for _, elt := range testCaseCompLit.Elts {
		if keyValueExpr, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := keyValueExpr.Key.(*ast.Ident); ok && ident.Name == "Steps" {
				if stepsCompLit, ok := keyValueExpr.Value.(*ast.CompositeLit); ok {
					return readStepsCompLitWithScope(stepsCompLit, funcDecls, varDecls, scope)
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to find Steps in %v", testCaseCompLit.Elts)
}

func unwrapStepCompLit(elt ast.Expr) *ast.CompositeLit {
	if compLit, ok := elt.(*ast.CompositeLit); ok {
		return compLit
	}
	if unaryExpr, ok := elt.(*ast.UnaryExpr); ok && unaryExpr.Op == token.AND {
		if compLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
			return compLit
		}
	}
	return nil
}

func readStepsCompLit(stepsCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (*Test, error) {
	return readStepsCompLitWithScope(stepsCompLit, funcDecls, varDecls, nil)
}

func readStepsCompLitWithScope(stepsCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, scope *evalScope) (*Test, error) {
	test := &Test{}
	errs := make([]error, 0)
	var configExprs []ast.Expr
	for _, elt := range stepsCompLit.Elts {
		eltCompLit := unwrapStepCompLit(elt)
		if eltCompLit == nil || shouldSkipStep(eltCompLit) {
			continue
		}
		for _, eltCompLitElt := range eltCompLit.Elts {
			if keyValueExpr, ok := eltCompLitElt.(*ast.KeyValueExpr); ok {
				if ident, ok := keyValueExpr.Key.(*ast.Ident); ok && ident.Name == "Config" {
					// A config is any expression that evaluates to a config
					// string, and reads the same way here as it does in a
					// config func's return statement.
					configStr, err := readConfigFuncResult(keyValueExpr.Value, funcDecls, varDecls)
					if err != nil {
						errs = append(errs, err)
					}
					step, err := readConfigStr(configStr)
					if err != nil {
						errs = append(errs, err)
					}
					test.Steps = append(test.Steps, step)
					configExprs = append(configExprs, keyValueExpr.Value)
				}
			}
		}
	}
	if len(errs) > 0 {
		return test, fmt.Errorf("errors reading test steps: %v", errs)
	}
	if len(test.Steps) > 1 && scope != nil {
		resolved := make([]Step, 0, len(configExprs))
		for idx, cfgExpr := range configExprs {
			if resStr, ok := readResolvedConfigExpr(cfgExpr, funcDecls, varDecls, scope); ok {
				if resStep, resErr := readConfigStr(resStr); resErr == nil {
					resolved = append(resolved, resStep)
					continue
				}
			}
			resolved = append(resolved, test.Steps[idx])
		}
		if !reflect.DeepEqual(resolved, test.Steps) {
			test.resolvedSteps = resolved
		}
	}
	return test, nil
}

// Report whether a test step composite literal does not apply its config, so
// its fields should not count towards coverage: steps with ExpectError set
// expect the config to be rejected, and steps with PlanOnly, Destroy,
// ImportState or RefreshState set to true run a plan, destroy, import or
// refresh rather than applying the config.
func shouldSkipStep(stepCompLit *ast.CompositeLit) bool {
	for _, elt := range stepCompLit.Elts {
		if keyValueExpr, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := keyValueExpr.Key.(*ast.Ident); ok {
				switch ident.Name {
				case "ExpectError":
					if valIdent, ok := keyValueExpr.Value.(*ast.Ident); !ok || valIdent.Name != "nil" {
						return true
					}
				case "PlanOnly", "Destroy", "ImportState", "RefreshState":
					if valIdent, ok := keyValueExpr.Value.(*ast.Ident); ok && valIdent.Name == "true" {
						return true
					}
				}
			}
		}
	}
	return false
}

// Read a call expression that produces a config and return the config string.
// The call is either to a config func declared in the same package, or to a
// formatting function that builds a config out of its arguments.
func readConfigCallExpr(configCallExpr *ast.CallExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (string, error) {
	// Dispatch on the callee rather than on the shape of the arguments: a
	// config func can take a string as its first argument, and that string is a
	// value to interpolate, not a config.
	if ident, ok := configCallExpr.Fun.(*ast.Ident); ok {
		if configFunc, ok := funcDecls[ident.Name]; ok {
			return readConfigFunc(configFunc, funcDecls, varDecls)
		}
		return "", fmt.Errorf("failed to find function declaration %s", ident.Name)
	}
	selectorExpr, ok := configCallExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", fmt.Errorf("failed to read config from call to %v (%T)", configCallExpr.Fun, configCallExpr.Fun)
	}
	// A function from another package can't be followed, so only the formatting
	// functions below, whose result is known from their arguments, can be read.
	// Reading the arguments of anything else would report fields as covered
	// that the test may never apply.
	switch selectorExpr.Sel.Name {
	case "Sprintf", "Nprintf":
		// The first argument is the config template and the rest are values
		// interpolated into it.
		if len(configCallExpr.Args) == 0 {
			return "", fmt.Errorf("failed to find a config template in call to %s.%s", selectorExpr.X, selectorExpr.Sel.Name)
		}
		configStr, err := readConfigFuncResult(configCallExpr.Args[0], funcDecls, varDecls)
		if err != nil || selectorExpr.Sel.Name != "Sprintf" {
			// Nprintf only replaces %{key}, so "%%" is not an escape there.
			return configStr, err
		}
		return markEscapedPercents(configStr), nil
	case "Sprint":
		// Every argument is part of the config. This is used to join a shared
		// setup config to the config under test.
		var configStr strings.Builder
		for _, arg := range configCallExpr.Args {
			argConfigStr, err := readConfigFuncResult(arg, funcDecls, varDecls)
			if err != nil {
				return "", err
			}
			configStr.WriteString(argConfigStr)
		}
		return configStr.String(), nil
	}
	return "", fmt.Errorf("failed to read config from call to %s.%s, which is not a config formatting function", selectorExpr.X, selectorExpr.Sel.Name)
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
	} else if ident, ok := result.(*ast.Ident); ok {
		// A config declared as a package-level string, often a base config
		// shared by several config funcs.
		if configVar, ok := varDecls[ident.Name]; ok {
			return strconv.Unquote(configVar.Value)
		}
		return "", fmt.Errorf("failed to find variable declaration %s", ident.Name)
	} else if callExpr, ok := result.(*ast.CallExpr); ok {
		return readConfigCallExpr(callExpr, funcDecls, varDecls)
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

// A substitution: an acctest.Nprintf key or a fmt verb. An Nprintf key may
// also follow a percent sign escaped in a fmt.Sprintf template, as in
// acctest.Nprintf(fmt.Sprintf("%%{key}"), ...), which formats to "%{key}".
var subPattern = regexp.MustCompile("(?:%|" + escapedPercent + "){[^{}]*}|%[vTtbcspqxXUeEfFgGdo]")

// A substitution alone on its line injects whole HCL statements rather than a
// value: most often a shared "setup" block, sometimes a conditionally present
// attribute or nested block. Substituting a value there leaves a bare token
// where HCL requires an argument or block definition.
var wholeLineSubPattern = regexp.MustCompile("(?m)^[ \t]*(?:" + subPattern.String() + ")[ \t]*$")

// A substitution, or a percent sign escaped as "%%" in a fmt.Sprintf template.
// Matching both in a single pass keeps an escaped percent sign followed by a
// verb letter, such as "%%d", from being read as a percent sign followed by
// the verb "%d".
var sprintfEscapedPercentPattern = regexp.MustCompile("%%|" + subPattern.String())

// Stands in for a percent sign escaped in a fmt.Sprintf template until the
// substitutions have been removed, so it isn't read as part of one. It's a
// private use character, so it doesn't occur in configs.
const escapedPercent = "\uE000"

// Mark the escaped percent signs in a fmt.Sprintf template, such as the one in
// HCL's modulo operator written as "%%". readConfigStr turns them back into
// percent signs.
func markEscapedPercents(template string) string {
	return sprintfEscapedPercentPattern.ReplaceAllStringFunc(template, func(match string) string {
		if match == "%%" {
			return escapedPercent
		}
		return match
	})
}

// Read the config string and return a test step.
func readConfigStr(configStr string) (Step, error) {
	// Remove fmt substitutions because they interfere with hcl parsing.
	// Drop whole-line substitutions entirely; whatever they inject is not
	// counted as covered, but the rest of the config still parses instead of
	// the whole test step being discarded.
	configStr = wholeLineSubPattern.ReplaceAllString(configStr, "")
	// Replace with a value that can be parsed outside quotation marks.
	configStr = subPattern.ReplaceAllString(configStr, "true")
	configStr = strings.ReplaceAll(configStr, escapedPercent, "%")
	parser := hclparse.NewParser()
	file, diagnostics := parser.ParseHCL([]byte(configStr), "config.hcl")
	if diagnostics.HasErrors() {
		return nil, fmt.Errorf("errors parsing hcl: %v", diagnostics.Errs())
	}
	// Iterate over the top-level blocks directly rather than validating the
	// body against a schema of known block types. Only blocks with a type and
	// a name label (resource, data, ephemeral, list, ...) are of interest, and
	// every other block is skipped below. Enumerating block types here would
	// mean a config using any block Terraform adds later fails to parse, which
	// silently drops the whole test step.
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("couldn't get hclsyntax body from %v", file.Body)
	}
	step := make(Step)
	errs := make([]error, 0)
	for _, block := range body.Blocks {
		if len(block.Labels) != 2 {
			// Not a type-and-name block: terraform, locals, variable, check,
			// output, provider, module and moved blocks all have zero or one
			// label and declare nothing whose attributes could be under test.
			continue
		}
		blockType := BlockType(block.Type)
		typeLabel, nameLabel := block.Labels[0], block.Labels[1]
		if _, ok := step[blockType]; !ok {
			// Create an empty map for this block type.
			step[blockType] = make(TypeLabels)
		}
		if _, ok := step[blockType][typeLabel]; !ok {
			// Create an empty map for this type label.
			step[blockType][typeLabel] = make(Blocks)
		}
		blockConfig, err := readHCLBlockBody(block.Body, file.Bytes)
		if err != nil {
			errs = append(errs, err)
		}
		blockConfig = flattenBlock(blockConfig, "")
		step[blockType][typeLabel][nameLabel] = blockConfig
	}
	if len(errs) > 0 {
		return step, fmt.Errorf("errors reading hcl blocks: %v", errs)
	}
	return step, nil
}

func readHCLBlockBody(body hcl.Body, fileBytes []byte) (Block, error) {
	var m Block
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
			// Merge the attributes from the current block into the existing one.
			if existingBlock, ok := existing.(Block); ok {
				mergeBlocks(existingBlock, blockConfig)
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
func mergeBlocks(a, b Block) {
	for k, bv := range b {
		if av, ok := a[k]; ok {
			if avr, ok := av.(Block); ok {
				if bvr, ok := bv.(Block); ok {
					mergeBlocks(avr, bvr)
				}
			}
		} else {
			a[k] = bv
		}
	}
}

// Flatten nested blocks into dotted attribute paths, e.g. a boot_disk block
// containing an initialize_params block containing an image attribute becomes
// "boot_disk.initialize_params.image".
func flattenBlock(b Block, parent string) Block {
	flattened := make(Block)

	if parent != "" {
		parent += "."
	}

	for fieldName, fieldValue := range b {
		key := parent + fieldName
		nestedObject, _ := fieldValue.(Block)
		if nestedObject != nil {
			for childKey, childField := range flattenBlock(nestedObject, key) {
				flattened[childKey] = childField
			}
		} else {
			flattened[key] = fieldValue
		}
	}

	return flattened
}

type evalScope struct {
	vars    map[string]string
	maps    map[string]map[string]string
	pkgVars map[string]*ast.BasicLit
}

func newEvalScope(pkgVars map[string]*ast.BasicLit) *evalScope {
	return &evalScope{
		vars:    make(map[string]string),
		maps:    make(map[string]map[string]string),
		pkgVars: pkgVars,
	}
}

func cloneStringMap(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func stringLitValue(expr ast.Expr) (string, bool) {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		s, err := strconv.Unquote(lit.Value)
		return s, err == nil
	}
	return "", false
}

func isStringKeyMapLit(compLit *ast.CompositeLit) bool {
	mapType, ok := compLit.Type.(*ast.MapType)
	if !ok {
		return false
	}
	keyIdent, ok := mapType.Key.(*ast.Ident)
	return ok && keyIdent.Name == "string" && !isSerialTestMap(compLit)
}

func evalMapExpr(expr ast.Expr, scope *evalScope) map[string]string {
	if scope == nil || expr == nil {
		return nil
	}
	if ident, ok := expr.(*ast.Ident); ok {
		if m, exists := scope.maps[ident.Name]; exists {
			return m
		}
		return nil
	}
	if compLit, ok := expr.(*ast.CompositeLit); ok && isStringKeyMapLit(compLit) {
		out := make(map[string]string, len(compLit.Elts))
		for _, elt := range compLit.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				if key, ok := stringLitValue(kv.Key); ok {
					out[key] = evalSymbolicValue(kv.Value, scope)
				}
			}
		}
		return out
	}
	return nil
}

func evalSymbolicValue(expr ast.Expr, scope *evalScope) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			if s, err := strconv.Unquote(e.Value); err == nil {
				return "str:" + s
			}
		}
		return "lit:" + e.Value
	case *ast.Ident:
		if e.Name == "true" || e.Name == "false" || e.Name == "nil" {
			return "kw:" + e.Name
		}
		if scope != nil {
			if v, ok := scope.vars[e.Name]; ok {
				return v
			}
			if pkgLit, ok := scope.pkgVars[e.Name]; ok {
				return evalSymbolicValue(pkgLit, scope)
			}
		}
		return "ident:" + e.Name
	case *ast.IndexExpr:
		if m := evalMapExpr(e.X, scope); m != nil {
			if key, ok := stringLitValue(e.Index); ok {
				if v, exists := m[key]; exists {
					return v
				}
			}
		}
		return fmt.Sprintf("index@%d", int(e.Pos()))
	case *ast.UnaryExpr:
		return e.Op.String() + evalSymbolicValue(e.X, scope)
	case *ast.BinaryExpr:
		return evalSymbolicValue(e.X, scope) + e.Op.String() + evalSymbolicValue(e.Y, scope)
	default:
		return fmt.Sprintf("dynamic@%d", int(expr.Pos()))
	}
}

func recordAssignInScope(assignStmt *ast.AssignStmt, scope *evalScope) {
	if scope == nil || len(assignStmt.Lhs) != len(assignStmt.Rhs) {
		return
	}
	for i, lhs := range assignStmt.Lhs {
		rhs := assignStmt.Rhs[i]
		if ident, ok := lhs.(*ast.Ident); ok {
			if m := evalMapExpr(rhs, scope); m != nil {
				scope.maps[ident.Name] = cloneStringMap(m)
			}
			scope.vars[ident.Name] = evalSymbolicValue(rhs, scope)
		} else if idxExpr, ok := lhs.(*ast.IndexExpr); ok {
			if mapIdent, ok := idxExpr.X.(*ast.Ident); ok {
				if key, ok := stringLitValue(idxExpr.Index); ok {
					if existing, exists := scope.maps[mapIdent.Name]; exists {
						cloned := cloneStringMap(existing)
						cloned[key] = evalSymbolicValue(rhs, scope)
						scope.maps[mapIdent.Name] = cloned
					}
				}
			}
		}
	}
}

func hclSubToken(symbolicVal string) string {
	sum := sha256.Sum256([]byte(symbolicVal))
	return "__mm_sub_" + hex.EncodeToString(sum[:6]) + "__"
}

func isWholeLineMatch(s string, start, end int) bool {
	lineStart := strings.LastIndexByte(s[:start], '\n') + 1
	lineEndOffset := strings.IndexByte(s[end:], '\n')
	lineEnd := len(s)
	if lineEndOffset != -1 {
		lineEnd = end + lineEndOffset
	}
	return strings.TrimLeft(s[lineStart:start], " \t") == "" &&
		strings.TrimRight(s[end:lineEnd], " \t\r") == ""
}

var (
	nprintfKeyPattern  = regexp.MustCompile("(?:%|" + escapedPercent + `)\{([^{}]+)\}`)
	sprintfVerbPattern = regexp.MustCompile("%[vTtbcspqxXUeEfFgGdo]")
)

func readResolvedConfigExpr(expr ast.Expr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, scope *evalScope) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s, err := strconv.Unquote(e.Value)
			return s, err == nil
		}
	case *ast.Ident:
		if scope != nil {
			if v, ok := scope.vars[e.Name]; ok && strings.HasPrefix(v, "str:") {
				return strings.TrimPrefix(v, "str:"), true
			}
		}
		if configVar, ok := varDecls[e.Name]; ok {
			s, err := strconv.Unquote(configVar.Value)
			return s, err == nil
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			left, okL := readResolvedConfigExpr(e.X, funcDecls, varDecls, scope)
			right, okR := readResolvedConfigExpr(e.Y, funcDecls, varDecls, scope)
			if okL && okR {
				return left + right, true
			}
		}
	case *ast.CallExpr:
		return readResolvedConfigCallExpr(e, funcDecls, varDecls, scope)
	}
	return "", false
}

func readResolvedConfigCallExpr(callExpr *ast.CallExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, scope *evalScope) (string, bool) {
	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
		configFunc, exists := funcDecls[ident.Name]
		if !exists || configFunc.Body == nil {
			return "", false
		}
		childScope := newEvalScope(varDecls)
		if configFunc.Type != nil && configFunc.Type.Params != nil {
			argIdx := 0
			for _, field := range configFunc.Type.Params.List {
				for _, paramIdent := range field.Names {
					if argIdx < len(callExpr.Args) {
						arg := callExpr.Args[argIdx]
						if m := evalMapExpr(arg, scope); m != nil {
							childScope.maps[paramIdent.Name] = cloneStringMap(m)
						}
						childScope.vars[paramIdent.Name] = evalSymbolicValue(arg, scope)
					}
					argIdx++
				}
			}
		}
		for _, stmt := range configFunc.Body.List {
			if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
				recordAssignInScope(assignStmt, childScope)
			} else if returnStmt, ok := stmt.(*ast.ReturnStmt); ok && len(returnStmt.Results) > 0 {
				return readResolvedConfigExpr(returnStmt.Results[0], funcDecls, varDecls, childScope)
			}
		}
		return "", false
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	switch selExpr.Sel.Name {
	case "Nprintf":
		if len(callExpr.Args) == 0 {
			return "", false
		}
		tmplStr, ok := readResolvedConfigExpr(callExpr.Args[0], funcDecls, varDecls, scope)
		if !ok {
			return "", false
		}
		mergedMap := make(map[string]string)
		for _, arg := range callExpr.Args[1:] {
			for k, v := range evalMapExpr(arg, scope) {
				mergedMap[k] = v
			}
		}
		if len(mergedMap) == 0 {
			return tmplStr, true
		}
		matches := nprintfKeyPattern.FindAllStringSubmatchIndex(tmplStr, -1)
		if len(matches) == 0 {
			return tmplStr, true
		}
		var out strings.Builder
		last := 0
		for _, loc := range matches {
			start, end := loc[0], loc[1]
			key := tmplStr[loc[2]:loc[3]]
			out.WriteString(tmplStr[last:start])
			if val, exists := mergedMap[key]; exists && !isWholeLineMatch(tmplStr, start, end) {
				out.WriteString(hclSubToken(val))
			} else {
				out.WriteString(tmplStr[start:end])
			}
			last = end
		}
		out.WriteString(tmplStr[last:])
		return out.String(), true

	case "Sprintf":
		if len(callExpr.Args) == 0 {
			return "", false
		}
		tmplStr, ok := readResolvedConfigExpr(callExpr.Args[0], funcDecls, varDecls, scope)
		if !ok {
			return "", false
		}
		tmplStr = markEscapedPercents(tmplStr)
		matches := sprintfVerbPattern.FindAllStringIndex(tmplStr, -1)
		if len(matches) == 0 {
			return tmplStr, true
		}
		var out strings.Builder
		last := 0
		for argOffset, loc := range matches {
			start, end := loc[0], loc[1]
			out.WriteString(tmplStr[last:start])
			if argOffset+1 < len(callExpr.Args) && !isWholeLineMatch(tmplStr, start, end) {
				sym := evalSymbolicValue(callExpr.Args[argOffset+1], scope)
				out.WriteString(hclSubToken(sym))
			} else {
				out.WriteString(tmplStr[start:end])
			}
			last = end
		}
		out.WriteString(tmplStr[last:])
		return out.String(), true

	case "Sprint":
		var out strings.Builder
		for _, arg := range callExpr.Args {
			part, ok := readResolvedConfigExpr(arg, funcDecls, varDecls, scope)
			if !ok {
				return "", false
			}
			out.WriteString(part)
		}
		return out.String(), true
	}
	return "", false
}
