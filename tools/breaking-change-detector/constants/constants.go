package constants

const BreakingChangeRelativeLocation = "/.github/"
const BreakingChangeFileName = "BREAKING_CHANGES.md"

var providerUrls = map[string]string{
	"google":      "https://github.com/hashicorp/terraform-provider-google/blob/main",
	"google-beta": "https://github.com/hashicorp/terraform-provider-google-beta/blob/main",
}

func GetFileUrl(version, identifier string) string {
	return providerUrls[version] + BreakingChangeRelativeLocation + BreakingChangeFileName + "#" + identifier
}
