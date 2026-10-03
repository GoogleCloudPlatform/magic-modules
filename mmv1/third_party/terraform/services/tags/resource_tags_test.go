package tags_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/cloudrun"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/google/services/tags"
	"github.com/hashicorp/terraform-provider-google/google/services/tagslocation"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Tags tests cannot be run in parallel without running into Error Code 10: ABORTED
// See https://github.com/hashicorp/terraform-provider-google/issues/8637

func TestAccTags(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"tagKeyBasic":                              testAccTagsTagKey_tagKeyBasic,
		"tagKeyBasicWithPurposeGceFirewall":        testAccTagsTagKey_tagKeyBasicWithPurposeGceFirewall,
		"tagKeyBasicWithPurposeDataGovernance":     testAccTagsTagKey_tagKeyBasicWithPurposeDataGovernance,
		"tagKeyBasicWithAllowedValuesRegex":        testAccTagsTagKey_tagKeyBasicWithAllowedValuesRegex,
		"tagKeyUpdate":                             testAccTagsTagKey_tagKeyUpdate,
		"tagKeyUpdateAllowedValuesRegex":           testAccTagsTagKey_tagKeyUpdateAllowedValuesRegex,
		"tagKeyIamBinding":                         testAccTagsTagKeyIamBinding,
		"tagKeyIamMember":                          testAccTagsTagKeyIamMember,
		"tagKeyIamPolicy":                          testAccTagsTagKeyIamPolicy,
		"tagValueBasic":                            testAccTagsTagValue_tagValueBasic,
		"tagValueUpdate":                           testAccTagsTagValue_tagValueUpdate,
		"tagBindingBasic":                          testAccTagsTagBinding_tagBindingBasic,
		"tagBindingBasicDynamic":                   testAccTagsTagBinding_tagBindingBasicDynamic,
		"tagBindingNamespaced":                     testAccTagsTagBinding_tagBindingNamespaced,
		"tagValueIamBinding":                       testAccTagsTagValueIamBinding,
		"tagValueIamMember":                        testAccTagsTagValueIamMember,
		"tagValueIamPolicy":                        testAccTagsTagValueIamPolicy,
		"tagsLocationTagBindingBasic":              testAccTagsLocationTagBinding_locationTagBindingbasic,
		"tagsLocationTagBindingBasicDynamic":       testAccTagsLocationTagBinding_locationTagBindingBasicDynamic,
		"tagsLocationTagBindingBasicWithProjectId": testAccTagsLocationTagBinding_locationTagBindingBasicWithProjectId,
		"tagsLocationTagBindingZonal":              testAccTagsLocationTagBinding_locationTagBindingzonal,
		"tagsLocationTagBindingZonalDynamic":       testAccTagsLocationTagBinding_locationTagBindingZonalDynamic,
		"tagsLocationTagBindingZonalNamespaced":    testAccTagsLocationTagBinding_locationTagBindingZonalNamespaced,
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

func testAccTagsTagKey_tagKeyBasic(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKey_tagKeyBasicExample(context),
			},
		},
	})
}

func testAccTagsTagKey_tagKeyBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "tf-test-foo%{random_suffix}"
  description = "For foo%{random_suffix} resources."
}
`, context)
}

func testAccTagsTagKey_tagKeyBasicWithPurposeGceFirewall(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKey_tagKeyBasicWithPurposeGceFirewallExample(context),
			},
			// Regression test for https://github.com/hashicorp/terraform-provider-google/issues/20073:
			// purpose_data was ignore_read, so importing left it empty and the (immutable) field
			// forced a replacement on the next plan.
			{
				ResourceName:            "google_tags_tag_key.key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccTagsTagKey_tagKeyBasicWithPurposeGceFirewallExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "tag_network" {
	name = "tf-test-vpc-%{random_suffix}"
	auto_create_subnetworks = false
}

resource "google_tags_tag_key" "key" {
	  parent = "organizations/%{org_id}"
	  short_name = "tf-test-foo%{random_suffix}"
	  description = "For foo%{random_suffix} resources."
	  purpose = "GCE_FIREWALL"
	  # purpose_data.network must be a Compute network self link containing the numeric network id,
	  # which is the form the API returns. google_compute_network.self_link is
	  # "{version}/projects/{project}/global/networks/{network_name}", so the id-bearing self link is
	  # built from its parts.
	  purpose_data = {network = "https://www.googleapis.com/compute/v1/projects/${google_compute_network.tag_network.project}/global/networks/${google_compute_network.tag_network.network_id}"}
	}

`, context)
}

