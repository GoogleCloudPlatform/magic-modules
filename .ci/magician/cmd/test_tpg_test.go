package cmd

import (
	"testing"
)

func TestExecTestTPG(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string]bool),
	}

	execTestTPG("beta", "sha1", "pr1", gh)

	if !gh.calledMethods["CreateWorkflowDispatchEvent"] {
		t.Fatal("workflow dispatch event not created")
	}
}
