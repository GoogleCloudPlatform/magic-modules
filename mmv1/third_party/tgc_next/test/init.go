package test

import (
	"os"
)

var (
	data   *testData
	tmpDir = os.TempDir()
)

// testData represents the full dataset that is used for templating terraform
// configs. It contains Google API resources that are expected to be returned
// after converting the terraform plan.
type testData struct {
	// provider "google"
	Provider map[string]string
	Project  map[string]string
	Time     map[string]string
	OrgID    string
	FolderID string
	Ancestry string
}

// init initializes the variables used for testing. As tests rely on
// environment variables, the parsing of those are only done once.
func init() {
	// credentials := getTestCredsFromEnv()
	// org := getTestOrgFromEnv(nil)
	// billingAccount := getTestBillingAccountFromEnv(nil)
	// folder, ok := os.LookupEnv("TEST_FOLDER_ID")
	// if !ok {
	// 	log.Printf("Missing required env var TEST_FOLDER_ID. Default (%s) will be used.", defaultFolder)
	// 	folder = defaultFolder
	// }
	// ancestry, ok := os.LookupEnv("TEST_ANCESTRY")
	// if !ok {
	// 	log.Printf("Missing required env var TEST_ANCESTRY. Default (%s) will be used.", defaultAncestry)
	// 	ancestry = defaultAncestry
	// }
	// providerVersion := defaultProviderVersion
	//As time is not information in terraform resource data, time is fixed for testing purposes
	// fixedTime := time.Date(2021, time.April, 14, 15, 16, 17, 0, time.UTC)
	// data = &testData{
	// 	Provider: map[string]string{
	// 		"version":     providerVersion,
	// 		"project":     defaultProject,
	// 		"credentials": credentials,
	// 	},
	// 	Time: map[string]string{
	// 		"RFC3339Nano": fixedTime.Format(time.RFC3339Nano),
	// 	},
	// 	Project: map[string]string{
	// 		"Name":      "My Project Name",
	// 		"ProjectId": "my-project-id",
	// 		// "BillingAccountName": billingAccount,
	// 		"Number": "1234567890",
	// 	},
	// 	// OrgID:    org,
	// 	// FolderID: folder,
	// 	Ancestry: defaultAncestry,
	// }
}
