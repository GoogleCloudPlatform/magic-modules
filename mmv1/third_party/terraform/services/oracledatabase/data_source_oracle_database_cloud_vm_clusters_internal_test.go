package oracledatabase

import (
	"reflect"
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestFlattenOracleDatabaseCloudVmClusters_identityConnectorFields(t *testing.T) {
	t.Parallel()

	d := DataSourceOracleDatabaseCloudVmClusters().Data(nil)
	cfg := &transport_tpg.Config{}

	rawClusters := []interface{}{
		map[string]interface{}{
			"name":        "projects/p/locations/us-east4/cloudVmClusters/c1",
			"displayName": "c1",
			"identityConnector": map[string]interface{}{
				"serviceAgentEmail":     "p123-sa@gcp-sa-oci.iam.gserviceaccount.com",
				"connectionState":       "PARTIALLY_CONNECTED",
				"workloadIdentityPool":  "projects/123/locations/global/workloadIdentityPools/pool-1",
				"identityConnectorId":   "ocid1.oracledbgcpidentityconnector.oc1.iad.abc",
				"lifecycleStateDetails": "Waiting for IAM binding",
			},
		},
		map[string]interface{}{
			"name":        "projects/p/locations/us-east4/cloudVmClusters/c2",
			"displayName": "c2",
		},
	}

	got := flattenOracleDatabaseCloudVmClusters(rawClusters, d, cfg)
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}

	wantConnector := []interface{}{
		map[string]interface{}{
			"service_agent_email":     "p123-sa@gcp-sa-oci.iam.gserviceaccount.com",
			"connection_state":        "PARTIALLY_CONNECTED",
			"workload_identity_pool":  "projects/123/locations/global/workloadIdentityPools/pool-1",
			"identity_connector_id":   "ocid1.oracledbgcpidentityconnector.oc1.iad.abc",
			"lifecycle_state_details": "Waiting for IAM binding",
		},
	}
	if !reflect.DeepEqual(got[0]["identity_connector"], wantConnector) {
		t.Fatalf("got[0][\"identity_connector\"] = %#v, want %#v", got[0]["identity_connector"], wantConnector)
	}
	if got[1]["identity_connector"] != nil {
		t.Fatalf("got[1][\"identity_connector\"] = %#v, want nil", got[1]["identity_connector"])
	}
}
