# GitHub Integration Tests

## Overview
This directory contains an interface for the GitHub client that make real API calls to GitHub's API. The tests in `integration_test.go` are isolated with build tags to prevent accidental execution.

## Build Tags
This file uses Go build tags (`//go:build integration`) which:
- Exclude these tests from normal test execution (`go test ./...`)
- Require explicit opt-in (`go test -tags=integration`)
- Prevent accidental execution of tests that make real API calls and may have side effects

## Usage

### Requirements
- GitHub API token with appropriate permissions
- Token set as environment variable: `GITHUB_API_TOKEN`

### Running Tests
```bash
# Run all integration tests
GITHUB_API_TOKEN=your_token_here go test -v -tags=integration ./github

# Run specific test
GITHUB_API_TOKEN=your_token_here go test -v -tags=integration -run TestIntegrationGetPullRequest ./github