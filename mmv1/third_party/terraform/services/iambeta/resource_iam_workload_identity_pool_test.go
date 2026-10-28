package iambeta_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/iambeta"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

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

func TestAccIAMBetaWorkloadIdentityPool_undeleteOnCreate(t *testing.T) {
	t.Parallel()

	context := map[string]any{
		"project":       envvar.GetTestProjectNumberFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMBetaWorkloadIdentityPoolUndeleteInitial(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccIAMBetaWorkloadIdentityPoolUndeleteInitial(context),
				Destroy: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPoolUndeleteRecreated(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIAMBetaWorkloadIdentityPoolUndeleteInitial(context map[string]any) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
  display_name              = "Initial pool name"
  description               = "Initial pool description"
  disabled                  = true
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPoolUndeleteRecreated(context map[string]any) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
  display_name              = "Recreated pool name"
  description               = "Recreated pool description"
  disabled                  = false
}
`, context)
}

type recordedRequest struct {
	Method     string
	Path       string
	UpdateMask string
	Body       map[string]any
}

func TestResourceIAMBetaWorkloadIdentityPoolCreate_undeleteWhenDeleted(t *testing.T) {
	_ = t.Context()

	var requests []recordedRequest
	undeleted := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var body map[string]any
		if r.Body != nil {
			raw, _ := io.ReadAll(r.Body)
			if len(raw) > 0 {
				_ = json.Unmarshal(raw, &body)
			}
		}

		if strings.HasPrefix(r.URL.Path, "/v1/operations/") {
			_, _ = w.Write([]byte(`{"name": "` + strings.TrimPrefix(r.URL.Path, "/v1/") + `", "done": true}`))
			return
		}

		if strings.HasSuffix(r.URL.Path, ":listAttestationRules") {
			_, _ = w.Write([]byte(`{}`))
			return
		}

		requests = append(requests, recordedRequest{
			Method:     r.Method,
			Path:       r.URL.Path,
			UpdateMask: r.URL.Query().Get("updateMask"),
			Body:       body,
		})

		switch {
		case r.Method == "GET" && !undeleted:
			_, _ = w.Write([]byte(`{
				"name": "projects/test-project/locations/global/workloadIdentityPools/my-pool",
				"state": "DELETED",
				"displayName": "Old Name",
				"description": "Old Description",
				"disabled": true,
				"mode": "FEDERATION_ONLY"
			}`))
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, ":undelete"):
			undeleted = true
			_, _ = w.Write([]byte(`{"name": "operations/undelete-pool-op", "done": true}`))
		case r.Method == "PATCH":
			_, _ = w.Write([]byte(`{"name": "operations/patch-pool-op", "done": true}`))
		case r.Method == "GET" && undeleted:
			_, _ = w.Write([]byte(`{
				"name": "projects/test-project/locations/global/workloadIdentityPools/my-pool",
				"state": "ACTIVE",
				"displayName": "Recreated Name",
				"description": "Recreated Description",
				"disabled": false,
				"mode": "FEDERATION_ONLY"
			}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "unexpected request"}`))
		}
	}))
	defer ts.Close()

	res := iambeta.ResourceIAMBetaWorkloadIdentityPool()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]any{
		"project":                   "test-project",
		"workload_identity_pool_id": "my-pool",
		"display_name":              "Recreated Name",
		"description":               "Recreated Description",
		"disabled":                  false,
		"mode":                      "FEDERATION_ONLY",
	})

	config := &transport_tpg.Config{
		Client:       ts.Client(),
		Project:      "test-project",
		UserAgent:    "test-agent",
		PollInterval: 1 * time.Millisecond,
		CustomEndpoints: map[string]string{
			"iam_beta_custom_endpoint": ts.URL + "/v1/",
		},
	}

	if err := res.Create(d, config); err != nil {
		t.Fatalf("ResourceIAMBetaWorkloadIdentityPool().Create(%q) returned unexpected error: %v", "my-pool", err)
	}

	wantRequests := []recordedRequest{
		{
			Method: "GET",
			Path:   "/v1/projects/test-project/locations/global/workloadIdentityPools/my-pool",
		},
		{
			Method: "POST",
			Path:   "/v1/projects/test-project/locations/global/workloadIdentityPools/my-pool:undelete",
			Body:   map[string]any{},
		},
		{
			Method:     "PATCH",
			Path:       "/v1/projects/test-project/locations/global/workloadIdentityPools/my-pool",
			UpdateMask: "displayName,description,disabled",
			Body: map[string]any{
				"displayName": "Recreated Name",
				"description": "Recreated Description",
				"disabled":    false,
			},
		},
		{
			Method: "GET",
			Path:   "/v1/projects/test-project/locations/global/workloadIdentityPools/my-pool",
		},
	}

	if diff := cmp.Diff(wantRequests, requests); diff != "" {
		t.Errorf("ResourceIAMBetaWorkloadIdentityPool().Create(%q) HTTP requests mismatch (-want +got):\n%s", "my-pool", diff)
	}
}
