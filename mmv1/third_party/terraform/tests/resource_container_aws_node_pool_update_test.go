package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccContainerAwsNodePool_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"service_account":        GetTestServiceAccountFromEnv(t),
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
		CheckDestroy:             testAccCheckContainerAwsNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAwsNodePool_containerAwsNodePool_full(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAwsNodePool_containerAwsNodePool_update(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAwsNodePool_containerAwsNodePool_destroy(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func TestAccContainerAwsNodePool_betaUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"service_account":        GetTestServiceAccountFromEnv(t),
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
		CheckDestroy:             testAccCheckContainerAwsNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAwsNodePool_containerAwsNodePool_betaFull(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAwsNodePool_containerAwsNodePool_betaUpdate(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAwsNodePool_containerAwsNodePool_betaDestroy(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccContainerAwsNodePool_containerAwsNodePool_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_container_aws_versions" "versions" {
	project = data.google_project.project.project_id
	location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
	location = "us-west1"
	name     = "full%{random_suffix}-cp"
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
		}

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		database_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
		subnet_ids           = ["%{aws_subnet}"]
		version              = data.google_container_aws_versions.versions.valid_versions[0]
		security_group_ids   = ["%{aws_sg}"]

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
		pod_address_cidr_blocks     = ["10.2.0.0/16"]
		service_address_cidr_blocks = ["10.1.0.0/16"]
		vpc_id                      = "%{aws_vpc}"
	}
}

resource "google_container_aws_node_pool" "primary" {
	location = "us-west1"
	name     = "full%{random_suffix}-np"
	project  = data.google_project.project.project_id
	cluster  = google_container_aws_cluster.primary.name
	version  = data.google_container_aws_versions.versions.valid_versions[1]

	autoscaling {
		min_node_count = 1
		max_node_count = 5
	}

	subnet_id = "%{aws_subnet}"

	max_pods_constraint {
		max_pods_per_node = 110
	}

	config {
		instance_type        = "m5.large"
		iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool"

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		root_volume {
			size_gib    = 10
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
			volume_type = "GP2"
		}

		taints {
			key    = "taint-key"
			value  = "taint-value"
			effect = "PREFER_NO_SCHEDULE"
		}

		labels = {
			label-one = "value-one"
		}

		tags = {
			tag-one = "value-one"
		}

		ssh_config {
			ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
		}

		proxy_config {
			secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
			secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
		}

		autoscaling_metrics_collection {
			granularity = "1Minute"
			metrics     = ["GroupMinSize"]
		}

		security_group_ids = ["%{aws_sg}"]
	}

	annotations = {
		label-one = "value-one"
	}

	lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

func testAccContainerAwsNodePool_containerAwsNodePool_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_container_aws_versions" "versions" {
	project = data.google_project.project.project_id
	location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
	location = "us-west1"
	name     = "full%{random_suffix}-cp"
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
		}

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		database_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
		subnet_ids           = ["%{aws_subnet}"]
		version              = data.google_container_aws_versions.versions.valid_versions[0]
		security_group_ids   = ["%{aws_sg}"]

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
		pod_address_cidr_blocks     = ["10.2.0.0/16"]
		service_address_cidr_blocks = ["10.1.0.0/16"]
		vpc_id                      = "%{aws_vpc}"
	}
}

resource "google_container_aws_node_pool" "primary" {
	location = "us-west1"
	name     = "full%{random_suffix}-np"
	project  = data.google_project.project.project_id
	cluster  = google_container_aws_cluster.primary.name
	version  = data.google_container_aws_versions.versions.valid_versions[0]

	autoscaling {
		min_node_count = 2
		max_node_count = 4
	}

	subnet_id = "%{aws_subnet}"

	max_pods_constraint {
		max_pods_per_node = 110
	}

	config {
		instance_type        = "m5.large"
		iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool-update"

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
		}

		root_volume {
			size_gib    = 15
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
			volume_type = "GP3"
			iops        = 3000
			throughput  = 500
		}

		taints {
			key    = "taint-key"
			value  = "taint-value"
			effect = "PREFER_NO_SCHEDULE"
		}

		labels = {
			label-two = "value-two"
		}

		tags = {
			tag-two = "value-two"
		}

		ssh_config {
			ec2_key_pair = "%{byo_prefix}-1p-dev-ssh-update"
		}

		proxy_config {
			secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-FEDCBA"
			secret_version = "87654321-ABCD-EFGH-IJKL-987654321098"
		}

		autoscaling_metrics_collection {
			granularity = "1Minute"
			metrics     = ["GroupMaxSize"]
		}

		security_group_ids = ["%{aws_sg_update}"]
	}

	annotations = {
		label-two = "value-two"
	}

	lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

// Duplicate of testAccContainerAwsNodePool_containerAwsNodePool_update without lifecycle.prevent_destroy set
// so the test can clean up the resource after the update.
func testAccContainerAwsNodePool_containerAwsNodePool_destroy(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_container_aws_versions" "versions" {
	project = data.google_project.project.project_id
	location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
	location = "us-west1"
	name     = "full%{random_suffix}-cp"
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
		}

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		database_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
		subnet_ids           = ["%{aws_subnet}"]
		version              = data.google_container_aws_versions.versions.valid_versions[0]
		security_group_ids   = ["%{aws_sg}"]

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
		pod_address_cidr_blocks     = ["10.2.0.0/16"]
		service_address_cidr_blocks = ["10.1.0.0/16"]
		vpc_id                      = "%{aws_vpc}"
	}
}

