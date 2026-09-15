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

package google

import (
	"log"
	"strings"
	"sync/atomic"
	"time"
)

var (
	VerboseLogging         bool
	totalResourceCount     atomic.Int32
	completedResourceCount atomic.Int32
	lastReportNano         atomic.Int64
)

const progressInterval = 2 * time.Second

// LogVerbose prints to log only when VerboseLogging is enabled.
func LogVerbose(format string, v ...any) {
	if VerboseLogging {
		log.Printf(format, v...)
	}
}

// InitProgress sets up the progress tracker for resource generation.
func InitProgress(total int) {
	totalResourceCount.Store(int32(total))
	completedResourceCount.Store(0)
	lastReportNano.Store(time.Now().UnixNano())
}

// IncrementResourceGenerated records progress and prints log messages periodically with UX progress bar and percentage.
func IncrementResourceGenerated() {
	completed := completedResourceCount.Add(1)
	total := totalResourceCount.Load()

	if total < 10 || VerboseLogging {
		return
	}

	now := time.Now().UnixNano()
	last := lastReportNano.Load()

	if time.Duration(now-last) >= progressInterval {
		if lastReportNano.CompareAndSwap(last, now) {
			percent := float64(completed) / float64(total) * 100.0
			filled := int((percent / 100.0) * 20)
			if filled > 20 {
				filled = 20
			}

			var bar string
			if filled == 0 {
				bar = strings.Repeat(".", 20)
			} else if filled == 20 {
				bar = strings.Repeat("=", 20)
			} else {
				bar = strings.Repeat("=", filled-1) + ">" + strings.Repeat(".", 20-filled)
			}

			log.Printf("[%s] %.1f%% (%d/%d resources)", bar, percent, completed, total)
		}
	}
}
