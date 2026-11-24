// Copyright 2025 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"maps"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template" // NOLINT

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"

	"gopkg.in/yaml.v3"
)

const nullVarSentinel = "\x00NULL\x00"

var nullVarLineRegexp = regexp.MustCompile(`(?m)^[ \t]*[a-zA-Z0-9_-]+[ \t]*=[ \t]*"?` + regexp.QuoteMeta(nullVarSentinel) + `"?[ \t]*\r?\n`)

func replaceNullVarSentinels(hcl string) string {
	hcl = nullVarLineRegexp.ReplaceAllString(hcl, "")
	hcl = strings.ReplaceAll(hcl, `"`+nullVarSentinel+`"`, "null")
	hcl = strings.ReplaceAll(hcl, nullVarSentinel, "null")
	return hcl
}

var hclResourceRegexp = regexp.MustCompile(`(?:resource|data|list|ephemeral) "(?P<resource>google_[^"]+)"`)

// Regular expressions used when rendering sample/step HCL. Compiled once at
// init as they are applied to every step of every sample, several times each.
var (
	trailingBlankLineRegexp = regexp.MustCompile(`\n\n$`)
	regionTagLineRegexp     = regexp.MustCompile(`# \[[a-zA-Z_ ]+\]\n`)
	regionTagTrailingRegexp = regexp.MustCompile(`\n# \[[a-zA-Z_ ]+\]`)
	testPlaceholderRegexp   = regexp.MustCompile(`%\{([^{}]+)\}`)
)

type Step struct {
	Name string `yaml:"name,omitempty"`

	// The path to this step's Terraform config.
	// Defaults to `templates/terraform/samples/{{service}}/{{name}}.tf.erb`
	ConfigPath string `yaml:"config_path,omitempty"`

	// ResourceIdVars is a Hash from template variable names to output variable names.
	// It is used for values that must be unique across test runs (e.g., resource names).
	// In generated tests, the value will be prefixed (e.g., "tf-test-") and have a
	// random suffix appended. In documentation, the value is used verbatim.
	ResourceIdVars map[string]string `yaml:"resource_id_vars,omitempty"`

	// Vars is a Hash from template variable names to output variable names.
	// It is used for values that should be inserted into the template as literal,
	// unmodified strings (e.g., labels, descriptions, or other non-identifier fields).
	// These values are used verbatim in both generated tests and documentation.
	// A value of `null` (or `~` / `"null"`) sets the test context variable to `nil`
	// and omits the attribute line (or renders unquoted HCL `null` when inline) in
	// documentation and OiCS configs.
	Vars StepVars `yaml:"vars,omitempty"`

	// Some variables need to hold special values during tests, and cannot
	// be inferred by Open in Cloud Shell.  For instance, org_id
	// needs to be the correct value during integration tests, or else
	// org tests cannot pass. Other examples include an existing project_id,
	// a zone, a service account name, etc.
	//
	// test_env_vars is a Hash from template variable names to one of the
	// following symbols:
	//  - PROJECT_NAME
	//  - CREDENTIALS
	//  - REGION
	//  - ORG_ID
	//  - ORG_TARGET
	//  - BILLING_ACCT
	//  - MASTER_BILLING_ACCT
	//  - SERVICE_ACCT
	//  - CUST_ID
	//  - IDENTITY_USER
	//  - CHRONICLE_ID
	//  - VMWAREENGINE_PROJECT
	//  - ORG_UNIT_ID
	// This list corresponds to the `get*FromEnv` methods in provider_test.go.
	TestEnvVars map[string]string `yaml:"test_env_vars,omitempty"`

	// Hash to provider custom override values for generating test config
	// If field my-var is set in this hash, it will replace vars[my-var] in
	// tests. i.e. if vars["network"] = "my-vpc", without override:
	//   - doc config will have `network = "my-vpc"`
	//   - tests config will have `"network = my-vpc%{random_suffix}"`
	//     with context
	//       map[string]interface{}{
	//         "random_suffix": acctest.RandString()
	//       }
	//
	// If test_vars_overrides["network"] = "nameOfVpc()"
	//   - doc config will have `network = "my-vpc"`
	//   - tests will replace with `"network = %{network}"` with context
	//       map[string]interface{}{
	//         "network": nameOfVpc
	//         ...
	//       }
	// NOTE: Keys in TestVarsOverrides will apply to a matching key in EITHER Vars or ResourceIdVars.
	TestVarsOverrides map[string]string `yaml:"test_vars_overrides,omitempty"`

	// Hash to provider custom override values for generating oics config
	// See test_vars_overrides for more details
	// NOTE: Keys in OicsVarsOverrides will apply to a matching key in EITHER Vars or ResourceIdVars for OiCS generation.
	OicsVarsOverrides map[string]string `yaml:"oics_vars_overrides,omitempty"`

	// The version name of the test step's version if it's different than the
	// test version, eg. `beta`
	MinVersion string `yaml:"min_version,omitempty"`

	// Extra properties to ignore read on during import.
	// These properties will likely be custom code.
	IgnoreReadExtra []string `yaml:"ignore_read_extra,omitempty"`

	// Whether to skip import tests for this test step
	ExcludeImportTest bool `yaml:"exclude_import_test,omitempty"`

	// Whether to generate docs for this test step (override sample's exclude_basic_doc)
	IncludeStepDoc bool `yaml:"include_step_doc,omitempty"`

	DocumentationHCLText string            `yaml:"-"`
	TestHCLText          string            `yaml:"-"`
	IamHCLText           string            `yaml:"-"`
	OicsHCLText          string            `yaml:"-"`
	PrimaryResourceId    string            `yaml:"-"`
	TestContextVars      map[string]string `yaml:"-"`
	StepSlug             string            `yaml:"-"`
	templateErrs         []error           `yaml:"-"`
	hclLoaded            bool
}

