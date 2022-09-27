package rules

import (
	"fmt"

	"github.com/GoogleCloudPlatform/magic-modules/.ci/breaking-change-detector/constants"
)

func documentationReference(version, identifier string) string {
	return fmt.Sprintf(" - [reference](%s)", constants.GetFileUrl(version, identifier))
}