func testAccTagsTagKey_tagKeyBasicWithPurposeDataGovernance(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKey_tagKeyBasicWithPurposeDataGovernanceExample(context),
			},
		},
	})
}

func testAccTagsTagKey_tagKeyBasicWithPurposeDataGovernanceExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "tf-test-data-gov-%{random_suffix}"
	description = "For data governance purposes."
	purpose = "DATA_GOVERNANCE"
}
`, context)
}

func testAccTagsTagKey_tagKeyBasicWithAllowedValuesRegex(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKey_tagKeyBasicWithAllowedValuesRegexExample(context),
			},
		},
	})
}

func testAccTagsTagKey_tagKeyBasicWithAllowedValuesRegexExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "tf-test-foo%{random_suffix}"
  description = "For foo%{random_suffix} resources."
  allowed_values_regex = "^[a-z]+$"
}
`, context)
}

func testAccTagsTagKey_tagKeyUpdate(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKey_basic(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTagsTagKey_basicUpdated(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsTagKey_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "tf-test-foo%{random_suffix}"
  description = "For foo%{random_suffix} resources."
}
`, context)
}

func testAccTagsTagKey_basicUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "tf-test-foo%{random_suffix}"
  description = "Anything related to foo%{random_suffix}"
}
`, context)
}

func testAccTagsTagKey_tagKeyUpdateAllowedValuesRegex(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKey_basicWithAllowedValuesRegex(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTagsTagKey_basicWithAllowedValuesRegexUpdated(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsTagKey_basicWithAllowedValuesRegex(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "tf-test-foo%{random_suffix}"
  description = "For foo%{random_suffix} resources."
  allowed_values_regex = "^[a-z]+$"
}
`, context)
}

func testAccTagsTagKey_basicWithAllowedValuesRegexUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "tf-test-foo%{random_suffix}"
  description = "For foo%{random_suffix} resources."
  allowed_values_regex = ".*"
}
`, context)
}

func testAccCheckTagsTagKeyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_tags_tag_key" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(tags.Product, config)+"tagKeys/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("TagsTagKey still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccTagsTagValue_tagValueBasic(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagValueDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagValue_tagValueBasicExample(context),
			},
		},
	})
}

func testAccTagsTagValue_tagValueBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "tf-test-foobarbaz%{random_suffix}"
  description = "For foo/bar/baz resources."
}

resource "google_tags_tag_value" "value" {

  parent      = google_tags_tag_key.key.id
  short_name  = "tf-test-foo%{random_suffix}"
  description = "For foo resources."
}
`, context)
}

func testAccTagsTagValue_tagValueUpdate(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagValueDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagValue_basic(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTagsTagValue_basicUpdated(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsTagValue_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "tf-test-foobarbaz%{random_suffix}"
  description = "For foo/bar/baz resources."
}

resource "google_tags_tag_value" "value" {

  parent      = google_tags_tag_key.key.id
  short_name  = "tf-test-foo%{random_suffix}"
  description = "For foo resources."
}
`, context)
}

func testAccTagsTagValue_basicUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "tf-test-foobarbaz%{random_suffix}"
  description = "For foo/bar/baz resources."
}

