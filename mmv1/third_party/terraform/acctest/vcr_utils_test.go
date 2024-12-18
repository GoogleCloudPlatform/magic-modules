package acctest_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestNewVcrMatcherFunc_canDetectMatches(t *testing.T) {

	// Everything should be determined as a match in this test
	cases := map[string]struct {
		httpRequest     requestDescription
		cassetteRequest requestDescription
	}{
		"matches POST requests with empty bodies": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{}",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{}",
			},
		},
		"matches POST requests with exact matching bodies": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value\"}",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value\"}",
			},
		},
		"matches GET requests when SendRequest adds alt=json but cassette has alt=json&prettyPrint=false": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "compute.googleapis.com",
				path:   "compute/v1/projects/my-project/zones/us-central1-a/instances/my-instance",
				query:  "alt=json",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "compute.googleapis.com",
				path:   "compute/v1/projects/my-project/zones/us-central1-a/instances/my-instance",
				query:  "alt=json&prettyPrint=false",
			},
		},
		"matches GET requests when both have alt and prettyPrint in different order": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "compute.googleapis.com",
				path:   "compute/v1/projects/my-project/zones/us-central1-a/instances/my-instance",
				query:  "prettyPrint=false&alt=json",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "compute.googleapis.com",
				path:   "compute/v1/projects/my-project/zones/us-central1-a/instances/my-instance",
				query:  "alt=json&prettyPrint=false",
			},
		},
		"matches POST requests with matching but re-ordered bodies, but only if Content-Type contains 'application/json'": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: "{\"field1\":\"value1\",\"field2\":\"value2\"}", // 1 before 2
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: "{\"field2\":\"value2\",\"field1\":\"value1\"}", // 2 before 1
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Make matcher
			ctx := context.Background()
			req := prepareHttpRequest(tc.httpRequest)
			cassetteReq := prepareCassetteRequest(tc.cassetteRequest)
			matcher := acctest.NewVcrMatcherFunc(ctx)

			// Act - use matcher
			matchDetected := matcher(req, cassetteReq)

			// Assert match
			if !matchDetected {
				t.Fatalf("expected matcher to match the requests")
			}
		})
	}
}

func TestNewVcrMatcherFunc_canDetectMismatches(t *testing.T) {

	// All these cases are expected to end with no match detected
	cases := map[string]struct {
		httpRequest     requestDescription
		cassetteRequest requestDescription
	}{
		"requests using different schemes": {
			httpRequest: requestDescription{
				scheme: "http",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
			},
		},
		"requests using different hosts": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "google.com",
				path:   "foobar",
			},
		},
		"requests using different paths": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar1",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar2",
			},
		},
		"requests with different methods": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{}",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "PUT",
				host:   "example.com",
				path:   "foobar",
				body:   "{}",
			},
		},
		"POST requests with different bodies": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value is ABCDEFG\"}",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value is MNLOP\"}",
			},
		},
		"requests with different non-standard query params do not match": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "compute.googleapis.com",
				path:   "compute/v1/projects/my-project/zones/us-central1-a/instances/my-instance/detachDisk",
				query:  "deviceName=disk-1&alt=json",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "compute.googleapis.com",
				path:   "compute/v1/projects/my-project/zones/us-central1-a/instances/my-instance/detachDisk",
				query:  "deviceName=disk-2&alt=json&prettyPrint=false",
			},
		},
		"POST requests with matching but re-ordered bodies aren't matching if Content-Type header is not 'application/json'": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"Content-Type": "foobar",
				},
				body: "{\"field1\":\"value1\",\"field2\":\"value2\"}", // 1 before 2
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"Content-Type": "foobar",
				},
				body: "{\"field2\":\"value2\",\"field1\":\"value1\"}", // 2 before 1
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Make matcher
			ctx := context.Background()
			req := prepareHttpRequest(tc.httpRequest)
			cassetteReq := prepareCassetteRequest(tc.cassetteRequest)
			matcher := acctest.NewVcrMatcherFunc(ctx)

			// Act - use matcher
			matchDetected := matcher(req, cassetteReq)

			// Assert match
			if matchDetected {
				t.Fatalf("expected matcher to not match the requests")
			}
		})
	}
}

