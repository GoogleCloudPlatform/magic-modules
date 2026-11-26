package compute_test

import (
	"fmt"
	"net/url"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Invokes the action from an action_trigger on create and again on update, and checks that one
// invalidateCache operation is recorded per invocation. The no-op plan check in step 2 shows the action
// does not drift the URL map.
func TestAccGoogleComputeUrlMapInvalidateCacheAction_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: testAccCheckComputeUrlMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleComputeUrlMapInvalidateCacheAction_basic(context, "first"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeUrlMapInvalidateCacheCount(t, "google_compute_url_map.test", 1),
				),
			},
			{
				Config: testAccGoogleComputeUrlMapInvalidateCacheAction_basic(context, "second"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_url_map.test", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("terraform_data.trigger", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeUrlMapInvalidateCacheCount(t, "google_compute_url_map.test", 2),
				),
			},
		},
	})
}

// Path rules are enforced during validation, before any API call, so this test records no requests.
func TestAccGoogleComputeUrlMapInvalidateCacheAction_invalidPath(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				action "google_compute_url_map_invalidate_cache" "test" {
					config {
						url_map = "projects/p/global/urlMaps/m"
						path = "/images/*.png"
					}
				}
				`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid path`),
			},
		},
	})
}

// testAccCheckComputeUrlMapInvalidateCacheCount counts completed invalidateCache operations that target the
// URL map. An invalidation does not change the URL map, so the operation log is the only observable evidence
// that the action ran.
func testAccCheckComputeUrlMapInvalidateCacheCount(t *testing.T, n string, want int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found in the state", n)
		}
		project := rs.Primary.Attributes["project"]
		selfLink := rs.Primary.Attributes["self_link"]

		config := acctest.GoogleProviderConfig(t)
		baseURL := fmt.Sprintf("%sprojects/%s/global/operations", transport_tpg.BaseUrl(tpgcompute.Product, config), project)
		query := url.Values{}
		query.Set("filter", `operationType = "invalidateCache"`)

		got := 0
		for {
			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   project,
				RawURL:    baseURL + "?" + query.Encode(),
				UserAgent: config.UserAgent,
			})
			if err != nil {
				return fmt.Errorf("error listing global operations: %s", err)
			}

			items, _ := res["items"].([]interface{})
			for _, raw := range items {
				op, ok := raw.(map[string]interface{})
				if !ok {
					continue
				}
				targetLink, _ := op["targetLink"].(string)
				if tpgresource.CompareSelfLinkOrResourceName("", targetLink, selfLink, nil) && op["status"] == "DONE" && op["error"] == nil {
					got++
				}
			}

			token, _ := res["nextPageToken"].(string)
			if token == "" {
				break
			}
			query.Set("pageToken", token)

		}
		if got != want {
			return fmt.Errorf("expected %d invalidateCache operation(s) on URL map %s, found %d", want, rs.Primary.ID, got)
		}
		return nil
	}
}

func testAccGoogleComputeUrlMapInvalidateCacheAction_basic(context map[string]interface{}, trigger string) string {
	context["trigger"] = trigger
	return acctest.Nprintf(`
resource "google_storage_bucket" "test" {
  name     = "tf-test-cdn-%{random_suffix}"
  location = "US"
  uniform_bucket_level_access = true
  force_destroy = true
}

resource "google_compute_backend_bucket" "test" {
  name        = "tf-test-cdn-%{random_suffix}"
  bucket_name = google_storage_bucket.test.name
  enable_cdn = true
}

resource "google_compute_url_map" "test" {
  name = "tf-test-cdn-%{random_suffix}"
  default_service = google_compute_backend_bucket.test.id
}

action "google_compute_url_map_invalidate_cache" "test" {
	config {
		url_map = google_compute_url_map.test.id
		path = "/*"
		timeout = "10m"
	}
}

resource "terraform_data" "trigger" {
  input = "${google_compute_url_map.test.map_id}-%{trigger}"

  lifecycle {
    action_trigger {
		events = [after_create, after_update]
		actions = [action.google_compute_url_map_invalidate_cache.test]
	}
  }
}
`, context)
}
