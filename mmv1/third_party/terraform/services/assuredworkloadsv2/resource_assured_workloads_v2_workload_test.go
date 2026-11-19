package assuredworkloadsv2_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/assuredworkloadsv2"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanagerv3"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// testAccAssuredWorkloadsV2AllowedServices is the allowedServices parameter of
// the builtin-enforce-service-usage-restrictions cloud control in revisions 60
// and 62 of the builtin-us-data-boundary framework.
const testAccAssuredWorkloadsV2AllowedServices = `["accessapproval.googleapis.com", "accesscontextintelligence.googleapis.com", "accesscontextmanager.googleapis.com", "aiplatform.googleapis.com", "alloydb.googleapis.com", "anthosconfigmanagement.googleapis.com", "anthosidentityservice.googleapis.com", "apigee.googleapis.com", "apphub.googleapis.com", "artifactregistry.googleapis.com", "assuredworkloads.googleapis.com", "auditmanager.googleapis.com", "backupdr.googleapis.com", "bigquery.googleapis.com", "bigqueryconnection.googleapis.com", "bigquerydatapolicy.googleapis.com", "bigquerydatatransfer.googleapis.com", "bigquerymigration.googleapis.com", "bigqueryreservation.googleapis.com", "bigquerystorage.googleapis.com", "bigtable.googleapis.com", "bigtableadmin.googleapis.com", "binaryauthorization.googleapis.com", "certificatemanager.googleapis.com", "chronicle.googleapis.com", "chronicleservicemanager.googleapis.com", "cloudasset.googleapis.com", "cloudbuild.googleapis.com", "clouddeploy.googleapis.com", "cloudfunctions.googleapis.com", "cloudkms.googleapis.com", "cloudresourcemanager.googleapis.com", "cloudscheduler.googleapis.com", "cloudsecuritycompliance.googleapis.com", "cloudtasks.googleapis.com", "cloudtrace.googleapis.com", "composer.googleapis.com", "compute.googleapis.com", "config.googleapis.com", "connectgateway.googleapis.com", "connectors.googleapis.com", "contactcenterinsights.googleapis.com", "container.googleapis.com", "containerfilesystem.googleapis.com", "containersecurity.googleapis.com", "containerthreatdetection.googleapis.com", "databasecenter.googleapis.com", "dataflow.googleapis.com", "dataform.googleapis.com", "datafusion.googleapis.com", "datalineage.googleapis.com", "datapipelines.googleapis.com", "dataplex.googleapis.com", "dataproc-control.googleapis.com", "dataproc.googleapis.com", "datastore.googleapis.com", "datastudio.googleapis.com", "dialogflow.googleapis.com", "discoveryengine.googleapis.com", "dlp.googleapis.com", "dns.googleapis.com", "documentai.googleapis.com", "domains.googleapis.com", "essentialcontacts.googleapis.com", "eventarc.googleapis.com", "file.googleapis.com", "firebaserules.googleapis.com", "firestore.googleapis.com", "gkebackup.googleapis.com", "gkeconnect.googleapis.com", "gkehub.googleapis.com", "iam.googleapis.com", "iap.googleapis.com", "logging.googleapis.com", "looker.googleapis.com", "managedkafka.googleapis.com", "mesh.googleapis.com", "meshca.googleapis.com", "meshconfig.googleapis.com", "modelarmor.googleapis.com", "monitoring.googleapis.com", "netapp.googleapis.com", "networkconnectivity.googleapis.com", "networksecurity.googleapis.com", "networkservices.googleapis.com", "notebooks.googleapis.com", "opsconfigmonitoring.googleapis.com", "orgpolicy.googleapis.com", "osconfig.googleapis.com", "oslogin.googleapis.com", "policytroubleshooter.googleapis.com", "privateca.googleapis.com", "privilegedaccessmanager.googleapis.com", "pubsub.googleapis.com", "redis.googleapis.com", "resourcesettings.googleapis.com", "run.googleapis.com", "secretmanager.googleapis.com", "securesourcemanager.googleapis.com", "securitycenter.googleapis.com", "securitycentermanagement.googleapis.com", "securityposture.googleapis.com", "servicedirectory.googleapis.com", "servicehealth.googleapis.com", "servicemanagement.googleapis.com", "servicenetworking.googleapis.com", "spanner.googleapis.com", "speech.googleapis.com", "sqladmin.googleapis.com", "stackdriver.googleapis.com", "storage.googleapis.com", "storagetransfer.googleapis.com", "texttospeech.googleapis.com", "trafficdirector.googleapis.com", "vision.googleapis.com", "vpcaccess.googleapis.com", "webrisk.googleapis.com", "websecurityscanner.googleapis.com", "workflowexecutions.googleapis.com", "workflows.googleapis.com", "workstations.googleapis.com"]`