// Currently there is no code to actively force the matcher to overlook differing User-Agent values.
// It isn't checked at any point in the matcher logic.
func TestNewVcrMatcherFunc_ignoresDifferentUserAgents(t *testing.T) {

	cases := map[string]struct {
		httpRequest     requestDescription
		cassetteRequest requestDescription
	}{
		"GET requests with different useragents are matched": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"User-Agent": "user-agent-HTTP",
				},
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"User-Agent": "user-agent-CASSETTE",
				},
			},
		},
		"POST requests with identical bodies and different useragents are matched": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value\"}",
				headers: map[string]string{
					"User-Agent": "user-agent-HTTP",
				},
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value\"}",
				headers: map[string]string{
					"User-Agent": "user-agent-CASSETTE",
				},
			},
		},
		"POST requests with reordered but matching bodies and different useragents are matched if Content-Type contains 'application/json'": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field1\":\"value1\",\"field2\":\"value2\"}",
				headers: map[string]string{
					"User-Agent":   "user-agent-HTTP",
					"Content-Type": "application/json",
				},
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field2\":\"value2\",\"field1\":\"value1\"}",
				headers: map[string]string{
					"User-Agent":   "user-agent-CASSETTE",
					"Content-Type": "application/json",
				},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Make matcher
			ctx := context.Background()
			req := prepareHttpRequest(tc.httpRequest)
			cassetteReq := prepareCassetteRequest(tc.cassetteRequest)
			matcher := acctest.NewVcrMatcherFunc(ctx)

			// Act - use matcher
			matchDetected := matcher(req, cassetteReq)

			// Assert match
			if !matchDetected {
				t.Fatalf("expected matcher to match the requests")
			}
		})
	}
}

type requestDescription struct {
	scheme  string
	method  string
	host    string
	path    string
	body    string
	query   string
	headers map[string]string
}

func prepareHttpRequest(d requestDescription) *http.Request {
	url := &url.URL{
		Scheme:   d.scheme,
		Host:     d.host,
		Path:     d.path,
		RawQuery: d.query,
	}

	req := &http.Request{
		Method: d.method,
		URL:    url,
	}

	// Conditionally set a body
	if d.body != "" {
		body := io.NopCloser(bytes.NewBufferString(d.body))
		req.Body = body
	}
	// Conditionally set headers
	if len(d.headers) > 0 {
		req.Header = http.Header{}
		for k, v := range d.headers {
			req.Header.Set(k, v)
		}
	}

	return req
}

func prepareCassetteRequest(d requestDescription) cassette.Request {
	fullUrl := fmt.Sprintf("%s://%s/%s", d.scheme, d.host, d.path)
	if d.query != "" {
		fullUrl = fmt.Sprintf("%s?%s", fullUrl, d.query)
	}

	req := cassette.Request{
		Method: d.method,
		URL:    fullUrl,
	}

	// Conditionally set a body
	if d.body != "" {
		req.Body = d.body
	}
	// Conditionally set headers
	if len(d.headers) > 0 {
		req.Headers = http.Header{}
		for k, v := range d.headers {
			req.Headers.Add(k, v)
		}
	}

	return req
}

func TestIsVcrReplaying(t *testing.T) {
	cases := []struct {
		name     string
		vcrPath  string
		vcrMode  string
		expected bool
	}{
		{"empty", "", "", false},
		{"path only", "/tmp", "", false},
		{"mode recording", "/tmp", "RECORDING", false},
		{"mode replaying no path", "", "REPLAYING", false},
		{"mode replaying with path", "/tmp", "REPLAYING", true},
		{"mode other", "/tmp", "OTHER", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("VCR_PATH", tc.vcrPath)
			t.Setenv("VCR_MODE", tc.vcrMode)
			if got := acctest.IsVcrReplaying(); got != tc.expected {
				t.Errorf("IsVcrReplaying() = %v, expected %v", got, tc.expected)
			}
		})
	}
}

