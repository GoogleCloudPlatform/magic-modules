/*
* Copyright 2026 Google LLC. All Rights Reserved.
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
	"fmt"
	"magician/provider"
	utils "magician/utility"
	"os"
	"path/filepath"
	"strings"
)

// Classification of a PR test failure against the nightly test history.
const (
	NightlyStatusFailing  = "Failing"
	NightlyStatusFlaky    = "Flaky"
	NightlyStatusPassing  = "Passing"
	NightlyStatusNotFound = "Not found"
)

// Minimum nightly failures within the history window for a test to be considered
// consistently failing in nightly.
const nightlyFailingThreshold = 3

// loadNightlyTestHistory downloads the rolling nightly test history for pVersion.
// It returns nil (without error) when gcs is nil so callers can treat the history as optional.
func loadNightlyTestHistory(pVersion provider.Version, gcs CloudstorageClient) (map[string]*NightlyTestHistory, error) {
	if gcs == nil {
		return nil, nil
	}
	localPath := filepath.Join(os.TempDir(), fmt.Sprintf("nightly-test-history-%s.json", pVersion.String()))
	defer os.Remove(localPath)

	if err := gcs.DownloadFile(nightlyDataBucket, nightlyTestHistoryObjectName(pVersion), localPath); err != nil {
		return nil, fmt.Errorf("failed to download nightly test history: %w", err)
	}
	var report NightlyTestHistoryReport
	if err := utils.ReadFromJson(&report, localPath); err != nil {
		return nil, fmt.Errorf("failed to read nightly test history: %w", err)
	}
	return report.Tests, nil
}

// classifyNightlyStatus returns how the given test behaves in recent nightly runs.
// The test name may be a VCR subtest name (Parent__sub); the parent test is used as a fallback.
func classifyNightlyStatus(testName string, history map[string]*NightlyTestHistory) string {
	h, ok := history[strings.ReplaceAll(testName, "__", "/")]
	if !ok {
		h, ok = history[compoundTest(testName)]
	}
	if !ok || h == nil {
		return NightlyStatusNotFound
	}
	if h.Failures > 0 && h.LastStatus == "FAILURE" && (h.Failures >= nightlyFailingThreshold || h.Passes == 0) {
		return NightlyStatusFailing
	}
	if h.Failures > 0 && h.Passes > 0 {
		return NightlyStatusFlaky
	}
	if h.Passes > 0 {
		return NightlyStatusPassing
	}
	return NightlyStatusNotFound
}

// nightlySymbol renders a nightly status for the PR comment table.
func nightlySymbol(status string) string {
	switch status {
	case NightlyStatusFailing:
		return "🔴 Failing in nightly"
	case NightlyStatusFlaky:
		return "🟡 Flaky in nightly"
	case NightlyStatusPassing:
		return "🟢 Passing in nightly"
	case NightlyStatusNotFound:
		return "Not run in nightly"
	default:
		return "-"
	}
}
