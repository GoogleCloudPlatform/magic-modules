package gotemplate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// helper to create a temp .go.tmpl file and return its path
func createTestFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go.tmpl")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	return path
}

// =============================================================================
// TESTS: Catches invalid functions
// =============================================================================

func TestFuncCheck_CatchesBareUppercaseIdentifier(t *testing.T) {
	path := createTestFile(t, `{{BigQueryBasePath}}`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one error, got none")
	}
	if !strings.Contains(results[0], "BigQueryBasePath") {
		t.Errorf("expected error to mention BigQueryBasePath, got: %s", results[0])
	}
}

func TestFuncCheck_CatchesMultipleInvalidFunctions(t *testing.T) {
	path := createTestFile(t, `{{BigQueryBasePath}}
{{ComputeBasePath}}
{{DNSBasePath}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 errors, got %d: %v", len(results), results)
	}
	if !strings.Contains(results[0], "BigQueryBasePath") {
		t.Errorf("expected first error to mention BigQueryBasePath, got: %s", results[0])
	}
	if !strings.Contains(results[1], "ComputeBasePath") {
		t.Errorf("expected second error to mention ComputeBasePath, got: %s", results[1])
	}
	if !strings.Contains(results[2], "DNSBasePath") {
		t.Errorf("expected third error to mention DNSBasePath, got: %s", results[2])
	}
}

func TestFuncCheck_CatchesTypoInFunctionName(t *testing.T) {
	path := createTestFile(t, `{{camalize .Name}}`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected error for typo 'camalize', got none")
	}
	if !strings.Contains(results[0], "camalize") {
		t.Errorf("expected error to mention 'camalize', got: %s", results[0])
	}
}

func TestFuncCheck_CatchesPipedInvalidFunction(t *testing.T) {
	path := createTestFile(t, `{{.Name | badFunc}}`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected error for piped 'badFunc', got none")
	}
	if !strings.Contains(results[0], "badFunc") {
		t.Errorf("expected error to mention 'badFunc', got: %s", results[0])
	}
}

func TestFuncCheck_CatchesInvalidFunctionWithArgs(t *testing.T) {
	path := createTestFile(t, `{{notAFunc .Arg1 .Arg2}}`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected error for 'notAFunc', got none")
	}
	if !strings.Contains(results[0], "notAFunc") {
		t.Errorf("expected error to mention 'notAFunc', got: %s", results[0])
	}
}

func TestFuncCheck_CatchesInvalidFunctionWithTrimMarkers(t *testing.T) {
	path := createTestFile(t, `{{- invalidFunc .X -}}`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected error for 'invalidFunc', got none")
	}
	if !strings.Contains(results[0], "invalidFunc") {
		t.Errorf("expected error to mention 'invalidFunc', got: %s", results[0])
	}
}

func TestFuncCheck_CatchesInvalidFunctionInMultiActionLine(t *testing.T) {
	path := createTestFile(t, `{{if .X}}{{badFunc .Y}}{{end}}`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected error for 'badFunc', got none")
	}
	if !strings.Contains(results[0], "badFunc") {
		t.Errorf("expected error to mention 'badFunc', got: %s", results[0])
	}
}

func TestFuncCheck_CatchesMixedValidAndInvalid(t *testing.T) {
	path := createTestFile(t, `{{- if ne $.TargetVersionName "ga" }}
{{camelize .Name}}
{{BigQueryBasePath}}
{{lower .S}}
{{StorageBasePath}}
{{- end }}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(results), results)
	}
	foundBQ, foundStorage := false, false
	for _, r := range results {
		if strings.Contains(r, "BigQueryBasePath") {
			foundBQ = true
		}
		if strings.Contains(r, "StorageBasePath") {
			foundStorage = true
		}
	}
	if !foundBQ {
		t.Error("expected error to mention BigQueryBasePath")
	}
	if !foundStorage {
		t.Error("expected error to mention StorageBasePath")
	}
}

// =============================================================================
// TESTS: Valid constructs that should NOT be flagged
// =============================================================================

