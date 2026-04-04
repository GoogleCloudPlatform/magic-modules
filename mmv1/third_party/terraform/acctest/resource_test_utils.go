package acctest

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// General test utils

var _ plancheck.PlanCheck = expectNoDelete{}

type expectNoDelete struct{}

func (e expectNoDelete) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result error
	for _, rc := range req.Plan.ResourceChanges {
		if slices.Contains(rc.Change.Actions, tfjson.ActionDelete) {
			result = errors.Join(result, fmt.Errorf("expected no deletion of resources, but %s has planned deletion", rc.Address))
		}
	}
	resp.Error = result
}

func ExpectNoDelete() plancheck.PlanCheck {
	return expectNoDelete{}
}

// TestExtractResourceAttr navigates a test's state to find the specified resource (or data source) attribute and makes the value
// accessible via the attributeValue string pointer.
func TestExtractResourceAttr(resourceName string, attributeName string, attributeValue *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName] // To find a datasource, include `data.` at the start of the resourceName value

		if !ok {
			return fmt.Errorf("resource name %s not found in state", resourceName)
		}

		attrValue, ok := rs.Primary.Attributes[attributeName]

		if !ok {
			return fmt.Errorf("attribute %s not found in resource %s state", attributeName, resourceName)
		}

		*attributeValue = attrValue

		return nil
	}
}

// GetTestRegion has the same logic as the provider's GetRegion, to be used in tests.
func GetTestRegion(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// GetTestProject has the same logic as the provider's GetProject, to be used in tests.
func GetTestProject(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["project"]; ok {
		return res, nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "project")
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func SkipIfVcr(t *testing.T) {
	if IsVcrEnabled() {
		t.Skipf("VCR enabled, skipping test: %s", t.Name())
	}
}

func SleepInSecondsForTest(t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Assume we never want to sleep when we're in replaying mode.
		if IsVcrEnabled() && os.Getenv("VCR_MODE") == "REPLAYING" {
			return nil
		}
		time.Sleep(time.Duration(t) * time.Second)
		return nil
	}
}

// TestCheckAttributeValuesEqual compares two string pointers, which have been used to retrieve attribute values from the test's state.
func TestCheckAttributeValuesEqual(i *string, j *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if testStringValue(i) != testStringValue(j) {
			return fmt.Errorf("attribute values are different, got %s and %s", testStringValue(i), testStringValue(j))
		}

		return nil
	}
}

// ConditionTitleIfPresent returns empty string if condition is not preset and " {condition.0.title}" if it is.
func BuildIAMImportId(name, role, member, condition string) string {
	ret := name
	if role != "" {
		ret += " " + role
	}
	if member != "" {
		ret += " " + member
	}
	if condition != "" {
		ret += " " + condition
	}
	return ret
}

// TagBindingCheckConfig configures a CheckTagBindings assertion.
// BuildParent must return the full resource name used as the tagBindings parent.
// If GetLocation is nil, the global tagBindings endpoint is used.
// If GetLocation is set, the location-scoped endpoint is used.
type TagBindingCheckConfig struct {
	ResourceName                string
	ExpectedTagValueResources   []string
	UnexpectedTagValueResources []string
	BuildParent                 func(rs *terraform.ResourceState) (string, error)
	GetLocation                 func(rs *terraform.ResourceState) (string, error)
}

// CheckTagBindings verifies that the target resource has the expected tag value
// bindings and does not have the unexpected ones.
func CheckTagBindings(t *testing.T, cfg TagBindingCheckConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cfg.ResourceName == "" {
			return fmt.Errorf("resource name must be set for CheckTagBindings")
		}
		if cfg.BuildParent == nil {
			return fmt.Errorf("BuildParent must be set for CheckTagBindings on resource %s", cfg.ResourceName)
		}

		rs, err := getResourceState(s, cfg.ResourceName)
		if err != nil {
			return err
		}

		expectedTagValues := make([]string, 0, len(cfg.ExpectedTagValueResources))
		for _, resourceName := range cfg.ExpectedTagValueResources {
			tagValueID, err := getResourceID(s, resourceName)
			if err != nil {
				return err
			}
			expectedTagValues = append(expectedTagValues, tagValueID)
		}

		unexpectedTagValues := make([]string, 0, len(cfg.UnexpectedTagValueResources))
		for _, resourceName := range cfg.UnexpectedTagValueResources {
			tagValueID, err := getResourceID(s, resourceName)
			if err != nil {
				return err
			}
			unexpectedTagValues = append(unexpectedTagValues, tagValueID)
		}

		parent, err := cfg.BuildParent(rs)
		if err != nil {
			return err
		}

		config := GoogleProviderConfig(t)
		basePath := config.TagsBasePath
		if cfg.GetLocation != nil {
			location, err := cfg.GetLocation(rs)
			if err != nil {
				return err
			}
			basePath = strings.Replace(config.TagsLocationBasePath, "{{location}}", location, 1)
		}

		listBindingsURL := fmt.Sprintf("%stagBindings/?parent=%s&pageSize=300", basePath, url.QueryEscape(parent))
		resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    listBindingsURL,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return fmt.Errorf("error calling tagBindings API for resource %s: %v", rs.Primary.ID, err)
		}

		tagBindingsVal, exists := resp["tagBindings"]
		if !exists {
			tagBindingsVal = []interface{}{}
		}

		tagBindings, ok := tagBindingsVal.([]interface{})
		if !ok {
			return fmt.Errorf("'tagBindings' is not a slice in response for resource %s. response: %v", rs.Primary.ID, resp)
		}

		foundExpected := make(map[string]bool, len(expectedTagValues))
		foundUnexpected := make(map[string]bool, len(unexpectedTagValues))

		for _, binding := range tagBindings {
			bindingMap, ok := binding.(map[string]interface{})
			if !ok {
				continue
			}

			tagValue, _ := bindingMap["tagValue"].(string)
			for _, expectedTagValue := range expectedTagValues {
				if tagValue == expectedTagValue {
					foundExpected[expectedTagValue] = true
				}
			}
			for _, unexpectedTagValue := range unexpectedTagValues {
				if tagValue == unexpectedTagValue {
					foundUnexpected[unexpectedTagValue] = true
				}
			}
		}

		for _, expectedTagValue := range expectedTagValues {
			if !foundExpected[expectedTagValue] {
				return fmt.Errorf("expected tag value %s not found in tag bindings for resource %s. bindings: %v", expectedTagValue, rs.Primary.ID, tagBindings)
			}
		}

		for _, unexpectedTagValue := range unexpectedTagValues {
			if foundUnexpected[unexpectedTagValue] {
				return fmt.Errorf("unexpected tag value %s found in tag bindings for resource %s. bindings: %v", unexpectedTagValue, rs.Primary.ID, tagBindings)
			}
		}

		return nil
	}
}

func getResourceState(s *terraform.State, resourceName string) (*terraform.ResourceState, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("terraform resource not found: %s", resourceName)
	}
	return rs, nil
}

func getResourceID(s *terraform.State, resourceName string) (string, error) {
	rs, err := getResourceState(s, resourceName)
	if err != nil {
		return "", err
	}
	if rs.Primary.ID == "" {
		return "", fmt.Errorf("terraform resource %s has no id", resourceName)
	}
	return rs.Primary.ID, nil
}

// testStringValue returns string values from string pointers, handling nil pointers.
func testStringValue(sPtr *string) string {
	if sPtr == nil {
		return ""
	}

	return *sPtr
}
