package iambeta_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/iambeta"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccIAMBetaWorkloadIdentityPool_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectNumberFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMBetaWorkloadIdentityPool_full(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPool_update(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIAMBetaWorkloadIdentityPool_minimal(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectNumberFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMBetaWorkloadIdentityPool_minimal(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPool_update(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIAMBetaWorkloadIdentityPool_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectNumberFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMBetaWorkloadIdentityPool_full(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPool_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_iam_workload_identity_pool.my_pool", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPool_update_use_default_ca(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_iam_workload_identity_pool.my_pool", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPool_minimal(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_iam_workload_identity_pool.my_pool", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIAMBetaWorkloadIdentityPool_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
  display_name              = "Name of the pool"
  description               = "Identity pool operates in TRUST_DOMAIN mode"
  disabled                  = true
  mode                      = "TRUST_DOMAIN"
  inline_certificate_issuance_config {
    ca_pools = {      
      "us-central1" : "projects/project-bar/locations/us-central1/caPools/ca-pool-bar"
      "asia-east2" : "projects/project-foo/locations/asia-east2/caPools/ca-pool-foo"
    }
    lifetime                   = "86400s"
    rotation_window_percentage = 50
    key_algorithm              = "ECDSA_P256"
  }
  inline_trust_config {
    additional_trust_bundles {
      trust_domain            = "ca-pool-foo.global.project-foo.workload.id.goog"
      trust_default_shared_ca = false
      trust_anchors {
        pem_certificate = file("test-fixtures/trust_anchor_1.pem")
      }
      trust_anchors {
        pem_certificate = file("test-fixtures/trust_anchor_2.pem")
      }
    }
    additional_trust_bundles {
      trust_domain            = "ca-pool-bar.global.project-bar.workload.id.goog"
      trust_default_shared_ca = false
      trust_anchors {
        pem_certificate = file("test-fixtures/trust_anchor_3.pem")
      }
      trust_anchors {
        pem_certificate = file("test-fixtures/trust_anchor_4.pem")
      }
    }
  }
  attestation_rules {
    google_cloud_resource = "//run.googleapis.com/projects/%{project}/type/Service/*"
  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPool_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
  display_name              = "Updated name of the pool"
  description               = "Updated identity pool operates in TRUST_DOMAIN mode"
  disabled                  = false
  mode                      = "TRUST_DOMAIN"
  inline_certificate_issuance_config {
    ca_pools = {      
      "us-central2" : "projects/project-bar/locations/us-central2/caPools/ca-pool-bar"
      "asia-east1" : "projects/project-foo/locations/asia-east1/caPools/ca-pool-foo"
    }
    lifetime                   = "172800s"
    rotation_window_percentage = 75
    key_algorithm              = "RSA_4096"
  }
  inline_trust_config {
    additional_trust_bundles {
      trust_domain            = "ca-pool-baz.global.project-baz.workload.id.goog"
	  trust_default_shared_ca = false
      trust_anchors {
        pem_certificate = file("test-fixtures/trust_anchor_updated.pem")
      }
    }
  }
  attestation_rules {
    google_cloud_resource = "//run.googleapis.com/projects/%{project}/type/Service/*"
  }
  attestation_rules {
    google_cloud_resource = "//run.googleapis.com/projects/%{project}/type/Service/*"
  }
  attestation_rules {
    google_cloud_resource = "//run.googleapis.com/projects/%{project}/type/Service/*"
  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPool_update_use_default_ca(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
  display_name              = "Updated name of the pool"
  description               = "Updated identity pool operates in TRUST_DOMAIN mode"
  disabled                  = false
  mode                      = "TRUST_DOMAIN"
  inline_certificate_issuance_config {
    use_default_shared_ca      = true
    lifetime                   = "172800s"
    rotation_window_percentage = 75
    key_algorithm              = "RSA_4096"
  }
  inline_trust_config {
    additional_trust_bundles {
      trust_domain            = "ca-pool-baz.global.project-baz.workload.id.goog"
	  trust_default_shared_ca = true
      trust_anchors {
        pem_certificate = file("test-fixtures/trust_anchor_updated.pem")
      }
    }
  }
  attestation_rules {
    google_cloud_resource = "//run.googleapis.com/projects/%{project}/type/Service/*"
  }
  attestation_rules {
    google_cloud_resource = "//run.googleapis.com/projects/%{project}/type/Service/*"
  }
  attestation_rules {
    google_cloud_resource = "//run.googleapis.com/projects/%{project}/type/Service/*"
  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPool_minimal(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
  mode                      = "TRUST_DOMAIN"
}
`, context)
}

func TestAccIAMWorkloadIdentityPool_meshCaCompatibilityLifecycle(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Step 1: BACKWARD COMPATIBILITY
				// Create the resource completely WITHOUT your new field.
				// This ensures existing Terraform manifests do not break after your code change.
				Config: testAccIAMWorkloadIdentityPool_withoutMeshCa(context),
				Check: resource.ComposeTestCheckFunc(
					// Assert that the pool creates successfully and the field defaults to false.
					resource.TestCheckResourceAttr("google_iam_workload_identity_pool.example", "enable_mesh_ca_compatibility", "false"),
				),
			},
			{
				// Step 2: ENABLE THE FIELD
				// Apply an update configuration where the field is explicitly set to true.
				Config: testAccIAMWorkloadIdentityPool_withMeshCa(context, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_iam_workload_identity_pool.example", "enable_mesh_ca_compatibility", "true"),
				),
			},
			{
				// Step 3: DISABLE THE FIELD
				// Apply another update configuration turning the field back to false.
				Config: testAccIAMWorkloadIdentityPool_withMeshCa(context, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_iam_workload_identity_pool.example", "enable_mesh_ca_compatibility", "false"),
				),
			},
		},
	})
}

// Config 1: Simulates an existing / legacy resource configuration
func testAccIAMWorkloadIdentityPool_withoutMeshCa(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "example" {
  workload_identity_pool_id = "tf-test-pool-%{random_suffix}"
  display_name              = "Mesh CA Lifecycle Pool"
}
`, context)
}

// Config 2 & 3: Simulates explicitly interacting with the new field
func testAccIAMWorkloadIdentityPool_withMeshCa(context map[string]interface{}, enabled bool) string {
	context["enabled"] = enabled
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "example" {
  workload_identity_pool_id = "tf-test-pool-%{random_suffix}"
  display_name              = "Mesh CA Lifecycle Pool"
  enable_mesh_ca_compatibility = %{enabled}
}
`, context)
}