const testAccAssuredWorkloadsV2FrameworkParam = `
    parameters {
      name = "framework"
      parameter_value {
        string_value = "US_REGIONAL_DATA_BOUNDARY"
      }
    }`

const testAccAssuredWorkloadsV2LoggingLocationsParam = `
    parameters {
      name = "cloudLoggingStorageLocations"
      parameter_value {
        string_list_value {
          values = ["us-central1", "us-central2", "us-east4", "us-east1", "us-west1", "us-west4", "us-west3", "us-west2", "us-south1", "us-east5"]
        }
      }
    }`

const testAccAssuredWorkloadsV2ServiceUsageParam = `
    parameters {
      name = "serviceUsage"
      parameter_value {
        oneof_value {
          name = "allowedServices"
          parameter_value {
            string_list_value {
              values = ` + testAccAssuredWorkloadsV2AllowedServices + `
            }
          }
        }
      }
    }`

const testAccAssuredWorkloadsV2ResourceLocationsParam = `
    parameters {
      name = "locations"
      parameter_value {
        oneof_value {
          name = "allowedLocations"
          parameter_value {
            string_list_value {
              values = ["us-locations", "us-central1", "us-central2", "us-east4", "us-east1", "us-west1", "us-west4", "us-west3", "us-west2", "us-south1", "us-east5"]
            }
          }
        }
      }
    }`

const testAccAssuredWorkloadsV2TlsVersionsParam = `
    parameters {
      name = "deniedTlsVersions"
      parameter_value {
        string_list_value {
          values = ["TLS_VERSION_1", "TLS_VERSION_1_1"]
        }
      }
    }`

type testAccAssuredWorkloadsV2Control struct {
	name           string
	rev60          int
	rev62          int
	preventiveOnly bool
	params         string
}

