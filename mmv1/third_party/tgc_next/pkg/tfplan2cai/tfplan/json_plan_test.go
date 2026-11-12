package tfplan

import (
	"encoding/json"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/require"
)

func newPlan(t *testing.T) []byte {
	t.Helper()
	return []byte(`
{
  "format_version": "0.1",
  "planned_values": {
    "root_module": {
      "child_modules": [
        {
          "address": "module.foo",
          "resources": [
            {
              "address": "module.foo.google_compute_instance.quz1",
              "values": {
								"key1": "value1",
								"nestedKey1": { "insideKey1": "insideValue1"}
              }
            }
          ],
          "child_modules": [
            {
              "address": "module.foo.bar",
              "resources": [
                {
                  "address": "module.foo.bar.google_compute_instance.quz2",
									"values": {"key2": "value2"}
                },
								{
									"address": "module.foo.bar.google_compute_instance.quz4",
									"values": {"key4": "value4"}
								}
              ]
            }
          ]
        }
      ]
    }
  },
	"resource_changes": [
		{
			"address": "module.foo.google_compute_instance.quz1",
			"mode": "managed",
			"type": "google_compute_instance",
			"name": "quz1",
			"provider_name": "google",
			"change": {
				"actions": ["delete", "create"],
				"before": {"key1": "value1"},
				"after": {
					"key1": "value1",
					"nestedKey1": { "insideKey1": "insideValue1"}
				}
			}
		},
		{
			"address": "module.foo.bar.google_compute_instance.quz2",
			"mode": "managed",
			"type": "google_compute_instance",
			"name": "quz2",
			"provider_name": "google",
			"change": {
				"actions": ["noop"],
				"before": {"key2": "value2"},
				"after": {"key2": "value2"}
			}
		},
		{
			"address": "module.foo.bar.google_compute_instance.quz3",
			"mode": "managed",
			"type": "google_compute_instance",
			"name": "quz3",
			"provider_name": "google",
			"change": {
				"actions": ["delete"],
				"before": {"key3": "value3"},
				"after": {}
			}
		},
		{
			"address": "module.foo.bar.google_compute_instance.quz4",
			"mode": "managed",
			"type": "google_compute_instance",
			"name": "quz4",
			"provider_name": "google",
			"change": {
				"actions": ["create"],
				"before": {},
				"after": {"key4": "value4"}
			}
		}
	]
}
`)
}

func TestReadResourceChanges(t *testing.T) {
	wantJSON := []byte(`
[
	{
		"address": "module.foo.google_compute_instance.quz1",
		"mode": "managed",
		"type": "google_compute_instance",
		"name": "quz1",
		"provider_name": "google",
		"change": {
			"actions": ["delete", "create"],
			"before": {"key1": "value1"},
			"after": {
				"key1": "value1",
				"nestedKey1": { "insideKey1": "insideValue1"}
			}
		}
	},
	{
		"address": "module.foo.bar.google_compute_instance.quz2",
		"mode": "managed",
		"type": "google_compute_instance",
		"name": "quz2",
		"provider_name": "google",
		"change": {
			"actions": ["noop"],
			"before": {"key2": "value2"},
			"after": {"key2": "value2"}
		}
	},
	{
		"address": "module.foo.bar.google_compute_instance.quz3",
		"mode": "managed",
		"type": "google_compute_instance",
		"name": "quz3",
		"provider_name": "google",
		"change": {
			"actions": ["delete"],
			"before": {"key3": "value3"},
			"after": {}
		}
	},
	{
		"address": "module.foo.bar.google_compute_instance.quz4",
		"mode": "managed",
		"type": "google_compute_instance",
		"name": "quz4",
		"provider_name": "google",
		"change": {
			"actions": ["create"],
			"before": {},
			"after": {"key4": "value4"}
		}
	}
]
`)
	data := newPlan(t)
	rcs, err := ReadResourceChanges(data)
	if err != nil {
		t.Fatalf("parsing %s: %v", string(data), err)
	}
	gotJSON, err := json.Marshal(rcs)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}
	require.JSONEq(t, string(wantJSON), string(gotJSON))
}

func TestIsReplace(t *testing.T) {
	cases := []struct {
		name    string
		actions tfjson.Actions
		want    bool
	}{
		{name: "destroy_before_create", actions: tfjson.Actions{"delete", "create"}, want: true},
		{name: "create_before_destroy", actions: tfjson.Actions{"create", "delete"}, want: true},
		{name: "create", actions: tfjson.Actions{"create"}, want: false},
		{name: "update", actions: tfjson.Actions{"update"}, want: false},
		{name: "delete", actions: tfjson.Actions{"delete"}, want: false},
		{name: "no-op", actions: tfjson.Actions{"no-op"}, want: false},
		{name: "read", actions: tfjson.Actions{"read"}, want: false},
		{name: "forget", actions: tfjson.Actions{"forget"}, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rc := &tfjson.ResourceChange{Change: &tfjson.Change{Actions: tc.actions}}
			require.Equal(t, tc.want, IsReplace(rc))
			require.Equal(t, tc.want, IsDeleteCreate(rc))
		})
	}
}
