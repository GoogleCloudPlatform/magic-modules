package eventarc

import (
	"testing"
)

func TestFlattenEventarcTriggerConditions(t *testing.T) {
	mockApiResponse := map[string]interface{}{
		"transport.pubsub.topic": map[string]interface{}{
			"code":    "UNKNOWN",
			"message": "Pub/Sub topic status is unknown. Try requesting the trigger description again.",
		},
	}

	result := flattenEventarcTriggerConditions(mockApiResponse, nil, nil)

	flatResult, ok := result.(map[string]string)
	if !ok {
		t.Fatalf("Flattener returned %T, expected map[string]string", result)
	}

	expectedMessage := "Pub/Sub topic status is unknown. Try requesting the trigger description again."
	if flatResult["transport.pubsub.topic"] != expectedMessage {
		t.Fatalf("Expected message %q, got %q", expectedMessage, flatResult["transport.pubsub.topic"])
	}
}
