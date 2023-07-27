package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

var flagBackfill = flag.String("backfill-date", "", "run in backfill mode to apply labels to issues filed after given date")
var flagDryRun = flag.Bool("backfill-dry-run", false, "when combined with backfill-date, perform a dry run of backfill mode")

func main() {
	flag.Parse()

	file, err := os.ReadFile("enrolled_teams.yaml")
	if err != nil {
		glog.Exitf("Error reading enrolled teams yaml: %v", err)
	}
	enrolledTeams := make(map[string][]string)
	if err := yaml.Unmarshal(file, &enrolledTeams); err != nil {
		glog.Exitf("Error unmarshalling enrolled teams yaml: %v", err)
	}

	if *flagBackfill == "" {
		issueBody := os.Getenv("ISSUE_BODY")
		desired := labels(issueBody, enrolledTeams, true)
		if len(desired) > 0 {
			desired = append(desired, `"forward/review"`)
			sort.Strings(desired)
			fmt.Println("[" + strings.Join(desired, ", ") + "]")
		}
	} else {
		backfill(*flagBackfill, enrolledTeams, *flagDryRun)
	}
}
