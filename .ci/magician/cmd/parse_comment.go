package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"

	"magician/github"

	"github.com/spf13/cobra"
)

// This regex captures the entire line starting with @modular-magician
// Example: "@modular-magician reassign-reviewer user1" or "@modular-magician assign review @user2"
var magicianInvocationRegex = regexp.MustCompile(`@modular-magician\s+([^\n\r]+)`)

// Command patterns for reassign-reviewer with flexible syntax
// Supports: assign-reviewer, reassign-reviewer, assign reviewer, reassign review, etc.
// Captures only valid GitHub usernames: [a-zA-Z0-9-_]
var reassignReviewerRegex = regexp.MustCompile(`^(?:re)?assign[- ]?review(?:er)?\s*@?([a-zA-Z0-9-_]*)`)

var parseCommentCmd = &cobra.Command{
	Use:   "parse-comment PR_NUMBER COMMENT_AUTHOR BASE64_COMMENT",
	Short: "Parses a base64 encoded comment to execute magician commands",
	Long: `This command parses GitHub PR comments for @modular-magician invocations.
	
	It supports flexible command syntax including:
	- Commands with hyphens: reassign-reviewer
	- Commands with spaces: reassign reviewer
	- Optional prefixes and suffixes: assign-review, reassign-reviewer
	- Optional @ prefix for usernames
	
	The command expects:
	1. PR_NUMBER - The pull request number
	2. COMMENT_AUTHOR - The GitHub username who made the comment
	3. BASE64_COMMENT - The base64 encoded comment body`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		prNumber := args[0]
		author := args[1]
		base64Comment := args[2]

		githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN environment variable")
		}
		gh := github.NewClient(githubToken)

		if gh.GetUserType(author) != github.CoreContributorUserType {
			return fmt.Errorf("comment author %s is not a core contributor", author)
		}

		commentBytes, err := base64.StdEncoding.DecodeString(base64Comment)
		if err != nil {
			return fmt.Errorf("failed to decode base64 comment: %w", err)
		}
		comment := string(commentBytes)

		return execParseComment(prNumber, comment, gh)
	},
}

// execParseComment is the main router that finds and executes the first command
func execParseComment(prNumber, comment string, gh GithubClient) error {
	// Find the first @modular-magician invocation in the comment
	match := magicianInvocationRegex.FindStringSubmatch(comment)

	if match == nil {
		fmt.Println("No @modular-magician invocation found. Ignoring comment.")
		return nil
	}

	if len(match) < 2 {
		fmt.Printf("Invalid match structure. Ignoring.\n")
		return nil
	}

	commandLine := strings.TrimSpace(match[1])
	if commandLine == "" {
		fmt.Printf("Empty command after @modular-magician. Ignoring.\n")
		return nil
	}

	fmt.Printf("Processing command: %q\n", commandLine)

	// Route to appropriate handler based on command pattern
	return routeCommand(prNumber, commandLine, gh)
}

// routeCommand determines which command handler to call based on the command pattern
func routeCommand(prNumber, commandLine string, gh GithubClient) error {
	// Check for reassign-reviewer command variants
	if matches := reassignReviewerRegex.FindStringSubmatch(commandLine); matches != nil {
		reviewer := strings.TrimSpace(matches[1])
		return handleReassignReviewer(prNumber, reviewer, gh)
	}

	// Add more command patterns here as needed
	// Example for future commands:
	// if matches := cherryPickRegex.FindStringSubmatch(commandLine); matches != nil {
	//     return handleCherryPick(prNumber, matches[1:], gh)
	// }

	fmt.Printf("Unknown command format: %q\n", commandLine)
	return nil
}

// handleReassignReviewer processes the reassign-reviewer command
func handleReassignReviewer(prNumber, reviewer string, gh GithubClient) error {
	// The regex already extracted just the username without @
	// and only allows valid GitHub username characters [a-zA-Z0-9-_]

	fmt.Printf("Reassigning reviewer for PR #%s", prNumber)
	if reviewer != "" {
		fmt.Printf(" to @%s", reviewer)
	} else {
		fmt.Printf(" (selecting random reviewer)")
	}
	fmt.Println()

	// Call the existing reassign reviewer logic
	return execReassignReviewer(prNumber, reviewer, gh)
}

func init() {
	rootCmd.AddCommand(parseCommentCmd)
}