var testAccAssuredWorkloadsV2Controls = []testAccAssuredWorkloadsV2Control{
	{name: "builtin-apigee-regulatory-settings", rev60: 2, rev62: 2},
	{name: "builtin-backup-dr-regulatory-settings", rev60: 5, rev62: 6, params: testAccAssuredWorkloadsV2FrameworkParam},
	{name: "builtin-big-query-regulatory-settings", rev62: 9, preventiveOnly: true, params: testAccAssuredWorkloadsV2FrameworkParam},
	{name: "builtin-cloud-databases-contextstore-regulatory-settings", rev62: 1, preventiveOnly: true},
	{name: "builtin-cloud-logging-regulatory-settings", rev60: 9, rev62: 10, params: testAccAssuredWorkloadsV2FrameworkParam},
	{name: "builtin-cloud-logging-storage-locations", rev60: 3, rev62: 5, params: testAccAssuredWorkloadsV2LoggingLocationsParam},
	{name: "builtin-cloud-monitoring-regulatory-settings", rev60: 2, rev62: 2},
	{name: "builtin-dataplex-regulatory-settings", rev60: 8, rev62: 10, params: testAccAssuredWorkloadsV2FrameworkParam},
	{name: "builtin-dataproc-regulatory-settings", rev62: 1, preventiveOnly: true},
	{name: "builtin-disable-global-cloud-armor", rev60: 4, rev62: 4},
	{name: "builtin-disable-pubsub-subscription-transforms", rev60: 6, rev62: 6},
	{name: "builtin-disable-pubsub-topic-transforms", rev60: 6, rev62: 6},
	{name: "builtin-disable-spanner-multi-region", rev60: 4, rev62: 4},
	{name: "builtin-enable-axt", rev62: 1, preventiveOnly: true},
	{name: "builtin-enable-spanner-advanced-service-controls", rev60: 1, rev62: 2},
	{name: "builtin-enforce-service-usage-restrictions", rev60: 5, rev62: 5, params: testAccAssuredWorkloadsV2ServiceUsageParam},
	{name: "builtin-firestore-disable-enhanced-text-search", rev60: 2, rev62: 3},
	{name: "builtin-gemini-enterprise-agent-platform-allow-all-compliant-apis", rev60: 2, rev62: 3},
	{name: "builtin-gemini-enterprise-regulatory-settings", rev62: 4, preventiveOnly: true, params: testAccAssuredWorkloadsV2FrameworkParam},
	{name: "builtin-google-agent-platform-regulatory-settings", rev62: 1, preventiveOnly: true},
	{name: "builtin-looker-regulatory-settings", rev60: 6, rev62: 7, params: testAccAssuredWorkloadsV2FrameworkParam},
	{name: "builtin-netapp-volumes-regulatory-settings", rev62: 1, preventiveOnly: true},
	{name: "builtin-restrict-resource-locations", rev60: 9, rev62: 9, params: testAccAssuredWorkloadsV2ResourceLocationsParam},
	{name: "builtin-restrict-technical-support-services", rev62: 12, preventiveOnly: true, params: testAccAssuredWorkloadsV2FrameworkParam},
	{name: "builtin-restrict-tls-versions", rev60: 1, rev62: 2, params: testAccAssuredWorkloadsV2TlsVersionsParam},
}

func testAccAssuredWorkloadsV2CloudControlConfigs(orgID string, frameworkRev int, preventive bool) string {
	mode := "DETECTIVE"
	if preventive {
		mode = "PREVENTIVE"
	}
	var b strings.Builder
	for _, c := range testAccAssuredWorkloadsV2Controls {
		if c.preventiveOnly && !preventive {
			continue
		}
		rev := c.rev62
		if frameworkRev == 60 {
			rev = c.rev60
		}
		fmt.Fprintf(&b, `
  cloud_control_configs {
    cloud_control          = "organizations/%s/locations/us-east4/cloudControls/%s"
    cloud_control_revision = %d
    enforcement_mode       = "%s"%s
  }
`, orgID, c.name, rev, mode, c.params)
	}
	return b.String()
}

func testAccAssuredWorkloadsV2Context(t *testing.T) map[string]interface{} {
	orgID := envvar.GetTestOrgFromEnv(t)
	return map[string]interface{}{
		"org_id":                 orgID,
		"billing_account":        envvar.GetTestBillingAccountFromEnv(t),
		"billing_project":        envvar.GetTestProjectFromEnv(),
		"random_suffix":          acctest.RandString(t, 10),
		"controls_rev60_monitor": testAccAssuredWorkloadsV2CloudControlConfigs(orgID, 60, false),
		"controls_rev62_monitor": testAccAssuredWorkloadsV2CloudControlConfigs(orgID, 62, false),
		"controls_rev62_prevent": testAccAssuredWorkloadsV2CloudControlConfigs(orgID, 62, true),
	}
}

