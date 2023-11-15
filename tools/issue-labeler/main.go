package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"
	"github.com/golang/glog"
)

var flagBackfillDate = flag.String("backfill-date", "", "run in backfill mode to apply labels to issues filed after given date")
var flagDryRun = flag.Bool("backfill-dry-run", false, "when combined with backfill-date, perform a dry run of backfill mode")

func main() {
	flag.Parse()

	regexpLabels, err := labeler.BuildRegexLabels(labeler.EnrolledTeamsYaml)
	if err != nil {
		glog.Exitf("Error building regex labels: %v", err)
	}

	if *flagBackfillDate == "" {
		issueBody := os.Getenv("ISSUE_BODY")
		affectedResources := labeler.ExtractAffectedResources(issueBody)
		labels := labeler.ComputeLabels(affectedResources, regexpLabels)

		if len(labels) > 0 {
			labels = append(labels, "forward/review")
			sort.Strings(labels)
			fmt.Println(`["` + strings.Join(labels, `", "`) + `"]`)
		}
	} else {
		repository := "hashicorp/terraform-provider-google"
		issues := labeler.GetIssues(repository, *flagBackfillDate)
		issueUpdates := labeler.ComputeIssueUpdates(issues, regexpLabels)
		labeler.UpdateIssues(repository, issueUpdates, *flagDryRun)
	}
}
