// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
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

	if err := entry.Validate(); err != nil {
		log.Printf("error parsing changelog entry in %s: %s", entry.Issue, err)
		switch err.Code {
		case changelog.EntryErrorNotFound:
			body := "Oops! It looks like no changelog entry is attached to" +
				" this PR. Please include a release note block" +
				" in the PR body, as described in https://github.com/GoogleCloudPlatform/magic-modules/blob/master/.ci/RELEASE_NOTES_GUIDE.md:" +
				"\n\n~~~\n```release-note:TYPE\nRelease note" +
				"\n```\n~~~"
			_, _, err := client.Issues.CreateComment(ctx, owner, repo,
				prNo, &github.IssueComment{
					Body: &body,
				})
			if err != nil {
				log.Fatalf("Error creating pull request comment on"+
					" github.com/%s/%s/%d: %s", owner, repo, prNo,
					err)
			}
			os.Exit(1)
		case changelog.EntryErrorUnknownTypes:
			unknownTypes := err.Details["unknownTypes"].([]string)

			body := "Oops! It looks like you're using"
			if len(unknownTypes) == 1 {
				body += " an"
			}
			body += " unknown release-note type"
			if len(unknownTypes) > 1 {
				body += "s"
			}
			body += " in your changelog entries:"
			for _, t := range unknownTypes {
				body += "\n* " + t
			}
			body += "\n\nPlease only use the types listed in https://github.com/GoogleCloudPlatform/magic-modules/blob/master/.ci/RELEASE_NOTES_GUIDE.md."
			_, _, err := client.Issues.CreateComment(ctx, owner, repo,
				prNo, &github.IssueComment{
					Body: &body,
				})
			if err != nil {
				log.Fatalf("Error creating pull request comment on"+
					" github.com/%s/%s/%d: %s", owner, repo, prNo,
					err)
			}
			os.Exit(1)
		}
	}
}
