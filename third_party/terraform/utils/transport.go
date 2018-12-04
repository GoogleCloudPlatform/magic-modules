package google

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/googleapi"
)

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

var retriableStatuses = [...]int{
	http.StatusTooManyRequests,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
}

// The google api client doesn't throw an error for non 200 responses, instead
// the response is returned with .StatusCode set and the error in the body. This
// inspects the response for common retriable errors such as 503 (stack driver throttling)
func isRetriableResponse(res *http.Response) bool {
	// Return early for any 2**/3** responses
	if res.StatusCode < 400 {
		return false
	}

	for _, code := range retriableStatuses {
		if res.StatusCode == code {
			return true
		}
	}
	return false
}

func sendRequest(config *Config, method, rawurl string, body map[string]interface{}) (map[string]interface{}, error) {
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

	u, err := addQueryParams(rawurl, map[string]string{"alt": "json"})
	if err != nil {
		return nil, err
	}

	var res *http.Response
	// Retry each request up to 3 times for retriable failures
	for i := 0; i < 3; i++ {
		if i > 0 {
			// Sleep on all requests other then the first
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		}
		req, err := http.NewRequest(method, u, &buf)
		if err != nil {
			return nil, err
		}
		req.Header = reqHeaders
		res, err = config.client.Do(req)
		if err != nil {
			return nil, err
		}

		if !isRetriableResponse(res) {
			break
		}
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}

	// 204 responses will have no body, so we're going to error with "EOF" if we
	// try to parse it. Instead, we can just return nil.
	if res.StatusCode == 204 {
		return nil, nil
	}

	result := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func addQueryParams(rawurl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func replaceVars(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string
	var err error

	if strings.Contains(linkTmpl, "{{project}}") {
		project, err = getProject(d, config)
		if err != nil {
			return "", err
		}
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region, err = getRegion(d, config)
		if err != nil {
			return "", err
		}
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone, err = getZone(d, config)
		if err != nil {
			return "", err
		}
	}

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}
		v, ok := d.GetOk(m)
		if ok {
			return v.(string)
		}
		return ""
	}

	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}
