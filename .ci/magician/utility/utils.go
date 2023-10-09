package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/exp/slices"
)

func RequestCall(url, method, credentials string, result any, body any) (int, error) {
	client := &http.Client{}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return 1, fmt.Errorf("error marshaling JSON: %s", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 2, fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials))
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("")
	fmt.Println("request url: ", url)
	fmt.Println("request body: ", string(jsonBody)) // Convert to string
	fmt.Println("")

	resp, err := client.Do(req)
	if err != nil {
		return 3, err
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 5, err
	}

	fmt.Println("response status-code: ", resp.StatusCode)
	fmt.Println("response body: ", string(respBodyBytes)) // Convert to string
	fmt.Println("")

	// Decode the response, if needed
	if result != nil {
		if err = json.Unmarshal(respBodyBytes, &result); err != nil {
			return 4, err
		}
	}

	return resp.StatusCode, nil
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
