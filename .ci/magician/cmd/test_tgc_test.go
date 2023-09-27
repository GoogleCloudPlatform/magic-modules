package cmd

import (
	"reflect"
	"testing"
)

func TestExecTestTGC(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string][]any),
	}

	execTestTGC("sha1", "pr1", gh)

	method := "CreateWorkflowDispatchEvent"
	expected := []any{"test-tgc.yml", map[string]any{"branch": "auto-pr-pr1", "owner": "modular-magician", "repo": "terraform-google-conversion", "sha": "sha1"}}
	if params, ok := gh.calledMethods[method]; !ok {
		t.Fatal("workflow dispatch event not created")
	} else if !reflect.DeepEqual(params, expected) {
		t.Fatalf("wrong params for %s, got %v, expected %v", method, params, expected)
	}
}
