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
		body := "\nOops! Some errors are detected for your changelog entries:\n"
		for i, err := range errors {
			body += fmt.Sprintf("\n* Issue %d\n", i+1)
			if err.Details != nil {
				body += fmt.Sprintf("Changelog:\n```release-note:%v\n%v\n```\n", err.Details["type"].(string), err.Details["note"].(string))
			}
			body += "Errors:\n"
			switch err.Code {
			case changelog.EntryErrorNotFound:
				body += "- It looks like no changelog entry is attached to this PR. Please include a release note block in the PR body, as described in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/.\n\n"
			case changelog.EntryErrorUnknownTypes:
				body += "- Unknown changelog types\nPlease only use the types listed in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/.\n\n"
			case changelog.EntryErrorMultipleLines:
				body += "- Multiple lines are found in changelog entry \nPlease only have one CONTENT line per release note block. Use multiple blocks if there are multiple related changes in a single PR.\n\n"
			case changelog.EntryErrorInvalidNewReourceOrDatasourceFormat:
				body += "- Invalid resource/datasource format\nPlease follow format in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/#type-specific-guidelines-and-examples.\n\n"
			case changelog.EntryErrorInvalidEnhancementOrBugFixFormat:
				body += "- Invalid enhancement/bug fix format\nPlease follow format in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/#type-specific-guidelines-and-examples.\n\n"
			}
		}
		log.Fatal(body)
	}
}
