package alloydb

import (
	"fmt"
	"testing"

	"google.golang.org/api/googleapi"
)

func TestIsAlloyDbInstanceUpdatingError(t *testing.T) {
	t.Parallel()

	updatingMessage := `Invalid resource state for "projects/653407317329/locations/us-central1/clusters/tf-test-alloydb-clusterna2xzcniou/instances/tf-test-alloydb-instancena2xzcniou": The instance is currently undergoing another update operation. Please try again later.`
	updatingMixedCaseMessage := `Invalid resource state: the Instance is currently undergoing another update operation`
	backgroundSetupMessage := `Invalid resource state for "projects/794360760068/locations/asia-southeast1/clusters/tf-test-alloydb-clustersd8rrxe71e/instances/tf-test-alloydb-instancesd8rrxe71e": background setup in progress, retry in a few minutes`
	queueMessage := "unable to queue the operation"
	expectedUpdatingReason := "Instance is undergoing another update operation"
	expectedBackgroundSetupReason := "Waiting for background setup to finish"
	expectedQueueReason := "Waiting for other enqueued operations to finish"

	cases := map[string]struct {
		err            error
		expected       bool
		expectedReason string
	}{
		"nil error": {
			err:            nil,
			expected:       false,
			expectedReason: "",
		},
		"googleapi 400 updating with message and body": {
			err: &googleapi.Error{
				Code:    400,
				Message: updatingMessage,
				Body:    fmt.Sprintf(`{"error":{"code":400,"message":%q,"status":"FAILED_PRECONDITION"}}`, updatingMessage),
			},
			expected:       true,
			expectedReason: expectedUpdatingReason,
		},
		"googleapi 400 updating case-insensitive": {
			err: &googleapi.Error{
				Code:    400,
				Message: updatingMixedCaseMessage,
			},
			expected:       true,
			expectedReason: expectedUpdatingReason,
		},
		"googleapi 400 updating with body only": {
			err: &googleapi.Error{
				Code: 400,
				Body: updatingMessage,
			},
			expected:       true,
			expectedReason: expectedUpdatingReason,
		},
		"googleapi 400 background setup with message and body": {
			err: &googleapi.Error{
				Code:    400,
				Message: backgroundSetupMessage,
				Body:    fmt.Sprintf(`{"error":{"code":400,"message":%q,"status":"FAILED_PRECONDITION"}}`, backgroundSetupMessage),
			},
			expected:       true,
			expectedReason: expectedBackgroundSetupReason,
		},
		"googleapi 400 background setup with body only": {
			err: &googleapi.Error{
				Code: 400,
				Body: backgroundSetupMessage,
			},
			expected:       true,
			expectedReason: expectedBackgroundSetupReason,
		},
		"wrapped standard error 400 background setup": {
			err:            fmt.Errorf("Error updating Cluster: googleapi: Error 400: %s", backgroundSetupMessage),
			expected:       true,
			expectedReason: expectedBackgroundSetupReason,
		},
		"googleapi 409 unable to queue with message and body": {
			err: &googleapi.Error{
				Code:    409,
				Message: queueMessage,
				Body:    fmt.Sprintf(`{"error":{"code":409,"message":%q,"status":"ABORTED"}}`, queueMessage),
			},
			expected:       true,
			expectedReason: expectedQueueReason,
		},
		"googleapi 409 unable to queue with body only": {
			err: &googleapi.Error{
				Code: 409,
				Body: queueMessage,
			},
			expected:       true,
			expectedReason: expectedQueueReason,
		},
		"googleapi 400 other error": {
			err: &googleapi.Error{
				Code:    400,
				Message: "Some other 400 error",
			},
			expected:       false,
			expectedReason: "",
		},
		"googleapi 409 other error": {
			err: &googleapi.Error{
				Code:    409,
				Message: "Some other 409 error",
			},
			expected:       false,
			expectedReason: "",
		},
		"wrapped standard error 400": {
			err:            fmt.Errorf("Error updating Instance: googleapi: Error 400: %s", updatingMessage),
			expected:       true,
			expectedReason: expectedUpdatingReason,
		},
		"wrapped standard error 409": {
			err:            fmt.Errorf("Error updating Instance: googleapi: Error 409: %s", queueMessage),
			expected:       true,
			expectedReason: expectedQueueReason,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			isRetryable, reason := IsAlloyDbInstanceUpdatingError(tc.err)
			if isRetryable != tc.expected || reason != tc.expectedReason {
				t.Errorf("IsAlloyDbInstanceUpdatingError(%v) = (%v, %q), want (%v, %q)", tc.err, isRetryable, reason, tc.expected, tc.expectedReason)
			}
		})
	}
}
