//go:build integration
// +build integration

/*
* Copyright 2025 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
* Integration tests - makes real GitHub API calls.
* NOT run during normal test execution (go test).
* Requires:
*   - GITHUB_API_TOKEN environment variable
*   - Run with: go test -tags=integration
 */

package github

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// https://github.com/GoogleCloudPlatform/magic-modules
const (
	testNonMember     = "bananaman5000"
	testRepo          = "magic-modules"
	testOwner         = "GoogleCloudPlatform"
	testPRNumber      = "13969"                                    // replace this with an actual PR Number
	testPRCommitSha   = "4a8409686810551655eea2533e939cc5344e83e2" // replace this with an actual SHA
	testMainCommitSha = "fd910977cf24595d2c04e3f0a369a82c79fdb8f8" // replace this with an actual SHA
	testLabel         = "terraform-3.0"
	testOrg           = "GoogleCloudPlatform"
	testTeam          = "terraform"
	workflowFileName  = "test-tpg.yml"
)

func skipIfNoToken(t *testing.T) *Client {
	token := os.Getenv("GITHUB_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: GITHUB_API_TOKEN environment variable not set")
	}

	return NewClient(token)
}

func TestIntegrationGetPullRequest(t *testing.T) {
	client := skipIfNoToken(t)

	pr, err := client.GetPullRequest(testPRNumber)
	if err != nil {
		t.Fatalf("GetPullRequest failed: %v", err)
	}

	t.Logf("PR Title: %s", pr.Title)
	if pr.Number == 0 {
		t.Error("Expected PR number to be non-zero")
	}
	if pr.Title == "" {
		t.Error("Expected PR title to be non-empty")
	}
}

func TestIntegrationGetPullRequests(t *testing.T) {
	client := skipIfNoToken(t)

	prs, err := client.GetPullRequests("open", "main", "created", "desc")
	if err != nil {
		t.Fatalf("GetPullRequests failed: %v", err)
	}

	t.Logf("Found %d PRs", len(prs))
}

func TestIntegrationGetCommitMessage(t *testing.T) {
	client := skipIfNoToken(t)

	// You'll need a valid commit SHA for this test
	if testMainCommitSha == "HEAD" {
		t.Skip("Skipping GetCommitMessage test: need a valid commit SHA")
	}

	message, err := client.GetCommitMessage(testOwner, testRepo, testMainCommitSha)
	if err != nil {
		t.Fatalf("GetCommitMessage failed: %v", err)
	}

	t.Logf("Commit message: %s", message)
	if message == "" {
		t.Error("Expected commit message to be non-empty")
	}
}

func TestIntegrationGetPullRequestComments(t *testing.T) {
	client := skipIfNoToken(t)

	comments, err := client.GetPullRequestComments(testPRNumber)
	if err != nil {
		t.Fatalf("GetPullRequestComments failed: %v", err)
	}

	t.Logf("Found %d comments", len(comments))
	for i, comment := range comments {
		t.Logf("Comment %d: %s by %s", i+1, comment.Body[:min(len(comment.Body), 50)], comment.User.Login)
	}
}

func TestIntegrationGetTeamMembers(t *testing.T) {
	client := skipIfNoToken(t)

	members, err := client.GetTeamMembers(testOrg, testTeam)
	if err != nil {
		t.Logf("GetTeamMembers failed: %v", err)
		t.Skip("Skipping team member test - might not have access to the specified team")
	}

	t.Logf("Found %d team members", len(members))
	for i, member := range members {
		t.Logf("Member %d: %s", i+1, member.Login)
	}
}

func TestIntegrationIsOrgMember(t *testing.T) {
	client := skipIfNoToken(t)

	isMember := client.IsOrgMember(testOwner, testOrg)
	t.Logf("Is %s a member of %s: %v", testOwner, testOrg, isMember)

	if !isMember {
		t.Errorf("Note: Expected %s to be a member of %s, but they're not", testOwner, testOrg)
	}

	isMember = client.IsOrgMember(testNonMember, testOrg)
	if isMember {
		t.Errorf("Expected %s to not be a member of %s, but they are", testNonMember, testOrg)
	} else {
		t.Logf("Is %s not a member of %s: %v", testNonMember, testOrg, isMember)
	}
}

func TestIntegrationIsTeamMember(t *testing.T) {
	client := skipIfNoToken(t)

	isMember := client.IsTeamMember(testOrg, testTeam, testOwner)
	if !isMember {
		t.Errorf("Expected %s to be a member of team %s in org %s, but they're not", testOwner, testTeam, testOrg)
	} else {
		t.Logf("Is %s a member of team %s in org %s: %v", testOwner, testTeam, testOrg, isMember)
	}

	isMember = client.IsTeamMember(testOrg, testTeam, testNonMember)
	if isMember {
		t.Errorf("Expected %s to not be a member of team %s in org %s, but they are", testNonMember, testTeam, testOrg)
	} else {
		t.Logf("Is %s not a member of team %s in org %s: %v", testNonMember, testTeam, testOrg, isMember)
	}
}

