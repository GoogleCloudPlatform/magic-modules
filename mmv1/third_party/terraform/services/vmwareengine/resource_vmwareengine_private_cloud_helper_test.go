package vmwareengine_test

import (
	"os"
	"testing"
)

func getTestNodeType() string {
	if nt := os.Getenv("GOOGLE_VMWAREENGINE_NODE_TYPE"); nt != "" {
		return nt
	}
	return "standard-72"
}

func skipIfLimitedNodes(t *testing.T) {
	if isLimitedNodes() {
		t.Skip("Skipping test because it requires more nodes than available in the limited nodes environment")
	}
}

func isLimitedNodes() bool {
	return os.Getenv("GOOGLE_VMWAREENGINE_LIMITED_NODES") == "true"
}


func getTestRegion() string {
	if r := os.Getenv("GOOGLE_VMWAREENGINE_REGION"); r != "" {
		return r
	}
	return "me-west1"
}

func isProjectCreationDisabled() bool {
	return os.Getenv("GOOGLE_VMWAREENGINE_PROJECT_CREATION_DISABLED") == "true"
}
