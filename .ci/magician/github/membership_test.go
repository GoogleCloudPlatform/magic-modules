package github

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestTrustedContributors(t *testing.T) {
	for _, member := range trustedContributors {
		if slices.Contains(reviewerRotation, member) {
			t.Fatalf(`%v should not be on reviewerRotation list`, member)
		}
	}
}

func TestOnVacationReviewers(t *testing.T) {
	for _, member := range onVacationReviewers {
		if !slices.Contains(reviewerRotation, member) {
			t.Fatalf(`%v is not on reviewerRotation list`, member)
		}
	}
}
