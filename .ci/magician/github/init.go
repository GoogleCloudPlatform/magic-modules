/*
* Copyright 2023 Google LLC. All Rights Reserved.
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
	"bytes"
	"context"
	"io"
	"net/http"

	utils "magician/utility"

	gh "github.com/google/go-github/v68/github"
)

// Client for GitHub interactions.
type Client struct {
	token string
	gh    *gh.Client
	ctx   context.Context
}

// retryTransport is a custom RoundTripper that adds retry and logging
type retryTransport struct {
	underlyingTransport http.RoundTripper
	token               string
}

func NewClient(token string) *Client {
	ctx := context.Background()

	// Create a custom transport with retry logic
	rt := &retryTransport{
		underlyingTransport: http.DefaultTransport,
		token:               token,
	}

	// Use this custom transport with OAuth2
	tc := &http.Client{Transport: rt}

	// Create the GitHub client with our custom transport
	ghClient := gh.NewClient(tc)

	return &Client{
		gh:    ghClient,
		token: token,
		ctx:   ctx,
	}
}

// RoundTrip implements the http.RoundTripper interface
func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Extract information from the request
	method := req.Method
	urlStr := req.URL.String()

	// Read and log the request body if present
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	resp, respBody, err := utils.RequestCallWithRetryRaw(urlStr, method, rt.token, bodyBytes)
	if err != nil {
		return nil, err
	}

	// Replace the response body with our captured body
	resp.Body.Close() // Close the original body
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	resp.ContentLength = int64(len(respBody))

	return resp, nil
}
