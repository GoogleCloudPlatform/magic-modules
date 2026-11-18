package spanner

import (
	"testing"
)

func TestSpannerBackupReferencesInstancePartition(t *testing.T) {
	t.Parallel()

	backup := func(partitions ...string) map[string]interface{} {
		refs := []interface{}{}
		for _, p := range partitions {
			refs = append(refs, map[string]interface{}{
				"instancePartition": "projects/p/instances/i/instancePartitions/" + p,
			})
		}
		return map[string]interface{}{
			"name":               "projects/p/instances/i/backups/b",
			"instancePartitions": refs,
		}
	}

	cases := map[string]struct {
		backup map[string]interface{}
		want   bool
	}{
		"only partition":              {backup("part"), true},
		"one of several partitions":   {backup("default", "part"), true},
		"other partitions":            {backup("default", "other"), false},
		"partition ID is a prefix":    {backup("part-2"), false},
		"partition ID is a suffix":    {backup("my-part"), false},
		"no partitions":               {backup(), false},
		"missing instancePartitions":  {map[string]interface{}{"name": "projects/p/instances/i/backups/b"}, false},
		"malformed instancePartition": {map[string]interface{}{"instancePartitions": []interface{}{nil, map[string]interface{}{}}}, false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := spannerBackupReferencesInstancePartition(tc.backup, "part"); got != tc.want {
				t.Errorf("spannerBackupReferencesInstancePartition() = %t, want %t", got, tc.want)
			}
		})
	}
}
