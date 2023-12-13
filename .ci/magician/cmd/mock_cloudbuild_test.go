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
package cmd

type mockCloudBuild struct {
	calledMethods map[string][][]any
}

func (m *mockCloudBuild) ApproveCommunityChecker(prNumber, commitSha string) error {
	m.calledMethods["ApproveCommunityChecker"] = append(m.calledMethods["ApproveCommunityChecker"], []any{prNumber, commitSha})
	return nil
}

func (m *mockCloudBuild) GetAwaitingApprovalBuildLink(prNumber, commitSha string) (string, error) {
	m.calledMethods["GetAwaitingApprovalBuildLink"] = append(m.calledMethods["GetAwaitingApprovalBuildLink"], []any{prNumber, commitSha})
	return "mocked_url", nil
}

func (m *mockCloudBuild) TriggerMMPresubmitRuns(commitSha string, substitutions map[string]string) error {
	m.calledMethods["TriggerMMPresubmitRuns"] = append(m.calledMethods["TriggerMMPresubmitRuns"], []any{commitSha, substitutions})
	return nil
}
