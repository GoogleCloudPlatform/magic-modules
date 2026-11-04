package oracledatabase

import (
	"fmt"
	"testing"

	"google.golang.org/api/googleapi"
)

func TestIsOracleDatabaseConnectionNetworkUnreachableError(t *testing.T) {
	t.Parallel()

	timeoutMessage := "Cannot reach the database host or network: the connection attempt timed out."
	expectedReason := "Waiting for Oracle database connection to be reachable"

	retryableErr := &googleapi.Error{
		Code:    400,
		Message: timeoutMessage,
		Body:    fmt.Sprintf(`{"error":{"code":400,"message":%q,"status":"FAILED_PRECONDITION","details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"DATABASE_NETWORK_UNREACHABLE","domain":"oracledatabase.googleapis.com","metadata":{"is_retryable":"true"}}]}}`, timeoutMessage),
	}

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
		"googleapi 400 with DATABASE_NETWORK_UNREACHABLE reason in body and timeout message": {
			err:            retryableErr,
			expected:       true,
			expectedReason: expectedReason,
		},
		"googleapi 400 with message only": {
			err: &googleapi.Error{
				Code:    400,
				Message: timeoutMessage,
			},
			expected:       false,
			expectedReason: "",
		},
		"googleapi 400 with body only": {
			err: &googleapi.Error{
				Code: 400,
				Body: "DATABASE_NETWORK_UNREACHABLE",
			},
			expected:       false,
			expectedReason: "",
		},
		"googleapi 400 other error": {
			err: &googleapi.Error{
				Code:    400,
				Message: "Invalid connection string syntax",
				Body:    `{"error":{"code":400,"message":"Invalid connection string syntax","status":"INVALID_ARGUMENT"}}`,
			},
			expected:       false,
			expectedReason: "",
		},
		"googleapi 404 error": {
			err: &googleapi.Error{
				Code:    404,
				Message: "Resource not found",
			},
			expected:       false,
			expectedReason: "",
		},
		"non-googleapi error": {
			err:            fmt.Errorf("Error updating DatabaseConnection: googleapi: Error 400: %s", timeoutMessage),
			expected:       false,
			expectedReason: "",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			predicate := IsOracleDatabaseConnectionNetworkUnreachableError(5)
			isRetryable, reason := predicate(tc.err)
			if isRetryable != tc.expected || reason != tc.expectedReason {
				t.Errorf("IsOracleDatabaseConnectionNetworkUnreachableError(5)(%v) = (%v, %q), want (%v, %q)", tc.err, isRetryable, reason, tc.expected, tc.expectedReason)
			}
		})
	}

	t.Run("exhaustion after maxAttempts", func(t *testing.T) {
		predicate := IsOracleDatabaseConnectionNetworkUnreachableError(3)
		for i := 1; i <= 3; i++ {
			isRetryable, reason := predicate(retryableErr)
			if !isRetryable || reason != expectedReason {
				t.Fatalf("attempt %d: got (%v, %q), want (true, %q)", i, isRetryable, reason, expectedReason)
			}
		}
		isRetryable, reason := predicate(retryableErr)
		if isRetryable || reason != "" {
			t.Fatalf("attempt 4: got (%v, %q), want (false, \"\")", isRetryable, reason)
		}
	})

	t.Run("maxAttempts <= 0", func(t *testing.T) {
		predicate := IsOracleDatabaseConnectionNetworkUnreachableError(0)
		isRetryable, reason := predicate(retryableErr)
		if isRetryable || reason != "" {
			t.Fatalf("attempt 1 with maxAttempts 0: got (%v, %q), want (false, \"\")", isRetryable, reason)
		}
	})
}
