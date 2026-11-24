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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"magician/provider"

	"github.com/stretchr/testify/assert"
)

// writeDailyStatus stores a daily test-metadata file in the fake bucket for the given day.
func writeDailyStatus(t *testing.T, gcs *fakeGCS, pVersion provider.Version, day string, tests []TestInfo) {
	t.Helper()
	data, err := json.Marshal(tests)
	if err != nil {
		t.Fatal(err)
	}
	object := fmt.Sprintf("test-metadata/%s/%s-%s.json", pVersion.String(), day, pVersion.String())
	gcs.objects[nightlyDataBucket+"/"+object] = data
}

// readHistory reads back the history the command uploaded to the fake bucket.
func readHistory(t *testing.T, gcs *fakeGCS, pVersion provider.Version) NightlyTestHistoryReport {
	t.Helper()
	data, ok := gcs.objects[nightlyDataBucket+"/"+nightlyTestHistoryObjectName(pVersion)]
	if !ok {
		t.Fatal("history was not uploaded")
	}
	var report NightlyTestHistoryReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	return report
}

func TestCreateNightlyTestHistoryCapturesTestNameId(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	assert.NoError(t, err)

	gcs := &fakeGCS{objects: map[string][]byte{}}
	// 3 days: the test fails, passes, then fails again.
	writeDailyStatus(t, gcs, provider.Beta, "2026-09-26", []TestInfo{
		{Name: "TestAccA", Status: "FAILURE", Service: "pubsub", TestNameId: "111"},
	})
	writeDailyStatus(t, gcs, provider.Beta, "2026-09-27", []TestInfo{
		{Name: "TestAccA", Status: "SUCCESS", Service: "pubsub", TestNameId: "111"},
	})
	writeDailyStatus(t, gcs, provider.Beta, "2026-09-28", []TestInfo{
		{Name: "TestAccA", Status: "FAILURE", Service: "pubsub", TestNameId: "111"},
	})

	err = createNightlyTestHistory(provider.Beta, nil, gcs, loc, "2026-09-28", 3)
	assert.NoError(t, err)

	report := readHistory(t, gcs, provider.Beta)
	assert.Equal(t, 3, report.DaysFound)

	h := report.Tests["TestAccA"]
	if assert.NotNil(t, h) {
		assert.Equal(t, 1, h.Passes)
		assert.Equal(t, 2, h.Failures)
		assert.Equal(t, "FAILURE", h.LastStatus)
		assert.Equal(t, "111", h.TestNameId)
	}
}

// Histories collected before test ids were captured must still aggregate cleanly.
func TestCreateNightlyTestHistoryWithoutTestNameId(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	assert.NoError(t, err)

	gcs := &fakeGCS{objects: map[string][]byte{}}
	writeDailyStatus(t, gcs, provider.Beta, "2026-09-28", []TestInfo{
		{Name: "TestAccA", Status: "SUCCESS", Service: "pubsub"},
	})

	err = createNightlyTestHistory(provider.Beta, nil, gcs, loc, "2026-09-28", 1)
	assert.NoError(t, err)

	h := readHistory(t, gcs, provider.Beta).Tests["TestAccA"]
	if assert.NotNil(t, h) {
		assert.Equal(t, "", h.TestNameId)
	}
}
