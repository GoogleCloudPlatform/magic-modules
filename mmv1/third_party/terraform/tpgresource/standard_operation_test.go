package tpgresource

import (
	"fmt"
	"strings"
	"testing"
)

const (
	standardTopLevelMsg             = "Top-level message."
	standardLocalizedMsgTmpl        = "LocalizedMessage%d message"
	standardHelpLinkDescriptionTmpl = "Help%dLink%d Description"
	standardHelpLinkUrlTmpl         = "https://help%d.com/link%d"
	standardQuotaExceededMsg        = "Quota DISKS_TOTAL_GB exceeded.  Limit: 1100.0 in region us-central1."
	standardQuotaExceededCode       = "QUOTA_EXCEEDED"
	standardQuotaMetricName         = "compute.googleapis.com/disks_total_storage"
	standardQuotaLimitName          = "DISKS-TOTAL-GB-per-project-region"
)

var standardOperationErrorLocales = []string{"en-US", "es-US", "es-ES", "es-MX", "de-DE"}

func buildStandardOperationError(numLocalizedMsg int, numHelpWithLinks []int) map[string]interface{} {
	errorDetails := []interface{}{}

	for n := 1; n <= numLocalizedMsg; n++ {
		errorDetails = append(errorDetails, map[string]interface{}{
			"localizedMessage": map[string]interface{}{
				"locale":  standardOperationErrorLocales[n-1%len(standardOperationErrorLocales)],
				"message": formatStandardLocalizedMsg(n),
			},
		})
	}

	for i := 0; i < len(numHelpWithLinks); i++ {
		links := []interface{}{}
		for nLinks := 1; nLinks <= numHelpWithLinks[i]; nLinks++ {
			desc, url := formatStandardLink(i+1, nLinks)
			links = append(links, map[string]interface{}{
				"description": desc,
				"url":         url,
			})
		}
		errorDetails = append(errorDetails, map[string]interface{}{
			"help": map[string]interface{}{
				"links": links,
			},
		})
	}

	return map[string]interface{}{
		"errors": []interface{}{
			map[string]interface{}{
				"message":      standardTopLevelMsg,
				"errorDetails": errorDetails,
			},
		},
	}
}

func buildStandardOperationErrorQuotaExceeded(withDetails bool, withDimensions bool, withFutureLimit bool) map[string]interface{} {
	opError := map[string]interface{}{
		"message": standardQuotaExceededMsg,
		"code":    standardQuotaExceededCode,
	}

	if withDetails {
		quotaInfo := map[string]interface{}{
			"metricName": standardQuotaMetricName,
			"limitName":  standardQuotaLimitName,
			"limit":      float64(1100),
		}
		if withFutureLimit {
			quotaInfo["futureLimit"] = float64(2200)
		}
		if withDimensions {
			quotaInfo["dimensions"] = map[string]interface{}{"region": "us-central1"}
		}
		opError["errorDetails"] = []interface{}{
			map[string]interface{}{
				"quotaInfo": quotaInfo,
			},
		}
	}

	return map[string]interface{}{
		"errors": []interface{}{opError},
	}
}

func omitStandardAlways(numLocalizedMsg int, numHelpWithLinks []int) []string {
	var omits []string

	for n := 2; n <= numLocalizedMsg; n++ {
		omits = append(omits, fmt.Sprintf("LocalizedMessage%d", n))
	}

	for i := 0; i < len(numHelpWithLinks); i++ {
		for j := maxStandardLinks(i); j < numHelpWithLinks[i]; j++ {
			desc, url := formatStandardLink(i+1, j+1)
			omits = append(omits, desc, url)
		}
	}

	return omits

}

func maxStandardLinks(helpIndex int) int {
	if helpIndex == 0 {
		return 1
	}

	return 0
}

func formatStandardLocalizedMsg(localizedMsgNum int) string {
	return fmt.Sprintf(standardLocalizedMsgTmpl, localizedMsgNum)
}

func formatStandardLink(helpNum, linkNum int) (string, string) {
	return fmt.Sprintf(standardHelpLinkDescriptionTmpl, helpNum, linkNum), fmt.Sprintf(standardHelpLinkUrlTmpl, helpNum, linkNum)
}