resource "google_tags_tag_value" "value" {

  parent      = google_tags_tag_key.key.id
  short_name  = "tf-test-foo%{random_suffix}"
  description = "For any foo resources."
}
`, context)
}

func testAccCheckTagsTagValueDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_tags_tag_key" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(tags.Product, config)+"tagValues/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("TagsTagValue still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccTagsTagBinding_tagBindingBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_id":    "tf-test-" + acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagBinding_tagBindingBasicExample(context),
			},
		},
	})
}

func testAccTagsTagBinding_tagBindingBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
	project_id = "%{project_id}"
	name       = "%{project_id}"
	org_id     = "%{org_id}"
	deletion_policy = "DELETE"
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "tf-test-keyname%{random_suffix}"
	description = "For a certain set of resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "tf-test-foo%{random_suffix}"
	description = "For foo%{random_suffix} resources."
}

resource "google_tags_tag_binding" "binding" {
	parent    = "//cloudresourcemanager.googleapis.com/projects/${google_project.project.number}"
	tag_value = google_tags_tag_value.value.id
}
`, context)
}

func testAccTagsTagBinding_tagBindingBasicDynamic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_id":    "tf-test-" + acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagBinding_tagBindingBasicDynamicExample(context),
			},
			{
				ResourceName: "google_tags_tag_binding.binding",
				ImportState:  true,
			},
		},
	})
}

func testAccTagsTagBinding_tagBindingBasicDynamicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
	project_id = "%{project_id}"
	name       = "%{project_id}"
	org_id     = "%{org_id}"
	deletion_policy = "DELETE"
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "tf-test-keyname%{random_suffix}"
	description = "For a certain set of resources."
	allowed_values_regex = "test-.*"
}

resource "google_tags_tag_binding" "binding" {
	parent    = "//cloudresourcemanager.googleapis.com/projects/${google_project.project.number}"
	tag_value = "${google_tags_tag_key.key.namespaced_name}/test-value"
}
`, context)
}

func testAccTagsTagBinding_tagBindingNamespaced(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_id":    "tf-test-" + acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagBinding_tagBindingNamespacedExample(context),
			},
			{
				ResourceName: "google_tags_tag_binding.binding",
				ImportState:  true,
			},
		},
	})
}

// Generates Terraform configuration for testing tag binding with namespaced value.
func testAccTagsTagBinding_tagBindingNamespacedExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
	project_id = "%{project_id}"
	name       = "%{project_id}"
	org_id     = "%{org_id}"
	deletion_policy = "DELETE"
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "tf-test-key-%{random_suffix}"
	description = "Key for namespaced test."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "tf-test-val-%{random_suffix}"
	description = "Value for namespaced test."
}

resource "google_tags_tag_binding" "binding" {
	parent    = "//cloudresourcemanager.googleapis.com/projects/${google_project.project.number}"
	tag_value = google_tags_tag_value.value.namespaced_name
}
`, context)
}

func testAccCheckTagsTagBindingDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_tags_tag_binding" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(tags.Product, config)+"tagBindings/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("TagsTagBinding still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccTagsTagKeyIamBinding(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/resourcemanager.tagAdmin",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"short_name": "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKeyIamBinding_basicGenerated(context),
			},
			{
				Config: testAccTagsTagKeyIamBinding_withCondition(context),
			},
			{
				// Test Iam Binding update
				Config: testAccTagsTagKeyIamBinding_updateGenerated(context),
			},
		},
	})
}

func testAccTagsTagKeyIamMember(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/resourcemanager.tagAdmin",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"short_name": "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccTagsTagKeyIamMember_basicGenerated(context),
			},
			{
				Config: testAccTagsTagKeyIamMember_withCondition(context),
			},
		},
	})
}

func testAccTagsTagKeyIamPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/resourcemanager.tagAdmin",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"short_name": "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKeyIamPolicy_basicGenerated(context),
			},
			{
				Config: testAccTagsTagKeyIamPolicy_emptyBinding(context),
			},
			{
				Config: testAccTagsTagKeyIamPolicy_withCondition(context),
			},
		},
	})
}

func testAccTagsTagKeyIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

