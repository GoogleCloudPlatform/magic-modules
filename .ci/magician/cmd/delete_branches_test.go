package cmd

import (
	"magician/exec"
	"magician/github"
	"os"
	"testing"
)

func TestFetchPRNumber(t *testing.T) {
	rnr, err := exec.NewRunner()
	if err != nil {
		t.Errorf("error creating Runner: %s", err)
	}

	githubToken, ok := os.LookupEnv("GITHUB_TOKEN_CLASSIC")
	if !ok {
		t.Errorf("did not provide GITHUB_TOKEN_CLASSIC environment variable")
	}

	gh := github.NewClient(githubToken)

	prNumber, err := fetchPRNumber("8c6e61bb62d52c950008340deafc1e2a2041898a", "main", rnr, gh)

	if err != nil {
		t.Errorf("error fetching PR number: %s", err)
	}

	if prNumber != "6504" {
		t.Errorf("PR number is %s, expected 6504", prNumber)
	}
}