func TestFuncCheck_PassesValidMmv1Functions(t *testing.T) {
	path := createTestFile(t, `{{title .Name}}
{{replace .S "old" "new" 1}}
{{replaceAll .S "old" "new"}}
{{camelize .Name}}
{{underscore .Name}}
{{plural .Name}}
{{contains .S "sub"}}
{{join .List ","}}
{{lower .Name}}
{{upper .Name}}
{{hasSuffix .S "suffix"}}
{{dict "key" "value"}}
{{format2regex .Fmt}}
{{hasPrefix .S "prefix"}}
{{sub .A .B}}
{{plus .A .B}}
{{firstSentence .Text}}
{{trimTemplate "path" .}}
{{customTemplate . "path" true}}
{{TemplatePath}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for valid mmv1 functions, got: %v", results)
	}
}

func TestFuncCheck_PassesGoBuiltinFunctions(t *testing.T) {
	path := createTestFile(t, `{{if eq .A .B}}equal{{end}}
{{if ne .A .B}}not equal{{end}}
{{if lt .A .B}}less{{end}}
{{if le .A .B}}less or equal{{end}}
{{if gt .A .B}}greater{{end}}
{{if ge .A .B}}greater or equal{{end}}
{{and .A .B}}
{{or .A .B}}
{{not .A}}
{{len .List}}
{{index .Map "key"}}
{{slice .List 1 3}}
{{print .A}}
{{printf "%s" .A}}
{{println .A}}
{{call .Func}}
{{html .S}}
{{js .S}}
{{urlquery .S}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for Go built-in functions, got: %v", results)
	}
}

func TestFuncCheck_PassesTemplateKeywords(t *testing.T) {
	path := createTestFile(t, `{{if .Condition}}
{{else if .Other}}
{{else}}
{{end}}
{{range .Items}}
{{end}}
{{with .Context}}
{{end}}
{{block "name" .}}{{end}}
{{define "helper"}}content{{end}}
{{template "helper" .}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for template keywords, got: %v", results)
	}
}

func TestFuncCheck_PassesDotAccess(t *testing.T) {
	path := createTestFile(t, `{{.Name}}
{{.Resource.Name}}
{{$.TargetVersionName}}
{{$.Vars}}
{{$.PrimaryResourceId}}
{{.Resource.Properties}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for dot access, got: %v", results)
	}
}

func TestFuncCheck_PassesVariableAssignment(t *testing.T) {
	path := createTestFile(t, `{{$name := .Resource.Name}}
{{$version := $.TargetVersionName}}
{{$items := .List}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for variable assignments, got: %v", results)
	}
}

func TestFuncCheck_PassesPipedValidFunctions(t *testing.T) {
	path := createTestFile(t, `{{.Name | camelize}}
{{.Name | underscore | upper}}
{{.Name | lower | title}}
{{.Items | len}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for piped valid functions, got: %v", results)
	}
}

func TestFuncCheck_PassesBareInterpolationVariables(t *testing.T) {
	path := createTestFile(t, `{{project}}
{{location}}
{{region}}
{{zone}}
{{name}}
{{anywhere_cache_id}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for bare interpolation variables, got: %v", results)
	}
}

func TestFuncCheck_PassesBooleanLiterals(t *testing.T) {
	path := createTestFile(t, `{{if true}}yes{{end}}
{{if false}}no{{end}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for boolean literals, got: %v", results)
	}
}

func TestFuncCheck_PassesNil(t *testing.T) {
	path := createTestFile(t, `{{if eq .Value nil}}is nil{{end}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for nil keyword, got: %v", results)
	}
}

func TestFuncCheck_PassesStringLiteralContainingIdentifier(t *testing.T) {
	path := createTestFile(t, `{{join ($.PropertyNamesToStrings (index $CustomUpdateProps $group)) "\") || d.HasChange(\""}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors when identifier is inside string literal, got: %v", results)
	}
}

func TestFuncCheck_PassesVersionGuard(t *testing.T) {
	path := createTestFile(t, `{{- if ne $.TargetVersionName "ga" }}
// beta-only code
{{- end }}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for version guard, got: %v", results)
	}
}

// =============================================================================
// TESTS: Line number reporting
// =============================================================================

func TestFuncCheck_ReportsCorrectLineNumbers(t *testing.T) {
	path := createTestFile(t, `line 1
line 2
{{badFunc .X}}
line 4
{{anotherBad .Y}}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(results), results)
	}
	if !strings.Contains(results[0], "line 3") {
		t.Errorf("expected first error on line 3, got: %s", results[0])
	}
	if !strings.Contains(results[1], "line 5") {
		t.Errorf("expected second error on line 5, got: %s", results[1])
	}
}

// =============================================================================
// TESTS: Edge cases
// =============================================================================

func TestFuncCheck_EmptyFile(t *testing.T) {
	path := createTestFile(t, ``)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for empty file, got: %v", results)
	}
}

func TestFuncCheck_FileWithNoTemplateActions(t *testing.T) {
	path := createTestFile(t, `package main
import "fmt"
func hello() {
	fmt.Println("no template actions here")
}
`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no errors for file with no actions, got: %v", results)
	}
}

func TestFuncCheck_NonexistentFile(t *testing.T) {
	_, err := CheckInvalidFuncsForFile("/nonexistent/path/test.go.tmpl")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

// =============================================================================
// TESTS: ValidFuncs map completeness
// =============================================================================

func TestValidFuncs_ContainsGoBuiltins(t *testing.T) {
	builtins := []string{
		"and", "call", "eq", "ge", "gt", "html", "index",
		"js", "le", "len", "lt", "ne", "not", "or",
		"print", "printf", "println", "slice", "urlquery",
	}
	for _, fn := range builtins {
		if !ValidFuncs[fn] {
			t.Errorf("Go built-in function %q missing from ValidFuncs", fn)
		}
	}
}

func TestValidFuncs_ContainsTemplateKeywords(t *testing.T) {
	keywords := []string{
		"if", "else", "end", "range", "with", "block", "define", "template", "nil",
	}
	for _, kw := range keywords {
		if !ValidFuncs[kw] {
			t.Errorf("Template keyword %q missing from ValidFuncs", kw)
		}
	}
}

func TestValidFuncs_ContainsMmv1Functions(t *testing.T) {
	mmv1Funcs := []string{
		"title", "replace", "replaceAll", "camelize", "underscore", "plural",
		"contains", "join", "lower", "upper", "hasSuffix", "dict",
		"format2regex", "hasPrefix", "sub", "plus", "firstSentence",
		"trimTemplate", "customTemplate", "TemplatePath",
	}
	for _, fn := range mmv1Funcs {
		if !ValidFuncs[fn] {
			t.Errorf("mmv1 registered function %q missing from ValidFuncs", fn)
		}
	}
}

func TestValidFuncs_DoesNotContainProductBasePaths(t *testing.T) {
	basePaths := []string{
		"BigQueryBasePath", "ComputeBasePath", "DNSBasePath",
		"StorageBasePath", "ContainerBasePath", "PubsubBasePath",
	}
	for _, fn := range basePaths {
		if ValidFuncs[fn] {
			t.Errorf("Product base path %q should NOT be in ValidFuncs", fn)
		}
	}
}

// =============================================================================
// TESTS: Real magic-modules templates (if available)
// =============================================================================

func TestFuncCheck_RealTemplates_ZeroErrors(t *testing.T) {
	repoRoot := "../../../mmv1/third_party/terraform/"
	if _, err := os.Stat(repoRoot); os.IsNotExist(err) {
		t.Skip("magic-modules templates not found — skipping real repo test")
	}
	var tmplFiles []string
	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go.tmpl") {
			tmplFiles = append(tmplFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk templates: %v", err)
	}
	if len(tmplFiles) == 0 {
		t.Skip("no .go.tmpl files found")
	}
	var allErrors []string
	for _, f := range tmplFiles {
		results, err := CheckInvalidFuncsForFile(f)
		if err != nil {
			t.Fatalf("error checking %s: %v", f, err)
		}
		for _, r := range results {
			allErrors = append(allErrors, f+": "+r)
		}
	}
	if len(allErrors) > 0 {
		t.Errorf("expected zero errors on real templates, got %d:\n%s", len(allErrors), strings.Join(allErrors, "\n"))
	}
}

func TestFuncCheck_RealMmv1Templates_ZeroErrors(t *testing.T) {
	repoRoot := "../../../mmv1/templates/terraform/"
	if _, err := os.Stat(repoRoot); os.IsNotExist(err) {
		t.Skip("magic-modules mmv1 templates not found — skipping real repo test")
	}
	var tmplFiles []string
	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go.tmpl") {
			tmplFiles = append(tmplFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk templates: %v", err)
	}
	if len(tmplFiles) == 0 {
		t.Skip("no .go.tmpl files found")
	}
	var allErrors []string
	for _, f := range tmplFiles {
		results, err := CheckInvalidFuncsForFile(f)
		if err != nil {
			t.Fatalf("error checking %s: %v", f, err)
		}
		for _, r := range results {
			allErrors = append(allErrors, f+": "+r)
		}
	}
	if len(allErrors) > 0 {
		t.Errorf("expected zero errors on real mmv1 templates, got %d:\n%s", len(allErrors), strings.Join(allErrors, "\n"))
	}
}
