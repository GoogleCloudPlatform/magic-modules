/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"magician/github"
 	"math/rand"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"
)

// requestServiceReviewersCmd represents the requestServiceReviewers command
var requestServiceReviewersCmd = &cobra.Command{
	Use:   "request-service-reviewers PR_NUMBER",
	Short: "Assigns reviewers based on the PR's service labels.",
	Long: `This command requests (or re-requests) review based on the PR's service labels.

	If a PR has more than 3 service labels, the command will not do anything.
	`,
	Args: cobra.ExactArgs(6),
	Run: func(cmd *cobra.Command, args []string) {
		prNumber := args[0]
		fmt.Println("PR Number: ", prNumber)

		gh := github.NewGithubService()
		execRequestServiceReviewers(prNumber, gh, labeler.EnrolledTeamsYaml)
	},
}

func execRequestServiceReviewers(prNumber string, gh github.GithubService, enrolledTeamsYaml []byte) {
	pullRequest, err := gh.GetPullRequest(prNumber)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	enrolledTeams := make(map[string]labeler.LabelData)
	if err := yaml.Unmarshal(enrolledTeamsYaml, &enrolledTeams); err != nil {
		fmt.Printf("Error unmarshalling enrolled teams yaml: %s", err)
		os.Exit(1)
	}

	previousReviewers, err := gh.GetPullRequestPreviousAssignedReviewers(prNumber)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// If more than three service labels are impacted, don't request reviews.
	// Only request reviews from unique service teams.
	githubTeams := make(map[string]struct{})
	teamCount := 0
	for _, label := range pullRequest.Labels {
		if !strings.HasPrefix(label.Name, "service/") || label.Name == "service/terraform" {
			continue
		}
		teamCount += 1
		if labelData, ok := enrolledTeams[label]; ok {
			githubTeams[labelData] = struct{}{}
		}
	}

	if teamCount > 3 {
		fmt.Println("Provider-wide change (>3 services impacted); not requesting service team reviews")
		return
	}

	// For each service team, check if there is already a reviewer in the team list. Rerequest
	// review if there is and choose a random reviewer from the list if there isn't.
	previousReviewersMap = make(map[string]struct{})
	toRequest := []string{}
	for _, reviewer := range previousReviewers {
		previousReviewersMap[reviewer] = struct{}{}
	}

	exitCode := 0
	for githubTeam, _ := range githubTeams {
		members, err := gh.GetTeamMembers("GoogleCloudPlatform", githubTeam)
		if err != nil {
			fmt.Printf("Error fetching members for GoogleCloudPlatform/%s: %s", githubTeam, err)
			exitCode = 1
			continue
		}
		hasReviewer := false
		for _, member := range members {
			if previousReviewersMap[member] {
				hasReviewer = true
				toRequest = append(toRequest, member)
			}
		}

		if !hasReviewer {
			toRequest = append(toRequest, members[rand.Intn(len(members))])
		}
	}

	for _, reviewer := range toRequest {
		err = gh.RequestPullRequestReviewer(prNumber, reviewer)
		if err != nil {
			fmt.Println(err)
			exitCode = 1
		}
	}
	if exitCode != 0 {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(requestServiceReviewersCmd)
}
