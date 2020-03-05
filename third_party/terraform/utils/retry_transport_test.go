package google

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/api/googleapi"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const testRetryTransportErrorMessageRetry = "retry error"
const testRetryTransportErrorMessageSuccess = "success"
const testRetryTransportErrorMessageFailure = "fail the request"

// Check for no errors if the request succeeds the first time
func TestRetryTransport_SingleRequestSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			if _, err := w.Write([]byte(testRetryTransportErrorMessageSuccess)); err != nil {
				t.Errorf("unable to write to response writer: %s", err)
			}
		}))
	defer ts.Close()

	client := ts.Client()
	client.Transport = &retryTransport{
		internal:        http.DefaultTransport,
		retryPredicates: []RetryErrorPredicateFunc{testRetryTransportRetryPredicate},
	}

	resp, err := client.Get(ts.URL)
	testRetryTransportCheckResponseForSuccess(t, resp, err)
}

// Check for error if the request fails the first time
func TestRetryTransport_SingleRequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
			w.Write([]byte(testRetryTransportErrorMessageFailure))
		}))
	defer ts.Close()

	client := ts.Client()
	client.Transport = &retryTransport{
		internal:        http.DefaultTransport,
		retryPredicates: []RetryErrorPredicateFunc{testRetryTransportRetryPredicate},
	}

	resp, err := client.Get(ts.URL)
	testRetryTransportCheckResponseForFailure(t, resp, err, 400, testRetryTransportErrorMessageFailure)
}

// Check for no errors if the request succeeds after a certain amount of time
func TestRetryTransport_SuccessAfterRetries(t *testing.T) {
	ts := httptest.NewServer(testRetryTransportSucceedAfterHandler(t, time.Second*1))
	defer ts.Close()

	client := ts.Client()
	client.Transport = &retryTransport{
		internal:        http.DefaultTransport,
		retryPredicates: []RetryErrorPredicateFunc{testRetryTransportRetryPredicate},
	}

	ctx, cc := context.WithTimeout(context.Background(), time.Second*2)
	defer cc()
	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("unable to construct err: %v", err)
	}

	resp, err := client.Do(req)
	testRetryTransportCheckResponseForSuccess(t, resp, err)
}

func TestRetryTransport_FailAfterRetries(t *testing.T) {
	ts := httptest.NewServer(testRetryTransportFailAfterHandler(t, time.Second*1))
	defer ts.Close()

	client := ts.Client()
	client.Transport = &retryTransport{
		internal:        http.DefaultTransport,
		retryPredicates: []RetryErrorPredicateFunc{testRetryTransportRetryPredicate},
	}

	ctx, cc := context.WithTimeout(context.Background(), time.Second*2)
	defer cc()
	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("unable to construct err: %v", err)
	}

	resp, err := client.Do(req)
	testRetryTransportCheckResponseForFailure(t, resp, err, 400, testRetryTransportErrorMessageFailure)
}

func TestRetryTransport_ContextTimeout(t *testing.T) {
	ts := httptest.NewServer(testRetryTransportSucceedAfterHandler(t, time.Second*4))
	defer ts.Close()

	client := ts.Client()
	client.Transport = &retryTransport{
		internal:        http.DefaultTransport,
		retryPredicates: []RetryErrorPredicateFunc{testRetryTransportRetryPredicate},
	}

	ctx, cc := context.WithTimeout(context.Background(), time.Second*2)
	defer cc()
	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("unable to construct err: %v", err)
	}
	resp, err := client.Do(req)
	// Last failure should have been a retryable error since we timed out
	testRetryTransportCheckResponseForFailure(t, resp, err, 500, testRetryTransportErrorMessageRetry)
}

// Check for no errors if the request succeeds after a certain amount of time
func TestRetryTransport_SuccessAfterRetriesWithBody(t *testing.T) {
	ts := httptest.NewServer(testRetryTransportSuccessCheckBodyHandler(t, time.Second*1))
	defer ts.Close()

	client := ts.Client()
	client.Transport = &retryTransport{
		internal:        http.DefaultTransport,
		retryPredicates: []RetryErrorPredicateFunc{testRetryTransportRetryPredicate},
	}

	msg := "body for successful request"
	ctx, cc := context.WithTimeout(context.Background(), time.Second*2)
	defer cc()

	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, bytes.NewReader([]byte(msg)))
	if err != nil {
		t.Fatalf("unable to construct err: %v", err)
	}

	resp, err := client.Do(req)
	testRetryTransportCheckResponseForSuccess(t, resp, err)
}

