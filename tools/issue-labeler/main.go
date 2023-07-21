package main

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

func main() {
	issueBody := os.Getenv("ISSUE_BODY")
	file, err := os.ReadFile("enrolled_teams.yaml")
	if err != nil {
		glog.Exitf("Error reading enrolled teams yaml: %v", err)
	}
	enrolledTeams := make(map[string][]string)
	err = yaml.Unmarshal(file, &enrolledTeams)
	if err != nil {
		glog.Exitf("Error unmarshalling enrolled teams yaml: %v", err)
	}
	fmt.Println(labels(issueBody, enrolledTeams))
}