resource "google_container_aws_node_pool" "primary" {
	location = "us-west1"
	name     = "full%{random_suffix}-np"
	project  = data.google_project.project.project_id
	cluster  = google_container_aws_cluster.primary.name
	version  = data.google_container_aws_versions.versions.valid_versions[0]

	autoscaling {
		min_node_count = 2
		max_node_count = 4
	}

	subnet_id = "%{aws_subnet}"

	max_pods_constraint {
		max_pods_per_node = 110
	}

	config {
		instance_type        = "m5.large"
		iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool-update"

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
		}

		root_volume {
			size_gib    = 15
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
			volume_type = "GP3"
			iops        = 3000
			throughput  = 500
		}

		taints {
			key    = "taint-key"
			value  = "taint-value"
			effect = "PREFER_NO_SCHEDULE"
		}

		labels = {
			label-two = "value-two"
		}

		tags = {
			tag-two = "value-two"
		}

		ssh_config {
			ec2_key_pair = "%{byo_prefix}-1p-dev-ssh-update"
		}

		proxy_config {
			secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-FEDCBA"
			secret_version = "87654321-ABCD-EFGH-IJKL-987654321098"
		}

		autoscaling_metrics_collection {
			granularity = "1Minute"
			metrics     = ["GroupMaxSize"]
		}

		security_group_ids = ["%{aws_sg_update}"]
	}

	annotations = {
		label-two = "value-two"
	}
}
`, context)
}

func testAccContainerAwsNodePool_containerAwsNodePool_betaFull(context map[string]interface{}) string {
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
	name     = "full%{random_suffix}-cp"
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
		}

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		database_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
		subnet_ids           = ["%{aws_subnet}"]
		version              = data.google_container_aws_versions.versions.valid_versions[0]
		security_group_ids   = ["%{aws_sg}"]

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
		pod_address_cidr_blocks     = ["10.2.0.0/16"]
		service_address_cidr_blocks = ["10.1.0.0/16"]
		vpc_id                      = "%{aws_vpc}"
	}
}

resource "google_container_aws_node_pool" "primary" {
	provider = google-beta
	location = "us-west1"
	name     = "full%{random_suffix}-np"
	project  = data.google_project.project.project_id
	cluster  = google_container_aws_cluster.primary.name
	version  = data.google_container_aws_versions.versions.valid_versions[1]

	autoscaling {
		min_node_count = 1
		max_node_count = 5
	}

	subnet_id = "%{aws_subnet}"

	max_pods_constraint {
		max_pods_per_node = 110
	}

	config {
		instance_type        = "m5.large"
		iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool"

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		root_volume {
			size_gib    = 10
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
			volume_type = "GP2"
		}

		taints {
			key    = "taint-key"
			value  = "taint-value"
			effect = "PREFER_NO_SCHEDULE"
		}

		labels = {
			label-one = "value-one"
		}

		tags = {
			tag-one = "value-one"
		}

		ssh_config {
			ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
		}

		proxy_config {
			secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
			secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
		}

		autoscaling_metrics_collection {
			granularity = "1Minute"
			metrics     = ["GroupMinSize"]
		}

		security_group_ids = ["%{aws_sg}"]

		instance_placement {
      tenancy = "DEFAULT"
    }
	}

	annotations = {
		label-one = "value-one"
	}

	lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

func testAccContainerAwsNodePool_containerAwsNodePool_betaUpdate(context map[string]interface{}) string {
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
	name     = "full%{random_suffix}-cp"
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
		}

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		database_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
		subnet_ids           = ["%{aws_subnet}"]
		version              = data.google_container_aws_versions.versions.valid_versions[0]
		security_group_ids   = ["%{aws_sg}"]

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
		pod_address_cidr_blocks     = ["10.2.0.0/16"]
		service_address_cidr_blocks = ["10.1.0.0/16"]
		vpc_id                      = "%{aws_vpc}"
	}
}

resource "google_container_aws_node_pool" "primary" {
	provider = google-beta
	location = "us-west1"
	name     = "full%{random_suffix}-np"
	project  = data.google_project.project.project_id
	cluster  = google_container_aws_cluster.primary.name
	version  = data.google_container_aws_versions.versions.valid_versions[0]

	autoscaling {
		min_node_count = 2
		max_node_count = 4
	}

	subnet_id = "%{aws_subnet}"

	max_pods_constraint {
		max_pods_per_node = 110
	}

	config {
		instance_type        = "m5.large"
		iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool-update"

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
		}

		root_volume {
			size_gib    = 15
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
			volume_type = "GP3"
			iops        = 3000
			throughput  = 500
		}

		taints {
			key    = "taint-key"
			value  = "taint-value"
			effect = "PREFER_NO_SCHEDULE"
		}

		labels = {
			label-two = "value-two"
		}

		tags = {
			tag-two = "value-two"
		}

		ssh_config {
			ec2_key_pair = "%{byo_prefix}-1p-dev-ssh-update"
		}

		proxy_config {
			secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-FEDCBA"
			secret_version = "87654321-ABCD-EFGH-IJKL-987654321098"
		}

		autoscaling_metrics_collection {
			granularity = "1Minute"
			metrics     = ["GroupMaxSize"]
		}

		security_group_ids = ["%{aws_sg_update}"]

		instance_placement {
      tenancy = "HOST"
    }
	}

	annotations = {
		label-two = "value-two"
	}

	management {
    auto_repair = true
  }

	lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

// Duplicate of testAccContainerAwsNodePool_containerAwsNodePool_betaUpdate without lifecycle.prevent_destroy set
// so the test can clean up the resource after the update.
func testAccContainerAwsNodePool_containerAwsNodePool_betaDestroy(context map[string]interface{}) string {
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
	name     = "full%{random_suffix}-cp"
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
		}

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		database_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key}"
		}

		iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
		subnet_ids           = ["%{aws_subnet}"]
		version              = data.google_container_aws_versions.versions.valid_versions[0]
		security_group_ids   = ["%{aws_sg}"]

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
		pod_address_cidr_blocks     = ["10.2.0.0/16"]
		service_address_cidr_blocks = ["10.1.0.0/16"]
		vpc_id                      = "%{aws_vpc}"
	}
}

resource "google_container_aws_node_pool" "primary" {
	provider = google-beta
	location = "us-west1"
	name     = "full%{random_suffix}-np"
	project  = data.google_project.project.project_id
	cluster  = google_container_aws_cluster.primary.name
	version  = data.google_container_aws_versions.versions.valid_versions[0]

	autoscaling {
		min_node_count = 2
		max_node_count = 4
	}

	subnet_id = "%{aws_subnet}"

	max_pods_constraint {
		max_pods_per_node = 110
	}

	config {
		instance_type        = "m5.large"
		iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool-update"

		config_encryption {
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
		}

		root_volume {
			size_gib    = 15
			kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_key_update}"
			volume_type = "GP3"
			iops        = 3000
			throughput  = 500
		}

		taints {
			key    = "taint-key"
			value  = "taint-value"
			effect = "PREFER_NO_SCHEDULE"
		}

		labels = {
			label-two = "value-two"
		}

		tags = {
			tag-two = "value-two"
		}

		ssh_config {
			ec2_key_pair = "%{byo_prefix}-1p-dev-ssh-update"
		}

		proxy_config {
			secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-FEDCBA"
			secret_version = "87654321-ABCD-EFGH-IJKL-987654321098"
		}

		autoscaling_metrics_collection {
			granularity = "1Minute"
			metrics     = ["GroupMaxSize"]
		}

		security_group_ids = ["%{aws_sg_update}"]

		instance_placement {
      tenancy = "HOST"
    }
	}

	annotations = {
		label-two = "value-two"
	}

	management {
    auto_repair = true
  }
}
`, context)
}