// StepVars is a map of template variable names to string values that preserves
// explicit YAML null scalars (`field: null` or `field: ~`) as `"null"` rather
// than decoding them as empty strings (`""`).
type StepVars map[string]string

func (v *StepVars) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]*string
	if err := value.Decode(&raw); err != nil {
		return err
	}
	if raw == nil {
		*v = nil
		return nil
	}
	out := make(StepVars, len(raw))
	for k, val := range raw {
		if val == nil {
			out[k] = "null"
		} else {
			out[k] = *val
		}
	}
	*v = out
	return nil
}

func (s *Step) TestStepSlug(productName, resourceName string) string {
	slug := s.StepSlug
	if slug == "" {
		slug = s.Name
	}
	return fmt.Sprintf("%s%s_%sExample", productName, resourceName, google.Camelize(slug, "lower"))
}

// TestDependencies returns a map of service names to import aliases that are required
// by this step.
func (s *Step) TestDependencies(resourcePrefixPkgMap map[string]string) map[string]string {
	deps := map[string]string{}
	for _, val := range s.TestContextVars {
		if strings.HasPrefix(val, "compute.") {
			deps["services/compute"] = ""
		}
		if strings.HasPrefix(val, "container.") {
			deps["services/container"] = ""
		}
		if strings.HasPrefix(val, "kms.") {
			deps["services/kms"] = ""
		}
		if strings.HasPrefix(val, "resourcemanager.") {
			deps["services/resourcemanager"] = ""
		}
		if strings.HasPrefix(val, "servicenetworking.") {
			deps["services/servicenetworking"] = ""
		}
		if strings.HasPrefix(val, "tags.") {
			deps["services/tags"] = ""
		}
	}
	matches := hclResourceRegexp.FindAllStringSubmatch(s.TestHCLText, -1)
	resources := map[string]struct{}{}
	for _, m := range matches {
		resources[m[1]] = struct{}{}
	}

	for r, _ := range resources {
		longestPrefix := ""
		for prefix, _ := range resourcePrefixPkgMap {
			if strings.HasPrefix(r, prefix) && len(prefix) > len(longestPrefix) {
				longestPrefix = prefix
			}
		}
		if longestPrefix != "" {
			pkg := resourcePrefixPkgMap[longestPrefix]
			if _, ok := deps[pkg]; !ok {
				deps[pkg] = "_"
			}
		}
	}
	return deps
}

func mergeStringMap(sampleMap, stepMap map[string]string) map[string]string {
	if len(sampleMap) == 0 {
		return stepMap
	}
	merged := make(map[string]string, len(sampleMap)+len(stepMap))
	maps.Copy(merged, sampleMap)
	maps.Copy(merged, stepMap)
	return merged
}