func TestRunVcrSubtest(t *testing.T) {
	origPath := os.Getenv("VCR_PATH")
	origMode := os.Getenv("VCR_MODE")
	defer func() {
		os.Setenv("VCR_PATH", origPath)
		os.Setenv("VCR_MODE", origMode)
	}()

	os.Setenv("VCR_PATH", "")
	os.Setenv("VCR_MODE", "")
	acctest.RunVcrSubtest(t, "disabled when not replaying", func(t *testing.T) {
		// Calling t.Parallel() here should not panic since RunVcrSubtest didn't call it
		t.Parallel()
	})

	os.Setenv("VCR_PATH", "/tmp")
	os.Setenv("VCR_MODE", "REPLAYING")
	acctest.RunVcrSubtest(t, "enabled when replaying", func(t *testing.T) {
		// In Go testing, calling t.Parallel() a second time on the same testing.T panics
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("expected panic when calling t.Parallel() a second time, but did not panic")
			}
		}()
		t.Parallel()
	})
}

func newFakeApiServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"test-resource"}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func setupWatcherTestProviderEnvs(t *testing.T) {
	t.Helper()
	acctest.UnsetTestProviderConfigEnvs(t)
	t.Setenv("VCR_PATH", "")
	t.Setenv("VCR_MODE", "")
	t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", "fake-access-token")
	t.Setenv("GOOGLE_PROJECT", "test-project")
	t.Setenv("GOOGLE_REGION", "us-central1")
	t.Setenv("GOOGLE_ZONE", "us-central1-a")
}

func TestWatchRequests_interceptsProviderAndGoogleProviderConfig(t *testing.T) {
	setupWatcherTestProviderEnvs(t)
	srv := newFakeApiServer(t)

	var seenPaths []string
	var seenVersions []string
	acctest.WatchRequests(t, func(req *http.Request) {
		seenPaths = append(seenPaths, req.URL.Path)
		seenVersions = append(seenVersions, req.Header.Get("X-Goog-Api-Version"))
	})

	// 1. Test via GetSDKProvider (used by MuxedProviders / VcrTest)
	prov := acctest.GetSDKProvider(t.Name())
	diags := prov.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if diags.HasError() {
		t.Fatalf("unexpected provider configure error: %v", diags)
	}
	sdkConfig := prov.Meta().(*transport_tpg.Config)

	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:     sdkConfig,
		Method:     "POST",
		Project:    "test-project",
		RawURL:     srv.URL + "/compute/v1/projects/test-project/global/urlMaps",
		UserAgent:  sdkConfig.UserAgent,
		ApiVersion: "2026-11-01",
	})
	if err != nil {
		t.Fatalf("SendRequest via SDK provider config failed: %v", err)
	}

	// 2. Test via GoogleProviderConfig (used by CheckDestroy and custom test checks)
	destroyConfig := acctest.GoogleProviderConfig(t)
	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:     destroyConfig,
		Method:     "GET",
		Project:    "test-project",
		RawURL:     srv.URL + "/compute/v1/projects/test-project/global/urlMaps/my-url-map",
		UserAgent:  destroyConfig.UserAgent,
		ApiVersion: "2026-11-01",
	})
	if err != nil {
		t.Fatalf("SendRequest via GoogleProviderConfig failed: %v", err)
	}

	if len(seenPaths) != 2 {
		t.Fatalf("expected 2 intercepted requests, got %d (%v)", len(seenPaths), seenPaths)
	}
	for i, v := range seenVersions {
		if v != "2026-11-01" {
			t.Errorf("request %d (%s) had X-Goog-Api-Version=%q, want %q", i, seenPaths[i], v, "2026-11-01")
		}
	}
}

