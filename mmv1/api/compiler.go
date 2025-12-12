// Copyright 2024 Google Inc.
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

package api

import (
	"bytes"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

func Compile(yamlPath string, obj interface{}) {
	objYaml, err := os.ReadFile(yamlPath)

	if err != nil {
		log.Fatalf("Cannot open the file: %s", yamlPath)
	}
	CompileContents(objYaml, obj, yamlPath)
}

func CompileContents(contents []byte, obj interface{}, yamlPath string) {
	// TODO: retire {{override_path}} from private overrides repositories,
	// and remove this later.
	contents = bytes.ReplaceAll(contents, []byte("{{override_path}}/"), []byte(""))

	yamlValidator := google.YamlValidator{}
	yamlValidator.Parse(contents, obj, yamlPath)
}
