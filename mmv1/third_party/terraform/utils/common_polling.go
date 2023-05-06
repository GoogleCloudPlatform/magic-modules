package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Helper functions to construct result of single pollRead as return result for a PollCheckResponseFunc
func ErrorPollResult(err error) transport_tpg.PollResult {
	return transport_tpg.ErrorPollResult(err)
}

func PendingStatusPollResult(status string) transport_tpg.PollResult {
	return transport_tpg.PendingStatusPollResult(status)
}

func SuccessPollResult() transport_tpg.PollResult {
	return transport_tpg.SuccessPollResult()
}

func PollingWaitTime(pollF transport_tpg.PollReadFunc, checkResponse transport_tpg.PollCheckResponseFunc, activity string,
	timeout time.Duration, targetOccurrences int) error {
	return transport_tpg.PollingWaitTime(pollF, checkResponse, activity, timeout, targetOccurrences)
}

// RetryWithTargetOccurrences is a basic wrapper around StateChangeConf that will retry
// a function until it returns the specified amount of target occurrences continuously.
// Adapted from the Retry function in the go SDK.
func RetryWithTargetOccurrences(timeout time.Duration, targetOccurrences int,
	f resource.RetryFunc) error {
	return transport_tpg.RetryWithTargetOccurrences(timeout, targetOccurrences, f)
}

/**
 * Common PollCheckResponseFunc implementations
 */

// PollCheckForExistence waits for a successful response, continues polling on 404,
// and returns any other error.
func PollCheckForExistence(_ map[string]interface{}, respErr error) transport_tpg.PollResult {
	return transport_tpg.PollCheckForExistence(nil, respErr)
}

// PollCheckForExistenceWith403 waits for a successful response, continues polling on 404 or 403,
// and returns any other error.
func PollCheckForExistenceWith403(_ map[string]interface{}, respErr error) transport_tpg.PollResult {
	return transport_tpg.PollCheckForExistenceWith403(nil, respErr)
}

// PollCheckForAbsence waits for a 404/403 response, continues polling on a successful
// response, and returns any other error.
func PollCheckForAbsenceWith403(_ map[string]interface{}, respErr error) transport_tpg.PollResult {
	return transport_tpg.PollCheckForAbsenceWith403(nil, respErr)
}

// PollCheckForAbsence waits for a 404 response, continues polling on a successful
// response, and returns any other error.
func PollCheckForAbsence(_ map[string]interface{}, respErr error) transport_tpg.PollResult {
	return transport_tpg.PollCheckForAbsence(nil, respErr)
}
