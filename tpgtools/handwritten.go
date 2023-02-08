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

		// Format file if ending in .go
		if strings.HasSuffix(f.Name(), ".go") {
			b, err = formatSource(bytes.NewBuffer(b))
			if err != nil {
				glog.Error("error formatting %s: %v", f.Name(), err)
				continue
			}
		}

		glog.Infof("copying handwritten file %s", path.Join(outPath, terraformResourceDirectory, f.Name()))
		// Write copied file.
		err = ioutil.WriteFile(path.Join(outPath, terraformResourceDirectory, f.Name()), b, 0644)
		if err != nil {
			glog.Exit(err)
		}
	}
}
