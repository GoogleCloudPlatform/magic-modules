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
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Classification of a PR test failure against the nightly test history.
const (
	NightlyStatusFailing       = "Failing"
	NightlyStatusFlaky         = "Flaky"
	NightlyStatusRecentlyFixed = "Recently fixed"
	NightlyStatusPassing       = "Passing"
	NightlyStatusNotFound      = "Not found"
)

// Minimum nightly failures within the history window for a test to be considered
// consistently failing in nightly.
const nightlyFailingThreshold = 3

// Consecutive most-recent nights without a failure after which a test with failures earlier in
// the window is treated as fixed rather than flakey. A fix merged to main mid-window leaves stale
// failures behind, so a run of green nights is better evidence than the window as a whole.
const nightlyRecentlyFixedNights = 3

// Minimum share of nightly runs that must fail for a test to be called flakey. Below this a
// single stale failure in a 30-night window would otherwise permanently label a healthy test.
const nightlyMinFlakyRate = 0.1

// loadNightlyTestHistory downloads the rolling nightly test history for pVersion.
// It returns nil (without error) when gcs is nil so callers can treat the history as optional.
func loadNightlyTestHistory(pVersion provider.Version, gcs CloudstorageClient) (*NightlyTestHistoryReport, error) {
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
	return &report, nil
}

// nightsSinceLastFailure returns how many nights have passed between a test's last nightly failure
// and the end of the history window. It returns -1 when this cannot be determined, either because
// the history predates last_failure_date or because a date failed to parse.
//
// This counts calendar days, which overstates the number of nightly runs on days where no nightly
// finished. That only ever makes the "recently fixed" check more conservative in the sense of
// needing more elapsed time, so it is acceptable for a reviewer-facing hint.
func nightsSinceLastFailure(h *NightlyTestHistory, historyEndDate string) int {
	if h.LastFailureDate == "" || historyEndDate == "" {
		return -1
	}
	const layout = "2006-01-02"
	last, err := time.Parse(layout, h.LastFailureDate)
	if err != nil {
		return -1
	}
	end, err := time.Parse(layout, historyEndDate)
	if err != nil {
		return -1
	}
	return int(end.Sub(last).Hours() / 24)
}

// classifyNightlyStatus returns how the given test behaves in recent nightly runs.
// A nil entry means the test was never seen in the history window. historyEndDate is the last
// night covered by the history and may be empty, which disables the recently-fixed check.
func classifyNightlyStatus(h *NightlyTestHistory, historyEndDate string) string {
	if h == nil {
		return NightlyStatusNotFound
	}
	// Skipped runs neither pass nor fail, so a test that only ever skipped was effectively not run.
	runs := h.Passes + h.Failures
	if runs == 0 {
		return NightlyStatusNotFound
	}
	if h.Failures == 0 {
		return NightlyStatusPassing
	}

	// A fix merged to main partway through the window leaves failures that no longer describe the
	// current state. Trust a streak of green nights over the window totals.
	if h.LastStatus != "FAILURE" {
		if n := nightsSinceLastFailure(h, historyEndDate); n >= nightlyRecentlyFixedNights {
			return NightlyStatusRecentlyFixed
		}
	}

	if h.LastStatus == "FAILURE" && (h.Failures >= nightlyFailingThreshold || h.Passes == 0) {
		return NightlyStatusFailing
	}
	if h.Passes == 0 {
		return NightlyStatusFailing
	}
	// Rare stale failures are not a useful signal that a PR failure is pre-existing, but a failure
	// on the most recent night is, regardless of rate.
	if float64(h.Failures)/float64(runs) >= nightlyMinFlakyRate || h.LastStatus == "FAILURE" {
		return NightlyStatusFlaky
	}
	return NightlyStatusPassing
}

// nightlyTestUrl links to a test's history page in the TeamCity nightly project. testNameId is
// TeamCity's cross-build test identifier; it cannot be derived from the test name, so it is
// captured when the nightly history is collected.
func nightlyTestUrl(testNameId string, pVersion provider.Version) string {
	if testNameId == "" {
		return ""
	}
	return fmt.Sprintf("https://hashicorp.teamcity.com/test/%s?currentProjectId=%s", testNameId, pVersion.TeamCityNightlyProjectName())
}

// lookupNightlyHistory finds the history entry for a test name. The name may be a VCR
// subtest name (Parent__sub); the parent test is used as a fallback.
func lookupNightlyHistory(testName string, history map[string]*NightlyTestHistory) *NightlyTestHistory {
	if h, ok := history[strings.ReplaceAll(testName, "__", "/")]; ok {
		return h
	}
	return history[compoundTest(testName)]
}

// nightlySymbol renders a nightly status as a finding for the PR comment table.
//
// Nightly findings only appear for tests that failed in this PR's recording, so the question they
// answer is "is this failure mine?". The emoji therefore signals how much the author needs to look
// at the row, not how healthy the test is in nightly: a test that always fails in nightly is very
// likely pre-existing (low alarm), while one that passes cleanly in nightly points at this PR.
//
// A test with no nightly runs yields no finding, since there is nothing to compare against.
func nightlySymbol(status string) string {
	switch status {
	case NightlyStatusFailing:
		return "⚪ Nightly fails"
	case NightlyStatusFlaky:
		return "🟡 Nightly flaky"
	case NightlyStatusRecentlyFixed:
		return "🔴 Nightly fixed"
	case NightlyStatusPassing:
		return "🔴 Nightly passes"
	default:
		return ""
	}
}

// nightlyDetail is the short qualifier shown next to a nightly status, kept terse so rows stay
// compact when a PR has many failures. The fuller explanation lives below the table.
//
// The rate always describes the verb it sits next to, so "passes 100% of 30" and "fails 100% of 26"
// read naturally. Flaky has no verb of its own, so it says "fails" explicitly.
//
// The denominator is kept because it carries the confidence of the rate: 100% of 26 runs is far
// stronger evidence than 100% of 2.
//
// A recently fixed test shows when it last failed rather than its rate, since its rate describes
// the period before the fix and would otherwise contradict the label.
func nightlyDetail(h *NightlyTestHistory, status string) string {
	if h == nil {
		return ""
	}
	if status == NightlyStatusRecentlyFixed && h.LastFailureDate != "" {
		return "last failed " + h.LastFailureDate
	}
	runs := h.Passes + h.Failures
	if runs == 0 {
		return ""
	}
	percentOf := func(n int) string {
		return fmt.Sprintf("%d%% of %d", int(math.Round(float64(n)/float64(runs)*100)), runs)
	}
	switch status {
	case NightlyStatusPassing:
		return percentOf(h.Passes)
	case NightlyStatusFlaky:
		return "fails " + percentOf(h.Failures)
	default:
		return percentOf(h.Failures)
	}
}

// nightlyFinding renders a nightly status and the detail backing it as a code-styled chip, e.g.
// "`⚪ Nightly fails 100% of 26`".
func nightlyFinding(row VCRTestTableRow) string {
	label := nightlySymbol(row.NightlyStatus)
	if label == "" {
		return ""
	}
	if row.NightlyDetail != "" {
		label += " " + row.NightlyDetail
	}
	return "`" + label + "`"
}
