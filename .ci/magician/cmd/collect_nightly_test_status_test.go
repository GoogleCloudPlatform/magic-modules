/*
* Copyright 2025 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertServiceName(t *testing.T) {
	cases := map[string]struct {
		servicePath string
		want        string
		wantError   bool
	}{
		"valid service path": {
			servicePath: "TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS_GOOGLE_PACKAGE_SECRETMANAGER",
			want:        "secretmanager",
			wantError:   false,
		},
		"invalid service path": {
			servicePath: "SECRETMANAGER",
			want:        "",
			wantError:   true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got, err := convertServiceName(tc.servicePath)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestConvertErrorMessage(t *testing.T) {
	cases := map[string]struct {
		rawErrorMessage string
		want            string
	}{
		"panic scenario keeps stderr": {
			rawErrorMessage: "Test ended in panic. ------- Stdout: ------- === RUN TestAccComputeProjectMetadata_modify_2 === PAUSE TestAccComputeProjectMetadata_modify_2 === CONT TestAccComputeProjectMetadata_modify_2 ------- Stderr: ------- panic: runtime error: invalid memory address or nil pointer dereference",
			want:            "=== RUN TestAccComputeProjectMetadata_modify_2 === PAUSE TestAccComputeProjectMetadata_modify_2 === CONT TestAccComputeProjectMetadata_modify_2\npanic: runtime error: invalid memory address or nil pointer dereference",
		},
		"standard failure with debug logs drops stderr": {
			rawErrorMessage: "=== RUN TestAccComputeInstance_Basic === PAUSE TestAccComputeInstance_Basic === CONT TestAccComputeInstance_Basic --- FAIL: TestAccComputeInstance_Basic (11.76s) FAIL ------- Stderr: ------- 2025/01/21 08:06:22 [DEBUG] [transport] Closing: EOF",
			want:            "=== RUN TestAccComputeInstance_Basic === PAUSE TestAccComputeInstance_Basic === CONT TestAccComputeInstance_Basic --- FAIL: TestAccComputeInstance_Basic (11.76s) FAIL",
		},
		"standard failure with start marker but no stderr": {
			rawErrorMessage: "------- Stdout: ------- === RUN TestAccComputeInstance_Basic === PAUSE TestAccComputeInstance_Basic === CONT TestAccComputeInstance_Basic --- FAIL: TestAccComputeInstance_Basic (11.76s) FAIL",
			want:            "=== RUN TestAccComputeInstance_Basic === PAUSE TestAccComputeInstance_Basic === CONT TestAccComputeInstance_Basic --- FAIL: TestAccComputeInstance_Basic (11.76s) FAIL",
		},
		"failure with no start and no stderr markers": {
			rawErrorMessage: "=== RUN TestAccComputeInstance_Basic === PAUSE TestAccComputeInstance_Basic === CONT TestAccComputeInstance_Basic --- FAIL: TestAccComputeInstance_Basic (11.76s) FAIL",
			want:            "=== RUN TestAccComputeInstance_Basic === PAUSE TestAccComputeInstance_Basic === CONT TestAccComputeInstance_Basic --- FAIL: TestAccComputeInstance_Basic (11.76s) FAIL",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got := convertErrorMessage(tc.rawErrorMessage)
			assert.Equal(t, tc.want, got)
		})
	}
}
