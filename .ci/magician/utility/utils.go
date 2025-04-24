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
package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/slices"
)

// retryConfig holds configuration for request retries
type retryConfig struct {
	MaxRetries       int           // Maximum number of retry attempts
	InitialBackoff   time.Duration // Initial backoff duration
	MaxBackoff       time.Duration // Maximum backoff duration
	BackoffFactor    float64       // Factor by which to multiply backoff on each retry
	RetryStatusCodes []int         // HTTP status codes that should trigger a retry
}

// defaultRetryConfig provides default retry configuration
func defaultRetryConfig() retryConfig {
	return retryConfig{
		MaxRetries:       3,
		InitialBackoff:   5000 * time.Millisecond,
		MaxBackoff:       60 * time.Second,
		BackoffFactor:    2.0,
		RetryStatusCodes: []int{408, 429, 500, 502, 503, 504}, // Common retry status codes
	}
}

// makeHTTPRequest performs the actual HTTP request and returns the response
func makeHTTPRequest(url, method, credentials string, body any) (*http.Response, []byte, error) {
	client := &http.Client{}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling JSON: %s", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	fmt.Println("")
	fmt.Println("request url: ", url)
	fmt.Println("request body: ", string(jsonBody))
	fmt.Println("")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("response status-code: ", resp.StatusCode)
	fmt.Println("response body: ", string(respBodyBytes))
	fmt.Println("")

	return resp, respBodyBytes, nil
}

// processResponse handles the response and unmarshals it to the result if provided
func processResponse(resp *http.Response, respBodyBytes []byte, result any) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResponse struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		}

		if err := json.Unmarshal(respBodyBytes, &errorResponse); err == nil {
			errorMsg := errorResponse.Message
			if errorMsg == "" {
				errorMsg = errorResponse.Error
			}

			if errorMsg != "" {
				return fmt.Errorf("got code %d from server: %s", resp.StatusCode, errorMsg)
			}
		}

		// Fall back to generic error if we couldn't parse the error message
		return fmt.Errorf("got code %d from server", resp.StatusCode)
	}

	// If no error status code, decode the response if needed
	if result != nil {
		if err := json.Unmarshal(respBodyBytes, &result); err != nil {
			return err
		}
	}

	return nil
}

// RequestCall makes a single HTTP request without retries
func RequestCall(url, method, credentials string, result any, body any) error {
	resp, respBodyBytes, err := makeHTTPRequest(url, method, credentials, body)
	if err != nil {
		return err
	}

	return processResponse(resp, respBodyBytes, result)
}

// shouldRetry determines if a retry should be attempted based on the status code
func shouldRetry(statusCode int, retryConfig retryConfig) bool {
	return slices.Contains(retryConfig.RetryStatusCodes, statusCode)
}

// calculateBackoff calculates the backoff duration for the current retry attempt
func calculateBackoff(attempt int, config retryConfig) time.Duration {
	backoff := config.InitialBackoff * time.Duration(math.Pow(config.BackoffFactor, float64(attempt)))
	if backoff > config.MaxBackoff {
		backoff = config.MaxBackoff
	}
	return backoff
}

// RequestCallWithRetry makes an HTTP request with retry capability
func requestCallWithRetry(url, method, credentials string, result any, body any, config retryConfig) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// If this is a retry attempt, wait before trying again
		if attempt > 0 {
			backoff := calculateBackoff(attempt-1, config)
			fmt.Printf("Retry attempt %d after %v\n", attempt, backoff)
			time.Sleep(backoff)
		}

		resp, respBodyBytes, err := makeHTTPRequest(url, method, credentials, body)
		if err != nil {
			lastErr = err
			continue // Network error, retry
		}

		// Process the response
		err = processResponse(resp, respBodyBytes, result)
		if err != nil {
			lastErr = err

			// Check if we should retry based on status code
			if shouldRetry(resp.StatusCode, config) {
				continue
			}
		}

		// If we got here with no error, return success
		return err
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// RequestCallWithRetry is a convenience function that uses default retry settings
func RequestCallWithRetry(url, method, credentials string, result any, body any) error {
	return requestCallWithRetry(url, method, credentials, result, body, defaultRetryConfig())
}

func Removes(s1 []string, s2 []string) []string {
	result := make([]string, 0, len(s1))

	for _, v := range s1 {
		if !slices.Contains(s2, v) {
			result = append(result, v)
		}
	}
	return result
}

func WriteToJson(data interface{}, path string) error {
	rsBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(path, rsBytes, 0644); err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

func ReadFromJson(data interface{}, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read data from file: %w", err)
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return nil
}