func TestAssertRequestApiVersion_passesWhenMatchingRequestsHaveExpectedHeader(t *testing.T) {
	setupWatcherTestProviderEnvs(t)
	srv := newFakeApiServer(t)

	acctest.AssertRequestApiVersion(t, "projects/{{project}}/global/urlMaps", "2026-11-01")

	config := acctest.GoogleProviderConfig(t)

	// Non-matching requests (e.g. backendServices, healthChecks, operations) without ApiVersion should be ignored.
	for _, nonMatchingPath := range []string{
		"/compute/v1/projects/test-project/global/backendServices",
		"/compute/v1/projects/test-project/global/healthChecks/my-hc",
		"/compute/v1/projects/test-project/global/operations/op-123",
	} {
		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   "test-project",
			RawURL:    srv.URL + nonMatchingPath,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			t.Fatalf("SendRequest to %s failed: %v", nonMatchingPath, err)
		}
	}

	// Matching requests (collection POST and item GET/PUT/DELETE across v1 and beta paths) with expected ApiVersion.
	for _, reqSpec := range []struct {
		method string
		path   string
	}{
		{"POST", "/compute/v1/projects/test-project/global/urlMaps"},
		{"GET", "/compute/v1/projects/test-project/global/urlMaps/my-url-map"},
		{"PUT", "/compute/beta/projects/test-project/global/urlMaps/my-url-map"},
		{"DELETE", "/compute/v1/projects/test-project/global/urlMaps/my-url-map"},
	} {
		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:     config,
			Method:     reqSpec.method,
			Project:    "test-project",
			RawURL:     srv.URL + reqSpec.path,
			UserAgent:  config.UserAgent,
			ApiVersion: "2026-11-01",
		})
		if err != nil {
			t.Fatalf("SendRequest %s %s failed: %v", reqSpec.method, reqSpec.path, err)
		}
	}
}

func TestAssertRequestApiVersion_failureModes(t *testing.T) {
	t.Run("fails when matching request has wrong or missing X-Goog-Api-Version", func(t *testing.T) {
		cmd := exec.Command(os.Args[0], "-test.run=^TestHelperAssertRequestApiVersionSubprocess$")
		cmd.Env = append(os.Environ(), "TEST_ASSERT_API_VERSION_MODE=mismatch")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("expected subprocess test to fail on missing X-Goog-Api-Version, but it succeeded. Output:\n%s", string(out))
		}
		if !strings.Contains(string(out), `had X-Goog-Api-Version="", want "2026-11-01"`) {
			t.Fatalf("expected failure output to mention X-Goog-Api-Version mismatch, got:\n%s", string(out))
		}
	})

	t.Run("fails in cleanup when zero matching requests are observed", func(t *testing.T) {
		cmd := exec.Command(os.Args[0], "-test.run=^TestHelperAssertRequestApiVersionSubprocess$")
		cmd.Env = append(os.Environ(), "TEST_ASSERT_API_VERSION_MODE=zero_matches")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("expected subprocess test to fail when 0 matching requests were observed, but it succeeded. Output:\n%s", string(out))
		}
		if !strings.Contains(string(out), "expected at least 1 matching request, observed 0") {
			t.Fatalf("expected failure output to mention 0 observed requests, got:\n%s", string(out))
		}
	})
}

func TestHelperAssertRequestApiVersionSubprocess(t *testing.T) {
	mode := os.Getenv("TEST_ASSERT_API_VERSION_MODE")
	if mode == "" {
		t.Skip("skipping subprocess helper test")
	}

	setupWatcherTestProviderEnvs(t)
	srv := newFakeApiServer(t)

	acctest.AssertRequestApiVersion(t, "projects/{{project}}/global/urlMaps", "2026-11-01")
	config := acctest.GoogleProviderConfig(t)

	switch mode {
	case "mismatch":
		_, _ = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   "test-project",
			RawURL:    srv.URL + "/compute/v1/projects/test-project/global/urlMaps",
			UserAgent: config.UserAgent,
			// Intentionally omit ApiVersion
		})
	case "zero_matches":
		_, _ = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   "test-project",
			RawURL:    srv.URL + "/compute/v1/projects/test-project/global/backendServices",
			UserAgent: config.UserAgent,
		})
	}
}
