package iap_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccIapBrand_Datasource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIapBrandDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// Check that brand is set automatically by the data source
					resource.TestCheckResourceAttrSet("data.google_iap_brand.project", "brand"),
					// Check that the automatically discovered brand ID matches the project's actual number
					resource.TestCheckResourceAttrPair("data.google_iap_brand.project", "brand", "google_project.project", "number"),
					// Check that the rest of the data source state matches the resource state
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_iap_brand.project",
						"google_iap_brand.project",
						map[string]struct{}{
							"project": {},
							"brand":   {},
						},
					),
				),
			},
		},
	})
}

func TestAccIapBrand_Datasource_WithBrand(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIapBrandDatasourceConfigWithBrand(context),
				Check: resource.ComposeTestCheckFunc(
					// Check that the explicitly provided brand ID is correctly used and matches
					resource.TestCheckResourceAttrPair("data.google_iap_brand.project", "brand", "google_project.project", "number"),
					// Check another attribute to ensure the read was successful
					resource.TestCheckResourceAttr("data.google_iap_brand.project", "application_title", "Cloud IAP protected Application"),
				),
			},
		},
	})
}

func testAccIapBrandDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id = "tf-test%{random_suffix}"
  name       = "tf-test%{random_suffix}"
  org_id     = "%{org_id}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "project_service" {
  project = google_project.project.project_id
  service = "iap.googleapis.com"
}
	  
resource "google_iap_brand" "project" {
  support_email     = "support@%{org_domain}"
  application_title = "Cloud IAP protected Application"
  project           = google_project_service.project_service.project
}

data "google_iap_brand" "project" {
  project = google_iap_brand.project.project
}
`, context)
}

func testAccIapBrandDatasourceConfigWithBrand(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id = "tf-test%{random_suffix}"
  name       = "tf-test%{random_suffix}"
  org_id     = "%{org_id}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "project_service" {
  project = google_project.project.project_id
  service = "iap.googleapis.com"
}

resource "google_iap_brand" "project" {
  support_email     = "support@%{org_domain}"
  application_title = "Cloud IAP protected Application"
  project           = google_project_service.project_service.project
}

data "google_iap_brand" "project" {
  project = google_iap_brand.project.project
  brand   = google_project.project.number
}
`, context)
}
