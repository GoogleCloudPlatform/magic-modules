package cmd

type mockCloudBuild struct {
	calledMethods map[string]bool
}

func (m *mockCloudBuild) ApproveCommunityChecker(prNumber, commitSha string) error {
	m.calledMethods["ApproveCommunityChecker"] = true
	return nil
}

func (m *mockCloudBuild) GetAwaitingApprovalBuildLink(prNumber, commitSha string) (string, error) {
	m.calledMethods["GetAwaitingApprovalBuildLink"] = true
	return "mocked_url", nil
}

func (m *mockCloudBuild) TriggerMMPresubmitRuns(commitSha string, substitutions map[string]string) error {
	m.calledMethods["TriggerMMPresubmitRuns"] = true
	return nil
}
