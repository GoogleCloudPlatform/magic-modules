package cmd

import (
	"magician/exec"
	"magician/source"
	"testing"
)

func TestFetchPRNumber(t *testing.T) {
	rnr, err := exec.NewRunner()
	if err != nil {
		t.Errorf("error creating Runner: %s", err)
	}

	ctlr := source.NewController("", "modular-magician", "", rnr)

	prNumber, err := fetchPRNumber("8c6e61bb62d52c950008340deafc1e2a2041898a", "main", rnr, ctlr)

	if err != nil {
		t.Errorf("error fetching PR number: %s", err)
	}

	if prNumber != "6504" {
		t.Errorf("PR number is %s, expected 6504", prNumber)
	}
}
