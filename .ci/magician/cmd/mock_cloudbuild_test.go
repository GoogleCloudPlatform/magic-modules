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
