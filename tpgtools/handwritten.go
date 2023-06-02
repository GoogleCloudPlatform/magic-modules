// Copyright 2021 Google LLC. All Rights Reserved.
//
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

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
)

func copyHandwrittenFiles(inPath string, outPath string) {
	if inPath == "" || outPath == "" {
		glog.Info("Skipping copying handwritten files, empty path specified")
		return
	}

	glog.Info("copying handwritten files")

	_, err := os.Stat(outPath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(outPath, 0755)
		if errDir != nil {
			glog.Fatal(err)
		}
	}

	// Log warning about unexpected outPath values before adding copyright headers
	// Matches equivalent in MMv1, see below:
	// https://github.com/hashicorp/magic-modules/blob/48ce1004bafd4b4ef1be7565eaa6727adabd0670/mmv1/provider/core.rb#L202-L206
	if !isOutputFolderExpected(outPath) {
		glog.Infof("Unexpected output folder (%s) detected when deciding to add HashiCorp copyright headers. Watch out for unexpected changes to copied files", outPath)
	}

	fs, err := ioutil.ReadDir(inPath)
	if err != nil {
		glog.Fatal(err)
	}
	for _, f := range fs {
		if f.IsDir() {
			copyHandwrittenFiles(path.Join(inPath, f.Name()), path.Join(outPath, f.Name()))
			return
		}

		// Ignore empty go.mod
		if f.Name() == "go.mod" {
			continue
		}

		b, err := ioutil.ReadFile(path.Join(inPath, f.Name()))
		if err != nil {
			if !os.IsNotExist(err) {
				glog.Exit(err)
			}
			// Ignore the error if the file just doesn't exist
			continue
		}

		// Add HashiCorp copyright header only if generating TPG/TPGB
		if strings.HasSuffix(outPath, "/terraform-provider-google") || strings.HasSuffix(outPath, "/terraform-provider-google-beta") {
			if strings.HasSuffix(f.Name(), ".go") {
				copyrightHeader := []byte("// Copyright (c) HashiCorp, Inc.\n// SPDX-License-Identifier: MPL-2.0\n")
				b = append(copyrightHeader, b...)
			}
		}

		// Format file if ending in .go
		if strings.HasSuffix(f.Name(), ".go") {
			b, err = formatSource(bytes.NewBuffer(b))
			if err != nil {
				glog.Error("error formatting %s: %v", f.Name(), err)
				continue
			}
		}

		// Write copied file.
		err = ioutil.WriteFile(path.Join(outPath, terraformResourceDirectory, f.Name()), b, 0644)
		if err != nil {
			glog.Exit(err)
		}
	}
}

// isOutputFolderExpected returns a boolean indicating if the output folder is present in an allow list.
// Intention is to warn users about unexpected diffs if they have renamed their cloned copy of downstream repos,
// as this affects detecting which downstream they're building and whether to add copyright headers.
// Written to match `expected_output_folder?` method in MMv1, see below:
// https://github.com/hashicorp/magic-modules/blob/48ce1004bafd4b4ef1be7565eaa6727adabd0670/mmv1/provider/core.rb#L266-L282
func isOutputFolderExpected(outPath string) bool {
	pathComponents := strings.Split(outPath, "/")
	folderName := pathComponents[len(pathComponents)-1] // last element

	switch folderName {
	case "terraform-provider-google",
		"terraform-provider-google-beta",
		"terraform-next",
		"terraform-google-conversion",
		"tfplan2cai":
		return true
	default:
		return false
	}
}
