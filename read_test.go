package test

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"
	"go.uber.org/zap/zaptest"
)

func TestReadPlannedAssetsCoverage(t *testing.T) {
	cases := []struct {
		name string
	}{
		// read-only, the following tests are not in cli_test or
		// have unique parameters that separate them
		{name: "example_folder_iam_binding"},
		{name: "example_folder_iam_member"},
		{name: "example_project_create"},
		{name: "example_project_update"},
		{name: "example_project_iam_binding"},
		{name: "example_project_iam_member"},
		{name: "example_storage_bucket"},
		{name: "example_storage_bucket_iam_binding"},
		{name: "example_storage_bucket_iam_member"},
		{name: "example_project_create_empty_project_id"},
		{name: "example_project_iam_member_empty_project"},
		// auto inserted tests that are not in list above or manually inserted in cli_test.go
		{name: "example_access_context_manager_access_policy"},
		{name: "example_access_context_manager_service_perimeter"},
		{name: "example_bigquery_dataset"},
		{name: "example_bigquery_dataset_iam_binding"},
		{name: "example_bigquery_dataset_iam_member"},
		{name: "example_bigquery_dataset_iam_policy"},
		{name: "example_bigquery_dataset_iam_policy_empty_policy_data"},
		{name: "example_bigquery_table"},
		{name: "example_bigtable_instance"},
		{name: "example_cloud_run_mapping"},
		{name: "example_cloud_run_service"},
		{name: "example_cloud_run_service_iam_binding"},
		{name: "example_cloud_run_service_iam_member"},
		{name: "example_cloud_run_service_iam_policy"},
		{name: "example_compute_address"},
		{name: "example_compute_disk"},
		{name: "example_compute_disk_empty_image"},
		{name: "example_compute_firewall"},
		{name: "example_compute_global_address"},
		{name: "example_compute_global_forwarding_rule"},
		{name: "example_compute_instance_iam_binding"},
		{name: "example_compute_instance_iam_member"},
		{name: "example_compute_instance_iam_policy"},
		{name: "example_compute_network"},
		{name: "example_compute_snapshot"},
		{name: "example_compute_ssl_policy"},
		{name: "example_compute_subnetwork"},
		{name: "example_compute_target_https_proxy"},
		{name: "example_compute_target_ssl_proxy"},
		{name: "example_container_cluster"},
		{name: "example_dns_managed_zone"},
		{name: "example_dns_policy"},
		{name: "example_filestore_instance"},
		{name: "example_folder_iam_member_empty_folder"},
		{name: "example_folder_iam_policy"},
		{name: "example_folder_organization_policy"},
		{name: "example_google_cloudfunctions_function"},
		{name: "example_google_sql_database"},
		{name: "example_kms_crypto_key"},
		{name: "example_kms_crypto_key_iam_binding"},
		{name: "example_kms_crypto_key_iam_member"},
		{name: "example_kms_crypto_key_iam_policy"},
		{name: "example_kms_key_ring"},
		{name: "example_kms_key_ring_iam_binding"},
		{name: "example_kms_key_ring_iam_member"},
		{name: "example_kms_key_ring_iam_policy"},
		{name: "example_logging_metric"},
		{name: "example_monitoring_notification_channel"},
		{name: "example_organization_iam_binding"},
		{name: "example_organization_iam_custom_role"},
		{name: "example_organization_iam_member"},
		{name: "example_organization_iam_policy"},
		{name: "example_organization_policy"},
		{name: "example_project_iam"},
		{name: "example_project_iam_custom_role"},
		{name: "example_project_iam_policy"},
		{name: "example_project_in_folder"},
		{name: "example_project_in_org"},
		{name: "example_project_organization_policy"},
		{name: "example_project_service"},
		{name: "example_pubsub_lite_reservation"},
		{name: "example_pubsub_lite_subscription"},
		{name: "example_pubsub_lite_topic"},
		{name: "example_pubsub_schema"},
		{name: "example_pubsub_subscription"},
		{name: "example_pubsub_subscription_iam_binding"},
		{name: "example_pubsub_subscription_iam_member"},
		{name: "example_pubsub_subscription_iam_policy"},
		{name: "example_pubsub_topic"},
		{name: "example_secret_manager_secret_iam_binding"},
		{name: "example_secret_manager_secret_iam_member"},
		{name: "example_secret_manager_secret_iam_policy"},
		{name: "example_service_account"},
		{name: "example_service_account_update"},
		{name: "example_spanner_database"},
		{name: "example_spanner_database_iam_binding"},
		{name: "example_spanner_database_iam_member"},
		{name: "example_spanner_database_iam_policy"},
		{name: "example_spanner_instance_iam_binding"},
		{name: "example_spanner_instance_iam_member"},
		{name: "example_spanner_instance_iam_policy"},
		{name: "example_sql_database_instance"},
		{name: "example_storage_bucket_iam_member_random_suffix"},
		{name: "example_storage_bucket_iam_policy"},
		{name: "example_vpc_access_connector"},
		{name: "full_compute_firewall"},
		{name: "full_compute_instance"},
		{name: "full_container_cluster"},
		{name: "full_container_node_pool"},
		{name: "full_spanner_instance"},
		{name: "full_sql_database_instance"},
		{name: "full_storage_bucket"},
	}
	for i := range cases {
		// Allocate a variable to make sure test can run in parallel.
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for running terraform.
			dir, err := ioutil.TempDir(tmpDir, "terraform")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)

			generateTestFiles(t, "../testdata/templates", dir, c.name+".json")
			generateTestFiles(t, "../testdata/templates", dir, c.name+".tfplan.json")

			// Unmarshal payload from testfile into `want` variable.
			f := filepath.Join(dir, c.name+".json")
			want, err := readExpectedTestFile(f)
			if err != nil {
				t.Fatal(err)
			}

			planfile := filepath.Join(dir, c.name+".tfplan.json")
			ctx := context.Background()
			ancestryCache := map[string]string{
				data.Provider["project"]: data.Ancestry,
			}
			got, err := tfgcv.ReadPlannedAssets(ctx, planfile, data.Provider["project"], "", "", ancestryCache, true, false, zaptest.NewLogger(t), "")
			if err != nil {
				t.Fatalf("ReadPlannedAssets(%s, %s, \"\", \"\", %s, %t): %v", planfile, data.Provider["project"], ancestryCache, true, err)
			}

			expectedAssets := normalizeAssets(t, want, true)
			actualAssets := normalizeAssets(t, got, true)
			require.ElementsMatch(t, actualAssets, expectedAssets)
		})
	}
}

