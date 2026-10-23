package compute

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"
)

func TestIsComputeInstanceStockoutError(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		err      error
		expected bool
	}{
		"zone capacity exhausted operation": {
			err: ComputeOperationError(map[string]interface{}{
				"errors": []interface{}{map[string]interface{}{"code": "ZONE_RESOURCE_POOL_EXHAUSTED"}},
			}),
			expected: true,
		},
		"zone capacity exhausted with details": {
			err: ComputeOperationError(map[string]interface{}{
				"errors": []interface{}{map[string]interface{}{"code": "ZONE_RESOURCE_POOL_EXHAUSTED_WITH_DETAILS"}},
			}),
			expected: true,
		},
		"zone capacity exhausted API reason": {
			err:      &googleapi.Error{Errors: []googleapi.ErrorItem{{Reason: "zoneResourcePoolExhausted"}}},
			expected: true,
		},
		"region capacity exhausted": {
			err:      &googleapi.Error{Errors: []googleapi.ErrorItem{{Reason: "REGION_RESOURCE_POOL_EXHAUSTED"}}},
			expected: false,
		},
		"GKE stockout": {
			err:      errors.New("GCE_STOCKOUT"),
			expected: false,
		},
		"generic stockout message": {
			err:      errors.New("a stockout occurred"),
			expected: false,
		},
		"unrelated error": {
			err:      errors.New("permission denied"),
			expected: false,
		},
		"nil error": {
			err:      nil,
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := isComputeInstanceStockoutError(tc.err); got != tc.expected {
				t.Errorf("isComputeInstanceStockoutError(%v) = %t, want %t", tc.err, got, tc.expected)
			}
		})
	}
}

func TestComputeInstanceStockoutRetrySchemaDefaults(t *testing.T) {
	schema := ResourceComputeInstance().Schema
	if got := schema["retry_on_stockout"].Default; got != false {
		t.Errorf("retry_on_stockout default = %v, want false", got)
	}
	if got := schema["stockout_retry_timeout_seconds"].Default; got != 300 {
		t.Errorf("stockout_retry_timeout_seconds default = %v, want 300", got)
	}
}

func TestComputeInstanceStockoutRetrySchemaConfig(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ResourceComputeInstance().Schema, map[string]interface{}{
		"retry_on_stockout":              true,
		"stockout_retry_timeout_seconds": 600,
	})

	if got := data.Get("retry_on_stockout"); got != true {
		t.Errorf("retry_on_stockout = %v, want true", got)
	}
	if got := data.Get("stockout_retry_timeout_seconds"); got != 600 {
		t.Errorf("stockout_retry_timeout_seconds = %v, want 600", got)
	}
}

func TestComputeInstanceStockoutRetryDelay(t *testing.T) {
	baseDelays := []time.Duration{
		5 * time.Second,
		10 * time.Second,
		20 * time.Second,
		30 * time.Second,
		30 * time.Second,
	}

	for retryNumber, baseDelay := range baseDelays {
		for sample := 0; sample < 100; sample++ {
			delay := computeInstanceStockoutRetryDelay(retryNumber)
			minDelay := time.Duration(float64(baseDelay) * 0.8)
			if delay < minDelay || delay > 30*time.Second {
				t.Errorf("retry %d delay %s outside expected range [%s, %s]", retryNumber, delay, minDelay, 30*time.Second)
			}
		}
	}
}

func TestHasTerminationTimeChanged(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		Old, New map[string]interface{}
		Expect   bool
	}{
		"empty": {
			Old:    map[string]interface{}{"termination_time": ""},
			New:    map[string]interface{}{"termination_time": ""},
			Expect: false,
		},
		"new": {
			Old:    map[string]interface{}{"termination_time": ""},
			New:    map[string]interface{}{"termination_time": "2025-01-31T15:04:05Z"},
			Expect: true,
		},
		"changed": {
			Old:    map[string]interface{}{"termination_time": "2025-01-30T15:04:05Z"},
			New:    map[string]interface{}{"termination_time": "2025-01-31T15:04:05Z"},
			Expect: true,
		},
		"same": {
			Old:    map[string]interface{}{"termination_time": "2025-01-30T15:04:05Z"},
			New:    map[string]interface{}{"termination_time": "2025-01-30T15:04:05Z"},
			Expect: false,
		},
	}
	for tn, tc := range cases {
		if hasTerminationTimeChanged(tc.Old, tc.New) != tc.Expect {
			t.Errorf("%s: expected %t for whether termination time matched for old = %q, new = %q", tn, tc.Expect, tc.Old, tc.New)
		}
	}
}
