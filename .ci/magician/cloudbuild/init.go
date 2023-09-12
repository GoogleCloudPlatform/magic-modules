package cloudbuild

type cloudBuild bool

type CloudBuild interface {
	ApproveCommunityChecker(prNumber, commitSha string) error
	GetAwaitingApprovalBuildLink(prNumber, commitSha string) (string, error)
	TriggerMMPresubmitRuns(commitSha string, substitutions map[string]string) error
}

func NewCloudBuildService() CloudBuild {
	var x cloudBuild = true
	return x
}
