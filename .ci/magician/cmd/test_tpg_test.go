package cmd

import (
	"reflect"
	"testing"
)

func TestExecTestTPG(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}

	execTestTPG("beta", "sha1", "pr1", gh)

	method := "CreateWorkflowDispatchEvent"
	expected := [][]any{{"test-tpg.yml", map[string]any{"branch": "auto-pr-pr1", "owner": "modular-magician", "repo": "terraform-provider-google-beta", "sha": "sha1"}}}
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("Workflow dispatch event not created")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

	gh.calledMethods = make(map[string][][]any)

	execTestTPG("ga", "sha1", "pr1", gh)

	method = "CreateWorkflowDispatchEvent"
	expected = [][]any{{"test-tpg.yml", map[string]any{"branch": "auto-pr-pr1", "owner": "modular-magician", "repo": "terraform-provider-google", "sha": "sha1"}}}
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("Workflow dispatch event not created")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}
}
