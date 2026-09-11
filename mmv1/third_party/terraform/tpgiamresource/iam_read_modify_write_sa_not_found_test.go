package tpgiamresource

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
)

// serviceAccountNotFoundErr constructs an error matching
// transport_tpg.IamServiceAccountNotFound's predicate, to trigger
// iamPolicyReadModifyWrite's "service account not found" retry branch.
func serviceAccountNotFoundErr() error {
	return &googleapi.Error{
		Code: 400,
		Body: "Service account foo@bar.iam.gserviceaccount.com does not exist.",
	}
}

// saNotFoundFakeUpdater is a ResourceIamUpdater whose SetResourceIamPolicy
// always fails with a "service account not found" error, used to exercise
// iamPolicyReadModifyWrite's rarely-hit recovery branch for that error.
type saNotFoundFakeUpdater struct {
	mutexKey string
	resource string

	getPolicy func(callNum int32) (*cloudresourcemanager.Policy, error)

	getCount int32
	setCount int32
}

func (u *saNotFoundFakeUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	n := atomic.AddInt32(&u.getCount, 1)
	return u.getPolicy(n)
}

func (u *saNotFoundFakeUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	atomic.AddInt32(&u.setCount, 1)
	return serviceAccountNotFoundErr()
}

func (u *saNotFoundFakeUpdater) GetMutexKey() string {
	return u.mutexKey
}

func (u *saNotFoundFakeUpdater) GetResourceId() string {
	return u.resource
}

func (u *saNotFoundFakeUpdater) DescribeResource() string {
	return u.resource
}

// runWithTimeout runs f in a goroutine and reports whether it completed
// within the given timeout, to detect a deadlock without hanging the test
// suite forever if the bug being tested for is present. A panic inside f is
// recovered (a panic in a spawned goroutine would otherwise crash the whole
// test binary, bypassing any recover() in the calling goroutine) and fails
// the test immediately via t.Fatalf.
func runWithTimeout(t *testing.T, timeout time.Duration, f func() error) (err error, completed bool) {
	t.Helper()
	done := make(chan error, 1)
	panicked := make(chan interface{}, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicked <- r
			}
		}()
		done <- f()
	}()
	select {
	case err := <-done:
		return err, true
	case r := <-panicked:
		t.Fatalf("function panicked: %v", r)
		return nil, true
	case <-time.After(timeout):
		return nil, false
	}
}

// TestIamPolicyReadModifyWrite_ServiceAccountNotFound_EtagUnchanged verifies
// that when SetResourceIamPolicy fails with a "service account not found"
// error and the policy is unchanged since it was first read,
// iamPolicyReadModifyWrite returns promptly with an error (rather than
// deadlocking, which is what the unfixed code does: it re-locks the same
// resource mutex it's already holding).
func TestIamPolicyReadModifyWrite_ServiceAccountNotFound_EtagUnchanged(t *testing.T) {
	updater := &saNotFoundFakeUpdater{
		mutexKey: "iam-project-my-project",
		resource: "project my-project",
		getPolicy: func(callNum int32) (*cloudresourcemanager.Policy, error) {
			// Etag never changes across calls.
			return &cloudresourcemanager.Policy{Etag: "etag-1"}, nil
		},
	}

	err, completed := runWithTimeout(t, 5*time.Second, func() error {
		return iamPolicyReadModifyWrite(updater, func(p *cloudresourcemanager.Policy) error { return nil })
	})

	if !completed {
		t.Fatal("iamPolicyReadModifyWrite did not return within timeout: deadlocked re-acquiring its own resource mutex")
	}
	if err == nil {
		t.Fatal("expected an error to be returned (service account not found, etag unchanged), got nil")
	}
	// It should have given up after a single recheck, not retried forever.
	if got := atomic.LoadInt32(&updater.getCount); got != 2 {
		t.Errorf("expected exactly 2 GetResourceIamPolicy calls (initial read + one recheck), got %d", got)
	}
	if got := atomic.LoadInt32(&updater.setCount); got != 1 {
		t.Errorf("expected exactly 1 SetResourceIamPolicy call, got %d", got)
	}
}

// TestIamPolicyReadModifyWrite_ServiceAccountNotFound_RecheckReadFails
// verifies that when the recheck read itself fails, iamPolicyReadModifyWrite
// falls through and returns the original error rather than panicking on a
// nil policy dereference (the unfixed code's condition is inverted, so it
// attempts to read .Etag off a nil *cloudresourcemanager.Policy in this
// case) or deadlocking.
func TestIamPolicyReadModifyWrite_ServiceAccountNotFound_RecheckReadFails(t *testing.T) {
	updater := &saNotFoundFakeUpdater{
		mutexKey: "iam-project-my-project",
		resource: "project my-project",
		getPolicy: func(callNum int32) (*cloudresourcemanager.Policy, error) {
			if callNum == 1 {
				// initial read at the top of the read-modify-write loop
				return &cloudresourcemanager.Policy{Etag: "etag-1"}, nil
			}
			// the recheck read fails
			return nil, fmt.Errorf("transient error rechecking policy")
		},
	}

	err, completed := runWithTimeout(t, 5*time.Second, func() error {
		return iamPolicyReadModifyWrite(updater, func(p *cloudresourcemanager.Policy) error { return nil })
	})

	if !completed {
		t.Fatal("iamPolicyReadModifyWrite did not return within timeout: deadlocked re-acquiring its own resource mutex")
	}
	if err == nil {
		t.Fatal("expected an error to be returned (service account not found, recheck read failed), got nil")
	}
}
