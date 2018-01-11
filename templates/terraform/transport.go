package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
)

func Post(config *Config, url string, body map[string]interface{}) (map[string]interface{}, error) {
	return sendRequest(config, "POST", url, body)
}

func Get(config *Config, url string) (map[string]interface{}, error) {
	return sendRequest(config, "GET", url, nil)
}

func Put(config *Config, url string, body map[string]interface{}) (map[string]interface{}, error) {
	return sendRequest(config, "PUT", url, body)
}

func Delete(config *Config, url string) (map[string]interface{}, error) {
	return sendRequest(config, "DELETE", url, nil)
}

func sendRequest(config *Config, method, url string, body map[string]interface{}) (map[string]interface{}, error) {
	reqHeaders := make(http.Header)
	reqHeaders.Set("User-Agent", config.userAgent)
	reqHeaders.Set("Content-Type", "application/json")

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url+"?alt=json", &buf)
	if err != nil {
		return nil, err
	}
	req.Header = reqHeaders
	res, err := config.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func constructUrl(d *schema.ResourceData, config *Config, linkTmpl string) (string, error) {
	re, err := regexp.Compile("{{([[:word:]]+)}}")
	if err != nil {
		return "", err
	}

	var replaceFunc func(string) string
	if strings.Contains(linkTmpl, "/global/") {
		replaceFunc, err = replaceVarsGlobal(d, config, re)
	} else if strings.Contains(linkTmpl, "/region/") {
		replaceFunc, err = replaceVarsRegional(d, config, re)
	} else if strings.Contains(linkTmpl, "/zone/") {
		replaceFunc, err = replaceVarsZonal(d, config, re)
	} else {
		return "", fmt.Errorf("Could not expand variables for URL template %s", linkTmpl)
	}
	if err != nil {
		return "", err
	}
	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}

func replaceVarsGlobal(d *schema.ResourceData, config *Config, re *regexp.Regexp) (func(string) string, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	return func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		return d.Get(m).(string)
	}, nil
}

func replaceVarsRegional(d *schema.ResourceData, config *Config, re *regexp.Regexp) (func(string) string, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}

	return func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "region" {
			return region
		}
		return d.Get(m).(string)
	}, nil
}

func replaceVarsZonal(d *schema.ResourceData, config *Config, re *regexp.Regexp) (func(string) string, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return nil, err
	}

	return func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "zone" {
			return zone
		}
		return d.Get(m).(string)
	}, nil
}