// InheritFromSample merges sample-level defaults into s without overwriting
// keys explicitly specified on s.
func (s *Step) InheritFromSample(sample *Sample) {
	if s.ConfigPath == "" && sample.ConfigPath != "" {
		s.ConfigPath = sample.ConfigPath
	}
	s.ResourceIdVars = mergeStringMap(sample.ResourceIdVars, s.ResourceIdVars)
	if len(sample.Vars) > 0 {
		merged := make(StepVars, len(sample.Vars)+len(s.Vars))
		maps.Copy(merged, sample.Vars)
		maps.Copy(merged, s.Vars)
		s.Vars = merged
	}
	s.TestEnvVars = mergeStringMap(sample.TestEnvVars, s.TestEnvVars)
	s.TestVarsOverrides = mergeStringMap(sample.TestVarsOverrides, s.TestVarsOverrides)
	s.OicsVarsOverrides = mergeStringMap(sample.OicsVarsOverrides, s.OicsVarsOverrides)
	if len(sample.IgnoreReadExtra) > 0 {
		combined := make([]string, 0, len(sample.IgnoreReadExtra)+len(s.IgnoreReadExtra))
		for _, item := range sample.IgnoreReadExtra {
			if !slices.Contains(combined, item) {
				combined = append(combined, item)
			}
		}
		for _, item := range s.IgnoreReadExtra {
			if !slices.Contains(combined, item) {
				combined = append(combined, item)
			}
		}
		s.IgnoreReadExtra = combined
	}
}

func (s *Step) usedPlaceholders() map[string]bool {
	used := make(map[string]bool)
	for _, match := range testPlaceholderRegexp.FindAllStringSubmatch(s.TestHCLText, -1) {
		used[match[1]] = true
	}
	for _, match := range testPlaceholderRegexp.FindAllStringSubmatch(s.IamHCLText, -1) {
		used[match[1]] = true
	}
	return used
}

func isQuoteWrapped(value string) bool {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < 2 {
		return false
	}
	quote := trimmed[0]
	if (quote != '"' && quote != '\'') || trimmed[len(trimmed)-1] != quote {
		return false
	}
	inner := trimmed[1 : len(trimmed)-1]
	for i := 0; i < len(inner); i++ {
		if inner[i] == '\\' {
			if i+1 >= len(inner) {
				return false
			}
			i++
			continue
		}
		if inner[i] == quote {
			return false
		}
	}
	return true
}

func (s *Step) Validate(rName, sName string) (es []error) {
	for _, k := range slices.Sorted(maps.Keys(s.Vars)) {
		if _, exists := s.ResourceIdVars[k]; exists {
			es = append(es, fmt.Errorf("variable key '%s' cannot exist in both 'vars' and 'resource_id_vars' for step '%s' in sample '%s' of resource '%s'", k, s.Name, sName, rName))
		}
	}
	if s.Name == "" && sName == "" {
		es = append(es, fmt.Errorf("missing `name` for one step in test sample %s in resource %s", sName, rName))
	}
	es = append(es, s.templateErrs...)

	for _, check := range []struct {
		field string
		vals  map[string]string
	}{
		{field: "vars", vals: s.Vars},
		{field: "resource_id_vars", vals: s.ResourceIdVars},
		{field: "oics_vars_overrides", vals: s.OicsVarsOverrides},
	} {
		for _, k := range slices.Sorted(maps.Keys(check.vals)) {
			if v := check.vals[k]; isQuoteWrapped(v) {
				es = append(es, fmt.Errorf("variable %q in '%s' has quote-wrapped value %q for step %q in sample %q of resource %q; values are substituted verbatim into HCL templates (unlike 'test_vars_overrides'), so remove the surrounding quotes", k, check.field, v, s.Name, sName, rName))
			}
		}
	}

	for _, k := range slices.Sorted(maps.Keys(s.OicsVarsOverrides)) {
		_, inVars := s.Vars[k]
		_, inResourceIdVars := s.ResourceIdVars[k]
		if !inVars && !inResourceIdVars {
			es = append(es, fmt.Errorf("unused variable %q in 'oics_vars_overrides' (not defined in 'vars' or 'resource_id_vars') for step %q in sample %q of resource %q", k, s.Name, sName, rName))
		}
	}

	// TestSampleSetUp skips steps above TargetVersionName (e.g. min_version: beta
	// during GA builds), whereas Validate runs on all steps; skip unloaded steps.
	if s.hclLoaded {
		used := s.usedPlaceholders()
		for _, k := range slices.Sorted(maps.Keys(s.Vars)) {
			if !used[k] {
				es = append(es, fmt.Errorf("unused variable %q in 'vars' for step %q in sample %q of resource %q", k, s.Name, sName, rName))
			}
		}
		for _, k := range slices.Sorted(maps.Keys(s.ResourceIdVars)) {
			if !used[k] {
				es = append(es, fmt.Errorf("unused variable %q in 'resource_id_vars' for step %q in sample %q of resource %q", k, s.Name, sName, rName))
			}
		}
		for _, k := range slices.Sorted(maps.Keys(s.TestEnvVars)) {
			if !used[k] {
				es = append(es, fmt.Errorf("unused variable %q in 'test_env_vars' for step %q in sample %q of resource %q", k, s.Name, sName, rName))
			}
		}
		for _, k := range slices.Sorted(maps.Keys(s.TestVarsOverrides)) {
			if !used[k] {
				es = append(es, fmt.Errorf("unused variable %q in 'test_vars_overrides' for step %q in sample %q of resource %q", k, s.Name, sName, rName))
			}
		}
		for _, k := range slices.Sorted(maps.Keys(s.OicsVarsOverrides)) {
			_, inVars := s.Vars[k]
			_, inResourceIdVars := s.ResourceIdVars[k]
			if (inVars || inResourceIdVars) && !used[k] {
				es = append(es, fmt.Errorf("unused variable %q in 'oics_vars_overrides' for step %q in sample %q of resource %q", k, s.Name, sName, rName))
			}
		}
	}

	return es
}

