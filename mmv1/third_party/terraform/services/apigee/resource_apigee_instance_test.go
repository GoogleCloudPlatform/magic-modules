package apigee_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/services/apigee"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestUnitApigeeInstance_projectListDiffSuppress(t *testing.T) {
	for _, tc := range apigeeInstanceDiffSuppressTestCases {
		tc.Test(t)
	}
}

type ApigeeInstanceDiffSuppressTestCase struct {
	Name           string
	KeysToSuppress []string
	Before         map[string]interface{}
	After          map[string]interface{}
}

var apigeeInstanceDiffSuppressTestCases = []ApigeeInstanceDiffSuppressTestCase{
	{
		Name:           "projects with the same length and one project entry is converted to project id",
		KeysToSuppress: []string{"consumer_accept_list", "consumer_accept_list.0"},
		Before: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"45796856818", "12345"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "45796856818",
			"consumer_accept_list.1": "12345",
		},
		After: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"tf-test8v1bd04pxa", "12345"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "tf-test8v1bd04pxa",
			"consumer_accept_list.1": "12345",
		},
	},
	{
		Name:           "projects with the same length and no project conversion",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"tf-test8v1bd04pxa", "12345"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "tf-test8v1bd04pxa",
			"consumer_accept_list.1": "12345",
		},
		After: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"tf-test8v1bd04pxa", "12345"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "tf-test8v1bd04pxa",
			"consumer_accept_list.1": "12345",
		},
	},
	{
		Name:           "projects are empty",
		KeysToSuppress: []string{},
		Before:         map[string]interface{}{},
		After:          map[string]interface{}{},
	},
	{
		Name:           "projects have the different length",
		KeysToSuppress: []string{},
		Before:         map[string]interface{}{},
		After: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"tf-test8v1bd04pxa", "12345"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "tf-test8v1bd04pxa",
			"consumer_accept_list.1": "12345",
		},
	},
	{
		Name:           "positive case: state contains configured projects plus auto-injected tenant project",
		KeysToSuppress: []string{"consumer_accept_list", "consumer_accept_list.#", "consumer_accept_list.1"},
		Before: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"my-project", "f0d55e5054793ca7d-tp"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "my-project",
			"consumer_accept_list.1": "f0d55e5054793ca7d-tp",
		},
		After: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"my-project"},
			"consumer_accept_list.#": 1,
			"consumer_accept_list.0": "my-project",
		},
	},
	{
		Name:           "positive case: state contains transformed project ID plus auto-injected tenant project",
		KeysToSuppress: []string{"consumer_accept_list", "consumer_accept_list.#", "consumer_accept_list.0", "consumer_accept_list.2"},
		Before: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"proj-1-id", "proj-2", "auto-injected-tp"},
			"consumer_accept_list.#": 3,
			"consumer_accept_list.0": "proj-1-id",
			"consumer_accept_list.1": "proj-2",
			"consumer_accept_list.2": "auto-injected-tp",
		},
		After: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"1234567890", "proj-2"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "1234567890",
			"consumer_accept_list.1": "proj-2",
		},
	},
	{
		Name:           "negative case: state contains different numbers of non-tenant projects",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"project-1", "tenant-tp"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "project-1",
			"consumer_accept_list.1": "tenant-tp",
		},
		After: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"project-1", "project-2"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "project-1",
			"consumer_accept_list.1": "project-2",
		},
	},
	{
		Name:           "negative case: state contains non-tenant project and tenant project vs empty config",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"project-1", "tenant-tp"},
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "project-1",
			"consumer_accept_list.1": "tenant-tp",
		},
		After: map[string]interface{}{
			"consumer_accept_list":   []interface{}{},
			"consumer_accept_list.#": 0,
		},
	},
	{
		Name:           "negative case: state contains only tenant project vs 1 configured project",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"tenant-tp"},
			"consumer_accept_list.#": 1,
			"consumer_accept_list.0": "tenant-tp",
		},
		After: map[string]interface{}{
			"consumer_accept_list":   []interface{}{"project-1"},
			"consumer_accept_list.#": 1,
			"consumer_accept_list.0": "project-1",
		},
	},
	{
		Name:           "fallback: same length without slice",
		KeysToSuppress: []string{"consumer_accept_list.0"},
		Before: map[string]interface{}{
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "45796856818",
			"consumer_accept_list.1": "12345",
		},
		After: map[string]interface{}{
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "tf-test8v1bd04pxa",
			"consumer_accept_list.1": "12345",
		},
	},
	{
		Name:           "fallback: different length without slice",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"consumer_accept_list.#": 1,
			"consumer_accept_list.0": "project-1",
		},
		After: map[string]interface{}{
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "project-1",
			"consumer_accept_list.1": "project-2",
		},
	},
}

func (tc *ApigeeInstanceDiffSuppressTestCase) Test(t *testing.T) {
	mockResourceDiff := &tpgresource.ResourceDiffMock{
		Before: tc.Before,
		After:  tc.After,
	}

	keysHavingDiff := map[string]bool{}

	for key, val1 := range tc.Before {
		val2, ok := tc.After[key]
		if !ok {
			keysHavingDiff[key] = true
		} else if !reflect.DeepEqual(val1, val2) {
			keysHavingDiff[key] = true
		}
	}

	for key, val1 := range tc.After {
		val2, ok := tc.Before[key]
		if !ok {
			keysHavingDiff[key] = true
		} else if !reflect.DeepEqual(val1, val2) {
			keysHavingDiff[key] = true
		}
	}

	keySuppressionMap := map[string]bool{}
	for key := range tc.Before {
		keySuppressionMap[key] = false
	}
	for key := range tc.After {
		keySuppressionMap[key] = false
	}

	for _, key := range tc.KeysToSuppress {
		keySuppressionMap[key] = true
	}

	for key := range keysHavingDiff {
		actual := apigee.ProjectListDiffSuppressFunc(mockResourceDiff)
		if actual != keySuppressionMap[key] {
			var expectation string
			if keySuppressionMap[key] {
				expectation = "be"
			} else {
				expectation = "not be"
			}
			t.Errorf("Test %s: expected key `%s` to %s suppressed", tc.Name, key, expectation)
		}
	}
}
