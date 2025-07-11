// Copyright 2025 Red Hat Inc.
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

package provider

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/golang/glog"
)

type AnsibleTemplateData struct {
	OutputFolder string
	VersionName  string

	OverWrite       bool
	PluginDirectory string
	ModuleDirectory string
	TestDirectories map[string]string
}

// NewAnsibleTemplateData returns a new AnsibleTemplateData struct.
func NewAnsibleTemplateData(outputFolder string, versionName string, overwrite bool) *AnsibleTemplateData {
	atd := AnsibleTemplateData{
		OutputFolder:    outputFolder,
		VersionName:     versionName,
		PluginDirectory: "plugins",
		ModuleDirectory: filepath.Join("plugins", "modules"),
		TestDirectories: make(map[string]string),
		OverWrite:       overwrite,
	}

	atd.TestDirectories["integration"] = path.Join("tests", "integration", "targets")

	return &atd
}

// fileExists returns true if the given path already exists, otherwise returns false
func fileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

// GenerateModuleFile renders the given template into a file in $OUTPUT_DIR/plugins/$MODULE/ directory.
// The directory must already exist.
func (atd *AnsibleTemplateData) GenerateModuleFile(outputFilePath string, resource *api.Resource) error {
	templateBaseDir := path.Join("templates", ANSIBLE_PROVIDER)
	templatePath := path.Join(templateBaseDir, atd.PluginDirectory, "module.py.tmpl")
	additionalTemplates := []string{
		filepath.Join(templateBaseDir, atd.PluginDirectory, "documentation.tmpl"),
	}
	return atd.writeFile(outputFilePath, templatePath, resource, additionalTemplates...)
}

// GenerateTestFile renders the given template into a file in $OUTPUT_DIR/tests/$TEST_TYPE/$MODULE/
// directory. The directory must already exist.
func (atd *AnsibleTemplateData) GenerateTestFile(outputFolder, targetFile, testType string, r *api.Resource) error {
	templateBaseDir := path.Join("templates", ANSIBLE_PROVIDER)
	templatePath := filepath.Join(templateBaseDir, "tests", testType, fmt.Sprintf("%s.tmpl", targetFile))
	filePath := path.Join(outputFolder, "tests", testType, "targets", r.AnsibleName(), targetFile)
	// glog.Infof("template %s", templatePath)
	// glog.Infof("output %s", outputFilePath)
	return atd.writeFile(filePath, templatePath, r)
}

// writeFile renders a given template (or templates) to an existing path.
func (atd *AnsibleTemplateData) writeFile(filePath, templatePath string, resource *api.Resource, additionalTemplates ...string) error {
	funcMap := template.FuncMap{
		"split":     strings.Split,
		"trimSpace": strings.TrimSpace,
		"trim":      strings.Trim,
		"now":       time.Now,
	}
	// add functions already defined in google/templates_utils.go
	maps.Copy(funcMap, google.TemplateFunctions)
	exists := fileExists(filePath)

	if atd.OverWrite || !exists {
		glog.Warningf("writing %s", filePath)

		templates := []string{
			templatePath,
			filepath.Join("templates", ANSIBLE_PROVIDER, "fragments.tmpl"),
		}
		templates = append(templates, additionalTemplates...)
		templateName := filepath.Base(templatePath)
		tmpl := template.New(templateName).Funcs(funcMap)
		render, err := tmpl.ParseFiles(templates...)
		//.Funcs(funcMap).ParseFiles(templates...)
		if err != nil {
			glog.Fatal(fmt.Sprintf("error parsing %s for filepath %s ", templatePath, filePath), err)
		}
		contents := bytes.Buffer{}
		if err = render.ExecuteTemplate(&contents, templateName, resource); err != nil {
			glog.Exit(fmt.Sprintf("error executing %s for filepath %s ", templatePath, filePath), err)
		}

		sourceByte := contents.Bytes()
		if len(sourceByte) == 0 {
			return nil
		}

		return os.WriteFile(filePath, sourceByte, 0644)
	}
	return nil
}