resource "google_tags_tag_key_iam_member" "foo" {
  tag_key = google_tags_tag_key.key.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccTagsTagKeyIamMember_withCondition(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

resource "google_tags_tag_key_iam_member" "foo" {
  tag_key = google_tags_tag_key.key.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
  condition {
    description = "Allow tagUser grant."
    expression  = "api.getAttribute('iam.googleapis.com/modifiedGrantsByRole', []).hasOnly([\"roles/resourcemanager.tagUser\"])"
    title       = "only_taguser_delegation"
  }
}
`, context)
}

func testAccTagsTagKeyIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_tags_tag_key_iam_policy" "foo" {
  tag_key = google_tags_tag_key.key.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagKeyIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

data "google_iam_policy" "foo" {
}

resource "google_tags_tag_key_iam_policy" "foo" {
  tag_key = google_tags_tag_key.key.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagKeyIamPolicy_withCondition(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
    condition {
      description = "Allow tagUser grant."
      expression  = "api.getAttribute('iam.googleapis.com/modifiedGrantsByRole', []).hasOnly([\"roles/resourcemanager.tagUser\"])"
      title       = "only_taguser_delegation"
    }
  }
}

resource "google_tags_tag_key_iam_policy" "foo" {
  tag_key = google_tags_tag_key.key.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagKeyIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

resource "google_tags_tag_key_iam_binding" "foo" {
  tag_key = google_tags_tag_key.key.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccTagsTagKeyIamBinding_withCondition(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

resource "google_tags_tag_key_iam_binding" "foo" {
  tag_key = google_tags_tag_key.key.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
  condition {
    description = "Allow tagUser grant."
    expression  = "api.getAttribute('iam.googleapis.com/modifiedGrantsByRole', []).hasOnly([\"roles/resourcemanager.tagUser\"])"
    title       = "only_taguser_delegation"
  }
}
`, context)
}

func testAccTagsTagKeyIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

resource "google_tags_tag_key_iam_binding" "foo" {
  tag_key = google_tags_tag_key.key.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}

func testAccTagsTagValueIamBinding(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/resourcemanager.tagAdmin",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"key_short_name":   "tf-test-key-" + acctest.RandString(t, 10),
		"value_short_name": "tf-test-value-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagValueIamBinding_basicGenerated(context),
			},
			{
				Config: testAccTagsTagValueIamBinding_withCondition(context),
			},
			{
				// Test Iam Binding update
				Config: testAccTagsTagValueIamBinding_updateGenerated(context),
			},
		},
	})
}

func testAccTagsTagValueIamMember(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/resourcemanager.tagAdmin",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"key_short_name":   "tf-test-key-" + acctest.RandString(t, 10),
		"value_short_name": "tf-test-value-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccTagsTagValueIamMember_basicGenerated(context),
			},
			{
				Config: testAccTagsTagValueIamMember_withCondition(context),
			},
		},
	})
}

func testAccTagsTagValueIamPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/resourcemanager.tagAdmin",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"key_short_name":   "tf-test-key-" + acctest.RandString(t, 10),
		"value_short_name": "tf-test-value-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagValueIamPolicy_basicGenerated(context),
			},
			{
				Config: testAccTagsTagValueIamPolicy_emptyBinding(context),
			},
			{
				Config: testAccTagsTagValueIamPolicy_withCondition(context),
			},
		},
	})
}

func testAccTagsTagValueIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

