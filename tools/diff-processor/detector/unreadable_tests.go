package detector

import (
	"sort"
	"strings"
)

// KnownUnreadableTests lists the acceptance tests in the beta provider that
// the test reader can't parse. The missing test detector skips these tests, so
// the fields they cover aren't counted.
//
// This list should only shrink. Rather than adding a new test here, rewrite
// it so that each step's Config is a call to a config function that returns a
// single fmt.Sprintf or acctest.Nprintf template. When a test on this list
// becomes readable, remove it.
var KnownUnreadableTests = []string{
	"TestAccComputeInstance_diskResourcePolicies_attachmentDiff",
	"TestAccComputeNetwork_bgpModeAndMedInteractions",
	"TestAccComputeNetwork_networkBgpStandardModeDeleteMed",
	"TestAccComputeRouterNamedSet_updatePrefix",
	"TestAccContainerCluster_resourceManagerTags",
	"TestAccDNSManagedZone_update",
	"TestAccEphemeralServiceAccountJwt_basic",
	"TestAccEphemeralServiceAccountJwt_withDelegates",
	"TestAccEphemeralServiceAccountJwt_withExpiresIn",
	"TestAccEphemeralServiceAccountKey_create",
	"TestAccEphemeralServiceAccountKey_fetch",
	"TestAccEphemeralServiceAccountKey_upload",
	"TestAccEphemeralServiceAccountToken_basic",
	"TestAccEphemeralServiceAccountToken_withCustomLifetime",
	"TestAccEphemeralServiceAccountToken_withDelegates",
	"TestAccPubsubTopic_schema",
	"TestAccStorageTransferJob_azureBlobStorageDataSourcePrivateNetwork",
}

// The longest error message reported for a single unreadable test.
const maxUnreadableTestErrorLength = 300

// UnreadableTestsReport compares the tests the test reader failed to read with
// the tests it is known to be unable to read.
type UnreadableTestsReport struct {
	// Tests the reader failed to read that aren't known to be unreadable,
	// mapped to the error the reader reported.
	NewUnreadableTests map[string]string
	// Tests known to be unreadable that the reader read without errors. These
	// should be removed from KnownUnreadableTests.
	StaleUnreadableTests []string
}

// DetectUnreadableTests compares the errors from reading a set of acceptance
// tests, as returned by reader.ReadTestFiles, with the tests known to be
// unreadable. testNames are the names of the acceptance test functions that
// were read.
func DetectUnreadableTests(testNames []string, readErrs map[string]error, known []string) UnreadableTestsReport {
	knownSet := make(map[string]struct{}, len(known))
	for _, test := range known {
		knownSet[test] = struct{}{}
	}
	report := UnreadableTestsReport{
		NewUnreadableTests:   make(map[string]string),
		StaleUnreadableTests: []string{},
	}
	for test, err := range readErrs {
		if _, ok := knownSet[test]; !ok {
			report.NewUnreadableTests[test] = truncateError(err)
		}
	}
	for _, test := range testNames {
		if _, failed := readErrs[test]; failed {
			continue
		}
		if _, ok := knownSet[test]; ok {
			report.StaleUnreadableTests = append(report.StaleUnreadableTests, test)
		}
	}
	sort.Strings(report.StaleUnreadableTests)
	return report
}

// Shorten an error so it can be shown in a PR comment. Reader errors can embed
// large printed syntax trees.
func truncateError(err error) string {
	msg := strings.Join(strings.Fields(err.Error()), " ")
	if runes := []rune(msg); len(runes) > maxUnreadableTestErrorLength {
		return string(runes[:maxUnreadableTestErrorLength]) + "…"
	}
	return msg
}