func TestIntegrationPostAndUpdateComment(t *testing.T) {
	client := skipIfNoToken(t)

	// First post a comment
	comment := fmt.Sprintf("Test comment from integration test at %s", time.Now().Format(time.RFC3339))
	err := client.PostComment(testPRNumber, comment)
	if err != nil {
		t.Fatalf("PostComment failed: %v", err)
	}

	// Get the comment ID
	comments, err := client.GetPullRequestComments(testPRNumber)
	if err != nil {
		t.Fatalf("GetPullRequestComments failed: %v", err)
	}

	var commentID int
	for _, c := range comments {
		if c.Body == comment {
			commentID = c.ID
			break
		}
	}

	if commentID == 0 {
		t.Fatal("Could not find the comment we just posted")
	}

	// Update the comment
	updatedComment := fmt.Sprintf("Updated test comment from integration test at %s", time.Now().Format(time.RFC3339))
	err = client.UpdateComment(testPRNumber, updatedComment, commentID)
	if err != nil {
		t.Fatalf("UpdateComment failed: %v", err)
	}

	t.Logf("Successfully posted and updated comment with ID: %d", commentID)
}

func TestIntegrationAddAndRemoveLabels(t *testing.T) {
	client := skipIfNoToken(t)

	err := client.AddLabels(testPRNumber, []string{testLabel})
	if err != nil {
		t.Fatalf("AddLabels failed: %v", err)
	}

	// Then remove the label
	err = client.RemoveLabel(testPRNumber, testLabel)
	if err != nil {
		t.Fatalf("RemoveLabel failed: %v", err)
	}

	t.Logf("Successfully added and removed label: %s", testLabel)
}

func TestIntegrationPostBuildStatus(t *testing.T) {
	client := skipIfNoToken(t)

	// You'll need a valid commit SHA for this test
	if testPRCommitSha == "HEAD" {
		t.Skip("Skipping PostBuildStatus test: need a valid commit SHA")
	}

	err := client.PostBuildStatus(
		testPRNumber,
		"integration-test",
		"success",
		"https://example.com/integration-test",
		testPRCommitSha,
	)
	if err != nil {
		t.Errorf("PostBuildStatus failed: %v", err)
	}

	err = client.PostBuildStatus(
		testPRNumber,
		"integration-test-failed",
		"failure",
		"https://example.com/integration-test-fail",
		testPRCommitSha,
	)
	if err != nil {
		t.Errorf("PostBuildStatus failed: %v", err)
	}

	t.Logf("Successfully posted build status")
}

func TestIntegrationCreateWorkflowDispatchEvent(t *testing.T) {
	client := skipIfNoToken(t)

	// Skip this test by default as it can have side effects
	if os.Getenv("RUN_WORKFLOW_DISPATCH_TEST") != "true" {
		t.Skip("Skipping workflow dispatch test: set RUN_WORKFLOW_DISPATCH_TEST=true to run")
	}

	if err := client.CreateWorkflowDispatchEvent("test-tpg.yml", map[string]any{
		"owner":     "modular-magician",
		"repo":      testRepo,
		"branch":    "main",
		"pr-number": testPRNumber,
		"sha":       testPRCommitSha,
	}); err != nil {
		t.Errorf("error creating workflow dispatch event: %v", err)
	}

	t.Logf("Successfully triggered workflow dispatch event")
}

// TestIntegrationMergePullRequest is commented out as it has permanent effects
// Uncomment and run only if you're sure you want to merge the PR
/*
 func TestIntegrationMergePullRequest(t *testing.T) {
	 client := skipIfNoToken(t)

	 // Skip this test by default as it has permanent effects
	 if os.Getenv("RUN_MERGE_PR_TEST") != "true" {
		 t.Skip("Skipping merge PR test: set RUN_MERGE_PR_TEST=true to run")
	 }

	 // You'll need a valid commit SHA for this test
	 if testPRCommitSha == "HEAD" {
		 t.Skip("Skipping MergePullRequest test: need a valid commit SHA")
	 }

	 err := client.MergePullRequest(testOwner, testRepo, testPRNumber, testPRCommitSha)
	 if err != nil {
		 t.Fatalf("MergePullRequest failed: %v", err)
	 }

	 t.Logf("Successfully merged pull request")
 }
*/

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