resource "google_tags_tag_value_iam_member" "foo" {
  tag_value = google_tags_tag_value.value.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccTagsTagValueIamMember_withCondition(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

resource "google_tags_tag_value_iam_member" "foo" {
  tag_value = google_tags_tag_value.value.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
  condition {
    description = "Allow tagUser grant."
    expression  = "api.getAttribute('iam.googleapis.com/modifiedGrantsByRole', []).hasOnly([\"roles/resourcemanager.tagUser\"])"
    title       = "only_taguser_delegation"
  }
}
`, context)
}

func testAccTagsTagValueIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_tags_tag_value_iam_policy" "foo" {
  tag_value = google_tags_tag_value.value.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagValueIamPolicy_withCondition(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
	condition {
        description = "Allow tagUser grant."
        expression  = "api.getAttribute('iam.googleapis.com/modifiedGrantsByRole', []).hasOnly([\"roles/resourcemanager.tagUser\"])"
        title       = "only_taguser_delegation"
    }
  }
}

resource "google_tags_tag_value_iam_policy" "foo" {
  tag_value = google_tags_tag_value.value.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagValueIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

data "google_iam_policy" "foo" {
}

resource "google_tags_tag_value_iam_policy" "foo" {
  tag_value = google_tags_tag_value.value.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagValueIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

resource "google_tags_tag_value_iam_binding" "foo" {
  tag_value = google_tags_tag_value.value.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccTagsTagValueIamBinding_withCondition(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

resource "google_tags_tag_value_iam_binding" "foo" {
  tag_value = google_tags_tag_value.value.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
  condition {
      description = "Allow tagUser grant."
      expression  = "api.getAttribute('iam.googleapis.com/modifiedGrantsByRole', []).hasOnly([\"roles/resourcemanager.tagUser\"])"
      title       = "only_taguser_delegation"
  }
}
`, context)
}

func testAccTagsTagValueIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

resource "google_tags_tag_value_iam_binding" "foo" {
  tag_value = google_tags_tag_value.value.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}

func testAccTagsLocationTagBinding_locationTagBindingbasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		// "org_id":        envvar.GetTestOrgFromEnv(t),
		// "project_id":    "tf-test-" + acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsLocationTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsLocationTagBinding_locationTagBindingBasicExample(context),
			},
			{
				ResourceName:      "google_tags_location_tag_binding.binding",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsLocationTagBinding_locationTagBindingBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_compute_instance" "vm" {
  name         = "tf-test-tagbinding-repro%{random_suffix}"
  machine_type = "e2-small"
  zone         = "us-east4-a"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
    }
  }
  network_interface {
    network = "default"
  }
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/${data.google_project.project.org_id}"
	short_name = "tf-test-keyname%{random_suffix}"
	description = "For a certain set of resources."

	# Setting purpose of GCE_FIREWALL exercises creation LRO logic
	purpose = "GCE_FIREWALL"
	purpose_data = {
		organization = "auto"
	}
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "tf-test-foo%{random_suffix}"
	description = "For foo%{random_suffix} resources."
}
  
resource "google_tags_location_tag_binding" "binding" {
	parent    = "//compute.googleapis.com/projects/${data.google_project.project.number}/zones/us-east4-a/instances/${google_compute_instance.vm.instance_id}"
	tag_value = google_tags_tag_value.value.id
	location  = "us-east4-a"
}
`, context)
}

func testAccTagsLocationTagBinding_locationTagBindingBasicDynamic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsLocationTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsLocationTagBinding_locationTagBindingBasicDynamicExample(context),
			},
			{
				ResourceName:      "google_tags_location_tag_binding.binding",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsLocationTagBinding_locationTagBindingBasicDynamicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/${data.google_project.project.org_id}"
	short_name = "tf-test-keyname%{random_suffix}"
	description = "For a certain set of resources."
	allowed_values_regex = "test-.*"
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-cloudrun-srv%{random_suffix}"
	location = "us-central1"
  
	template {
	  spec {
		containers {
		  image = "us-docker.pkg.dev/cloudrun/container/hello"
		}
	  }
	}
  
	traffic {
	  percent         = 100
	  latest_revision = true
	}
}
  
resource "google_tags_location_tag_binding" "binding" {
	parent    = "//run.googleapis.com/projects/${data.google_project.project.number}/locations/${google_cloud_run_service.default.location}/services/${google_cloud_run_service.default.name}"
	tag_value = "${google_tags_tag_key.key.namespaced_name}/test-value"
	location  = "us-central1"
}
`, context)
}

