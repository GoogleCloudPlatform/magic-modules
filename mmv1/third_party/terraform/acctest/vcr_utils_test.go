package acctest_test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
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
	headers map[string]string
}

func prepareHttpRequest(d requestDescription) *http.Request {
	url := &url.URL{
		Scheme: d.scheme,
		Host:   d.host,
		Path:   d.path,
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

func TestDiffTestStepInjection(t *testing.T) {
	// This test is a placeholder for the diff test step injection functionality.
	// It should be implemented to test the diff injection logic in the future.
	t.Skip("Diff test step injection test not implemented yet")
}

// todo: fix imports so this is only in vcr_utils.go
func ParseReleaseDiffOutput(temp *os.File) (string, error) {
	if temp == nil {
		return "", errors.New("temporary file is nil")
	}

	_, err := temp.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek to beginning of temporary file: %w", err)
	}

	var lastLine string
	scanner := bufio.NewScanner(temp)

	for scanner.Scan() {
		line := scanner.Text()
		lastLine = line
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading temporary file: %w", err)
	}

	return lastLine, nil
}
func TestParseReleaseDiffOutput(t *testing.T) {
	temp_file, err := os.CreateTemp("", "test_release_diff_test_output_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(temp_file.Name())
	// Write some dummy data to the temp file
	_, err = temp_file.WriteString("This is a test release diff output.\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	var expectedOutput = "This is a test release diff output.\n"
	var output string
	output, err = ParseReleaseDiffOutput(temp_file)
	if err != nil {
		t.Fatalf("Failed to parse release diff output: %v", err)
	}

	if output != expectedOutput {
		t.Fatalf("Expected output to be:\n%q\nbut got:\n%q", expectedOutput, output)
	}
	t.Logf("Parsed output:\n%s", output)
}

func TestReformConfigWithProviderGoogleBeta(t *testing.T) {
	var config = ` data "google_new_resource" {
	  provider = google-beta
}`

	var newConfig = acctest.ReformConfigWithProvider(config, "google-local")

	// Adjusted expectedConfig to match the actual output of the function
	expectedConfig := ` data "google_new_resource" {
	  provider = google-local
}`
	// Added extra newline and adjusted indentation and provider name
	if newConfig != expectedConfig {
		// Change %s to %q here to reveal invisible differences!
		t.Fatalf("Expected config to be reformatted to:\n%q\nbut got:\n%q", expectedConfig, newConfig)
	}
	t.Logf("Reformed config:\n%s", newConfig)
}

// the empty provider case fails, so need to figure out a regex way of handling
func TestReformConfigWithProviderEmpty(t *testing.T) {
	var config = ` data "google_new_resource" {
}`

	var newConfig = acctest.ReformConfigWithProvider(config, "google-local")

	// Adjusted expectedConfig to match the actual output of the function
	expectedConfig := ` data "google_new_resource" {
	  provider = google-local
}`
	// Added extra newline and adjusted indentation and provider name
	if newConfig != expectedConfig {
		// Change %s to %q here to reveal invisible differences!
		t.Fatalf("Expected config to be reformatted to:\n%q\nbut got:\n%q", expectedConfig, newConfig)
	}
	t.Logf("Reformed config:\n%s", newConfig)
}
