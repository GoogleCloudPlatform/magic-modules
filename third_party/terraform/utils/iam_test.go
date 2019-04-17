package google

import (
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
	"testing"
)

// Mock ResourceIamUpdater
type readTestUpdater struct {
	errorCode           int
	returnAfterAttempts int
	numAttempts         int
}

// Not being tested
func (u *readTestUpdater) SetResourceIamPolicy(p *cloudresourcemanager.Policy) error { return nil }
func (u *readTestUpdater) GetMutexKey() string                                       { return "a-resource-key" }
func (u *readTestUpdater) GetResourceId() string                                     { return "id" }
func (u *readTestUpdater) DescribeResource() string                                  { return "mock updater" }

// Fetch the existing IAM policy attached to a resource.
func (u *readTestUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	u.numAttempts++
	if u.numAttempts < u.returnAfterAttempts {
		// Indicate retry
		return nil, &googleapi.Error{Code: 429}
	}

	// Was testing successes or was testing retryable
	if u.errorCode == 200 {
		return &cloudresourcemanager.Policy{}, nil
	}
	return nil, &googleapi.Error{Code: u.errorCode}
}

func TestIamPolicyReadWithRetry_returnImmediately(t *testing.T) {
	mockUpdater := &readTestUpdater{
		returnAfterAttempts: 1,
		errorCode:           200,
	}
	p, err := iamPolicyReadWithRetry(mockUpdater)
	if err != nil || p == nil {
		t.Errorf("expected valid policy and no error, got nil policy and error %v", err)
	}
	if mockUpdater.numAttempts != 1 {
		t.Errorf("expected GetResourceIamPolicy to have been called once")
	}
}

func TestIamPolicyReadWithRetry_retry(t *testing.T) {
	mockUpdater := &readTestUpdater{
		returnAfterAttempts: 3,
		errorCode:           404,
	}
	p, err := iamPolicyReadWithRetry(mockUpdater)
	if err == nil || !isGoogleApiErrorWithCode(err, 404) {
		t.Errorf("expected googleapi error 404, got policy %v, err %v", p, err)
	}
	if mockUpdater.numAttempts != mockUpdater.returnAfterAttempts {
		t.Errorf("expected GetResourceIamPolicy to have been called %d times, was called %d", mockUpdater.numAttempts, mockUpdater.returnAfterAttempts)
	}
}
