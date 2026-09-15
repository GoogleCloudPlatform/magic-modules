package apphub

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestResourceApphubServiceCreate_LeaseConflictRetry(t *testing.T) {
	var postCount int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/services"):
			count := atomic.AddInt32(&postCount, 1)
			fmt.Fprintf(w, `{"name": "projects/test-project/locations/us-central1/operations/op-%d", "done": false}`, count)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/operations/op-1"):
			fmt.Fprint(w, `{"name": "projects/test-project/locations/us-central1/operations/op-1", "done": true, "error": {"code": 9, "message": "failed to add registration status to discovered service since the discovered service is under lease"}}`)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/operations/op-2"):
			fmt.Fprint(w, `{"name": "projects/test-project/locations/us-central1/operations/op-2", "done": true, "response": {"name": "projects/test-project/locations/us-central1/applications/test-app/services/test-svc"}}`)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/services/test-svc"):
			fmt.Fprint(w, `{"name": "projects/test-project/locations/us-central1/applications/test-app/services/test-svc", "uid": "test-uid"}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	config := &transport_tpg.Config{
		Project: "test-project",
		CustomEndpoints: map[string]string{
			Product.CustomEndpointField: ts.URL + "/",
		},
		PollInterval: 10 * time.Millisecond,
		Client:       ts.Client(),
	}

	d := ResourceApphubService().TestResourceData()
	_ = d.Set("project", "test-project")
	_ = d.Set("location", "us-central1")
	_ = d.Set("application_id", "test-app")
	_ = d.Set("service_id", "test-svc")
	_ = d.Set("discovered_service", "projects/12345/locations/us-central1/discoveredServices/ds-1")

	err := resourceApphubServiceCreate(d, config)
	if err != nil {
		t.Fatalf("expected nil error after lease conflict retry, got: %v", err)
	}
	if got := atomic.LoadInt32(&postCount); got != 2 {
		t.Errorf("expected POST to be called 2 times, got %d", got)
	}
	expectedID := "projects/test-project/locations/us-central1/applications/test-app/services/test-svc"
	if d.Id() != expectedID {
		t.Errorf("expected resource ID %q, got %q", expectedID, d.Id())
	}
}

func TestResourceApphubWorkloadCreate_LeaseConflictRetry(t *testing.T) {
	var postCount int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/workloads"):
			count := atomic.AddInt32(&postCount, 1)
			fmt.Fprintf(w, `{"name": "projects/test-project/locations/us-central1/operations/op-%d", "done": false}`, count)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/operations/op-1"):
			fmt.Fprint(w, `{"name": "projects/test-project/locations/us-central1/operations/op-1", "done": true, "error": {"code": 9, "message": "failed to add registration status to discovered workload since the discovered workload is under lease"}}`)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/operations/op-2"):
			fmt.Fprint(w, `{"name": "projects/test-project/locations/us-central1/operations/op-2", "done": true, "response": {"name": "projects/test-project/locations/us-central1/applications/test-app/workloads/test-workload"}}`)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/workloads/test-workload"):
			fmt.Fprint(w, `{"name": "projects/test-project/locations/us-central1/applications/test-app/workloads/test-workload", "uid": "test-uid"}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	config := &transport_tpg.Config{
		Project: "test-project",
		CustomEndpoints: map[string]string{
			Product.CustomEndpointField: ts.URL + "/",
		},
		PollInterval: 10 * time.Millisecond,
		Client:       ts.Client(),
	}

	d := ResourceApphubWorkload().TestResourceData()
	_ = d.Set("project", "test-project")
	_ = d.Set("location", "us-central1")
	_ = d.Set("application_id", "test-app")
	_ = d.Set("workload_id", "test-workload")
	_ = d.Set("discovered_workload", "projects/12345/locations/us-central1/discoveredWorkloads/dw-1")

	err := resourceApphubWorkloadCreate(d, config)
	if err != nil {
		t.Fatalf("expected nil error after lease conflict retry, got: %v", err)
	}
	if got := atomic.LoadInt32(&postCount); got != 2 {
		t.Errorf("expected POST to be called 2 times, got %d", got)
	}
	expectedID := "projects/test-project/locations/us-central1/applications/test-app/workloads/test-workload"
	if d.Id() != expectedID {
		t.Errorf("expected resource ID %q, got %q", expectedID, d.Id())
	}
}

func TestResourceApphubServiceCreate_NonLeaseErrorDoesNotRetry(t *testing.T) {
	var postCount int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/services"):
			count := atomic.AddInt32(&postCount, 1)
			fmt.Fprintf(w, `{"name": "projects/test-project/locations/us-central1/operations/op-%d", "done": false}`, count)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/operations/op-1"):
			fmt.Fprint(w, `{"name": "projects/test-project/locations/us-central1/operations/op-1", "done": true, "error": {"code": 3, "message": "invalid resource name"}}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	config := &transport_tpg.Config{
		Project: "test-project",
		CustomEndpoints: map[string]string{
			Product.CustomEndpointField: ts.URL + "/",
		},
		PollInterval: 10 * time.Millisecond,
		Client:       ts.Client(),
	}

	d := ResourceApphubService().TestResourceData()
	_ = d.Set("project", "test-project")
	_ = d.Set("location", "us-central1")
	_ = d.Set("application_id", "test-app")
	_ = d.Set("service_id", "test-svc")
	_ = d.Set("discovered_service", "projects/12345/locations/us-central1/discoveredServices/ds-1")

	err := resourceApphubServiceCreate(d, config)
	if err == nil {
		t.Fatalf("expected error for non-lease failure, got nil")
	}
	if got := atomic.LoadInt32(&postCount); got != 1 {
		t.Errorf("expected POST to be called exactly 1 time (no retry), got %d", got)
	}
	if d.Id() != "" {
		t.Errorf("expected empty resource ID on failure, got %q", d.Id())
	}
}
