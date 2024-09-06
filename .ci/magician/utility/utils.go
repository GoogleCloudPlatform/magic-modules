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
	"net/http"

	"golang.org/x/exp/slices"
)

func RequestCall(url, method, credentials string, result any, body any) error {
	client := &http.Client{}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %s", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials))
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("")
	fmt.Println("request url: ", url)
	fmt.Println("request body: ", string(jsonBody)) // Convert to string
	fmt.Println("")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("response status-code: ", resp.StatusCode)
	fmt.Println("response body: ", string(respBodyBytes)) // Convert to string
	fmt.Println("")

	// Decode the response, if needed
	if result != nil {
		if err = json.Unmarshal(respBodyBytes, &result); err != nil {
			return err
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("got code %d from server", resp.StatusCode)
	}

	return nil
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
