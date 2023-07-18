package constants

const BreakingChangeRelativeLocation = "reference/"
const BreakingChangeFileName = "breaking-change-detector"

var docsite = "https://googlecloudplatform.github.io/magic-modules/"

func GetFileUrl(identifier string) string {
	return docsite + BreakingChangeRelativeLocation + BreakingChangeFileName + "#" + identifier
}
