package cmd

import (
	"fmt"
	"os"
	"testing"
	"time"
)

const validApache2LicenseFormat string = `
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`
const validCopyright string = "Copyright %d Google Inc\n"
const invalidCopyright string = "Copyright 1900 Google Inc\n"
const invalidLicense string = "not a license"

func TestCheckLicenseType(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "success",
			content: fmt.Sprintf(validCopyright, time.Now().Year()) + validApache2LicenseFormat,
		},
		{
			name:    "invalid license",
			content: fmt.Sprintf(validCopyright, time.Now().Year()) + invalidLicense,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		f, err := os.CreateTemp(dir, "testfile")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString(tc.content); err != nil {
			t.Fatal(err)
		}
		err = checkLicenseType(f.Name())
		if err == nil && tc.wantErr {
			t.Errorf("checkLicenseType() want err, but got nil")
		}
		if err != nil && !tc.wantErr {
			t.Errorf("checkLicenseType() got error %s", err)
		}
	}
}

func TestCheckCopyright(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "success",
			content: fmt.Sprintf(validCopyright, time.Now().Year()) + validApache2LicenseFormat,
		},
		{
			name:    "wrong year copyright",
			content: invalidCopyright + validApache2LicenseFormat,
			wantErr: true,
		},
		{
			name:    "not matching the pattern",
			content: "random str",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.CreateTemp(dir, "testfile")
			if err != nil {
				t.Fatal(err)
			}
			if _, err := f.WriteString(tc.content); err != nil {
				t.Fatal(err)
			}
			err = checkCopyright(f.Name(), time.Now().Year())
			if err == nil && tc.wantErr {
				t.Errorf("checkCopyright() want err, but got nil")
			}
			if err != nil && !tc.wantErr {
				t.Errorf("checkCopyright() got error %s", err)
			}
		})
	}
}
