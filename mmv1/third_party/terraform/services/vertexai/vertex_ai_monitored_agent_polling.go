package vertexai

import (
	"fmt"
	"log"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// PollCheckForMonitoredAgentActive checks the agent's state during enable polling.
// ACTIVE → done, ENABLING → retry, DISABLED/error → error.
func PollCheckForMonitoredAgentActive(resp map[string]interface{}, respErr error) transport_tpg.PollResult {
	if respErr != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(respErr, 404) {
			log.Printf("[DEBUG] MonitoredAgent poll: not found yet, retrying...")
			return transport_tpg.PendingStatusPollResult("not found")
		}
		log.Printf("[DEBUG] MonitoredAgent poll: error: %s", respErr)
		return transport_tpg.ErrorPollResult(respErr)
	}

	state, ok := resp["state"].(string)
	if !ok {
		log.Printf("[DEBUG] MonitoredAgent poll: state field not found, retrying...")
		return transport_tpg.PendingStatusPollResult("state unknown")
	}

	log.Printf("[DEBUG] MonitoredAgent poll: state=%s", state)

	switch state {
	case "ACTIVE":
		return transport_tpg.SuccessPollResult()
	case "ENABLING":
		return transport_tpg.PendingStatusPollResult("ENABLING")
	case "DISABLED":
		return transport_tpg.ErrorPollResult(fmt.Errorf("MonitoredAgent is DISABLED after enable attempt"))
	default:
		return transport_tpg.PendingStatusPollResult(state)
	}
}

// PollCheckForMonitoredAgentDisabled checks the agent's state during disable polling.
// DISABLED → done, DISABLING → retry.
func PollCheckForMonitoredAgentDisabled(resp map[string]interface{}, respErr error) transport_tpg.PollResult {
	if respErr != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(respErr, 404) {
			log.Printf("[DEBUG] MonitoredAgent disable poll: not found (404), treating as disabled")
			return transport_tpg.SuccessPollResult()
		}
		log.Printf("[DEBUG] MonitoredAgent disable poll: error: %s", respErr)
		return transport_tpg.ErrorPollResult(respErr)
	}

	state, ok := resp["state"].(string)
	if !ok {
		log.Printf("[DEBUG] MonitoredAgent disable poll: state field not found, retrying...")
		return transport_tpg.PendingStatusPollResult("state unknown")
	}

	log.Printf("[DEBUG] MonitoredAgent disable poll: state=%s", state)

	switch state {
	case "DISABLED":
		return transport_tpg.SuccessPollResult()
	default:
		return transport_tpg.PendingStatusPollResult(state)
	}
}