func testAccTagsLocationTagBinding_locationTagBindingBasicWithProjectId(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsLocationTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsLocationTagBinding_locationTagBindingBasicExampleWithProjectId(context),
			},
			{
				ResourceName:      "google_tags_location_tag_binding.binding",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsLocationTagBinding_locationTagBindingBasicExampleWithProjectId(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/${data.google_project.project.org_id}"
	short_name = "tf-test-keyname%{random_suffix}"
	description = "For a certain set of resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "tf-test-foo%{random_suffix}"
	description = "For foo%{random_suffix} resources."
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-cloudrun-srv%{random_suffix}"
	location = "us-central1"
  
	template {
	  spec {
		containers {
		  image = "us-docker.pkg.dev/cloudrun/container/hello"
		}
	  }
	}
  
	traffic {
	  percent         = 100
	  latest_revision = true
	}
}
  
resource "google_tags_location_tag_binding" "binding" {
	parent    = "//run.googleapis.com/projects/${data.google_project.project.project_id}/locations/${google_cloud_run_service.default.location}/services/${google_cloud_run_service.default.name}"
	tag_value = google_tags_tag_value.value.id
	location  = "us-central1"
}
`, context)
}

func testAccTagsLocationTagBinding_locationTagBindingzonal(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsLocationTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsLocationTagBinding_locationTagBindingZonalExample(context),
			},
			{
				ResourceName:      "google_tags_location_tag_binding.binding",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsLocationTagBinding_locationTagBindingZonalExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "tf-test-keyname%{random_suffix}"
	description = "For a certain set of resources."
}
resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "tf-test-foo%{random_suffix}"
	description = "For foo%{random_suffix} resources."
}
resource "google_compute_instance" "default" {
	name         = "tf-test-%{random_suffix}"
	machine_type = "e2-medium"
	zone         = "us-central1-a"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-13"
    }
  }
  network_interface {
    network = "default"
  }
}
resource "google_tags_location_tag_binding" "binding" {
	parent    = "//compute.googleapis.com/projects/${data.google_project.project.number}/zones/us-central1-a/instances/${google_compute_instance.default.instance_id}"
	tag_value = google_tags_tag_value.value.id
	location  = "us-central1-a"
}
`, context)
}

func testAccTagsLocationTagBinding_locationTagBindingZonalDynamic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsLocationTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsLocationTagBinding_locationTagBindingZonalDynamicExample(context),
			},
			{
				ResourceName:      "google_tags_location_tag_binding.binding",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsLocationTagBinding_locationTagBindingZonalDynamicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "tf-test-keyname%{random_suffix}"
	description = "For a certain set of resources."
	allowed_values_regex = "test-.*"
}

resource "google_compute_instance" "default" {
	name         = "tf-test-%{random_suffix}"
	machine_type = "e2-medium"
	zone         = "us-central1-a"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-13"
    }
  }
  network_interface {
    network = "default"
  }
}