func TestStandardOperationError_Error(t *testing.T) {
	testCases := []struct {
		name           string
		input          map[string]interface{}
		expectContains []string
		expectOmits    []string
	}{
		{
			name:  "MessageOnly",
			input: buildStandardOperationError(0, []int{}),
			expectContains: []string{
				"Top-level",
			},
			expectOmits: append(omitStandardAlways(0, []int{}), []string{
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			}...),
		},
		{
			name:  "WithLocalizedMessageAndNoHelp",
			input: buildStandardOperationError(1, []int{}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
			},
			expectOmits: append(omitStandardAlways(1, []int{}), []string{
				"Help1Link1 Description",
				"https://help1.com/link1",
			}...),
		},
		{
			name:  "WithLocalizedMessageAndHelp",
			input: buildStandardOperationError(1, []int{1}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitStandardAlways(1, []int{1}), []string{}...),
		},
		{
			name:  "WithNoLocalizedMessageAndHelp",
			input: buildStandardOperationError(0, []int{1}),
			expectContains: []string{
				"Top-level",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitStandardAlways(0, []int{1}), []string{
				"LocalizedMessage1",
			}...),
		},
		{
			name:  "WithLocalizedMessageAndHelpWithTwoLinks",
			input: buildStandardOperationError(1, []int{2}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitStandardAlways(1, []int{2}), []string{}...),
		},
		// The case below should never happen because the server should just send multiple links
		// but the protobuf defition would allow it, so testing anyway.
		{
			name:  "WithLocalizedMessageAndTwoHelpsWithTwoLinks",
			input: buildStandardOperationError(1, []int{2, 2}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitStandardAlways(1, []int{2, 2}), []string{}...),
		},
		// This should never happen because the server should never respond with the messages for
		// two locales at once, but should rather take the locale as input to the API and serve
		// the appropriate message for that locale. However, the protobuf defition would allow it,
		// so we'll test for it. The second message in the list would be ignored.
		{
			name:  "WithTwoLocalizedMessageAndHelp",
			input: buildStandardOperationError(2, []int{1}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitStandardAlways(2, []int{1}), []string{}...),
		},
		{
			name:  "QuotaMessageOnly",
			input: buildStandardOperationErrorQuotaExceeded(false, false, false),
			expectContains: []string{
				"Quota DISKS_TOTAL_GB exceeded.  Limit: 1100.0 in region us-central1.",
			},
			expectOmits: append(omitStandardAlways(0, []int{}), []string{
				"metric name = compute.googleapis.com/disks_total_storage",
				"limit = 1100",
			}...),
		},
		{
			name:  "QuotaMessageWithDetailsNoDimensions",
			input: buildStandardOperationErrorQuotaExceeded(true, false, false),
			expectContains: []string{
				"Quota DISKS_TOTAL_GB exceeded.  Limit: 1100.0 in region us-central1.",
				"metric name = compute.googleapis.com/disks_total_storage",
				"limit name = DISKS-TOTAL-GB-per-project-region",
				"limit = 1100",
			},
			expectOmits: append(omitStandardAlways(0, []int{}), []string{
				"dimensions = map[region:us-central1]",
			}...),
		},
		{
			name:  "QuotaMessageWithDetailsWithDimensions",
			input: buildStandardOperationErrorQuotaExceeded(true, true, false),
			expectContains: []string{
				"Quota DISKS_TOTAL_GB exceeded.  Limit: 1100.0 in region us-central1.",
				"metric name = compute.googleapis.com/disks_total_storage",
				"limit name = DISKS-TOTAL-GB-per-project-region",
				"limit = 1100",
				"dimensions = map[region:us-central1]",
			},
			expectOmits: append(omitStandardAlways(0, []int{}), []string{
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			}...),
		},
		{
			name:  "QuotaMessageWithDetailsWithFutureLimit",
			input: buildStandardOperationErrorQuotaExceeded(true, false, true),
			expectContains: []string{
				"Quota DISKS_TOTAL_GB exceeded.  Limit: 1100.0 in region us-central1.",
				"metric name = compute.googleapis.com/disks_total_storage",
				"limit name = DISKS-TOTAL-GB-per-project-region",
				"limit = 1100",
				"future limit = 2200",
				"rollout status = in progress",
			},
			expectOmits: append(omitStandardAlways(0, []int{}), []string{
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			}...),
		},
		{
			name:  "QuotaMessageWithDetailsWithDimensionsWithFutureLimit",
			input: buildStandardOperationErrorQuotaExceeded(true, true, true),
			expectContains: []string{
				"Quota DISKS_TOTAL_GB exceeded.  Limit: 1100.0 in region us-central1.",
				"metric name = compute.googleapis.com/disks_total_storage",
				"limit name = DISKS-TOTAL-GB-per-project-region",
				"limit = 1100",
				"future limit = 2200",
				"rollout status = in progress",
				"dimensions = map[region:us-central1]",
			},
			expectOmits: append(omitStandardAlways(0, []int{}), []string{
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			}...),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := StandardOperationError(tc.input)
			str := err.Error()

			for _, contains := range tc.expectContains {
				if !strings.Contains(str, contains) {
					t.Errorf("expected\n%s\nto contain, %q, and did not", str, contains)
				}
			}

			for _, omits := range tc.expectOmits {
				if strings.Contains(str, omits) {
					t.Errorf("expected\n%s\nnot to contain, %q, and did not", str, omits)
				}
			}
		})
	}
}
