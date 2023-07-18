package constants

const BreakingChangeRelativeLocation = "develop/"
const BreakingChangeFileName = "make-a-breaking-change"

var docsite = "https://googlecloudplatform.github.io/magic-modules/"

func GetFileUrl(identifier string) string {
	return docsite + BreakingChangeRelativeLocation + BreakingChangeFileName + "#" + identifier
}
