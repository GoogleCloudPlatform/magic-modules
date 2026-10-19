/*
* Copyright 2026 Google LLC. All Rights Reserved.
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
 */
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"magician/exec"
	"magician/github"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultClassifierURL       = "https://us-central1-aiplatform.googleapis.com/v1/projects/673497134629/locations/us-central1/publishers/google/models/gemini-2.5-flash:generateContent"
	defaultMetadataTokenURL    = "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token"
	defaultCloudBuildProjectID = "graphite-docker-images"
	classifierStatusContext    = "pr-type-classifier"
	classifierConfidenceFloor  = 0.50
	maxChangedFilesInPrompt    = 25
	maxBodyCharsInPrompt       = 4000
	maxPatchCharsInPrompt      = 10000
)

var validPRTypeLabels = []string{
	"type/new-resource",
	"type/new-field",
	"type/ga-promotion",
	"type/bugfix",
	"type/test-fix",
	"type/sweeper",
	"type/refactor",
	"type/docs",
	"type/generator",
	"type/ci",
	"type/agents",
	"type/deprecation",
	"type/revert",
	"type/unknown",
}

var mergeCommitPRRegex = regexp.MustCompile(`\(#(\d+)\)\s*$`)

type geminiRequest struct {
	Contents         []geminiContent        `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	ResponseMimeType string               `json:"responseMimeType"`
	Temperature      float64              `json:"temperature"`
	ThinkingConfig   geminiThinkingConfig `json:"thinkingConfig"`
	ResponseSchema   map[string]any       `json:"responseSchema"`
}

type geminiThinkingConfig struct {
	ThinkingBudget int `json:"thinkingBudget"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type prClassificationResult struct {
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
}

var labelMergedPRCmd = &cobra.Command{
	Use:   "label-merged-pr",
	Short: "Classify and label a merged pull request by type",
	Long: `This command classifies a merged pull request into one of the 14 type/* taxonomy labels
using Vertex AI Serverless Gemini 2.5 Flash, applies the label on GitHub, and posts a commit status.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commitSha := args[0]

		githubToken, ok := os.LookupEnv("GITHUB_TOKEN_CLASSIC")
		if !ok || githubToken == "" {
			return fmt.Errorf("did not provide GITHUB_TOKEN_CLASSIC environment variable")
		}

		buildID := os.Getenv("BUILD_ID")
		projectID := os.Getenv("PROJECT_ID")
		targetURL := buildCloudBuildTargetURL(buildID, projectID)

		classifierURL := os.Getenv("CLASSIFIER_URL")
		if classifierURL == "" {
			classifierURL = defaultClassifierURL
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating Runner: %w", err)
		}

		gh := github.NewClient(githubToken)
		return execLabelMergedPR(commitSha, targetURL, classifierURL, rnr, gh)
	},
}

func buildCloudBuildTargetURL(buildID, projectID string) string {
	if buildID == "" {
		return ""
	}
	if projectID == "" {
		projectID = defaultCloudBuildProjectID
	}
	return fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s?project=%s", buildID, projectID)
}

func execLabelMergedPR(commitSha, targetURL, classifierURL string, runner ExecRunner, gh GithubClient) error {
	if classifierURL == "" {
		classifierURL = defaultClassifierURL
	}

	prNumber, err := extractMergedPRNumber(commitSha, runner, gh)
	if err != nil {
		fmt.Printf("Warning: could not extract PR number for commit %s: %v\n", commitSha, err)
		return nil
	}

	pr, err := gh.GetPullRequest(prNumber)
	if err != nil {
		fmt.Printf("Warning: failed to fetch PR #%s: %v\n", prNumber, err)
		return softFailClassification(gh, prNumber, targetURL, commitSha)
	}

	for _, label := range pr.Labels {
		if strings.HasPrefix(label.Name, "type/") {
			fmt.Printf("PR #%s already has label %q; skipping classification\n", prNumber, label.Name)
			return nil
		}
	}

	changedFilesOut, err := runner.Run("git", []string{"show", "--name-only", "--format=", commitSha}, nil)
	if err != nil {
		fmt.Printf("Warning: git show --name-only failed for %s: %v\n", commitSha, err)
	}
	var changedFiles []string
	for _, line := range strings.Split(changedFilesOut, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			changedFiles = append(changedFiles, trimmed)
		}
	}
	if len(changedFiles) > maxChangedFilesInPrompt {
		changedFiles = changedFiles[:maxChangedFilesInPrompt]
	}

	diffstatOut, err := runner.Run("git", []string{"show", "--shortstat", "--format=", commitSha}, nil)
	if err != nil {
		fmt.Printf("Warning: git show --shortstat failed for %s: %v\n", commitSha, err)
	}
	diffstat := strings.TrimSpace(diffstatOut)

	patchOut, err := runner.Run("git", []string{"show", "-p", "--format=", commitSha}, nil)
	if err != nil {
		fmt.Printf("Warning: git show -p failed for %s: %v\n", commitSha, err)
	}
	patchSample := truncateUTF8(patchOut, maxPatchCharsInPrompt)

	token, err := getGCPAccessToken(defaultMetadataTokenURL)
	if err != nil {
		fmt.Printf("Warning: failed to obtain GCP access token for PR #%s: %v\n", prNumber, err)
		return softFailClassification(gh, prNumber, targetURL, commitSha)
	}

	prompt := buildClassificationPrompt(pr, changedFiles, diffstat, patchSample)
	result, err := callGeminiClassifier(classifierURL, token, prompt)
	if err != nil {
		fmt.Printf("Warning: classifier request failed for PR #%s: %v\n", prNumber, err)
		return softFailClassification(gh, prNumber, targetURL, commitSha)
	}

	if !isValidNonUnknownTypeLabel(result.Label) || result.Confidence < classifierConfidenceFloor {
		fmt.Printf("Warning: low-confidence or unknown classification for PR #%s (label=%q, confidence=%.2f)\n", prNumber, result.Label, result.Confidence)
		return softFailClassification(gh, prNumber, targetURL, commitSha)
	}

	if err := gh.AddLabels(prNumber, []string{result.Label}); err != nil {
		fmt.Printf("Warning: failed to add label %q to PR #%s: %v\n", result.Label, prNumber, err)
	}
	if err := gh.PostBuildStatus(prNumber, classifierStatusContext, "success", targetURL, commitSha); err != nil {
		fmt.Printf("Warning: failed to post success build status for PR #%s: %v\n", prNumber, err)
	}
	return nil
}

func softFailClassification(gh GithubClient, prNumber, targetURL, commitSha string) error {
	if err := gh.AddLabels(prNumber, []string{"type/unknown"}); err != nil {
		fmt.Printf("Warning: failed to add fallback label type/unknown to PR #%s: %v\n", prNumber, err)
	}
	if err := gh.PostBuildStatus(prNumber, classifierStatusContext, "failure", targetURL, commitSha); err != nil {
		fmt.Printf("Warning: failed to post failure build status for PR #%s: %v\n", prNumber, err)
	}
	return nil
}

func extractMergedPRNumber(sha string, runner ExecRunner, gh GithubClient) (string, error) {
	if runner != nil {
		subject, err := runner.Run("git", []string{"show", "-s", "--format=%s", sha}, nil)
		if err == nil {
			if prNumber := parsePRNumberFromCommitMessage(subject); prNumber != "" {
				return prNumber, nil
			}
		}
	}

	message, err := gh.GetCommitMessage("GoogleCloudPlatform", "magic-modules", sha)
	if err != nil {
		return "", fmt.Errorf("failed to get commit message for %s: %w", sha, err)
	}
	if prNumber := parsePRNumberFromCommitMessage(message); prNumber != "" {
		return prNumber, nil
	}
	return "", fmt.Errorf("no PR number found in commit subject for %s", sha)
}

func parsePRNumberFromCommitMessage(message string) string {
	lines := strings.Split(strings.TrimSpace(message), "\n")
	if len(lines) == 0 {
		return ""
	}
	firstLine := strings.TrimSpace(lines[0])
	matches := mergeCommitPRRegex.FindStringSubmatch(firstLine)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func getGCPAccessToken(metadataURL string) (string, error) {
	if token := os.Getenv("GCP_ACCESS_TOKEN"); token != "" {
		return token, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, metadataURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating metadata token request: %w", err)
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("requesting metadata token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("metadata server returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decoding metadata token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("empty access_token in metadata response")
	}
	return tokenResp.AccessToken, nil
}

func truncateUTF8(s string, maxRunes int) string {
	s = strings.ToValidUTF8(s, "")
	runes := []rune(s)
	if len(runes) > maxRunes {
		return string(runes[:maxRunes])
	}
	return s
}

func buildClassificationPrompt(pr github.PullRequest, changedFiles []string, diffstat, patchSample string) string {
	body := truncateUTF8(pr.Body, maxBodyCharsInPrompt)
	patchSample = truncateUTF8(patchSample, maxPatchCharsInPrompt)
	return fmt.Sprintf(`Classify this Google Cloud Terraform provider (Magic Modules) pull request into exactly one type label.

Taxonomy (choose the single primary purpose):
- type/new-resource: Adds a brand new Terraform resource or data source
- type/new-field: Adds new fields, properties, or arguments to an existing resource
- type/ga-promotion: Promoting beta-only resources, fields, or data sources to GA (google provider)
- type/bugfix: Fixes a provider runtime bug, API error, diff suppression, timeout, or state issue
- type/test-fix: Fixes, stabilizes, skips, or adds acceptance tests, VCR cassettes, or test helpers
- type/sweeper: Adding or modifying resource sweepers to clean up leaked test infrastructure
- type/refactor: Internal code refactoring, cleanup, lint fixes, or style updates without behavioral change
- type/docs: Documentation, website markdown, or example code updates only
- type/generator: Changes to Magic Modules generator templates, compiler, or mmv1 engine
- type/ci: Changes to CI workflows, Cloud Build configs, magician tooling, or release scripts
- type/agents: Changes to repository agent skills, instructions, or developer tooling in .agents/ or AGENTS.md (NOT the GCP Agentic Applications product)
- type/deprecation: Deprecating or removing fields or resources for major releases
- type/revert: Reverting a previous pull request or commit
- type/unknown: Unknown type or cannot be determined from description

PR #%d Title: %s

PR Body:
%s

Changed Files (up to 25):
%s

Diffstat:
%s

Patch Sample:
%s
`, pr.Number, pr.Title, body, strings.Join(changedFiles, "\n"), diffstat, patchSample)
}

func callGeminiClassifier(classifierURL, token, prompt string) (prClassificationResult, error) {
	reqPayload := geminiRequest{
		Contents: []geminiContent{
			{
				Role: "user",
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: geminiGenerationConfig{
			ResponseMimeType: "application/json",
			Temperature:      0.0,
			ThinkingConfig: geminiThinkingConfig{
				ThinkingBudget: 0,
			},
			ResponseSchema: map[string]any{
				"type": "OBJECT",
				"properties": map[string]any{
					"label": map[string]any{
						"type": "STRING",
						"enum": validPRTypeLabels,
					},
					"confidence": map[string]any{
						"type": "NUMBER",
					},
				},
				"required": []string{"label", "confidence"},
			},
		},
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return prClassificationResult{}, fmt.Errorf("marshaling request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodPost, classifierURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return prClassificationResult{}, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return prClassificationResult{}, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return prClassificationResult{}, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return prClassificationResult{}, fmt.Errorf("unexpected HTTP status %d: %s", resp.StatusCode, string(respBytes))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBytes, &geminiResp); err != nil {
		return prClassificationResult{}, fmt.Errorf("unmarshaling gemini response: %w", err)
	}
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return prClassificationResult{}, fmt.Errorf("no candidate parts in gemini response")
	}

	rawJSON := geminiResp.Candidates[0].Content.Parts[0].Text
	var result prClassificationResult
	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		return prClassificationResult{}, fmt.Errorf("unmarshaling classification JSON: %w", err)
	}
	return result, nil
}

func isValidNonUnknownTypeLabel(label string) bool {
	if label == "" || label == "type/unknown" {
		return false
	}
	for _, valid := range validPRTypeLabels {
		if label == valid {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(labelMergedPRCmd)
}
