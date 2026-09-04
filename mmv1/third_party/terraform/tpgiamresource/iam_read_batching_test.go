package tpgiamresource

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// fakeResourceIamUpdater is a minimal ResourceIamUpdater used to exercise
// iamPolicyReadWithRetry's batching behavior without any real API calls.
type fakeResourceIamUpdater struct {
	mutexKey  string
	resource  string
	getCount  int32
	getPolicy func() (*cloudresourcemanager.Policy, error)
}

func (u *fakeResourceIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	atomic.AddInt32(&u.getCount, 1)
	return u.getPolicy()
}

func (u *fakeResourceIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	return fmt.Errorf("not implemented")
}

func (u *fakeResourceIamUpdater) GetMutexKey() string {
	return u.mutexKey
}

func (u *fakeResourceIamUpdater) GetResourceId() string {
	return u.resource
}

func (u *fakeResourceIamUpdater) DescribeResource() string {
	return u.resource
}

func newTestConfigForIamBatching() *transport_tpg.Config {
	ctx := context.Background()
	batchingConfig := &transport_tpg.BatchingConfig{
		SendAfter:      100 * time.Millisecond,
		EnableBatching: true,
	}
	return &transport_tpg.Config{
		Context:           ctx,
		BatchingConfig:    batchingConfig,
		RequestBatcherIam: transport_tpg.NewRequestBatcher("IAM", ctx, batchingConfig),
	}
}

func TestIamPolicyReadWithRetry_CollapsesConcurrentCallsIntoOneRequest(t *testing.T) {
	policy := &cloudresourcemanager.Policy{Etag: "etag-1"}
	updater := &fakeResourceIamUpdater{
		mutexKey: "iam-project-my-project",
		resource: "project my-project",
		getPolicy: func() (*cloudresourcemanager.Policy, error) {
			return policy, nil
		},
	}
	config := newTestConfigForIamBatching()

	const numCallers = 10
	var wg sync.WaitGroup
	wg.Add(numCallers)
	errs := make([]error, numCallers)
	results := make([]*cloudresourcemanager.Policy, numCallers)

	for i := 0; i < numCallers; i++ {
		go func(idx int) {
			defer wg.Done()
			p, err := iamPolicyReadWithRetry(updater, config)
			errs[idx] = err
			results[idx] = p
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("caller %d: unexpected error: %v", i, err)
		}
		if results[i] != policy {
			t.Errorf("caller %d: expected shared policy pointer %p, got %p", i, policy, results[i])
		}
	}

	if got := atomic.LoadInt32(&updater.getCount); got != 1 {
		t.Errorf("expected exactly 1 GetResourceIamPolicy call for %d concurrent callers, got %d", numCallers, got)
	}
}

func TestIamPolicyReadWithRetry_SeparateResourcesAreNotCombined(t *testing.T) {
	config := newTestConfigForIamBatching()

	updaters := []*fakeResourceIamUpdater{
		{
			mutexKey: "iam-project-project-a",
			resource: "project project-a",
			getPolicy: func() (*cloudresourcemanager.Policy, error) {
				return &cloudresourcemanager.Policy{Etag: "etag-a"}, nil
			},
		},
		{
			mutexKey: "iam-project-project-b",
			resource: "project project-b",
			getPolicy: func() (*cloudresourcemanager.Policy, error) {
				return &cloudresourcemanager.Policy{Etag: "etag-b"}, nil
			},
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(updaters))
	for _, u := range updaters {
		go func(u *fakeResourceIamUpdater) {
			defer wg.Done()
			if _, err := iamPolicyReadWithRetry(u, config); err != nil {
				t.Errorf("resource %s: unexpected error: %v", u.resource, err)
			}
		}(u)
	}
	wg.Wait()

	for _, u := range updaters {
		if got := atomic.LoadInt32(&u.getCount); got != 1 {
			t.Errorf("resource %s: expected exactly 1 GetResourceIamPolicy call, got %d", u.resource, got)
		}
	}
}
