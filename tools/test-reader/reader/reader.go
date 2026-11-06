package reader

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
	Name  string
	Steps []Step
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
	vars := make(map[string]*ast.CompositeLit, len(testFunc.Body.List))              // map of variable names to composite literal values in function body
	contextVars := make(map[string]map[string]contextValue, len(testFunc.Body.List)) // map of variable names to context map values
	for _, stmt := range testFunc.Body.List {
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				// This is a call expression.
				ident, isIdent := callExpr.Fun.(*ast.Ident)
				selExpr, isSelExpr := callExpr.Fun.(*ast.SelectorExpr)
				if isIdent && ident.Name == "VcrTest" || isSelExpr && selExpr.Sel.Name == "VcrTest" {
					test, err := readVcrTestCallWithContext(callExpr, funcDecls, varDecls, contextVars)
					if err != nil {
						errs = append(errs, err)
					}
					test.Name = testFunc.Name.Name
					tests = append(tests, test)
				}
			}
		} else if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
			if len(assignStmt.Lhs) == 1 && len(assignStmt.Rhs) == 1 {
				// For now, only allow single assignment variables for serial test maps
				// and local context maps.
				if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
					if rhsCompLit, ok := assignStmt.Rhs[0].(*ast.CompositeLit); ok {
						if isSerialTestMap(rhsCompLit) {
							vars[ident.Name] = rhsCompLit
						} else if ctxMap := parseContextMapCompLit(rhsCompLit); len(ctxMap) > 0 {
							contextVars[ident.Name] = ctxMap
						}
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

type contextValue struct {
	raw      string
	isString bool
}

// Extract key-value pairs from a composite literal of type map[string]interface{} or map[string]any.
func parseContextMapCompLit(compLit *ast.CompositeLit) map[string]contextValue {
	mapType, ok := compLit.Type.(*ast.MapType)
	if !ok {
		return nil
	}
	keyIdent, ok := mapType.Key.(*ast.Ident)
	if !ok || keyIdent.Name != "string" {
		return nil
	}
	if _, isFunc := mapType.Value.(*ast.FuncType); isFunc {
		return nil
	}
	result := make(map[string]contextValue)
	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		keyLit, ok := kv.Key.(*ast.BasicLit)
		if !ok || keyLit.Kind != token.STRING {
			continue
		}
		keyStr, err := strconv.Unquote(keyLit.Value)
		if err != nil {
			continue
		}
		if val, ok := extractLiteralContextValue(kv.Value); ok {
			result[keyStr] = val
		}
	}
	return result
}

func extractLiteralContextValue(expr ast.Expr) (contextValue, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		switch v.Kind {
		case token.STRING:
			unquoted, err := strconv.Unquote(v.Value)
			if err != nil {
				return contextValue{}, false
			}
			return contextValue{raw: unquoted, isString: true}, true
		case token.INT, token.FLOAT:
			return contextValue{raw: v.Value, isString: false}, true
		}
	case *ast.Ident:
		if v.Name == "true" || v.Name == "false" {
			return contextValue{raw: v.Name, isString: false}, true
		}
	case *ast.UnaryExpr:
		if v.Op == token.SUB {
			if basicLit, ok := v.X.(*ast.BasicLit); ok && (basicLit.Kind == token.INT || basicLit.Kind == token.FLOAT) {
				return contextValue{raw: "-" + basicLit.Value, isString: false}, true
			}
		}
	}
	return contextValue{}, false
}

func extractContextMapFromExpr(expr ast.Expr, contextVars map[string]map[string]contextValue) map[string]contextValue {
	switch v := expr.(type) {
	case *ast.Ident:
		if contextVars != nil {
			return contextVars[v.Name]
		}
	case *ast.CompositeLit:
		return parseContextMapCompLit(v)
	}
	return nil
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
	return readVcrTestCallWithContext(vcrTestCall, funcDecls, varDecls, nil)
}

