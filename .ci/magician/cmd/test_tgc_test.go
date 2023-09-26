package cmd

import (
	"testing"
)

func TestExecTestTGC(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string]bool),
	}

	execTestTGC("sha1", "pr1", gh)

	if !gh.calledMethods["CreateWorkflowDispatchEvent"] {
		t.Fatal("workflow dispatch event not created")
	}
}
