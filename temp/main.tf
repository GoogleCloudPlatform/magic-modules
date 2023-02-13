#
# This configuration is expected to be run locally by an administrator. It specifies the configuration needed for a
# test environment where the full set of acceptance tests can be run.
#
# Googlers can find record of internal requests at b/268353203.
#
# Prerequisites:
#   - An existing organization
#   - An existing billing account where charges can be applied
#   - A second existing billing account where charges can be applied (used only for TestAccProject_billing)
#   - An existing billing account where subaccounts can be created
#   - A BeyondCorp subscription on the organization
#
# After applying this configuration:
#   - Increase project quota for the new service account
#   - Increase project quota for the billing account
#   - Enable Workforce Identity Federation for new project
#   - Deploy "Hello World" app: https://cloud.google.com/appengine/docs/flexible/go/create-app
#     ```
#     gcloud components install app-engine-go
#     git clone https://github.com/GoogleCloudPlatform/golang-samples
#     cp -r golang-samples/appengine_flexible/helloworld ./.
#     cd helloworld
#     gcloud app deploy --project=<project>
#     ```
#   - Enable Multi-Tenancy
#     ```
#     curl --header "Authorization: Bearer $(gcloud auth print-access-token -q)" --header "X-Goog-User-Project: <project>" -X POST https://identitytoolkit.googleapis.com/v2/projects/<project>/identityPlatform:initializeAuth
#     curl --header "Content-Type: application/json" --header "Authorization: Bearer $(gcloud auth print-access-token -q)" --header "X-Goog-User-Project: <project>" -X PATCH https://identitytoolkit.googleapis.com/admin/v2/projects/<project>/config?updateMask=multiTenant -d '{"multiTenant": {"allowTenants": true}}'
#     ```
#   - Create a Service Account key for the new service account
#   - Add Group Admin role to new service account in Google Workspace: https://admin.google.com/ac/roles
#

provider "google" {}

data "google_organization" "org" {
  organization = var.org_id
}

data "google_billing_account" "master_acct" {
  billing_account = var.master_billing_account_id
}

resource "google_project" "proj" {
  name            = var.project_id
  project_id      = var.project_id
  org_id          = data.google_organization.org.org_id
  billing_account = var.billing_account_id
}

resource "google_service_account" "sa" {
  project      = google_project.proj.project_id
  account_id   = "hashicorp-test-runner"
  display_name = "HashiCorp Test Runner"
}

resource "google_organization_iam_member" "sa_assuredworkloads_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/assuredworkloads.admin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_billing_user" {
  org_id = data.google_organization.org.org_id
  role   = "roles/billing.user"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_billing_viewer" {
  org_id = data.google_organization.org.org_id
  role   = "roles/billing.viewer"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_cloudkms_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/cloudkms.admin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_compute_xpn_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/compute.xpnAdmin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_deny_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/iam.denyAdmin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_folder_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/resourcemanager.folderAdmin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_iap_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/iap.admin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_iap_settings_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/iap.settingsAdmin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_orgpolicy_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/orgpolicy.policyAdmin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_org_role_viewer" {
  org_id = data.google_organization.org.org_id
  role   = "roles/iam.organizationRoleViewer"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_owner" {
  org_id = data.google_organization.org.org_id
  role   = "roles/owner"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_billing_project_manager" {
  org_id = data.google_organization.org.org_id
  role   = "roles/billing.projectManager"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_project_creator" {
  org_id = data.google_organization.org.org_id
  role   = "roles/resourcemanager.projectCreator"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_project_deleter" {
  org_id = data.google_organization.org.org_id
  role   = "roles/resourcemanager.projectDeleter"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_service_account_token_creator" {
  org_id = data.google_organization.org.org_id
  role   = "roles/iam.serviceAccountTokenCreator"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_storage_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/storage.admin"
  member = google_service_account.sa.member
}

resource "google_billing_account_iam_member" "sa_master_billing_admin" {
  billing_account_id = data.google_billing_account.master_acct.id
  role               = "roles/billing.admin"
  member             = google_service_account.sa.member
}

resource "google_app_engine_application" "app" {
  project     = google_project.proj.project_id
  location_id = "us-central"
}

