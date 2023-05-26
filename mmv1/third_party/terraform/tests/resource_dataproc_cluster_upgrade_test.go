package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"google.golang.org/api/dataproc/v1"
)

// Tests schema version migration by creating a certificate with an old version of the provider (4.59.0)
// and then updating it with the current version the provider.
func TestAccDataprocClusterLabels_migration(t *testing.T) {
	SkipIfVcr(t)
	t.Parallel()

	rnd := RandString(t, 10)
	var cluster dataproc.Cluster
	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.65.0", // a version that doesn't support user-defined labels.
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccDataprocCluster_withLabels(rnd),
				ExternalProviders: oldVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_labels", &cluster),

					// We only provide one, but GCP adds three and we added goog-dataproc-autozone internally, so expect 5.
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.%", "5"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.key1", "value1"),
				),
			},
			{
				Config:                   testAccDataprocCluster_withLabels(rnd),
				ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_labels", &cluster),

					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.%", "5"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.key1", "value1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.%", "5"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.key1", "value1"),
				),
			},
		},
	})
}
