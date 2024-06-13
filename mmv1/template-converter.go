package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang/glog"
)

func find(root, ext string) []string {
	var a []string

	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ext {
			a = append(a, file.Name())
		}
	}
	return a
}

func convertTemplates() {
	// exculding iam
	folders := []string{"examples", "constants", "custom_check_destroy", "custom_create", "custom_delete", "custom_import", "custom_update", "decoders", "encoders", "extra_schema_entry", "post_create", "post_create_failure", "post_delete", "post_import", "post_update", "pre_create", "pre_delete", "pre_read", "pre_update", "state_migrations", "update_encoder", "custom_expand", "custom_flatten", "iam/example_config_body"}
	counts := 0
	for _, folder := range folders {
		counts += convertTemplate(folder)
	}
	log.Printf("%d template files in %d subfolders total", counts, len(folders))
}

func convertTemplate(folder string) int {
	rubyDir := fmt.Sprintf("templates/terraform/%s", folder)
	goDir := fmt.Sprintf("templates/terraform/%s/go", folder)

	if err := os.MkdirAll(goDir, os.ModePerm); err != nil {
		glog.Error(fmt.Errorf("error creating directory %v: %v", goDir, err))
	}

	templates := find(rubyDir, ".erb")
	log.Printf("%d template files in folder %s", len(templates), folder)

	for _, file := range templates {
		filePath := path.Join(rubyDir, file)
		if checkExceptionList(filePath) {
			continue
		}
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Cannot open the file: %v", file)
		}

		// Replace {{}}
		r, err := regexp.Compile(`{{(.*?)}}`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{"{{"}}$1{{"}}"}}`))

		// Replace primary_resource_id
		r, err = regexp.Compile(`<%=\s*ctx\[:primary_resource_id\]\s*-?%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte("{{$.PrimaryResourceId}}"))

		// Replace vars
		r, err = regexp.Compile(`<%=\s*ctx\[:vars\]\[('|")([a-zA-Z0-9_-]+)('|")\]\s*-?%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{index $.Vars "$2"}}`))

		// Replace test_env_vars
		r, err = regexp.Compile(`<%=\s*ctx\[:test_env_vars\]\[('|")([a-zA-Z0-9_-]+)('|")\]\s*-?%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{index $.TestEnvVars "$2"}}`))

		// Replace <% unless compiler == "terraformgoogleconversion-codegen" -%>
		r, err = regexp.Compile(`<% unless compiler == "terraformgoogleconversion-codegen" -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- if ne $.Compiler "terraformgoogleconversion-codegen" }}`))

		// Replace <% unless version == 'ga' -%>
		r, err = regexp.Compile(`<% unless version == ['|"]ga['|"] -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- if ne $.TargetVersionName "ga" }}`))

		// Replace <% if version == 'ga' -%>
		r, err = regexp.Compile(`<% if version == 'ga' -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- if eq $.TargetVersionName "ga" }}`))

		// Replace <% else -%>
		r, err = regexp.Compile(`<% else -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- else }}`))

		// Replace <%= object.name -%>
		r, err = regexp.Compile(`<%= object.name -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.Name}}`))

		// Replace <%= object.resource_name -%>
		r, err = regexp.Compile(`<%= object.resource_name -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.ResourceName}}`))

		// Replace <%=object.self_link_uri-%>
		r, err = regexp.Compile(`<%=object.self_link_uri-%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.SelfLinkUri}}`))

		// Replace <%=object.create_uri-%>
		r, err = regexp.Compile(`<%=object.create_uri-%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.CreateUri}}`))

		// Replace <%=object.base_url-%>
		r, err = regexp.Compile(`<%=object.base_url-%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.BaseUrl}}`))

		// Replace <%=object.__product.name-%>
		r, err = regexp.Compile(`<%=object.__product.name-%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.ProductMetadata.Name}}`))

		// Replace <% if object.name == 'Disk' -%>
		r, err = regexp.Compile(`<% if object.name == 'Disk' -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- if eq $.Name "Disk" }}`))

		// Replace <% elsif object.name == 'RegionDisk' -%>
		r, err = regexp.Compile(`<% elsif object.name == 'RegionDisk' -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- else if eq $.Name "RegionDisk" }}`))

		// Replace <% if object.properties.any?{ |p| p.name == "labels" } -%>
		r, err = regexp.Compile(`<% if object\.properties.any\?\{ \|p\| p\.name == "labels" \} -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- if $.HasLabelsField }}`))

		// Replace <% if object.error_retry_predicates -%>
		r, err = regexp.Compile(`<% if object.error_retry_predicates -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- if $.ErrorRetryPredicates }}`))

		// Replace <% if object.error_abort_predicates -%>
		r, err = regexp.Compile(`<% if object.error_abort_predicates -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- if $.ErrorAbortPredicates }}`))

		// Replace <%= object.error_retry_predicates.join(',') -%>
		r, err = regexp.Compile(`<%= object.error_retry_predicates.join\(','\) -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(` {{- join $.ErrorRetryPredicates "," -}} `))

		// Replace <%= object.error_abort_predicates.join(',') -%>
		r, err = regexp.Compile(`<%= object.error_abort_predicates.join\(','\) -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(` {{- join $.ErrorAbortPredicates "," -}} `))

		// Replace <%= object.name.camelize(:lower) -%>
		r, err = regexp.Compile(`<%= object.name.camelize\(:lower\) -?%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{camelize $.Name "lower"}}`))

		// Replace <%= object.name.plural.camelize(:lower) -%>
		r, err = regexp.Compile(`<%= object.name.plural.camelize\(:lower\) -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{camelize (plural $.Name) "lower"}}`))

		// Replace <%= id_format(object) -%>
		r, err = regexp.Compile(`<%= id_format\(object\) -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.GetIdFormat}}`))

		// Replace <%= prefix -%>
		r, err = regexp.Compile(`<%= prefix -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.GetPrefix}}`))

		// Replace <%= titlelize_property(property) -%>
		r, err = regexp.Compile(`<%= titlelize_property\(property\) -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.TitlelizeProperty}}`))

		// Replace <%= prop_path -%>
		r, err = regexp.Compile(`<%= prop_path -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.PropPath}}`))

		// Replace <%= go_literal(property.default_value) -%>
		r, err = regexp.Compile(`<%= go_literal\(property.default_value\) -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.DefaultValue}}`))

		// Replace <%= build_expand_resource_ref('v.(string)', property, pwd) %>
		r, err = regexp.Compile(`<%= build_expand_resource_ref\('v\.\(string\)', property, pwd\) %>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{ template "expandResourceRef" dict "VarName" "v.(string)" "ResourceRef" $.ResourceRef "ResourceType" $.ResourceType}}`))

		// Replace <%= build_expand_resource_ref('raw.(string)', property.item_type, pwd) %>
		r, err = regexp.Compile(`<%= build_expand_resource_ref\('raw\.\(string\)', property\.item_type, pwd\) %>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{ template "expandResourceRef" dict "VarName" "raw.(string)" "ResourceRef" $.ItemType.ResourceRef "ResourceType" $.ItemType.ResourceType}}`))

		// Replace <%- if property.is_a?(Api::Type::Integer) -%>
		r, err = regexp.Compile(`<%- if property.is_a\?\(Api::Type::Integer\) -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- if $.IsA "Integer" }}`))

		// Replace <%= property.name.underscore -%>
		r, err = regexp.Compile(`<%= property.name.underscore -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{underscore $.Name}}`))

		// Replace <%= resource_type -%>
		r, err = regexp.Compile(`<%= resource_type -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{$.ResourceType}}`))

		// Replace <%  if property.is_set -%>
		r, err = regexp.Compile(`<%  if property.is_set -%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- if $.IsSet }}`))

		// Replace <% end -%>
		r, err = regexp.Compile(`<%[\s-]*end[\s-]*%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- end }}`))

		copyRight := `{{- /*
	The license inside this block applies to this file
	Copyright 2024 Google Inc.
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/}}`
		// Replace copyright
		r, err = regexp.Compile(`(?s)<%[-\s#]*[tT]he license inside this.*?limitations under the License..*?%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(copyRight))

		// Replace comments
		r, err = regexp.Compile(`(?s)<%#-?\s?(.*?)\s?-?%>`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`{{- /* $1 */}}`))

		// Replace .erb
		r, err = regexp.Compile(`\.erb`)
		if err != nil {
			log.Fatalf("Cannot compile the regular expression: %v", err)
		}
		data = r.ReplaceAll(data, []byte(`.tmpl`))

		goTemplate := strings.Replace(file, "erb", "tmpl", 1)
		err = ioutil.WriteFile(path.Join(goDir, goTemplate), data, 0644)
		if err != nil {
			glog.Exit(err)
		}
	}

	return len(templates)
}

func checkExceptionList(filePath string) bool {
	exceptionPaths := []string{
		"custom_flatten/bigquery_table_ref_load_destinationtable.go",
		"custom_flatten/bigquery_table_ref.go",
		"custom_flatten/bigquery_table_ref_copy_destinationtable.go",
		"custom_flatten/bigquery_table_ref_extract_sourcetable.go",
		"custom_flatten/bigquery_table_ref_query_destinationtable.go",
	}

	for _, t := range exceptionPaths {
		if strings.Contains(filePath, t) {
			return true
		}
	}

	return false
}
