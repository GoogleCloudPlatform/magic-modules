package compute

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
		"zone capacity exhausted operation wrapped by operation wait": {
			err: fmt.Errorf("Error waiting for instance to create: %w", ComputeOperationError(map[string]interface{}{
				"errors": []interface{}{map[string]interface{}{"code": "ZONE_RESOURCE_POOL_EXHAUSTED"}},
			})),
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

func TestComputeInstanceStockoutRetrySchemaHasNoDefaults(t *testing.T) {
	s := ResourceComputeInstance().Schema
	for _, field := range []string{"retry_on_stockout", "stockout_retry_timeout_seconds"} {
		if got := s[field].Default; got != nil {
			t.Errorf("%s default = %v, want nil", field, got)
		}
	}
}

func TestComputeInstanceStockoutRetryTimeout(t *testing.T) {
	unset := schema.TestResourceDataRaw(t, ResourceComputeInstance().Schema, map[string]interface{}{})
	if got := computeInstanceStockoutRetryTimeout(unset); got != 300*time.Second {
		t.Errorf("unset timeout = %s, want 300s", got)
	}

	set := schema.TestResourceDataRaw(t, ResourceComputeInstance().Schema, map[string]interface{}{
		"stockout_retry_timeout_seconds": 600,
	})
	if got := computeInstanceStockoutRetryTimeout(set); got != 600*time.Second {
		t.Errorf("configured timeout = %s, want 600s", got)
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

func newTestComputeInstanceCreateRetry(ctx context.Context) computeInstanceCreateRetry {
	return computeInstanceCreateRetry{
		ctx:                  ctx,
		createDeadline:       time.Now().Add(time.Minute),
		retryOnStockout:      true,
		stockoutRetryTimeout: time.Minute,
		retryDelay:           func(int) time.Duration { return 0 },
		insert:               func() (map[string]interface{}, error) { return map[string]interface{}{}, nil },
		waitForOperation:     func(map[string]interface{}, time.Duration) error { return nil },
	}
}

func computeInstanceStockoutTestError() error {
	return &googleapi.Error{Errors: []googleapi.ErrorItem{{Reason: "zoneResourcePoolExhausted"}}}
}

func TestRunComputeInstanceCreateWithStockoutRetry_RetriesInsertStockout(t *testing.T) {
	t.Parallel()
	inserts := 0
	r := newTestComputeInstanceCreateRetry(context.Background())
	r.insert = func() (map[string]interface{}, error) {
		inserts++
		if inserts < 3 {
			return nil, computeInstanceStockoutTestError()
		}
		return map[string]interface{}{}, nil
	}
	if err := runComputeInstanceCreateWithStockoutRetry(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inserts != 3 {
		t.Errorf("insert attempts = %d, want 3", inserts)
	}
}

func TestRunComputeInstanceCreateWithStockoutRetry_RetriesOperationStockout(t *testing.T) {
	t.Parallel()
	inserts, waits := 0, 0
	r := newTestComputeInstanceCreateRetry(context.Background())
	r.insert = func() (map[string]interface{}, error) {
		inserts++
		return map[string]interface{}{}, nil
	}
	r.waitForOperation = func(map[string]interface{}, time.Duration) error {
		waits++
		if waits == 1 {
			return computeInstanceStockoutTestError()
		}
		return nil
	}
	if err := runComputeInstanceCreateWithStockoutRetry(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inserts != 2 {
		t.Errorf("insert attempts = %d, want 2", inserts)
	}
}

func TestRunComputeInstanceCreateWithStockoutRetry_NonStockoutNotRetried(t *testing.T) {
	t.Parallel()
	inserts := 0
	r := newTestComputeInstanceCreateRetry(context.Background())
	r.insert = func() (map[string]interface{}, error) {
		inserts++
		return nil, errors.New("permission denied")
	}
	err := runComputeInstanceCreateWithStockoutRetry(r)
	if err == nil || err.Error() != "Error creating instance: permission denied" {
		t.Errorf("error = %v, want Error creating instance: permission denied", err)
	}
	if inserts != 1 {
		t.Errorf("insert attempts = %d, want 1", inserts)
	}
}

func TestRunComputeInstanceCreateWithStockoutRetry_DisabledNotRetried(t *testing.T) {
	t.Parallel()
	inserts := 0
	r := newTestComputeInstanceCreateRetry(context.Background())
	r.retryOnStockout = false
	r.insert = func() (map[string]interface{}, error) {
		inserts++
		return nil, computeInstanceStockoutTestError()
	}
	err := runComputeInstanceCreateWithStockoutRetry(r)
	if err == nil || !strings.HasPrefix(err.Error(), "Error creating instance: ") {
		t.Errorf("error = %v, want Error creating instance prefix", err)
	}
	if inserts != 1 {
		t.Errorf("insert attempts = %d, want 1", inserts)
	}
}

func TestRunComputeInstanceCreateWithStockoutRetry_StockoutRetryTimeout(t *testing.T) {
	t.Parallel()
	inserts := 0
	r := newTestComputeInstanceCreateRetry(context.Background())
	r.stockoutRetryTimeout = 250 * time.Millisecond
	r.retryDelay = func(int) time.Duration { return 100 * time.Millisecond }
	r.insert = func() (map[string]interface{}, error) {
		inserts++
		return nil, computeInstanceStockoutTestError()
	}
	err := runComputeInstanceCreateWithStockoutRetry(r)
	if inserts < 2 {
		t.Errorf("insert attempts = %d, want at least 2 before timeout", inserts)
	}
	wantPrefix := fmt.Sprintf("Error creating instance after %d attempts due to zone resource exhaustion: ", inserts)
	if err == nil || !strings.HasPrefix(err.Error(), wantPrefix) {
		t.Errorf("error = %v, want prefix %q", err, wantPrefix)
	}
	var apiErr *googleapi.Error
	if !errors.As(err, &apiErr) {
		t.Errorf("error = %v, want wrapped *googleapi.Error", err)
	}
}

func TestRunComputeInstanceCreateWithStockoutRetry_OperationStockoutRetryTimeout(t *testing.T) {
	t.Parallel()
	inserts := 0
	opErr := ComputeOperationError(map[string]interface{}{
		"errors": []interface{}{map[string]interface{}{"code": "ZONE_RESOURCE_POOL_EXHAUSTED"}},
	})
	r := newTestComputeInstanceCreateRetry(context.Background())
	r.stockoutRetryTimeout = 50 * time.Millisecond
	r.retryDelay = func(int) time.Duration { return 20 * time.Millisecond }
	r.insert = func() (map[string]interface{}, error) {
		inserts++
		return map[string]interface{}{}, nil
	}
	r.waitForOperation = func(map[string]interface{}, time.Duration) error { return opErr }
	err := runComputeInstanceCreateWithStockoutRetry(r)
	wantPrefix := fmt.Sprintf("Error creating instance after %d attempts due to zone resource exhaustion: ", inserts)
	if err == nil || !strings.HasPrefix(err.Error(), wantPrefix) {
		t.Errorf("error = %v, want prefix %q", err, wantPrefix)
	}
	var gotOpErr ComputeOperationError
	if !errors.As(err, &gotOpErr) {
		t.Errorf("error = %v, want wrapped ComputeOperationError", err)
	}
}

func TestRunComputeInstanceCreateWithStockoutRetry_FinalAttemptAtRetryDeadline(t *testing.T) {
	t.Parallel()
	inserts := 0
	r := newTestComputeInstanceCreateRetry(context.Background())
	r.stockoutRetryTimeout = 50 * time.Millisecond
	r.retryDelay = func(int) time.Duration { return time.Hour }
	r.insert = func() (map[string]interface{}, error) {
		inserts++
		if inserts == 1 {
			return nil, computeInstanceStockoutTestError()
		}
		return map[string]interface{}{}, nil
	}
	if err := runComputeInstanceCreateWithStockoutRetry(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inserts != 2 {
		t.Errorf("insert attempts = %d, want 2", inserts)
	}
}

func TestRunComputeInstanceCreateWithStockoutRetry_CreateDeadlineBounds(t *testing.T) {
	t.Parallel()
	r := newTestComputeInstanceCreateRetry(context.Background())
	r.createDeadline = time.Now().Add(time.Millisecond)
	r.insert = func() (map[string]interface{}, error) {
		return nil, computeInstanceStockoutTestError()
	}
	if err := runComputeInstanceCreateWithStockoutRetry(r); err == nil {
		t.Error("expected an error once the create deadline passes")
	}
}

func TestRunComputeInstanceCreateWithStockoutRetry_CanceledBeforeInsert(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	inserts := 0
	r := newTestComputeInstanceCreateRetry(ctx)
	r.insert = func() (map[string]interface{}, error) {
		inserts++
		return map[string]interface{}{}, nil
	}
	if err := runComputeInstanceCreateWithStockoutRetry(r); !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
	if inserts != 0 {
		t.Errorf("insert attempts = %d, want 0", inserts)
	}
}

func TestRunComputeInstanceCreateWithStockoutRetry_CanceledDuringDelay(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inserts := 0
	r := newTestComputeInstanceCreateRetry(ctx)
	r.retryDelay = func(int) time.Duration { return time.Hour }
	r.insert = func() (map[string]interface{}, error) {
		inserts++
		cancel()
		return nil, computeInstanceStockoutTestError()
	}
	if err := runComputeInstanceCreateWithStockoutRetry(r); !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
	if inserts != 1 {
		t.Errorf("insert attempts = %d, want 1", inserts)
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
