package bigquery

import (
	"testing"
)

func TestBigqueryDatasetAccessHash_nilVsEmptyString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		a    map[string]interface{}
		b    map[string]interface{}
	}{
		{
			name: "role with nil vs empty string fields",
			a: map[string]interface{}{
				"role":           "OWNER",
				"user_by_email":  "alice@example.com",
				"group_by_email": nil,
				"domain":         nil,
				"special_group":  nil,
				"iam_member":     nil,
				"view":           nil,
				"dataset":        nil,
				"routine":        nil,
				"condition":      nil,
			},
			b: map[string]interface{}{
				"role":           "OWNER",
				"user_by_email":  "alice@example.com",
				"group_by_email": "",
				"domain":         "",
				"special_group":  "",
				"iam_member":     "",
				"view":           nil,
				"dataset":        nil,
				"routine":        nil,
				"condition":      nil,
			},
		},
		{
			name: "domain access with nil vs empty string",
			a: map[string]interface{}{
				"role":           "READER",
				"domain":         "hashicorp.com",
				"user_by_email":  nil,
				"group_by_email": nil,
				"special_group":  nil,
				"iam_member":     nil,
				"view":           nil,
				"dataset":        nil,
				"routine":        nil,
				"condition":      nil,
			},
			b: map[string]interface{}{
				"role":           "READER",
				"domain":         "hashicorp.com",
				"user_by_email":  "",
				"group_by_email": "",
				"special_group":  "",
				"iam_member":     "",
				"view":           nil,
				"dataset":        nil,
				"routine":        nil,
				"condition":      nil,
			},
		},
		{
			name: "iam_member access with nil vs empty string",
			a: map[string]interface{}{
				"role":           "READER",
				"iam_member":     "allUsers",
				"user_by_email":  nil,
				"group_by_email": nil,
				"domain":         nil,
				"special_group":  nil,
				"view":           nil,
				"dataset":        nil,
				"routine":        nil,
				"condition":      nil,
			},
			b: map[string]interface{}{
				"role":           "READER",
				"iam_member":     "allUsers",
				"user_by_email":  "",
				"group_by_email": "",
				"domain":         "",
				"special_group":  "",
				"view":           nil,
				"dataset":        nil,
				"routine":        nil,
				"condition":      nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hashA := resourceBigqueryDatasetAccessHash(tc.a)
			hashB := resourceBigqueryDatasetAccessHash(tc.b)
			if hashA != hashB {
				t.Errorf("hash mismatch: nil-fields=%d, empty-string-fields=%d", hashA, hashB)
			}
		})
	}
}
