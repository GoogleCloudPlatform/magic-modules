package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTagsLocation(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"tagsLocationTagBindingBasic": testAccTagsLocationTagBinding_locationTagBindingbasic,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccTagsLocationTagBinding_locationTagBindingbasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"project_id":    "tf-test-" + randString(t, 10),
		"random_suffix": randString(t, 10),

		"key_short_name":   "tf-test-key-" + randString(t, 10),
		"value_short_name": "tf-test-value-" + randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsLocationTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsLocationTagBinding_locationTagBindingBasicExample(context),
			},
		},
	})
}

func testAccTagsLocationTagBinding_locationTagBindingBasicExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "project" {
	project_id = "%{project_id}"
	name       = "%{project_id}"
	org_id     = "%{org_id}"
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "keyname%{random_suffix}"
	description = "For a certain set of resources."
}

resource "google_tags_tag_value" "value" {
	parent = google_tags_tag_key.key.id
	short_name = "foo%{random_suffix}"
	description = "For foo%{random_suffix} resources."
}

resource "google_sql_database_instance" "main" {
	name             = "tf-test-main-instance"
	database_version = "POSTGRES_14"
	region           = "us-central1"
	deletion_protection = false
	settings {
	  # Second-generation instance tiers are based on the machine
	  # type. See argument reference below.
	  tier = "db-f1-micro"
	}
}
resource "google_tags_location_tag_binding" "binding" {
	parent = "//sqladmin.googleapis.com/projects/${google_project.project.number}/instances/${google_sql_database_instance.main.id}"
	tag_value = google_tags_tag_value.value.id
	location = "us-central1"
}
`, context)
}

func testAccCheckTagsLocationTagBindingDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_tags_location_tag_binding" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{TagsLocationBasePath}}{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = sendRequest(config, "GET", billingProject, url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("TagsLocationTagBinding still exists at %s", url)
			}
		}

		return nil
	}
}
