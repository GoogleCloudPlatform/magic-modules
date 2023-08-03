package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccContainerAwsCluster_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"service_account":        GetTestServiceAccountFromEnv(t),
		"service_account_update": "fake_user",
		"aws_acct_id":            "111111111111",
		"aws_key":                "00000000-0000-0000-0000-17aad2f0f61f",
		"aws_key_update":         "00000000-0000-0000-0000-998877665544",
		"aws_region":             "us-west-2",
		"aws_sg":                 "sg-0b3f63cb91b247628",
		"aws_sg_update":          "sg-9EEEE00001111FFFF",
		"aws_subnet":             "subnet-0b3f63cb91b247628",
		"aws_vpc":                "vpc-0b3f63cb91b247628",
		"byo_prefix":             "mmv1",
		"random_suffix":          acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerAwsClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAwsCluster_containerAwsCluster_full(context),
			},
			{
				ResourceName:            "google_container_aws_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAwsCluster_containerAwsCluster_update(context),
			},
			{
				ResourceName:            "google_container_aws_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAwsCluster_containerAwsCluster_destroy(context),
			},
			{
				ResourceName:            "google_container_aws_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func TestAccContainerAwsCluster_betaUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"service_account":        GetTestServiceAccountFromEnv(t),
		"service_account_update": "fake_user",
		"aws_acct_id":            "111111111111",
		"aws_key":                "00000000-0000-0000-0000-17aad2f0f61f",
		"aws_key_update":         "00000000-0000-0000-0000-998877665544",
		"aws_region":             "us-west-2",
		"aws_sg":                 "sg-0b3f63cb91b247628",
		"aws_sg_update":          "sg-9EEEE00001111FFFF",
		"aws_subnet":             "subnet-0b3f63cb91b247628",
		"aws_vpc":                "vpc-0b3f63cb91b247628",
		"byo_prefix":             "mmv1",
		"random_suffix":          acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckContainerAwsClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAwsCluster_containerAwsCluster_beta(context),
			},
			{
				ResourceName:            "google_container_aws_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAwsCluster_containerAwsCluster_betaUpdate(context),
			},
			{
				ResourceName:            "google_container_aws_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAwsCluster_containerAwsCluster_betaDestroy(context),
			},
			{
				ResourceName:            "google_container_aws_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccContainerAwsCluster_containerAwsCluster_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_container_aws_versions" "versions" {
  project = data.google_project.project.project_id
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  location = "us-west1"
  name     = "full%{random_suffix}"
  description = "A sample aws cluster"
  project     = data.google_project.project.project_id

  authorization {
    admin_users {
      username = "%{service_account}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-oneplatform"
      role_session_name = "%{byo_prefix}-1p-dev-session"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[1]}"
    instance_type        = "t3.medium"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    root_volume {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
      size_gib    = 10
      volume_type = "GP2"
    }

    security_group_ids = ["%{aws_sg}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      owner = "%{service_account}"
    }
  }

  fleet {
    project = "projects/${data.google_project.project.number}"
  }

  networking {
    pod_address_cidr_blocks         = ["10.2.0.0/16"]
    service_address_cidr_blocks     = ["10.1.0.0/16"]
    vpc_id                          = "%{aws_vpc}"
    per_node_pool_sg_rules_disabled = true
  }

  annotations = {
    label-one = "value-one"
  }

  logging_config {
    component_config {
      enable_components = ["SYSTEM_COMPONENTS"]
    }
  }

  monitoring_config {
    managed_prometheus_config {
      enabled = false
    }
  }

  lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

func testAccContainerAwsCluster_containerAwsCluster_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_container_aws_versions" "versions" {
  project = data.google_project.project.project_id
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  location = "us-west1"
  name     = "full%{random_suffix}"
  description = "An updated sample aws cluster"
  project     = data.google_project.project.project_id

  authorization {
    admin_users {
      username = "%{service_account}"
    }
    admin_users {
      username = "%{service_account_update}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-update"
      role_session_name = "%{byo_prefix}-1p-dev-session-update"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-update"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"
    instance_type        = "m5.large"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "987654321098-ABCD-EFGH-IJKL-12345678"
    }

    root_volume {
      iops        = 4000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
      size_gib    = 15
      throughput  = 1000
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}", "%{aws_sg_update}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh-update"
    }

    tags = {
      owner = "%{service_account}"
      updated = "new tag"
    }
  }

  fleet {
    project = "projects/${data.google_project.project.number}"
  }

  networking {
    pod_address_cidr_blocks         = ["10.2.0.0/16"]
    service_address_cidr_blocks     = ["10.1.0.0/16"]
    vpc_id                          = "%{aws_vpc}"
    per_node_pool_sg_rules_disabled = false
  }

  annotations = {
    label-two = "value-two"
  }

  logging_config {
    component_config {
      enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
    }
  }

  monitoring_config {
    managed_prometheus_config {
      enabled = true
    }
  }

  lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

// Duplicate of testAccContainerAwsCluster_containerAwsCluster_update without lifecycle.prevent_destroy set
// so the test can clean up the resource after the update.
func testAccContainerAwsCluster_containerAwsCluster_destroy(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_container_aws_versions" "versions" {
  project = data.google_project.project.project_id
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  location = "us-west1"
  name     = "full%{random_suffix}"
  description = "An updated sample aws cluster"
  project     = data.google_project.project.project_id

  authorization {
    admin_users {
      username = "%{service_account}"
    }
    admin_users {
      username = "%{service_account_update}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-update"
      role_session_name = "%{byo_prefix}-1p-dev-session-update"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-update"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"
    instance_type        = "m5.large"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "987654321098-ABCD-EFGH-IJKL-12345678"
    }

    root_volume {
      iops        = 4000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
      size_gib    = 15
      throughput  = 1000
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}", "%{aws_sg_update}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh-update"
    }

    tags = {
      owner = "%{service_account}"
      updated = "new tag"
    }
  }

  fleet {
    project = "projects/${data.google_project.project.number}"
  }

  networking {
    pod_address_cidr_blocks         = ["10.2.0.0/16"]
    service_address_cidr_blocks     = ["10.1.0.0/16"]
    vpc_id                          = "%{aws_vpc}"
    per_node_pool_sg_rules_disabled = false
  }

  annotations = {
    label-two = "value-two"
  }

  logging_config {
    component_config {
      enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
    }
  }

  monitoring_config {
    managed_prometheus_config {
      enabled = true
    }
  }
}
`, context)
}

func testAccContainerAwsCluster_containerAwsCluster_beta(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  provider = google-beta
}

data "google_container_aws_versions" "versions" {
  provider = google-beta
  project = data.google_project.project.project_id
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  provider = google-beta
  location = "us-west1"
  name     = "full%{random_suffix}"
  description = "A sample aws cluster"
  project     = data.google_project.project.project_id

  authorization {
    admin_users {
      username = "%{service_account}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-oneplatform"
      role_session_name = "%{byo_prefix}-1p-dev-session"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[1]}"
    instance_type        = "t3.medium"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    root_volume {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
      size_gib    = 10
      volume_type = "GP2"
    }

    security_group_ids = ["%{aws_sg}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      owner = "%{service_account}"
    }

    instance_placement {
      tenancy = "DEFAULT"
    }
  }

  fleet {
    project = "projects/${data.google_project.project.number}"
  }

  networking {
    pod_address_cidr_blocks         = ["10.2.0.0/16"]
    service_address_cidr_blocks     = ["10.1.0.0/16"]
    vpc_id                          = "%{aws_vpc}"
    per_node_pool_sg_rules_disabled = true
  }

  annotations = {
    label-one = "value-one"
  }

  logging_config {
    component_config {
      enable_components = ["SYSTEM_COMPONENTS"]
    }
  }

  monitoring_config {
    managed_prometheus_config {
      enabled = false
    }
  }

  lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

func testAccContainerAwsCluster_containerAwsCluster_betaUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  provider = google-beta
}

data "google_container_aws_versions" "versions" {
  provider = google-beta
  project = data.google_project.project.project_id
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  provider = google-beta
  location = "us-west1"
  name     = "full%{random_suffix}"
  description = "An updated sample aws cluster"
  project     = data.google_project.project.project_id

  authorization {
    admin_users {
      username = "%{service_account}"
    }
    admin_users {
      username = "%{service_account_update}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-update"
      role_session_name = "%{byo_prefix}-1p-dev-session-update"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-update"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"
    instance_type        = "m5.large"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "987654321098-ABCD-EFGH-IJKL-12345678"
    }

    root_volume {
      iops        = 4000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
      size_gib    = 15
      throughput  = 1000
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}", "%{aws_sg_update}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh-update"
    }

    tags = {
      owner = "%{service_account}"
      updated = "new tag"
    }

    instance_placement {
      tenancy = "HOST"
    }
  }

  fleet {
    project = "projects/${data.google_project.project.number}"
  }

  networking {
    pod_address_cidr_blocks         = ["10.2.0.0/16"]
    service_address_cidr_blocks     = ["10.1.0.0/16"]
    vpc_id                          = "%{aws_vpc}"
    per_node_pool_sg_rules_disabled = false
  }

  annotations = {
    label-two = "value-two"
  }

  logging_config {
    component_config {
      enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
    }
  }

  monitoring_config {
    managed_prometheus_config {
      enabled = true
    }
  }

  lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

// Duplicate of testAccContainerAwsCluster_containerAwsCluster_update without lifecycle.prevent_destroy set
// so the test can clean up the resource after the update.
func testAccContainerAwsCluster_containerAwsCluster_betaDestroy(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  provider = google-beta
}

data "google_container_aws_versions" "versions" {
  provider = google-beta
  project = data.google_project.project.project_id
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  provider = google-beta
  location = "us-west1"
  name     = "full%{random_suffix}"
  description = "An updated sample aws cluster"
  project     = data.google_project.project.project_id

  authorization {
    admin_users {
      username = "%{service_account}"
    }
    admin_users {
      username = "%{service_account_update}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-update"
      role_session_name = "%{byo_prefix}-1p-dev-session-update"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-update"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"
    instance_type        = "m5.large"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "987654321098-ABCD-EFGH-IJKL-12345678"
    }

    root_volume {
      iops        = 4000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
      size_gib    = 15
      throughput  = 1000
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}", "%{aws_sg_update}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh-update"
    }

    tags = {
      owner = "%{service_account}"
      updated = "new tag"
    }

    instance_placement {
      tenancy = "HOST"
    }
  }

  fleet {
    project = "projects/${data.google_project.project.number}"
  }

  networking {
    pod_address_cidr_blocks         = ["10.2.0.0/16"]
    service_address_cidr_blocks     = ["10.1.0.0/16"]
    vpc_id                          = "%{aws_vpc}"
    per_node_pool_sg_rules_disabled = false
  }

  annotations = {
    label-two = "value-two"
  }

  logging_config {
    component_config {
      enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
    }
  }

  monitoring_config {
    managed_prometheus_config {
      enabled = true
    }
  }
}
`, context)
}
