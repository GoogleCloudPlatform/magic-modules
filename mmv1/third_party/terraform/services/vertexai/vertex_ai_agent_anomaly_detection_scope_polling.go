package vertexai

import (
	"fmt"
	"log"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// PollCheckForScopeActive checks the scope's state field during create polling.
// ACTIVE → done, CREATING/UPDATING → retry, FAILED → error.
func PollCheckForScopeActive(resp map[string]interface{}, respErr error) transport_tpg.PollResult {
	if respErr != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(respErr, 404) {
			log.Printf("[DEBUG] AgentAnomalyDetectionScope poll: not found yet, retrying...")
			return transport_tpg.PendingStatusPollResult("not found")
		}
		log.Printf("[DEBUG] AgentAnomalyDetectionScope poll: error: %s", respErr)
		return transport_tpg.ErrorPollResult(respErr)
	}

	state, ok := resp["state"].(string)
	if !ok {
		log.Printf("[DEBUG] AgentAnomalyDetectionScope poll: state field not found, retrying...")
		return transport_tpg.PendingStatusPollResult("state unknown")
	}

	log.Printf("[DEBUG] AgentAnomalyDetectionScope poll: state=%s", state)

	switch state {
	case "ACTIVE":
		return transport_tpg.SuccessPollResult()
	case "CREATING", "UPDATING":
		return transport_tpg.PendingStatusPollResult(state)
	case "FAILED":
		return transport_tpg.ErrorPollResult(fmt.Errorf("AgentAnomalyDetectionScope reached FAILED state"))
	default:
		return transport_tpg.PendingStatusPollResult(state)
	}
}

// PollCheckForScopeDeleted checks the scope's state field during delete polling.
// 404 → done, DELETING → retry, FAILED → error.
func PollCheckForScopeDeleted(resp map[string]interface{}, respErr error) transport_tpg.PollResult {
	if respErr != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(respErr, 404) {
			log.Printf("[DEBUG] AgentAnomalyDetectionScope delete poll: scope deleted (404)")
			return transport_tpg.SuccessPollResult()
		}
		log.Printf("[DEBUG] AgentAnomalyDetectionScope delete poll: error: %s", respErr)
		return transport_tpg.ErrorPollResult(respErr)
	}

	state, ok := resp["state"].(string)
	if !ok {
		log.Printf("[DEBUG] AgentAnomalyDetectionScope delete poll: state field not found, retrying...")
		return transport_tpg.PendingStatusPollResult("state unknown")
	}

	log.Printf("[DEBUG] AgentAnomalyDetectionScope delete poll: state=%s", state)

	switch state {
	case "DELETING":
		return transport_tpg.PendingStatusPollResult("DELETING")
	case "FAILED":
		return transport_tpg.ErrorPollResult(fmt.Errorf("AgentAnomalyDetectionScope deletion reached FAILED state"))
	default:
		return transport_tpg.PendingStatusPollResult(state)
	}
}
