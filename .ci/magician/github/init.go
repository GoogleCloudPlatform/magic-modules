package github

import (
	"fmt"
	"os"
	"strings"
)

var github_token string

func init() {
	isTest := false
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.v") {
			isTest = true
			break
		}
	}

	if isTest {
		github_token = "dummyToken"
		return
	}

	GITHUB_TOKEN, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("Did not provide GITHUB_TOKEN environment variable")
		os.Exit(1)
	}

	github_token = GITHUB_TOKEN
}
