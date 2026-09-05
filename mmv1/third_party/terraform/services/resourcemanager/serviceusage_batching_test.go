package resourcemanager

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func newTestConfigForProjectBatching(ts *httptest.Server) *transport_tpg.Config {
	ctx := context.Background()
	batchingConfig := &transport_tpg.BatchingConfig{
		SendAfter:      100 * time.Millisecond,
		EnableBatching: true,
	}
	return &transport_tpg.Config{
		Context: ctx,
		Client:  ts.Client(),
		CustomEndpoints: map[string]string{
			// Matches resourcemanager's Product.CustomEndpointField.
			"resource_manager_custom_endpoint": ts.URL + "/",
		},
		BatchingConfig:             batchingConfig,
		RequestBatcherServiceUsage: transport_tpg.NewRequestBatcher("Service Usage", ctx, batchingConfig),
	}
}

func testProjectServiceResourceData(t *testing.T, project, service string) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, ResourceGoogleProjectService().Schema, map[string]interface{}{
		"project": project,
		"service": service,
	})
}

func TestBatchRequestReadProject_CollapsesConcurrentCallsIntoOneRequest(t *testing.T) {
	const project = "my-project"
	var requestCount int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"projectId": %q, "lifecycleState": "ACTIVE"}`, project)
	}))
	defer ts.Close()

	config := newTestConfigForProjectBatching(ts)

	const numCallers = 10
	var wg sync.WaitGroup
	wg.Add(numCallers)
	errs := make([]error, numCallers)
	lifecycleStates := make([]string, numCallers)

	for i := 0; i < numCallers; i++ {
		go func(idx int) {
			defer wg.Done()
			d := testProjectServiceResourceData(t, project, "foo.googleapis.com")
			p, err := BatchRequestReadProject(project, d, config)
			errs[idx] = err
			if err == nil {
				lifecycleStates[idx] = p.LifecycleState
			}
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("caller %d: unexpected error: %v", i, err)
		}
		if lifecycleStates[i] != "ACTIVE" {
			t.Errorf("caller %d: expected lifecycleState ACTIVE, got %q", i, lifecycleStates[i])
		}
	}

	if got := atomic.LoadInt32(&requestCount); got != 1 {
		t.Errorf("expected exactly 1 HTTP request to be sent for %d concurrent callers, got %d", numCallers, got)
	}
}

func TestBatchRequestReadProject_SeparateProjectsAreNotCombined(t *testing.T) {
	var requestCount int32
	seenProjects := make(map[string]bool)
	var mu sync.Mutex

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		mu.Lock()
		seenProjects[r.URL.Path] = true
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"lifecycleState": "ACTIVE"}`)
	}))
	defer ts.Close()

	config := newTestConfigForProjectBatching(ts)

	projects := []string{"project-a", "project-b"}
	var wg sync.WaitGroup
	wg.Add(len(projects))
	for _, project := range projects {
		go func(project string) {
			defer wg.Done()
			d := testProjectServiceResourceData(t, project, "foo.googleapis.com")
			if _, err := BatchRequestReadProject(project, d, config); err != nil {
				t.Errorf("project %s: unexpected error: %v", project, err)
			}
		}(project)
	}
	wg.Wait()

	if got := atomic.LoadInt32(&requestCount); got != int32(len(projects)) {
		t.Errorf("expected %d HTTP requests (one per distinct project), got %d", len(projects), got)
	}
}
