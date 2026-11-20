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
		backup      map[string]interface{}
		partitionId string
		want        bool
	}{
		"only partition":              {backup("part"), "part", true},
		"one of several partitions":   {backup("default", "part"), "part", true},
		"other partitions":            {backup("default", "other"), "part", false},
		"partition ID is a prefix":    {backup("part-2"), "part", false},
		"partition ID is a suffix":    {backup("my-part"), "part", false},
		"no partitions":               {backup(), "part", false},
		"missing instancePartitions":  {map[string]interface{}{"name": "projects/p/instances/i/backups/b"}, "part", false},
		"malformed instancePartition": {map[string]interface{}{"instancePartitions": []interface{}{nil, map[string]interface{}{}}}, "part", false},
		"empty partition ID":          {map[string]interface{}{"instancePartitions": []interface{}{map[string]interface{}{"instancePartition": ""}}}, "", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := spannerBackupReferencesInstancePartition(tc.backup, tc.partitionId); got != tc.want {
				t.Errorf("spannerBackupReferencesInstancePartition() = %t, want %t", got, tc.want)
			}
		})
	}
}
