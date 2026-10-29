package vertexai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"
)

func TestIsEvaluationRunOperationInProgressError(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		err            error
		expected       bool
		expectedReason string
	}{
		"nil error": {
			err:            nil,
			expected:       false,
			expectedReason: "",
		},
		"evaluationRun delete 400 with other operations running": {
			err: &googleapi.Error{
				Code: 400,
				Body: `{"error":{"code":400,"message":"There are other operations running on the EvaluationRun \"projects/123/locations/us-central1/evaluationRuns/eval_run\". Operation(s) are: projects/123/locations/us-central1/operations/op1","status":"FAILED_PRECONDITION"}}`,
			},
			expected:       true,
			expectedReason: "Waiting for operations on EvaluationRun to finish",
		},
		"wrapped error with other operations running": {
			err:            fmt.Errorf("Error when reading or editing EvaluationRun: googleapi: Error 400: There are other operations running on the EvaluationRun \"projects/123/locations/us-central1/evaluationRuns/eval_run\". Operation(s) are: projects/123/locations/us-central1/operations/op1"),
			expected:       true,
			expectedReason: "Waiting for operations on EvaluationRun to finish",
		},
		"non-400 error": {
			err: &googleapi.Error{
				Code: 500,
				Body: `{"error":{"code":500,"message":"Internal error"}}`,
			},
			expected:       false,
			expectedReason: "",
		},
		"unrelated 400 error": {
			err: &googleapi.Error{
				Code: 400,
				Body: `{"error":{"code":400,"message":"Invalid argument"}}`,
			},
			expected:       false,
			expectedReason: "",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			isRetryable, reason := IsEvaluationRunOperationInProgressError(tc.err)
			if isRetryable != tc.expected || reason != tc.expectedReason {
				t.Errorf("IsEvaluationRunOperationInProgressError(%v) = (%v, %q), want (%v, %q)", tc.err, isRetryable, reason, tc.expected, tc.expectedReason)
			}
		})
	}
}

func TestFlattenVertexAIEvaluationRunMetadata(t *testing.T) {
	t.Parallel()

	resourceSchema := map[string]*schema.Schema{
		"metadata": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	cases := map[string]struct {
		configMetadata string
		apiValue       interface{}
		expected       interface{}
	}{
		"nil API value": {
			configMetadata: "",
			apiValue:       nil,
			expected:       nil,
		},
		"server-injected pipeline_id only when metadata unset returns nil": {
			configMetadata: "",
			apiValue: map[string]interface{}{
				"pipeline_id": "12345",
			},
			expected: nil,
		},
		"server-injected pipeline_id stripped alongside user metadata": {
			configMetadata: `{"tracker":"eval-tracker"}`,
			apiValue: map[string]interface{}{
				"tracker":     "eval-tracker",
				"pipeline_id": "12345",
			},
			expected: `{"tracker":"eval-tracker"}`,
		},
		"user-configured pipeline_id preserved": {
			configMetadata: `{"pipeline_id":"12345","tracker":"eval-tracker"}`,
			apiValue: map[string]interface{}{
				"pipeline_id": "12345",
				"tracker":     "eval-tracker",
			},
			expected: `{"pipeline_id":"12345","tracker":"eval-tracker"}`,
		},
		"explicit empty JSON object preserved when server injects pipeline_id": {
			configMetadata: `{}`,
			apiValue: map[string]interface{}{
				"pipeline_id": "12345",
			},
			expected: `{}`,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			raw := map[string]interface{}{}
			if tc.configMetadata != "" {
				raw["metadata"] = tc.configMetadata
			}
			d := schema.TestResourceDataRaw(t, resourceSchema, raw)
			got := flattenVertexAIEvaluationRunMetadata(tc.apiValue, d, nil)
			if got != tc.expected {
				t.Errorf("flattenVertexAIEvaluationRunMetadata() = %#v, want %#v", got, tc.expected)
			}
		})
	}
}
