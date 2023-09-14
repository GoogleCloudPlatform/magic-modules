package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/exp/slices"
)

func RequestCall(url, method, credentials string, result interface{}, body interface{}) (int, error) {
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
	resp, err := client.Do(req)
	if err != nil {
		return 3, err
	}
	defer resp.Body.Close()

	if result != nil {
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
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