// Executes step configuration templates for documentation and tests
func (s *Step) SetHCLText(sysfs fs.FS) {
	originalResourceIdVars := s.ResourceIdVars
	originalVars := s.Vars
	originalTestEnvVars := s.TestEnvVars
	docTestEnvVars := make(map[string]string)
	docs_defaults := map[string]string{
		"PROJECT_NAME":         "my-project-name",
		"PROJECT_NUMBER":       "1111111111111",
		"CREDENTIALS":          "my/credentials/filename.json",
		"REGION":               "us-west1",
		"ORG_ID":               "123456789",
		"ORG_DOMAIN":           "example.com",
		"ORG_TARGET":           "123456789",
		"BILLING_ACCT":         "000000-0000000-0000000-000000",
		"MASTER_BILLING_ACCT":  "000000-0000000-0000000-000000",
		"SERVICE_ACCT":         "my@service-account.com",
		"CUST_ID":              "A01b123xz",
		"IDENTITY_USER":        "cloud_identity_user",
		"PAP_DESCRIPTION":      "description",
		"CHRONICLE_ID":         "00000000-0000-0000-0000-000000000000",
		"VMWAREENGINE_PROJECT": "my-vmwareengine-project",
		"ORG_UNIT_ID":          "03abcxyz",
	}

	// Apply doc defaults to test_env_vars from YAML
	for key := range s.TestEnvVars {
		docTestEnvVars[key] = docs_defaults[s.TestEnvVars[key]]
	}
	docVars := make(map[string]string, len(originalVars))
	for key, value := range originalVars {
		if value == "null" {
			docVars[key] = nullVarSentinel
		} else {
			docVars[key] = value
		}
	}
	s.Vars = docVars
	s.TestEnvVars = docTestEnvVars
	docHCL, err := s.ExecuteTemplate(sysfs)
	if err != nil {
		s.templateErrs = append(s.templateErrs, err)
		s.ResourceIdVars = originalResourceIdVars
		s.Vars = originalVars
		s.TestEnvVars = originalTestEnvVars
		return
	}
	s.DocumentationHCLText = replaceNullVarSentinels(docHCL)
	s.DocumentationHCLText = trailingBlankLineRegexp.ReplaceAllString(s.DocumentationHCLText, "\n")

	// Remove region tags
	re1 := regionTagLineRegexp
	re2 := regionTagTrailingRegexp
	s.DocumentationHCLText = re1.ReplaceAllString(s.DocumentationHCLText, "")
	s.DocumentationHCLText = re2.ReplaceAllString(s.DocumentationHCLText, "")

	testResourceIdVars := make(map[string]string)
	testVars := make(map[string]string)
	testTestEnvVars := make(map[string]string)
	testContextVars := make(map[string]string)
	// Override resource_id_vars to inject test values into configs - will have
	//   - "a-example-var-value%{random_suffix}""
	//   - "%{my_var}" for overrides that have custom Golang values
	for key, value := range originalResourceIdVars {
		var newVal string
		if strings.Contains(value, "_") {
			newVal = fmt.Sprintf("tf_test_%s", value)
		} else {
			newVal = fmt.Sprintf("tf-test-%s", value)
		}
		// Random suffix is 10 characters and standard name length <= 64
		if len(newVal) > 54 {
			newVal = newVal[:54]
		}
		testResourceIdVars[key] = fmt.Sprintf("%%{%s}", key)
		testContextVars[key] = fmt.Sprintf("\"%s\"+randomSuffix", newVal)
	}

	for key, value := range originalVars {
		testVars[key] = fmt.Sprintf("%%{%s}", key)
		if value == "null" {
			testContextVars[key] = "nil"
		} else if _, err := strconv.ParseBool(value); err == nil {
			testContextVars[key] = value
		} else if _, err := strconv.ParseInt(value, 10, 64); err == nil {
			testContextVars[key] = value
		} else {
			testContextVars[key] = fmt.Sprintf("%#v", value)
		}
	}

	// Apply overrides from YAML
	for key, value := range s.TestVarsOverrides {
		testContextVars[key] = fmt.Sprintf("%s", value)
	}

	for key := range originalTestEnvVars {
		testTestEnvVars[key] = fmt.Sprintf("%%{%s}", key)
	}

	s.ResourceIdVars = testResourceIdVars
	s.TestEnvVars = testTestEnvVars
	s.Vars = testVars
	s.TestContextVars = testContextVars
	testHCL, err := s.ExecuteTemplate(sysfs)
	if err != nil {
		s.templateErrs = append(s.templateErrs, err)
		s.ResourceIdVars = originalResourceIdVars
		s.Vars = originalVars
		s.TestEnvVars = originalTestEnvVars
		return
	}
	s.TestHCLText = trailingBlankLineRegexp.ReplaceAllString(testHCL, "\n")
	// Remove region tags
	s.TestHCLText = re1.ReplaceAllString(s.TestHCLText, "")
	s.TestHCLText = re2.ReplaceAllString(s.TestHCLText, "")
	s.TestHCLText = SubstituteTestPaths(s.TestHCLText)
	s.hclLoaded = true

	// Reset the step
	s.ResourceIdVars = originalResourceIdVars
	s.Vars = originalVars
	s.TestEnvVars = originalTestEnvVars
}

