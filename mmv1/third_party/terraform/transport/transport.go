package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"
)

var DefaultRequestTimeout = 5 * time.Minute

type SendRequestOptions struct {
	Config               *Config
	Method               string
	Project              string
	RawURL               string
	UserAgent            string
	Body                 map[string]any
	Timeout              time.Duration
	Headers              http.Header
	ErrorRetryPredicates []RetryErrorPredicateFunc
	ErrorAbortPredicates []RetryErrorPredicateFunc
}

func SendRequest(opt SendRequestOptions) (map[string]interface{}, error) {
	if opt.Config == nil || opt.Config.Client == nil {
		return nil, fmt.Errorf("client is nil for request to %s", opt.RawURL)
	}

	reqHeaders := opt.Headers
	if reqHeaders == nil {
		reqHeaders = make(http.Header)
	}
	reqHeaders.Set("User-Agent", opt.UserAgent)
	reqHeaders.Set("Content-Type", "application/json")

	if opt.Config.UserProjectOverride && opt.Project != "" {
		// When opt.Project is "NO_BILLING_PROJECT_OVERRIDE" in the function GetCurrentUserEmail,
		// set the header X-Goog-User-Project to be empty string.
		if opt.Project == "NO_BILLING_PROJECT_OVERRIDE" {
			reqHeaders.Set("X-Goog-User-Project", "")
		} else {
			// Pass the project into this fn instead of parsing it from the URL because
			// both project names and URLs can have colons in them.
			reqHeaders.Set("X-Goog-User-Project", opt.Project)
		}
	}

	if opt.Timeout == 0 {
		opt.Timeout = DefaultRequestTimeout
	}

	var res *http.Response
	err := Retry(RetryOptions{
		RetryFunc: func() error {
			var buf bytes.Buffer
			if opt.Body != nil {
				err := json.NewEncoder(&buf).Encode(opt.Body)
				if err != nil {
					return err
				}
			}

			u, err := AddQueryParams(opt.RawURL, map[string]string{"alt": "json"})
			if err != nil {
				return err
			}
			req, err := http.NewRequest(opt.Method, u, &buf)
			if err != nil {
				return err
			}

			req.Header = reqHeaders
			res, err = opt.Config.Client.Do(req)
			if err != nil {
				return err
			}

			if err := googleapi.CheckResponse(res); err != nil {
				googleapi.CloseBody(res)
				return err
			}

			return nil
		},
		Timeout:              opt.Timeout,
		ErrorRetryPredicates: opt.ErrorRetryPredicates,
		ErrorAbortPredicates: opt.ErrorAbortPredicates,
	})
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, fmt.Errorf("Unable to parse server response. This is most likely a terraform problem, please file a bug at https://github.com/hashicorp/terraform-provider-google/issues.")
	}

	// The defer call must be made outside of the retryFunc otherwise it's closed too soon.
	defer googleapi.CloseBody(res)

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

func AddQueryParams(rawurl string, params map[string]string) (string, error) {
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

func AddArrayQueryParams(rawurl string, param string, values []interface{}) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for _, v := range values {
		q.Add(param, v.(string))
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func HandleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore
		d.SetId("")

		return nil
	}

	return errwrap.Wrapf(
		fmt.Sprintf("Error when reading or editing %s: {{err}}", resource), err)
}

func HandleDataSourceNotFoundError(err error, d *schema.ResourceData, resource, url string) error {
	if IsGoogleApiErrorWithCode(err, 404) {
		return fmt.Errorf("%s not found", url)
	}

	return errwrap.Wrapf(
		fmt.Sprintf("Error when reading or editing %s: {{err}}", resource), err)
}

func IsGoogleApiErrorWithCode(err error, errCode int) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	return ok && gerr != nil && gerr.Code == errCode
}

func IsApiNotEnabledError(err error) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	if !ok {
		return false
	}
	if gerr == nil {
		return false
	}
	if gerr.Code != 403 {
		return false
	}
	for _, e := range gerr.Errors {
		if e.Reason == "accessNotConfigured" {
			return true
		}
	}
	return false
}

type GetPaginatedItemsOptions struct {
	ResourceData         *schema.ResourceData
	Config               *Config
	BillingProject       *string
	UserAgent            string
	URL                  string
	ListFlattener        func(config *Config, res []map[string]interface{}) ([]map[string]interface{}, error)
	Params               map[string]string
	ResourceToList       string
	ErrorRetryPredicates []RetryErrorPredicateFunc
}

func GetPaginatedItems(paginationOptions GetPaginatedItemsOptions) ([]map[string]interface{}, error) {
	if paginationOptions.Params == nil {
		paginationOptions.Params = make(map[string]string)
	}

	items := make([]map[string]interface{}, 0)
	for {
		url, err := AddQueryParams(paginationOptions.URL, paginationOptions.Params)
		if err != nil {
			return nil, err
		}

		headers := make(http.Header)
		opts := SendRequestOptions{
			Config:               paginationOptions.Config,
			Method:               "GET",
			RawURL:               url,
			UserAgent:            paginationOptions.UserAgent,
			Headers:              headers,
			ErrorRetryPredicates: paginationOptions.ErrorRetryPredicates,
		}
		if paginationOptions.BillingProject != nil {
			opts.Project = *paginationOptions.BillingProject
		}
		res, err := SendRequest(opts)
		if err != nil {
			return nil, HandleNotFoundError(err, paginationOptions.ResourceData, fmt.Sprintf("%s %q", paginationOptions.ResourceToList, paginationOptions.ResourceData.Id()))
		}

		var newItems []map[string]interface{}
		if res[paginationOptions.ResourceToList] != nil {
			itemsAsMap, err := InterfaceSliceToMapSlice(res[paginationOptions.ResourceToList])
			if err != nil {
				return nil, err
			}

			if paginationOptions.ListFlattener != nil {
				log.Printf("[DEBUG] res[paginationOptions.ResourceToList]: %#v", res[paginationOptions.ResourceToList])
				flattened, err := paginationOptions.ListFlattener(paginationOptions.Config, itemsAsMap)
				if err != nil {
					return nil, err
				}
				newItems = flattened
			} else {
				newItems = itemsAsMap
			}
		}
		items = append(items, newItems...)

		if v, ok := res["nextPageToken"]; ok && v != nil && v.(string) != "" {
			paginationOptions.Params["pageToken"] = v.(string)
		} else {
			break
		}
	}
	return items, nil
}

func InterfaceSliceToMapSlice(i interface{}) ([]map[string]interface{}, error) {
	itemsAsInterface, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot convert to slice of interfaces: got %T", i)
	}

	itemsAsMap := make([]map[string]interface{}, len(itemsAsInterface))
	for idx, item := range itemsAsInterface {
		var ok bool
		itemsAsMap[idx], ok = item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot convert item to map[string]interface{}: got %T", item)
		}
	}
	return itemsAsMap, nil
}
