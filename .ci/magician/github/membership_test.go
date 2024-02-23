/*
* Copyright 2023 Google LLC. All Rights Reserved.
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
package github

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestTrustedContributors(t *testing.T) {
	for _, member := range trustedContributors {
		if slices.Contains(reviewerRotation, member) {
			t.Fatalf(`%v should not be on reviewerRotation list`, member)
		}
	}
}

func TestOnVacationReviewers(t *testing.T) {
	for _, member := range onVacationReviewers {
		if !slices.Contains(reviewerRotation, member) {
			t.Fatalf(`%v is not on reviewerRotation list`, member)
		}
	}
}
