package cmd

import (
	"testing"
)

func TestFetchPRNumber(t *testing.T) {
	mr := NewMockRunner()
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
		commitMessage: "Add `additional_group_keys` attribute to `google_cloud_identity_group` resource (#9217) (#6504)\n\n* Add `additional_group_keys` attribute to `google_cloud_identity_group` resource\n\n* Update acceptance test to check for attribute\n\n* Fix test check\n\n* Add `output: true` to nested properties in output field\n[upstream:49d3741f9d4d810a0a4768363bb8498afa21c688]\n\nSigned-off-by: Modular Magician <magic-modules@google.com>",
	}

	// Call function with mocks
	prNumber, err := fetchPRNumber("8c6e61bb62d52c950008340deafc1e2a2041898a", "main", mr, gh)
	if err != nil {
		t.Errorf("error fetching PR number: %s", err)
	}

	if prNumber != "6504" {
		t.Errorf("PR number is %s, expected 6504", prNumber)
	}

	// Verify GitHub API was called
	if calls, ok := gh.calledMethods["GetCommitMessage"]; !ok || len(calls) == 0 {
		t.Errorf("Expected GetCommitMessage to be called")
	} else {
		args := calls[0]
		if len(args) != 3 {
			t.Errorf("Expected GetCommitMessage to be called with 3 arguments, got %d", len(args))
		} else {
			if args[0] != "hashicorp" {
				t.Errorf("Expected owner to be 'hashicorp', got '%s'", args[0])
			}
			if args[1] != "terraform-provider-google-beta" {
				t.Errorf("Expected repo to be 'terraform-provider-google-beta', got '%s'", args[1])
			}
			if args[2] != "8c6e61bb62d52c950008340deafc1e2a2041898a" {
				t.Errorf("Expected SHA to be '8c6e61bb62d52c950008340deafc1e2a2041898a', got '%s'", args[2])
			}
		}
	}
}
