/*
* Copyright 2026 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type rewriteTransport struct {
	base      http.RoundTripper
	serverURL *url.URL
}

func (rt *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqCopy := req.Clone(req.Context())
	reqCopy.URL.Scheme = rt.serverURL.Scheme
	reqCopy.URL.Host = rt.serverURL.Host
	return rt.base.RoundTrip(reqCopy)
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL: %v", err)
	}

	origTransport := http.DefaultTransport
	http.DefaultTransport = &rewriteTransport{
		base:      server.Client().Transport,
		serverURL: serverURL,
	}
	t.Cleanup(func() {
		http.DefaultTransport = origTransport
	})

	return NewClient("test-token")
}

func TestMergePullRequest(t *testing.T) {
	origDelay := mergeRetryDelay
	mergeRetryDelay = time.Millisecond
	t.Cleanup(func() {
		mergeRetryDelay = origDelay
	})

	t.Run("merges successfully on first attempt", func(t *testing.T) {
		mergeCalls := 0
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut && r.URL.Path == "/repos/hashicorp/terraform-provider-google-beta/pulls/12548/merge" {
				mergeCalls++
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}
				if body["merge_method"] != "squash" || body["sha"] != "abc1234" {
					t.Errorf("unexpected merge body: %v", body)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"sha":"mergedsha","merged":true,"message":"Pull Request successfully merged"}`))
				return
			}
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		})

		err := client.MergePullRequest("hashicorp", "terraform-provider-google-beta", "12548", "abc1234")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if mergeCalls != 1 {
			t.Errorf("expected 1 merge call, got %d", mergeCalls)
		}
	})

	t.Run("checks target repo and succeeds when PR is already merged", func(t *testing.T) {
		mergeCalls := 0
		getCalls := 0
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodPut && r.URL.Path == "/repos/GoogleCloudPlatform/terraform-google-conversion/pulls/4867/merge":
				mergeCalls++
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(`{"message":"Pull Request is not mergeable"}`))
			case r.Method == http.MethodGet && r.URL.Path == "/repos/GoogleCloudPlatform/terraform-google-conversion/pulls/4867":
				getCalls++
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"number":4867,"merged":true}`))
			default:
				t.Errorf("unexpected request (should not query magic-modules): %s %s", r.Method, r.URL.Path)
				w.WriteHeader(http.StatusNotFound)
			}
		})

		err := client.MergePullRequest("GoogleCloudPlatform", "terraform-google-conversion", "4867", "abc1234")
		if err != nil {
			t.Fatalf("expected no error when PR was already merged in target repo, got: %v", err)
		}
		if mergeCalls != 1 {
			t.Errorf("expected 1 merge call, got %d", mergeCalls)
		}
		if getCalls != 1 {
			t.Errorf("expected 1 get PR call on target repo, got %d", getCalls)
		}
	})

	t.Run("retries when PR is transiently not mergeable and then succeeds", func(t *testing.T) {
		mergeCalls := 0
		getCalls := 0
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodPut && r.URL.Path == "/repos/hashicorp/terraform-provider-google/pulls/999/merge":
				mergeCalls++
				if mergeCalls == 1 {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"message":"Pull Request is not mergeable"}`))
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"merged":true}`))
			case r.Method == http.MethodGet && r.URL.Path == "/repos/hashicorp/terraform-provider-google/pulls/999":
				getCalls++
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"number":999,"merged":false}`))
			default:
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				w.WriteHeader(http.StatusNotFound)
			}
		})

		err := client.MergePullRequest("hashicorp", "terraform-provider-google", "999", "abc1234")
		if err != nil {
			t.Fatalf("expected retry to succeed, got: %v", err)
		}
		if mergeCalls != 2 {
			t.Errorf("expected 2 merge calls, got %d", mergeCalls)
		}
		if getCalls != 1 {
			t.Errorf("expected 1 get PR call, got %d", getCalls)
		}
	})

	t.Run("fails after max retries when PR remains not mergeable and unmerged", func(t *testing.T) {
		mergeCalls := 0
		getCalls := 0
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodPut && r.URL.Path == "/repos/hashicorp/terraform-provider-google/pulls/999/merge":
				mergeCalls++
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(`{"message":"Pull Request is not mergeable"}`))
			case r.Method == http.MethodGet && r.URL.Path == "/repos/hashicorp/terraform-provider-google/pulls/999":
				getCalls++
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"number":999,"merged":false}`))
			default:
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				w.WriteHeader(http.StatusNotFound)
			}
		})

		err := client.MergePullRequest("hashicorp", "terraform-provider-google", "999", "abc1234")
		if err == nil || !strings.Contains(err.Error(), "Pull Request is not mergeable") {
			t.Fatalf("expected Pull Request is not mergeable error, got: %v", err)
		}
		if mergeCalls != maxMergeRetries+1 {
			t.Errorf("expected %d merge calls, got %d", maxMergeRetries+1, mergeCalls)
		}
		if getCalls != maxMergeRetries+1 {
			t.Errorf("expected %d get PR calls, got %d", maxMergeRetries+1, getCalls)
		}
	})

	t.Run("retries on Base branch was modified and succeeds", func(t *testing.T) {
		mergeCalls := 0
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut && r.URL.Path == "/repos/hashicorp/terraform-provider-google-beta/pulls/12548/merge" {
				mergeCalls++
				if mergeCalls < 3 {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"message":"Base branch was modified. Review and try the merge again."}`))
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"merged":true}`))
				return
			}
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		})

		err := client.MergePullRequest("hashicorp", "terraform-provider-google-beta", "12548", "abc1234")
		if err != nil {
			t.Fatalf("expected retry on base branch modified to succeed, got: %v", err)
		}
		if mergeCalls != 3 {
			t.Errorf("expected 3 merge calls, got %d", mergeCalls)
		}
	})

	t.Run("retries on Merge already in progress and succeeds", func(t *testing.T) {
		mergeCalls := 0
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut && r.URL.Path == "/repos/hashicorp/terraform-provider-google/pulls/100/merge" {
				mergeCalls++
				if mergeCalls == 1 {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"message":"Merge already in progress"}`))
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"merged":true}`))
				return
			}
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		})

		err := client.MergePullRequest("hashicorp", "terraform-provider-google", "100", "abc1234")
		if err != nil {
			t.Fatalf("expected retry on merge in progress to succeed, got: %v", err)
		}
		if mergeCalls != 2 {
			t.Errorf("expected 2 merge calls, got %d", mergeCalls)
		}
	})

	t.Run("does not retry on non-retryable error", func(t *testing.T) {
		mergeCalls := 0
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut && r.URL.Path == "/repos/hashicorp/terraform-provider-google/pulls/100/merge" {
				mergeCalls++
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(`{"message":"Head branch was modified. Review and try the merge again."}`))
				return
			}
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		})

		err := client.MergePullRequest("hashicorp", "terraform-provider-google", "100", "abc1234")
		if err == nil || !strings.Contains(err.Error(), "Head branch was modified") {
			t.Fatalf("expected Head branch was modified error, got: %v", err)
		}
		if mergeCalls != 1 {
			t.Errorf("expected 1 merge call without retry, got %d", mergeCalls)
		}
	})
}