func indirectInterface(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		return v.Elem()
	}
	return v
}

// strictIndex mirrors text/template's built-in index function, except that
// indexing a missing key on a map returns an explicit error instead of a zero
// value (which text/template's built-in index returns even with missingkey=error).
func strictIndex(item reflect.Value, indexes ...reflect.Value) (reflect.Value, error) {
	item = indirectInterface(item)
	if !item.IsValid() {
		return reflect.Value{}, fmt.Errorf("index of untyped nil")
	}
	for _, idx := range indexes {
		item = indirectInterface(item)
		for item.Kind() == reflect.Pointer {
			if item.IsNil() {
				return reflect.Value{}, fmt.Errorf("index of nil pointer")
			}
			item = item.Elem()
		}
		switch item.Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			idx = indirectInterface(idx)
			var x int
			switch idx.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				x = int(idx.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				x = int(idx.Uint())
			default:
				return reflect.Value{}, fmt.Errorf("cannot index slice/array with type %s", idx.Type())
			}
			if x < 0 || x >= item.Len() {
				return reflect.Value{}, fmt.Errorf("index out of range: %d", x)
			}
			item = item.Index(x)
		case reflect.Map:
			idx = indirectInterface(idx)
			var key reflect.Value
			if !idx.IsValid() {
				key = reflect.Zero(item.Type().Key())
			} else if idx.Type().AssignableTo(item.Type().Key()) {
				key = idx
			} else if idx.Type().ConvertibleTo(item.Type().Key()) {
				key = idx.Convert(item.Type().Key())
			} else {
				return reflect.Value{}, fmt.Errorf("value of type %s is not assignable to type %s", idx.Type(), item.Type().Key())
			}
			if item.IsNil() {
				return reflect.Value{}, fmt.Errorf("map has no entry for key %q", fmt.Sprint(key.Interface()))
			}
			elem := item.MapIndex(key)
			if !elem.IsValid() {
				return reflect.Value{}, fmt.Errorf("map has no entry for key %q", fmt.Sprint(key.Interface()))
			}
			item = elem
		default:
			return reflect.Value{}, fmt.Errorf("can't index item of type %s", item.Type())
		}
	}
	return item, nil
}