// SUCCESS handlers and check
func testRetryTransportSucceedAfterHandler(t *testing.T, successInterval time.Duration) http.Handler {
	var firstReqTime time.Time
	var testOnce sync.Once

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testOnce.Do(func() {
			firstReqTime = time.Now()
		})
		if time.Since(firstReqTime) >= successInterval {
			w.WriteHeader(200)
			if _, err := w.Write([]byte(testRetryTransportErrorMessageSuccess)); err != nil {
				t.Errorf("[ERROR] unable to write to response writer: %v", err)
			}
		} else {
			w.WriteHeader(500)
			if _, err := w.Write([]byte(testRetryTransportErrorMessageRetry)); err != nil {
				t.Errorf("[ERROR] unable to write to response writer: %v", err)
			}
		}
	})
}

func testRetryTransportSuccessCheckBodyHandler(t *testing.T, successInterval time.Duration) http.Handler {
	var firstReqTime time.Time
	var testOnce sync.Once

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testOnce.Do(func() {
			firstReqTime = time.Now()
		})

		slurp, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			if _, err := w.Write([]byte(fmt.Sprintf("unable to read request body: %v", err))); err != nil {
				t.Errorf("[ERROR] unable to write to response writer: %v", err)
			}
			return
		}

		if time.Since(firstReqTime) >= successInterval {
			w.WriteHeader(200)
			resp := fmt.Sprintf("%s\nRequest Body: %s", testRetryTransportErrorMessageSuccess, string(slurp))
			if _, err := w.Write([]byte(resp)); err != nil {
				t.Errorf("[ERROR] unable to write to response writer: %v", err)
			}
		} else {
			w.WriteHeader(500)
			resp := fmt.Sprintf("%s\nRequest Body: %s", testRetryTransportErrorMessageRetry, string(slurp))
			if _, err := w.Write([]byte(resp)); err != nil {
				t.Errorf("[ERROR] unable to write to response writer: %v", err)
			}
		}
	})
}

func testRetryTransportCheckResponseForSuccess(t *testing.T, resp *http.Response, respErr error) {
	if respErr != nil {
		t.Fatalf("expected no error, got: %v", respErr)
	}

	err := googleapi.CheckResponse(resp)
	if err != nil {
		t.Fatalf("expected no error, got response error: %v", err)
	}
}

// FAILURE handler and check
func testRetryTransportFailAfterHandler(t *testing.T, successInterval time.Duration) http.Handler {
	var firstReqTime time.Time
	var testOnce sync.Once

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testOnce.Do(func() {
			firstReqTime = time.Now()
		})
		if time.Since(firstReqTime) >= successInterval {
			w.WriteHeader(400)
			if _, err := w.Write([]byte(testRetryTransportErrorMessageFailure)); err != nil {
				t.Errorf("[ERROR] unable to write to response writer: %v", err)
			}
		} else {
			w.WriteHeader(500)
			if _, err := w.Write([]byte(testRetryTransportErrorMessageRetry)); err != nil {
				t.Errorf("[ERROR] unable to write to response writer: %v", err)
			}
		}
	})
}

func testRetryTransportCheckResponseForFailure(t *testing.T, resp *http.Response, respErr error, expectedCode int, expectedMsg string) {
	if respErr != nil {
		t.Fatalf("expected response error, got actual error for doing request: %v", respErr)
	}

	err := googleapi.CheckResponse(resp)
	if err == nil {
		t.Fatalf("expected googleapi error, got no error")
	}

	gerr, ok := err.(*googleapi.Error)
	if !ok {
		t.Fatalf("expected error to be googleapi error: %v", err)
	}

	if gerr.Code != expectedCode {
		t.Errorf("expected error code 400, got error: %v", err)
	}

	if !strings.Contains(gerr.Body, expectedMsg) {
		t.Errorf("expected error %q in %v", testRetryTransportErrorMessageFailure, err)
	}
}

// ERROR RETRY PREDICATE
// Retries 500.
func testRetryTransportRetryPredicate(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 500 && gerr.Message == "retry error" {
			return true, "retry error"
		}
	}
	return false, ""
}
