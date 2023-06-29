package constants

const BreakingChangeRelativeLocation = "/.github/"
const BreakingChangeFileName = "BREAKING_CHANGES.md"

var docsite = "https://googlecloudplatform.github.io/magic-modules"

func GetFileUrl(identifier string) string {
	return docsite + BreakingChangeRelativeLocation + BreakingChangeFileName + "#" + identifier
}
