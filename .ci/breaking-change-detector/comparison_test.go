package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestComparisonEngine(t *testing.T) {
	for _, tc := range comparisonEngineTestCases {
		tc.check(t)
	}
}

type comparisonEngineTestCase struct {
	name               string
	oldResourceMap     map[string]*schema.Resource
	newResourceMap     map[string]*schema.Resource
	expectedViolations int
}

var comparisonEngineTestCases = []comparisonEngineTestCase{
	{
		name: "control",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		newResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		expectedViolations: 0,
	},
	{
		name: "adding resources",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		newResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
			"google-y": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
				},
			},
		},
		expectedViolations: 0,
	},
	{
		name: "adding fields",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		newResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
					"field-c": {Description: "beep", Optional: true},
				},
			},
		},
		expectedViolations: 0,
	},
	{
		name: "resource missing",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep"},
					"field-b": {Description: "beep"},
				},
			},
		},
		newResourceMap:     map[string]*schema.Resource{},
		expectedViolations: 1,
	},
	{
		name: "field missing",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		newResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
				},
			},
		},
		expectedViolations: 1,
	},
	{
		name: "optional field to required",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		newResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Required: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		expectedViolations: 1,
	},
	{
		name: "field missing and optional to required",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
		},
		newResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Required: true},
				},
			},
		},
		expectedViolations: 2,
	},
	{
		name: "field missing, resource missing, and optional to required",
		oldResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
					"field-b": {Description: "beep", Optional: true},
				},
			},
			"google-y": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Optional: true},
				},
			},
		},
		newResourceMap: map[string]*schema.Resource{
			"google-x": {
				Schema: map[string]*schema.Schema{
					"field-a": {Description: "beep", Required: true},
				},
			},
		},
		expectedViolations: 3,
	},
}

func (tc *comparisonEngineTestCase) check(t *testing.T) {
	violations := compareResourceMaps(tc.oldResourceMap, tc.newResourceMap)
	if tc.expectedViolations != len(violations) {
		t.Errorf("Test `%s` failed: expected %d violations, got %d", tc.name, tc.expectedViolations, len(violations))
	}
}
