// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-changelog"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	if len(os.Args) < 2 {
		log.Fatalf("Usage: changelog-pr-body-check PR#\n")
	}
	pr := os.Args[1]
	prNo, err := strconv.Atoi(pr)
	if err != nil {
		log.Fatalf("Error parsing PR %q as a number: %s", pr, err)
	}

	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")
	token := os.Getenv("GITHUB_TOKEN")

	if owner == "" {
		log.Fatalf("GITHUB_OWNER not set")
	}
	if repo == "" {
		log.Fatalf("GITHUB_REPO not set")
	}
	if token == "" {
		log.Fatalf("GITHUB_TOKEN not set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	pullRequest, _, err := client.PullRequests.Get(ctx, owner, repo, prNo)
	if err != nil {
		log.Fatalf("Error retrieving pull request github.com/"+
			"%s/%s/%d: %s", owner, repo, prNo, err)
	}

	entry := changelog.Entry{
		Issue: pr,
		Body:  pullRequest.GetBody(),
	}

	if errors := entry.Validate(); errors != nil {
		log.Printf("error parsing changelog entry in %s: %s", entry.Issue, errors)
		body := "\nOops! Some errors are detected for your changelog entries:"
		for i, err := range errors {
			switch err.Code {
			case changelog.EntryErrorNotFound:
				body += fmt.Sprintf("\n\n Issue %d: It looks like no changelog entry is attached to this PR. Please include a release note block in the PR body, as described in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/", i+1)
			case changelog.EntryErrorUnknownTypes:
				body += fmt.Sprintf("\n\n Issue %d: unknown changelog types %v \nPlease only use the types listed in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/.", i+1, err.Details["type"].(string))
			case changelog.EntryErrorMultipleLines:
				body += fmt.Sprintf("\n\n Issue %d: multiple lines are found in changelog entry: %v \nPlease only have one CONTENT line per release note block. Use multiple blocks if there are multiple related changes in a single PR.", i+1, err.Details["note"].(string))
			case changelog.EntryErrorInvalidNewReourceFormat:
				body += fmt.Sprintf("\n\n Issue %d: invalid resource/datasource format in changelog entry: %v \nPlease follow format in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/#type-specific-guidelines-and-examples.", i+1, err.Details["note"].(string))
			case changelog.EntryErrorInvalidEnhancementOrBugFixFormat:
				body += fmt.Sprintf("\n\n Issue %d: invalid enhancement/bug fix format in changelog entry: %v \nPlease follow format in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/#type-specific-guidelines-and-examples.", i+1, err.Details["note"].(string))
			}
		}
		log.Fatal(body)
	}
}