func (s *Step) ExecuteTemplate(sysfs fs.FS) (string, error) {
	templateContent, err := fs.ReadFile(sysfs, s.ConfigPath)
	if err != nil {
		return "", err
	}

	fileContentString := string(templateContent)
	templateFileName := filepath.Base(s.ConfigPath)

	funcs := google.TemplateFunctions(sysfs)
	funcs["index"] = strictIndex

	tmpl, err := template.New(templateFileName).Option("missingkey=error").Funcs(funcs).Parse(fileContentString)
	if err != nil {
		return "", fmt.Errorf("parsing sample template %s: %w", s.ConfigPath, err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, templateFileName, s); err != nil {
		return "", fmt.Errorf("executing sample template %s: %w", s.ConfigPath, err)
	}

	rs := contents.String()

	if !strings.HasSuffix(rs, "\n") {
		rs = fmt.Sprintf("%s\n", rs)
	}

	return rs, nil
}

func (s *Step) OiCSLink() string {
	v := url.Values{}
	v.Add("cloudshell_git_repo", "https://github.com/terraform-google-modules/docs-examples.git")
	v.Add("cloudshell_working_dir", s.Name)
	v.Add("cloudshell_image", "gcr.io/cloudshell-images/cloudshell:latest")
	v.Add("open_in_editor", "main.tf")
	v.Add("cloudshell_print", "./motd")
	v.Add("cloudshell_tutorial", "./tutorial.md")
	u := url.URL{
		Scheme:   "https",
		Host:     "console.cloud.google.com",
		Path:     "/cloudshell/open",
		RawQuery: v.Encode(),
	}
	return u.String()
}

func SubstituteExamplePaths(config string) string {
	config = strings.ReplaceAll(config, "../static/img/header-logo.png", "../static/header-logo.png")
	config = strings.ReplaceAll(config, "path/to/private.key", "../static/ssl_cert/test.key")
	config = strings.ReplaceAll(config, "path/to/id_rsa.pub", "../static/ssh_rsa.pub")
	config = strings.ReplaceAll(config, "path/to/certificate.crt", "../static/ssl_cert/test.crt")
	return config
}

func SubstituteTestPaths(config string) string {
	config = strings.ReplaceAll(config, "../static/img/header-logo.png", "test-fixtures/header-logo.png")
	config = strings.ReplaceAll(config, "path/to/private.key", "test-fixtures/test.key")
	config = strings.ReplaceAll(config, "path/to/certificate.crt", "test-fixtures/test.crt")
	config = strings.ReplaceAll(config, "path/to/index.zip", "%{zip_path}")
	config = strings.ReplaceAll(config, "verified-domain.com", "tf-test-domain%{random_suffix}.gcp.tfacc.hashicorptest.com")
	config = strings.ReplaceAll(config, "path/to/id_rsa.pub", "test-fixtures/ssh_rsa.pub")
	return config
}

// Executes step configuration templates for documentation and tests
func (s *Step) SetOiCSHCLText(sysfs fs.FS) {
	originalResourceIdVars := s.ResourceIdVars
	originalVars := s.Vars

	// // Remove region tags
	re1 := regionTagLineRegexp
	re2 := regionTagTrailingRegexp

	testResourceIdVars := make(map[string]string)
	testVars := make(map[string]string)
	for key, value := range originalResourceIdVars {
		testResourceIdVars[key] = fmt.Sprintf("%s-${local.name_suffix}", value)
	}

	for key, value := range originalVars {
		if value == "null" {
			testVars[key] = nullVarSentinel
		} else {
			testVars[key] = value
		}
	}

	// Apply overrides from YAML
	for key, value := range s.OicsVarsOverrides {
		testResourceIdVars[key] = value
		if value == "null" {
			testVars[key] = nullVarSentinel
		} else {
			testVars[key] = value
		}
	}

	s.ResourceIdVars = testResourceIdVars
	s.Vars = testVars
	oicsHCL, err := s.ExecuteTemplate(sysfs)
	if err != nil {
		log.Fatal(err)
	}
	s.OicsHCLText = replaceNullVarSentinels(oicsHCL)
	s.OicsHCLText = trailingBlankLineRegexp.ReplaceAllString(s.OicsHCLText, "\n")

	// Remove region tags
	s.OicsHCLText = re1.ReplaceAllString(s.OicsHCLText, "")
	s.OicsHCLText = re2.ReplaceAllString(s.OicsHCLText, "")
	s.OicsHCLText = SubstituteExamplePaths(s.OicsHCLText)

	// Reset the step
	s.ResourceIdVars = originalResourceIdVars
	s.Vars = originalVars
}

func (s *Step) ShouldGenerateDoc(index int, sample *Sample) bool {
	return s.IncludeStepDoc || (index == 0 && !sample.ExcludeBasicDoc)
}
