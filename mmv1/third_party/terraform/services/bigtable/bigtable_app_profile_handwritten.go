package bigtable

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/bigtableadmin/v2"
)

func internalExpandBigtableAppProfileMultiClusterRoutingUseAny(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil || !v.(bool) {
		return nil, nil
	}

	obj := bigtableadmin.MultiClusterRoutingUseAny{}

	clusterIds := d.Get("multi_cluster_routing_cluster_ids").([]interface{})

	for _, id := range clusterIds {
		obj.ClusterIds = append(obj.ClusterIds, id.(string))
	}

	affinity, _ := d.GetOkExists("row_affinity")
	if affinity != nil && affinity == true {
		obj.RowAffinity = &bigtableadmin.RowAffinity{}
	} else {
		obj.RowAffinity = nil
	}

	return obj, nil
}