// testAccCheckAssuredWorkloadsV2WorkloadDestroyProducer checks that every
// workload in the state returns 404 and deletes any folder or project
// provisioned via new_folder_config or new_project_config.
func testAccCheckAssuredWorkloadsV2WorkloadDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_assured_workloads_v2_workload" || strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(assuredworkloadsv2.Product, config)+"organizations/{{organization}}/locations/{{location}}/workloads/{{workload_id}}")
			if err != nil {
				return err
			}

			headers := make(http.Header)
			headers.Set("X-Goog-User-Project", envvar.GetTestProjectFromEnv())
			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
				Headers:   headers,
			})
			if err == nil {
				return fmt.Errorf("AssuredWorkloadsV2Workload still exists at %s", url)
			}
			if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
				return fmt.Errorf("error checking that AssuredWorkloadsV2Workload %s is destroyed: %s", url, err)
			}

			// Deleting a workload does not delete a folder or project created
			// via new_folder_config or new_project_config, so clean it up here.
			target := rs.Primary.Attributes["computed_target_resource"]
			if rs.Primary.Attributes["resource_config.0.new_folder_config.#"] == "1" && strings.HasPrefix(target, "folders/") {
				if _, err := resourcemanagerv3.NewClient(config, config.UserAgent).Folders.Delete(target).Do(); err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
					return fmt.Errorf("error deleting provisioned folder %s: %s", target, err)
				}
			} else if rs.Primary.Attributes["resource_config.0.new_project_config.#"] == "1" && strings.HasPrefix(target, "projects/") {
				if _, err := resourcemanagerv3.NewClient(config, config.UserAgent).Projects.Delete(target).Do(); err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
					return fmt.Errorf("error deleting provisioned project %s: %s", target, err)
				}
			}
		}
		return nil
	}
}

// Creates a workload in monitor mode on an existing project, upgrades the
// framework revision in place, moves the workload to monitor and prevent mode,
// and changes the Key Access Transparency choice.
func TestAccAssuredWorkloadsV2Workload_update(t *testing.T) {
	t.Parallel()

	context := testAccAssuredWorkloadsV2Context(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAssuredWorkloadsV2WorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAssuredWorkloadsV2Workload_existingProject(context, 60, "%{controls_rev60_monitor}", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_assured_workloads_v2_workload.workload", "framework.0.major_revision_id", "60"),
					resource.TestCheckResourceAttr("google_assured_workloads_v2_workload.workload", "state", "ACTIVE"),
					resource.TestCheckResourceAttrSet("google_assured_workloads_v2_workload.workload", "name"),
					resource.TestCheckResourceAttrSet("google_assured_workloads_v2_workload.workload", "computed_target_resource"),
					resource.TestCheckResourceAttrSet("google_assured_workloads_v2_workload.workload", "target_resource_display_name"),
					resource.TestCheckResourceAttrSet("google_assured_workloads_v2_workload.workload", "create_time"),
					resource.TestCheckResourceAttrSet("google_assured_workloads_v2_workload.workload", "update_time"),
					resource.TestCheckResourceAttrSet("google_assured_workloads_v2_workload.workload", "etag"),
				),
			},
			{
				ResourceName:            "google_assured_workloads_v2_workload.workload",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_config", "deletion_policy"},
			},
			{
				Config: testAccAssuredWorkloadsV2Workload_existingProject(context, 62, "%{controls_rev62_monitor}", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_assured_workloads_v2_workload.workload", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr("google_assured_workloads_v2_workload.workload", "framework.0.major_revision_id", "62"),
			},
			{
				ResourceName:            "google_assured_workloads_v2_workload.workload",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_config", "deletion_policy"},
			},
			{
				Config: testAccAssuredWorkloadsV2Workload_existingProject(context, 62, "%{controls_rev62_prevent}", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_assured_workloads_v2_workload.workload", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr("google_assured_workloads_v2_workload.workload", "kaj_enrollment_type", "KEY_ACCESS_TRANSPARENCY_ENABLED"),
			},
			{
				Config: testAccAssuredWorkloadsV2Workload_existingProject(context, 62, "%{controls_rev62_prevent}", `kaj_enrollment_type = "KEY_ACCESS_TRANSPARENCY_DISABLED"`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_assured_workloads_v2_workload.workload", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr("google_assured_workloads_v2_workload.workload", "kaj_enrollment_type", "KEY_ACCESS_TRANSPARENCY_DISABLED"),
			},
			{
				ResourceName:            "google_assured_workloads_v2_workload.workload",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_config", "deletion_policy"},
			},
		},
	})
}

