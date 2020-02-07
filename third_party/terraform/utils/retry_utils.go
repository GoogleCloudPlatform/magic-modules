package google

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/googleapi"
)

func retry(retryFunc func() error) error {
	return retryTime(retryFunc, 1)
}

func retryTime(retryFunc func() error, minutes int) error {
	return retryTimeDuration(retryFunc, time.Duration(minutes)*time.Minute)
}

func retryTimeDuration(retryFunc func() error, duration time.Duration, errorRetryPredicates ...func(e error) (bool, string)) error {
	return resource.Retry(duration, func() *resource.RetryError {
		err := retryFunc()
		if err == nil {
			return nil
		}
		for _, e := range getAllTypes(err, &googleapi.Error{}, &url.Error{}) {
			if isRetryableError(e, errorRetryPredicates) {
				return resource.RetryableError(e)
			}
		}
		return resource.NonRetryableError(err)
	})
}

func getAllTypes(err error, args ...interface{}) []error {
	var result []error
	for _, v := range args {
		subResult := errwrap.GetAllType(err, v)
		if subResult != nil {
			result = append(result, subResult...)
		}
	}
	return result
}

func isRetryableError(err error, retryPredicates []func(e error) (bool, string)) bool {
	// These operations are always hitting googleapis.com - they should rarely
	// time out, and if they do, that timeout is retryable.
	if urlerr, ok := err.(*url.Error); ok && urlerr.Timeout() {
		log.Printf("[DEBUG] Dismissed an error as retryable based on googleapis.com target: %s", err)
		return true
	}

	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 429 || gerr.Code == 500 || gerr.Code == 502 || gerr.Code == 503 {
			log.Printf("[DEBUG] Dismissed an error as retryable based on error code: %s", err)
			return true
		}

		if gerr.Code == 409 && strings.Contains(gerr.Body, "operationInProgress") {
			// 409's are retried because cloud sql throws a 409 when concurrent calls are made.
			// The only way right now to determine it is a SQL 409 due to concurrent calls is to
			// look at the contents of the error message.
			// See https://github.com/terraform-providers/terraform-provider-google/issues/3279
			log.Printf("[DEBUG] Dismissed an error as retryable based on error code 409 and error reason 'operationInProgress': %s", err)
			return true
		}

		if gerr.Code == 412 && isFingerprintError(err) {
			log.Printf("[DEBUG] Dismissed an error as retryable as a fingerprint mismatch: %s", err)
			return true
		}

	}
	for _, pred := range retryPredicates {
		if retry, reason := (pred(err)); retry {
			log.Printf("[DEBUG] Dismissed an error as retryable. %s - %s", reason, err)
			return true
		}
	}

	return false
}

type ExistenceCheckFunc func(*schema.ResourceData, interface{}) (bool, error)

// tryPollForResourceExistence attempts to make sure a resource exists right after creation, before
// the final Read() after resource has been created.
// It is meant for use with resources/APIs that are eventually consistent
// or implement caching of read results.
//
// existenceCheckF should handle retries as needed for base resource/API (i.e. retry predicates, common retry)
func tryPollForResourceExistence(d *schema.ResourceData, meta interface{}, existenceCheckF ExistenceCheckFunc) {
	log.Printf("[DEBUG] Confirming successful read of resource post-create")
	err := resource.Retry(time.Second*10, func() *resource.RetryError {
		found, err := existenceCheckF(d, meta)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !found {
			return resource.RetryableError(fmt.Errorf("resource %q not found", d.Id()))
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] Unable to confirm existence of resource during create, deferring to read: %v", err)
	}
}
