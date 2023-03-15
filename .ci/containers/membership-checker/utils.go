package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func requestCall(url, method, GITHUB_TOKEN string, result interface{}, body interface{}) (int, error) {
	client := &http.Client{}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return 1, fmt.Errorf("Error marshaling JSON: %s", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 1, fmt.Errorf("Error creating request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", GITHUB_TOKEN))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 1, err
	}
	defer resp.Body.Close()

	if result != nil {
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return 1, err
		}
	}

	return resp.StatusCode, nil
}

func readFile(filename string) (string, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
