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
		"error message with start and end markers": {
			rawErrorMessage: "------- Stdout: ------- === RUN TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === PAUSE TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === CONT TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample --- PASS: TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample (11.76s) PASS ------- Stderr: ------- 2025/01/21 08:06:22 [DEBUG] [transport] [server-transport 0xc002614000] Closing: EOF 2025/01/21 08:06:22 [DEBUG] [transport] [server-transport 0xc002614000] loopyWriter exiting with error: transport closed by client 2025/01/21 08:06:22",
			want:            "=== RUN TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === PAUSE TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === CONT TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample --- PASS: TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample (11.76s) PASS",
		},
		"error message with start but no end markers": {
			rawErrorMessage: "------- Stdout: ------- === RUN TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === PAUSE TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === CONT TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample --- PASS: TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample (11.76s) PASS",
			want:            "=== RUN TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === PAUSE TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === CONT TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample --- PASS: TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample (11.76s) PASS",
		},
		"error message with no start but with end markers": {
			rawErrorMessage: "=== RUN TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === PAUSE TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === CONT TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample --- PASS: TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample (11.76s) PASS ------- Stderr: ------- 2025/01/21 08:06:22 [DEBUG] [transport] [server-transport 0xc002614000] Closing: EOF 2025/01/21 08:06:22 [DEBUG] [transport] [server-transport 0xc002614000] loopyWriter exiting with error: transport closed by client 2025/01/21 08:06:22",
			want:            "=== RUN TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === PAUSE TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === CONT TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample --- PASS: TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample (11.76s) PASS",
		},
		"error message with no start and no end markers": {
			rawErrorMessage: "=== RUN TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === PAUSE TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === CONT TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample --- PASS: TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample (11.76s) PASS",
			want:            "=== RUN TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === PAUSE TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample === CONT TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample --- PASS: TestAccColabRuntimeTemplate_colabRuntimeTemplateBasicExample (11.76s) PASS",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got := convertErrorMessage(tc.rawErrorMessage)
			assert.Equal(t, tc.want, got)
		})
	}
}