resource "google_tags_location_tag_binding" "binding" {
	parent    = "//compute.googleapis.com/projects/${data.google_project.project.number}/zones/us-central1-a/instances/${google_compute_instance.default.instance_id}"
	tag_value = "${google_tags_tag_key.key.namespaced_name}/test-value"
	location  = "us-central1-a"
}
`, context)
}

func testAccTagsLocationTagBinding_locationTagBindingZonalNamespaced(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsLocationTagBinding_locationTagBindingZonalNamespacedExample(context),
			},
			{
				ResourceName: "google_tags_location_tag_binding.binding",
				ImportState:  true,
			},
		},
	})
}

// Generates Terraform configuration for testing location tag binding with namespaced value.
func testAccTagsLocationTagBinding_locationTagBindingZonalNamespacedExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "tf-test-keyname%{random_suffix}"
	description = "For a certain set of resources."
}

resource "google_tags_tag_value" "value" {
	parent      = google_tags_tag_key.key.id
	short_name  = "tf-test-foo%{random_suffix}"
	description = "For foo%{random_suffix} resources."
}

resource "google_compute_instance" "default" {
	name         = "tf-test-%{random_suffix}"
	machine_type = "e2-medium"
	zone         = "us-central1-a"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-13"
    }
  }
  network_interface {
    network = "default"
  }
}

resource "google_tags_location_tag_binding" "binding" {
	parent    = "//compute.googleapis.com/projects/${data.google_project.project.number}/zones/us-central1-a/instances/${google_compute_instance.default.instance_id}"
	tag_value = google_tags_tag_value.value.namespaced_name
	location  = "us-central1-a"
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

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(tagslocation.Product, config)+"{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("TagsTagBinding still exists at %s", url)
			}
		}
		return nil
	}
}

func TestTagsTagKeyPurposeDataValidate(t *testing.T) {
	cases := map[string]struct {
		purposeData map[string]interface{}
		expectErr   bool
	}{
		"self link with a numeric id is valid": {
			purposeData: map[string]interface{}{
				"network": "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/123456789",
			},
		},
		"a custom compute endpoint host is valid": {
			purposeData: map[string]interface{}{
				"network": "https://compute.mirror.example.com/compute/v1/projects/my-project/global/networks/123456789",
			},
		},
		"no network key is valid": {
			purposeData: map[string]interface{}{"organization": "auto"},
		},
		"empty network is valid": {
			purposeData: map[string]interface{}{"network": ""},
		},
		"organization auto is valid": {
			purposeData: map[string]interface{}{"organization": "auto"},
		},
		"empty organization is valid": {
			purposeData: map[string]interface{}{"organization": ""},
		},
		"short form is rejected": {
			purposeData: map[string]interface{}{"network": "my-project/vpc-us-west1"},
			expectErr:   true,
		},
		"organization other than auto is rejected": {
			purposeData: map[string]interface{}{"organization": "123456789012"},
			expectErr:   true,
		},
		"self link with a network name is rejected": {
			purposeData: map[string]interface{}{
				"network": "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/vpc-us-west1",
			},
			expectErr: true,
		},
		"missing numeric id is rejected": {
			purposeData: map[string]interface{}{
				"network": "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/",
			},
			expectErr: true,
		},
		"regional network path is rejected": {
			purposeData: map[string]interface{}{
				"network": "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-west1/networks/123456789",
			},
			expectErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, errs := tags.TagsTagKeyPurposeDataValidate(tc.purposeData, "purpose_data")
			if tc.expectErr && len(errs) == 0 {
				t.Fatalf("expected an error for %#v, got none", tc.purposeData)
			}
			if !tc.expectErr && len(errs) > 0 {
				t.Fatalf("expected no error for %#v, got %v", tc.purposeData, errs)
			}
		})
	}
}

func TestTagsTagKeyPurposeDataDiffSuppress(t *testing.T) {
	const (
		orgID    = "123456789012"
		selfLink = "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/123456789"
	)

	cases := map[string]struct {
		k        string
		oldValue string
		newValue string
		want     bool
	}{
		// The API replaces "auto" (the config value) with the organization id it stores, so the
		// config side is what has to be suppressed.
		"config auto against the resolved organization id is suppressed": {
			k: "purpose_data.organization", oldValue: orgID, newValue: "auto", want: true,
		},
		// State never holds "auto", so this direction is not expected; suppress nothing.
		"state auto against a resolved id is not suppressed": {
			k: "purpose_data.organization", oldValue: "auto", newValue: orgID, want: false,
		},
		// network is validated to a single form, so it is compared normally.
		"identical network values are not suppressed here": {
			k: "purpose_data.network", oldValue: selfLink, newValue: selfLink, want: false,
		},
		// Any other key falls through.
		"other key is not suppressed": {
			k: "purpose_data.something_else", oldValue: "auto", newValue: orgID, want: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := tags.TagsTagKeyPurposeDataDiffSuppress(tc.k, tc.oldValue, tc.newValue, nil)
			if got != tc.want {
				t.Errorf("want suppress=%t, got %t for k=%q old=%q new=%q", tc.want, got, tc.k, tc.oldValue, tc.newValue)
			}
		})
	}
}
