package netapp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestNetappVolumeHybridReplicationParametersDiffSuppress(t *testing.T) {
	dExisting := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"hybrid_replication_parameters": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"peer_cluster_name": {
						Type: schema.TypeString,
					},
				},
			},
		},
	}, map[string]interface{}{})
	dExisting.SetId("projects/my-project/locations/us-central1/volumes/my-volume")

	dNew := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"hybrid_replication_parameters": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"peer_cluster_name": {
						Type: schema.TypeString,
					},
				},
			},
		},
	}, map[string]interface{}{})

	cases := map[string]struct {
		Data               *schema.ResourceData
		Key                string
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"existing resource with empty count and non-empty config (upgrade scenario)": {
			Data:               dExisting,
			Key:                "hybrid_replication_parameters.#",
			Old:                "0",
			New:                "1",
			ExpectDiffSuppress: true,
		},
		"existing resource with empty string and non-empty config": {
			Data:               dExisting,
			Key:                "hybrid_replication_parameters.0.peer_cluster_name",
			Old:                "",
			New:                "cluster1",
			ExpectDiffSuppress: true,
		},
		"existing resource with non-empty state and same config": {
			Data:               dExisting,
			Key:                "hybrid_replication_parameters.0.peer_cluster_name",
			Old:                "cluster1",
			New:                "cluster1",
			ExpectDiffSuppress: false,
		},
		"existing resource with non-empty state and different config": {
			Data:               dExisting,
			Key:                "hybrid_replication_parameters.0.peer_cluster_name",
			Old:                "cluster1",
			New:                "cluster2",
			ExpectDiffSuppress: false,
		},
		"new resource creation with non-empty config": {
			Data:               dNew,
			Key:                "hybrid_replication_parameters.#",
			Old:                "0",
			New:                "1",
			ExpectDiffSuppress: false,
		},
		"new resource creation with empty config": {
			Data:               dNew,
			Key:                "hybrid_replication_parameters.#",
			Old:                "0",
			New:                "0",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if netappVolumeHybridReplicationParametersDiffSuppress(tc.Key, tc.Old, tc.New, tc.Data) != tc.ExpectDiffSuppress {
			t.Errorf("failed case %s: Key='%s', Old='%s', New='%s', expected suppress=%t", tn, tc.Key, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}
