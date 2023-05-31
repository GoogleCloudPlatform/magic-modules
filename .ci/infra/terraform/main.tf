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

resource "google_organization_iam_member" "sa_access_boundary_admin" {
  org_id = data.google_organization.org.org_id
  role   = "roles/iam.accessBoundaryAdmin"
  member = google_service_account.sa.member
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

resource "google_billing_account_iam_member" "sa_master_billing_log_writer" {
  billing_account_id = data.google_billing_account.master_acct.id
  role               = "roles/logging.configWriter"
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
    "datamigration.googleapis.com",
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
    "gkeonprem.googleapis.com",
    "googlecloudmessaging.googleapis.com",
    "healthcare.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "iap.googleapis.com",
    "identitytoolkit.googleapis.com",
    "ids.googleapis.com",
    "logging.googleapis.com",
    "looker.googleapis.com",
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

resource "google_project_service_identity" "bigtable_sa" {
  provider = google-beta
  depends_on = [module.project-services]

  project = google_project.proj.project_id
  service = "bigtableadmin.googleapis.com"
}

resource "google_project_service_identity" "secretmanager_sa" {
  provider = google-beta
  depends_on = [module.project-services]

  project = google_project.proj.project_id
  service = "secretmanager.googleapis.com"
}

resource "google_project_service_identity" "sqladmin_sa" {
  provider = google-beta
  depends_on = [module.project-services]

  project = google_project.proj.project_id
  service = "sqladmin.googleapis.com"
}

# TODO: Replace these permissions with bootstrapped permissions

# TestAccComposerEnvironment_fixPyPiPackages
# TestAccComposerEnvironmentComposer2_private
# TestAccComposerEnvironment_withEncryptionConfigComposer1
# TestAccComposerEnvironment_withEncryptionConfigComposer2
# TestAccComposerEnvironment_ComposerV2
# TestAccComposerEnvironment_UpdateComposerV2
# TestAccComposerEnvironment_composerV2PrivateServiceConnect
# TestAccComposerEnvironment_composerV2MasterAuthNetworks
# TestAccComposerEnvironment_composerV2MasterAuthNetworksUpdate
# TestAccComposerEnvironmentAirflow2_withRecoveryConfig
resource "google_project_iam_member" "composer_agent_v2_ext" {
  project = google_project.proj.project_id
  role    = "roles/composer.ServiceAgentV2Ext"
  member  = "serviceAccount:service-${google_project.proj.number}@cloudcomposer-accounts.iam.gserviceaccount.com"
}

# TestAccComputeInstance_resourcePolicyUpdate
resource "google_project_iam_member" "compute_agent_instance_admin" {
  project = google_project.proj.project_id
  role    = "roles/compute.instanceAdmin"
  member  = "serviceAccount:service-${google_project.proj.number}@compute-system.iam.gserviceaccount.com"
}

# TestAccCloudfunctions2function_cloudfunctions2SecretEnvExample
# TestAccCloudfunctions2function_cloudfunctions2SecretVolumeExample
resource "google_project_iam_member" "compute_agent_secret_accessor" {
  project = google_project.proj.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_project.proj.number}-compute@developer.gserviceaccount.com"
}

# TestAccVertexAIEndpoint_vertexAiEndpointNetwork
# TestAccVertexAIFeaturestoreEntitytype_vertexAiFeaturestoreEntitytypeExample
# TestAccVertexAIFeaturestoreEntitytype_vertexAiFeaturestoreEntitytypeWithBetaFieldsExample
# TestAccVertexAIFeaturestore_vertexAiFeaturestoreExample
# TestAccVertexAIFeaturestore_vertexAiFeaturestoreScalingExample
# TestAccVertexAIFeaturestore_vertexAiFeaturestoreWithBetaFieldsExample
# TestAccVertexAIMetadataStore_vertexAiMetadataStoreExample
# TestAccVertexAITensorboard_vertexAiTensorboardFullExample
resource "google_project_iam_member" "aiplatform_agent_encrypter_decrypter" {
  project = google_project.proj.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${google_project.proj.number}@gcp-sa-aiplatform.iam.gserviceaccount.com"
}

data "google_organization" "org2" {
  organization = var.org2_id
}

resource "google_organization_iam_member" "sa_org2_admin" {
  org_id = data.google_organization.org2.org_id
  role   = "roles/resourcemanager.organizationAdmin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_org2_owner" {
  org_id = data.google_organization.org2.org_id
  role   = "roles/owner"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_org2_policy_admin" {
  org_id = data.google_organization.org2.org_id
  role   = "roles/orgpolicy.policyAdmin"
  member = google_service_account.sa.member
}

resource "google_organization_iam_member" "sa_org2_resource_settings_admin" {
  org_id = data.google_organization.org2.org_id
  role   = "roles/resourcesettings.admin"
  member = google_service_account.sa.member
}

resource "google_project" "firestore_proj" {
  name            = var.firestore_project_id
  project_id      = var.firestore_project_id
  org_id          = data.google_organization.org.org_id
  billing_account = var.billing_account_id
}

module "firestore-project-services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  version = "~> 14.1"

  project_id = google_project.firestore_proj.project_id

  activate_apis = [
    "firestore.googleapis.com",
  ]
}

resource "google_firestore_database" "firestore_db" {
  provider = google-beta
  depends_on = [module.firestore-project-services]

  project     = google_project.firestore_proj.project_id
  name        = "(default)"
  location_id = "nam5"
  type        = "FIRESTORE_NATIVE"
}

output "service_account" {
  value = google_service_account.sa.email
}

output "project_number" {
  value = google_project.proj.number
}