func readVcrTestCallWithContext(vcrTestCall *ast.CallExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, contextVars map[string]map[string]contextValue) (*Test, error) {
	for _, arg := range vcrTestCall.Args {
		if vcrTestArgCompLit, ok := arg.(*ast.CompositeLit); ok {
			if selExpr, ok := vcrTestArgCompLit.Type.(*ast.SelectorExpr); ok {
				if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "resource" && selExpr.Sel.Name == "TestCase" {
					return readTestCaseCompLitWithContext(vcrTestArgCompLit, funcDecls, varDecls, contextVars)
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to find TestCase in %v", vcrTestCall.Args)
}

func readTestCaseCompLit(testCaseCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (*Test, error) {
	return readTestCaseCompLitWithContext(testCaseCompLit, funcDecls, varDecls, nil)
}

func readTestCaseCompLitWithContext(testCaseCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, contextVars map[string]map[string]contextValue) (*Test, error) {
	for _, elt := range testCaseCompLit.Elts {
		if keyValueExpr, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := keyValueExpr.Key.(*ast.Ident); ok && ident.Name == "Steps" {
				if stepsCompLit, ok := keyValueExpr.Value.(*ast.CompositeLit); ok {
					return readStepsCompLitWithContext(stepsCompLit, funcDecls, varDecls, contextVars)
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to find Steps in %v", testCaseCompLit.Elts)
}

func readStepsCompLit(stepsCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (*Test, error) {
	return readStepsCompLitWithContext(stepsCompLit, funcDecls, varDecls, nil)
}

func readStepsCompLitWithContext(stepsCompLit *ast.CompositeLit, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, contextVars map[string]map[string]contextValue) (*Test, error) {
	test := &Test{}
	errs := make([]error, 0)
	for _, elt := range stepsCompLit.Elts {
		if eltCompLit, ok := elt.(*ast.CompositeLit); ok {
			if shouldSkipStep(eltCompLit) {
				continue
			}
			for _, eltCompLitElt := range eltCompLit.Elts {
				if keyValueExpr, ok := eltCompLitElt.(*ast.KeyValueExpr); ok {
					if ident, ok := keyValueExpr.Key.(*ast.Ident); ok && ident.Name == "Config" {
						// A config is any expression that evaluates to a config
						// string, and reads the same way here as it does in a
						// config func's return statement.
						configStr, err := readConfigFuncResultWithContext(keyValueExpr.Value, funcDecls, varDecls, contextVars)
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
	return readConfigCallExprWithContext(configCallExpr, funcDecls, varDecls, nil)
}

func readConfigCallExprWithContext(configCallExpr *ast.CallExpr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, contextVars map[string]map[string]contextValue) (string, error) {
	// Dispatch on the callee rather than on the shape of the arguments: a
	// config func can take a string as its first argument, and that string is a
	// value to interpolate, not a config.
	if ident, ok := configCallExpr.Fun.(*ast.Ident); ok {
		if configFunc, ok := funcDecls[ident.Name]; ok {
			return readConfigFuncWithContext(configFunc, configCallExpr.Args, funcDecls, varDecls, contextVars)
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
		configStr, err := readConfigFuncResultWithContext(configCallExpr.Args[0], funcDecls, varDecls, contextVars)
		if err != nil {
			return "", err
		}
		if selectorExpr.Sel.Name == "Sprintf" {
			return markEscapedPercents(configStr), nil
		}
		// Nprintf only replaces %{key}, so "%%" is not an escape there.
		if len(configCallExpr.Args) >= 2 {
			if ctxMap := extractContextMapFromExpr(configCallExpr.Args[1], contextVars); len(ctxMap) > 0 {
				configStr = applyContextMap(configStr, ctxMap)
			}
		}
		return configStr, nil
	case "Sprint":
		// Every argument is part of the config. This is used to join a shared
		// setup config to the config under test.
		var configStr strings.Builder
		for _, arg := range configCallExpr.Args {
			argConfigStr, err := readConfigFuncResultWithContext(arg, funcDecls, varDecls, contextVars)
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
	return readConfigFuncWithContext(configFunc, nil, funcDecls, varDecls, nil)
}

func readConfigFuncWithContext(configFunc *ast.FuncDecl, callArgs []ast.Expr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, callerContextVars map[string]map[string]contextValue) (string, error) {
	calleeContextVars := make(map[string]map[string]contextValue)
	for k, v := range callerContextVars {
		calleeContextVars[k] = v
	}
	if configFunc.Type != nil && configFunc.Type.Params != nil {
		argIdx := 0
		for _, field := range configFunc.Type.Params.List {
			for _, nameIdent := range field.Names {
				if argIdx < len(callArgs) {
					if ctxMap := extractContextMapFromExpr(callArgs[argIdx], callerContextVars); len(ctxMap) > 0 {
						calleeContextVars[nameIdent.Name] = ctxMap
					}
				}
				argIdx++
			}
		}
	}
	for _, stmt := range configFunc.Body.List {
		if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
			if len(assignStmt.Lhs) == 1 && len(assignStmt.Rhs) == 1 {
				if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
					if rhsCompLit, ok := assignStmt.Rhs[0].(*ast.CompositeLit); ok {
						if ctxMap := parseContextMapCompLit(rhsCompLit); len(ctxMap) > 0 {
							calleeContextVars[ident.Name] = ctxMap
						}
					}
				}
			}
		} else if returnStmt, ok := stmt.(*ast.ReturnStmt); ok {
			if len(returnStmt.Results) > 0 {
				return readConfigFuncResultWithContext(returnStmt.Results[0], funcDecls, varDecls, calleeContextVars)
			}
			return "", fmt.Errorf("failed to find a config string in results %v", returnStmt.Results)
		}
	}
	return "", fmt.Errorf("failed to find a return statement in %v", configFunc.Body.List)
}

// Read the return result of a config func and return the config string.
func readConfigFuncResult(result ast.Expr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit) (string, error) {
	return readConfigFuncResultWithContext(result, funcDecls, varDecls, nil)
}

func readConfigFuncResultWithContext(result ast.Expr, funcDecls map[string]*ast.FuncDecl, varDecls map[string]*ast.BasicLit, contextVars map[string]map[string]contextValue) (string, error) {
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
		return readConfigCallExprWithContext(callExpr, funcDecls, varDecls, contextVars)
	} else if binaryExpr, ok := result.(*ast.BinaryExpr); ok {
		xConfigStr, err := readConfigFuncResultWithContext(binaryExpr.X, funcDecls, varDecls, contextVars)
		if err != nil {
			return "", err
		}
		yConfigStr, err := readConfigFuncResultWithContext(binaryExpr.Y, funcDecls, varDecls, contextVars)
		if err != nil {
			return "", err
		}
		return xConfigStr + yConfigStr, nil
	}
	return "", fmt.Errorf("unknown config func result %v (%T)", result, result)
}

var namedPlaceholderPattern = regexp.MustCompile("(?:%|" + escapedPercent + `)\{([^{}]+)\}`)
var nprintfPlaceholderOrPercentPattern = regexp.MustCompile(`%\{[^{}]*\}|%`)

func escapeLiteralPercents(s string) string {
	return nprintfPlaceholderOrPercentPattern.ReplaceAllStringFunc(s, func(match string) string {
		if match == "%" {
			return escapedPercent
		}
		return match
	})
}

func isValidHCLValueExpr(s string) bool {
	trimmed := strings.TrimSpace(s)
	if trimmed == "true" || trimmed == "false" {
		return true
	}
	if _, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
		return true
	}
	if _, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return true
	}
	if strings.HasPrefix(trimmed, "\"") || strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{") {
		return true
	}
	return false
}

// Substitute known %{key} placeholders from ctxMap into templateStr.
func applyContextMap(templateStr string, ctxMap map[string]contextValue) string {
	matches := namedPlaceholderPattern.FindAllStringSubmatchIndex(templateStr, -1)
	if len(matches) == 0 {
		return templateStr
	}
	var out strings.Builder
	lastEnd := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		key := templateStr[m[2]:m[3]]
		cv, ok := ctxMap[key]
		if !ok {
			out.WriteString(templateStr[lastEnd:end])
			lastEnd = end
			continue
		}
		// If a literal '%' immediately precedes %{key} in an acctest.Nprintf
		// template (e.g. "name-%%{suffix}"), protect that '%' from matching a
		// fmt verb in subPattern once the substituted value follows it.
		if start > lastEnd && templateStr[start-1] == '%' {
			out.WriteString(templateStr[lastEnd : start-1])
			out.WriteString(escapedPercent)
		} else {
			out.WriteString(templateStr[lastEnd:start])
		}
		lineStart := strings.LastIndexByte(templateStr[:start], '\n') + 1
		inQuotes := false
		for i := lineStart; i < start; i++ {
			if templateStr[i] == '"' && (i == lineStart || templateStr[i-1] != '\\') {
				inQuotes = !inQuotes
			}
		}
		escapedRaw := escapeLiteralPercents(cv.raw)
		if inQuotes {
			out.WriteString(escapedRaw)
		} else {
			trimmedBefore := strings.TrimSpace(templateStr[lineStart:start])
			if trimmedBefore == "" {
				// Whole-line placeholder (e.g. injecting or omitting a nested HCL block).
				out.WriteString(escapedRaw)
			} else if cv.raw == "" {
				out.WriteString(`""`)
			} else if !cv.isString || isValidHCLValueExpr(cv.raw) {
				out.WriteString(escapedRaw)
			} else {
				out.WriteString(strconv.Quote(escapedRaw))
			}
		}
		lastEnd = end
	}
	out.WriteString(templateStr[lastEnd:])
	return out.String()
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
		if blockType == ResourceBlock {
			blockConfig = flattenBlockWithEmptyBlocks(blockConfig, "")
		} else {
			blockConfig = flattenBlock(blockConfig, "")
		}
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
			} else if avStr, aOk := av.(string); aOk {
				if bvStr, bOk := bv.(string); bOk {
					a[k] = avStr + "\x00" + bvStr
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
	return flattenBlockInternal(b, parent, false)
}

func flattenBlockWithEmptyBlocks(b Block, parent string) Block {
	return flattenBlockInternal(b, parent, true)
}

func flattenBlockInternal(b Block, parent string, preserveEmptyBlocks bool) Block {
	flattened := make(Block)

	if parent != "" {
		parent += "."
	}

	for fieldName, fieldValue := range b {
		key := parent + fieldName
		nestedObject, _ := fieldValue.(Block)
		if nestedObject != nil {
			if preserveEmptyBlocks && len(nestedObject) == 0 {
				flattened[key] = "{}"
			} else {
				for childKey, childField := range flattenBlockInternal(nestedObject, key, preserveEmptyBlocks) {
					flattened[childKey] = childField
				}
			}
		} else {
			flattened[key] = fieldValue
		}
	}

	return flattened
}