module "project-services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  version = "~> 14.1"

  project_id = google_project.proj.project_id

  activate_apis = [
    "accessapproval.googleapis.com",
    "accesscontextmanager.googleapis.com",
    "aiplatform.googleapis.com",
    "alloydb.googleapis.com",
    "analyticshub.googleapis.com",
    "apigateway.googleapis.com",
    "apikeys.googleapis.com",
    "appengine.googleapis.com",
    "appengineflex.googleapis.com",
    "artifactregistry.googleapis.com",
    "assuredworkloads.googleapis.com",
    "autoscaling.googleapis.com",
    "beyondcorp.googleapis.com",
    "bigquery.googleapis.com",
    "bigqueryconnection.googleapis.com",
    "bigquerydatapolicy.googleapis.com",
    "bigquerydatatransfer.googleapis.com",
    "bigquerymigration.googleapis.com",
    "bigqueryreservation.googleapis.com",
    "bigquerystorage.googleapis.com",
    "bigtable.googleapis.com",
    "bigtableadmin.googleapis.com",
    "billingbudgets.googleapis.com",
    "binaryauthorization.googleapis.com",
    "certificatemanager.googleapis.com",
    "cloudapis.googleapis.com",
    "cloudasset.googleapis.com",
    "cloudbilling.googleapis.com",
    "cloudbuild.googleapis.com",
    "clouddebugger.googleapis.com",
    "clouddeploy.googleapis.com",
    "cloudfunctions.googleapis.com",
    "cloudidentity.googleapis.com",
    "cloudiot.googleapis.com",
    "cloudkms.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "cloudscheduler.googleapis.com",
    "cloudtasks.googleapis.com",
    "cloudtrace.googleapis.com",
    "composer.googleapis.com",
    "compute.googleapis.com",
    "container.googleapis.com",
    "containeranalysis.googleapis.com",
    "containerfilesystem.googleapis.com",
    "containerregistry.googleapis.com",
    "daily-serviceconsumermanagement.sandbox.googleapis.com",
    "daily-serviceusage.sandbox.googleapis.com",
    "datacatalog.googleapis.com",
    "dataflow.googleapis.com",
    "dataform.googleapis.com",
    "datafusion.googleapis.com",
    "dataplex.googleapis.com",
    "dataproc.googleapis.com",
    "datastore.googleapis.com",
    "datastream.googleapis.com",
    "deploymentmanager.googleapis.com",
    "dialogflow.googleapis.com",
    "dlp.googleapis.com",
    "dns.googleapis.com",
    "documentai.googleapis.com",
    "edgecache.googleapis.com",
    "essentialcontacts.googleapis.com",
    "eventarc.googleapis.com",
    "eventarcpublishing.googleapis.com",
    "fcm.googleapis.com",
    "fcmregistrations.googleapis.com",
    "file.googleapis.com",
    "firebase.googleapis.com",
    "firebaseappdistribution.googleapis.com",
    "firebasedatabase.googleapis.com",
    "firebasedynamiclinks.googleapis.com",
    "firebasehosting.googleapis.com",
    "firebaseinstallations.googleapis.com",
    "firebaseremoteconfig.googleapis.com",
    "firebaserules.googleapis.com",
    "firebasestorage.googleapis.com",
    "firestore.googleapis.com",
    "firestorekeyvisualizer.googleapis.com",
    "gameservices.googleapis.com",
    "gkebackup.googleapis.com",
    "gkeconnect.googleapis.com",
    "gkehub.googleapis.com",
    "gkemulticloud.googleapis.com",
    "googlecloudmessaging.googleapis.com",
    "healthcare.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "iap.googleapis.com",
    "identitytoolkit.googleapis.com",
    "ids.googleapis.com",
    "logging.googleapis.com",
    "managedidentities.googleapis.com",
    "memcache.googleapis.com",
    "metastore.googleapis.com",
    "ml.googleapis.com",
    "mobilecrashreporting.googleapis.com",
    "monitoring.googleapis.com",
    "multiclustermetering.googleapis.com",
    "networkconnectivity.googleapis.com",
    "networkmanagement.googleapis.com",
    "networkservices.googleapis.com",
    "notebooks.googleapis.com",
    "orgpolicy.googleapis.com",
    "osconfig.googleapis.com",
    "oslogin.googleapis.com",
    "privateca.googleapis.com",
    "pubsub.googleapis.com",
    "pubsublite.googleapis.com",
    "recaptchaenterprise.googleapis.com",
    "redis.googleapis.com",
    "replicapool.googleapis.com",
    "replicapoolupdater.googleapis.com",
    "resourcesettings.googleapis.com",
    "resourceviews.googleapis.com",
    "run.googleapis.com",
    "runtimeconfig.googleapis.com",
    "secretmanager.googleapis.com",
    "securetoken.googleapis.com",
    "securitycenter.googleapis.com",
    "serviceconsumermanagement.googleapis.com",
    "servicecontrol.googleapis.com",
    "servicedirectory.googleapis.com",
    "servicemanagement.googleapis.com",
    "servicenetworking.googleapis.com",
    "serviceusage.googleapis.com",
    "sourcerepo.googleapis.com",
    "spanner.googleapis.com",
    "sql-component.googleapis.com",
    "sqladmin.googleapis.com",
    "stackdriver.googleapis.com",
    "storage-api.googleapis.com",
    "storage-component.googleapis.com",
    "storage.googleapis.com",
    "storagetransfer.googleapis.com",
    "test-file.sandbox.googleapis.com",
    "testing.googleapis.com",
    "tpu.googleapis.com",
    "trafficdirector.googleapis.com",
    "vpcaccess.googleapis.com",
    "websecurityscanner.googleapis.com",
    "workflowexecutions.googleapis.com",
    "workflows.googleapis.com",
    "workstations.googleapis.com"
  ]
}

output "service_account" {
  value = google_service_account.sa.email
}

output "project_number" {
  value = google_project.proj.number
}