// Covers existing_folder, new_folder_config (with latest framework revision),
// and new_project_config targets in a single test.
func TestAccAssuredWorkloadsV2Workload_resourceConfigs(t *testing.T) {
	t.Parallel()

	context := testAccAssuredWorkloadsV2Context(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAssuredWorkloadsV2WorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAssuredWorkloadsV2Workload_resourceConfigs(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("google_assured_workloads_v2_workload.existing_folder", "computed_target_resource", regexp.MustCompile("^folders/")),
					resource.TestCheckResourceAttrSet("google_assured_workloads_v2_workload.new_folder", "framework.0.major_revision_id"),
					resource.TestMatchResourceAttr("google_assured_workloads_v2_workload.new_folder", "computed_target_resource", regexp.MustCompile("^folders/")),
					resource.TestMatchResourceAttr("google_assured_workloads_v2_workload.new_project", "computed_target_resource", regexp.MustCompile("^projects/")),
				),
			},
			{
				ResourceName:            "google_assured_workloads_v2_workload.existing_folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_config", "deletion_policy"},
			},
			{
				ResourceName:            "google_assured_workloads_v2_workload.new_folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_config", "deletion_policy"},
			},
			{
				ResourceName:            "google_assured_workloads_v2_workload.new_project",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_config", "deletion_policy"},
			},
		},
	})
}

func testAccAssuredWorkloadsV2Workload_existingProject(context map[string]interface{}, majorRevisionID int, controlsPlaceholder, extraAttrs string) string {
	return acctest.Nprintf(fmt.Sprintf(`
provider "google" {
  user_project_override = true
  billing_project       = "%%{billing_project}"
}

resource "google_project" "target" {
  project_id      = "tf-test-aw-%%{random_suffix}"
  name            = "tf-test-aw-%%{random_suffix}"
  org_id          = "%%{org_id}"
  billing_account = "%%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_assured_workloads_v2_workload" "workload" {
  organization = "%%{org_id}"
  location     = "us-east4"
  workload_id  = "tf-test-%%{random_suffix}"
  description  = "Terraform test workload on an existing project"
  %s

  resource_config {
    existing_project = "projects/${google_project.target.number}"
  }

  framework {
    framework         = "organizations/%%{org_id}/locations/us-east4/frameworks/builtin-us-data-boundary"
    major_revision_id = %d
  }
%s
}
`, extraAttrs, majorRevisionID, controlsPlaceholder), context)
}

func testAccAssuredWorkloadsV2Workload_resourceConfigs(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  user_project_override = true
  billing_project       = "%{billing_project}"
}

resource "google_folder" "target" {
  display_name        = "tf-test-aw-folder-%{random_suffix}"
  parent              = "organizations/%{org_id}"
  deletion_protection = false
}

resource "google_assured_workloads_v2_workload" "existing_folder" {
  organization = "%{org_id}"
  location     = "us-east4"
  workload_id  = "tf-test-ef-%{random_suffix}"

  resource_config {
    existing_folder = google_folder.target.name
  }

  framework {
    framework         = "organizations/%{org_id}/locations/us-east4/frameworks/builtin-us-data-boundary"
    major_revision_id = 62
  }
%{controls_rev62_monitor}
}

resource "google_assured_workloads_v2_workload" "new_folder" {
  organization = "%{org_id}"
  location     = "us-east4"
  workload_id  = "tf-test-nf-%{random_suffix}"

  resource_config {
    new_folder_config {
      parent       = "organizations/%{org_id}"
      display_name = "tf-test-awf%{random_suffix}"
    }
  }

  framework {
    framework = "organizations/%{org_id}/locations/us-east4/frameworks/builtin-us-data-boundary"
  }
%{controls_rev62_monitor}
}

resource "google_assured_workloads_v2_workload" "new_project" {
  organization = "%{org_id}"
  location     = "us-east4"
  workload_id  = "tf-test-np-%{random_suffix}"

  resource_config {
    new_project_config {
      parent          = "organizations/%{org_id}"
      display_name    = "tf-test-awp%{random_suffix}"
      billing_account = "billingAccounts/%{billing_account}"
    }
  }

  framework {
    framework         = "organizations/%{org_id}/locations/us-east4/frameworks/builtin-us-data-boundary"
    major_revision_id = 62
  }
%{controls_rev62_monitor}
}
`, context)
}
