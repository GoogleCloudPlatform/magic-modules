package github

import (
	"fmt"
	"os"
)

// Client for GitHub interactions.
type Client struct {
	token string
}

func NewClient() *Client {
	githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("Did not provide GITHUB_TOKEN environment variable")
		os.Exit(1)
	}

	return &Client{token: githubToken}
}