func TestReadPlannedAssetsCoverage_WithoutDefaultProject(t *testing.T) {
	cases := []struct {
		name string
	}{
		{name: "example_project_create_empty_project_id"},
		{name: "example_storage_bucket"},
		{name: "example_project_iam_member_empty_project"},
	}
	for i := range cases {
		// Allocate a variable to make sure test can run in parallel.
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for running terraform.
			dir, err := ioutil.TempDir(tmpDir, "terraform")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)

			generateTestFiles(t, "../testdata/templates", dir, c.name+"_without_default_project.json")
			generateTestFiles(t, "../testdata/templates", dir, c.name+".tfplan.json")

			// Unmarshal payload from testfile into `want` variable.
			f := filepath.Join(dir, c.name+"_without_default_project.json")
			want, err := readExpectedTestFile(f)
			if err != nil {
				t.Fatal(err)
			}

			planfile := filepath.Join(dir, c.name+".tfplan.json")
			ctx := context.Background()
			ancestryCache := map[string]string{
				// data.Provider["project"]: data.Ancestry,
			}
			got, err := tfgcv.ReadPlannedAssets(ctx, planfile, "", "", "", ancestryCache, true, false, zaptest.NewLogger(t), "")
			if err != nil {
				t.Fatalf("ReadPlannedAssets(%s, %s, \"\", \"\", %s, %t): %v", planfile, data.Provider["project"], ancestryCache, true, err)
			}

			expectedAssets := normalizeAssets(t, want, true)
			actualAssets := normalizeAssets(t, got, true)
			require.ElementsMatch(t, actualAssets, expectedAssets)
		})
	}
}
