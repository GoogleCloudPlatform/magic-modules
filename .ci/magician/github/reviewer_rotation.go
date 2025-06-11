package github

import (
	"fmt"
	utils "magician/utility"
	"math/rand"
	"slices"
	"time"

	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type ReviewerRotation map[string]*ReviewerConfig

func (rr *ReviewerRotation) setStartEnd() {
	for _, config := range *rr {
		config.setStartEnd()
	}
}

func (rr *ReviewerRotation) read(data []byte) error {
	if err := yaml.Unmarshal(data, rr); err != nil {
		return err
	}
	return nil
}

func (rr *ReviewerRotation) write() ([]byte, error) {
	data, err := yaml.Marshal(rr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (rr *ReviewerRotation) UnmarshalYAML(value *yaml.Node) error {
	if len(value.Content)%2 != 0 {
		return fmt.Errorf("reviewer rotation map content must be even, got %d", len(value.Content))
	}
	partial := make(map[string]*yaml.Node, len(value.Content)/2)
	// Iterate through content of value.
	for i := 0; i < len(value.Content); i += 2 {
		keyNode := value.Content[i]
		valueNode := value.Content[i+1]
		partial[keyNode.Value] = valueNode
	}

	for reviewer, node := range partial {
		if err := node.Decode((*rr)[reviewer]); err != nil {
			return err
		}
	}
	return nil
}

// isCoreReviewer returns true if the given user is a core reviewer
func (rr ReviewerRotation) isCoreReviewer(user string) bool {
	_, isCoreReviewer := rr[user]
	return isCoreReviewer
}

// getRandomReviewer returns a random available reviewer (optionally excluding some people from the reviewer pool)
func (rr ReviewerRotation) getRandomReviewer(excludedReviewers []string) string {
	availableReviewers := rr.availableReviewers(excludedReviewers)
	reviewer := availableReviewers[rand.Intn(len(availableReviewers))]
	return reviewer
}

func (rr ReviewerRotation) availableReviewers(excludedReviewers []string) []string {
	return rr.available(time.Now(), excludedReviewers)
}

func (rr ReviewerRotation) available(nowTime time.Time, excludedReviewers []string) []string {
	excludedReviewers = append(excludedReviewers, rr.onVacation(nowTime)...)
	ret := utils.Removes(maps.Keys(rr), excludedReviewers)
	slices.Sort(ret)
	return ret
}

// onVacation returns a list of reviewers who are on vacation at the given time
func (rr ReviewerRotation) onVacation(now time.Time) []string {
	var onVacationList []string
	for reviewer, config := range rr {
		if config.onVacation(now) {
			onVacationList = append(onVacationList, reviewer)
		}
	}
	return onVacationList
}
