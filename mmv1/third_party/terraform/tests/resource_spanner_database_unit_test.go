package google

import (
	"fmt"
	"testing"
)

// Unit Tests for type SpannerDatabaseId
func TestDatabaseNameForApi(t *testing.T) {
	id := SpannerDatabaseId{
		Project:  "project123",
		Instance: "instance456",
		Database: "db789",
	}
	actual := id.databaseUri()
	expected := "projects/project123/instances/instance456/databases/db789"
	expectEquals(t, expected, actual)
}

// Unit Tests for ForceNew when the change in ddl
func TestSpannerDatabase_resourceSpannerDBDdlCustomDiffFuncForceNew(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		before   interface{}
		after    interface{}
		forcenew bool
	}{
		"remove_old_statements": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)"},
			after: []interface{}{
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)"},
			forcenew: true,
		},
		"append_new_statements": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)"},
			after: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
			},
			forcenew: false,
		},
		"no_change": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)"},
			after: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)"},
			forcenew: false,
		},
		"order_of_statments_change": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
				"CREATE TABLE t3 (t3 INT64 NOT NULL,) PRIMARY KEY(t3)",
			},
			after: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t3 (t3 INT64 NOT NULL,) PRIMARY KEY(t3)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
			},
			forcenew: true,
		},
		"missing_an_old_statement": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
				"CREATE TABLE t3 (t3 INT64 NOT NULL,) PRIMARY KEY(t3)",
			},
			after: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
			},
			forcenew: true,
		},
	}

	for tn, tc := range cases {
		d := &ResourceDiffMock{
			Before: map[string]interface{}{
				"ddl": tc.before,
			},
			After: map[string]interface{}{
				"ddl": tc.after,
			},
		}
		err := resourceSpannerDBDdlCustomDiffFunc(d)
		if err != nil {
			t.Errorf("failed, expected no error but received - %s for the condition %s", err, tn)
		}
		if d.IsForceNew != tc.forcenew {
			t.Errorf("ForceNew not setup correctly for the condition-'%s', expected:%v;actual:%v", tn, tc.forcenew, d.IsForceNew)
		}
	}
}

// Unit Tests for validation of retention period argument
func TestValidateDatabaseRetentionPeriod(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input       string
		expectError bool
	}{
		// Not valid input
		"empty_string": {
			input:       "",
			expectError: true,
		},
		"number_with_no_unit": {
			input:       "1",
			expectError: true,
		},
		"less_than_1h": {
			input:       "59m",
			expectError: true,
		},
		"more_than_7days": {
			input:       "8d",
			expectError: true,
		},
		// Valid input
		"1_hour_in_secs": {
			input:       "3600s",
			expectError: false,
		},
		"1_hour_in_mins": {
			input:       "60m",
			expectError: false,
		},
		"1_hour_in_hours": {
			input:       "1h",
			expectError: false,
		},
		"7_days_in_secs": {
			input:       fmt.Sprintf("%ds", 7*24*60*60),
			expectError: false,
		},
		"7_days_in_mins": {
			input:       fmt.Sprintf("%dm", 7*24*60),
			expectError: false,
		},
		"7_days_in_hours": {
			input:       fmt.Sprintf("%dh", 7*24),
			expectError: false,
		},
		"7_days_in_days": {
			input:       "7d",
			expectError: false,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			_, errs := ValidateDatabaseRetentionPeriod(tc.input, "foobar")
			var wantErrCount string
			if tc.expectError {
				wantErrCount = "1+"
			} else {
				wantErrCount = "0"
			}
			if (len(errs) > 0 && tc.expectError == false) || (len(errs) == 0 && tc.expectError == true) {
				t.Errorf("failed, expected `%s` test case validation to have %s errors", tn, wantErrCount)
			}
		})
	}
}
