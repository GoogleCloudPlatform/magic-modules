package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerUpdateExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerBasic(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerUpdateExample2(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerIdentifyingFields(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.identifying_fields",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerIdentifyingFieldsUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.identifying_fields_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerPubsub(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": GetTestProjectFromEnv(),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_publishToPubSub(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.pubsub",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerDeidentifyUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerDeidentifyBasic(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.actions",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerDeidentifyUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.actions",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerChangingActions(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerJobNotificationEmails(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.actions",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerDeidentifyBasic(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.actions",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerJobNotificationEmails(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.actions",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerInspect(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": GetTestProjectFromEnv(),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_inspectBasic(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.inspect",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccDataLossPreventionJobTrigger_inspectUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.inspect",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerInspectCustomInfoTypes(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": GetTestProjectFromEnv(),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_inspectCustomInfoTypes(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.inspect",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "basic" {
	parent = "projects/%{project}"
	description = "Starting description"
	display_name = "display"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "project"
						dataset_id = "dataset123"
					}
				}
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerIdentifyingFields(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "identifying_fields" {
	parent = "projects/%{project}"
	description = "Starting description"
	display_name = "display"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "project"
						dataset_id = "dataset123"
					}
				}
			}
		}
		storage_config {
			big_query_options {
				table_reference {
					project_id = "project"
					dataset_id = "dataset"
					table_id = "table_to_scan"
				}
				rows_limit = 1000
				sample_method = "RANDOM_START"
				identifying_fields {
					name = "field"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "basic" {
	parent = "projects/%{project}"
	description = "An updated description"
	display_name = "Different"

	triggers {
		schedule {
			recurrence_period_duration = "86500s"
		}
	}

	inspect_job {
		inspect_template_name = "other"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "different"
						dataset_id = "asdf"
					}
				}
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerIdentifyingFieldsUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "identifying_fields_update" {
	parent = "projects/%{project}"
	description = "An updated description"
	display_name = "Different"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "project"
						dataset_id = "dataset123"
					}
				}
			}
		}
		storage_config {
			big_query_options {
				table_reference {
					project_id = "project"
					dataset_id = "dataset"
					table_id = "table_to_scan"
				}
				rows_limit = 1000
				sample_method = "RANDOM_START"
				identifying_fields {
					name = "different"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_publishToPubSub(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "pubsub" {
	parent = "projects/%{project}"
	description = "Starting description"
	display_name = "display"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			pub_sub {
				topic = "projects/%{project}/topics/bar"
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerDeidentifyBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "actions" {
	parent       = "projects/%{project}"
	description  = "Description for the job_trigger created by terraform"
	display_name = "TerraformDisplayName"
	
	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}
	
	inspect_job {
		inspect_template_name = "sample-inspect-template"
		actions {
			deidentify {
				cloud_storage_output    = "gs://samplebucket/dir/"
				file_types_to_transform = ["CSV", "IMAGE", "TSV"]
				transformation_details_storage_config {
					table {
						project_id = "%{project}"
						dataset_id = google_bigquery_dataset.default.dataset_id
						table_id   = google_bigquery_table.default.table_id
					}
				}
				transformation_config {
					deidentify_template            = "sample-deidentify-template"
					image_redact_template          = "sample-image-redact-template"
					structured_deidentify_template = "sample-structured-deidentify-template"
				}
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
	}
}
	
resource "google_bigquery_dataset" "default" {
	dataset_id                  = "tf_test_%{random_suffix}"
	friendly_name               = "terraform-test"
	description                 = "Description for the dataset created by terraform"
	location                    = "US"
	default_table_expiration_ms = 3600000
	
	labels = {
		env = "default"
	}
}
	
resource "google_bigquery_table" "default" {
	dataset_id          = google_bigquery_dataset.default.dataset_id
	table_id            = "tf_test_%{random_suffix}"
	deletion_protection = false
	
	time_partitioning {
		type = "DAY"
	}
	
	labels = {
		env = "default"
	}
	
	schema = <<EOF
		[
		{
			"name": "quantity",
			"type": "NUMERIC",
			"mode": "NULLABLE",
			"description": "The quantity"
		},
		{
			"name": "name",
			"type": "STRING",
			"mode": "NULLABLE",
			"description": "Name of the object"
		}
		]
	EOF
}
`, context)
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerDeidentifyUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "actions" {
	parent       = "projects/%{project}"
	description  = "Description for the job_trigger created by terraform"
	display_name = "TerraformDisplayName"
	
	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}
	
	inspect_job {
		inspect_template_name = "sample-inspect-template"
		actions {
			deidentify {
				cloud_storage_output    = "gs://samplebucketnew/dir/"
				file_types_to_transform = ["TEXT_FILE", "TSV"]
				transformation_details_storage_config {
					table {
						project_id = "%{project}"
						dataset_id = google_bigquery_dataset.default.dataset_id
						table_id   = google_bigquery_table.default.table_id
					}
				}
				transformation_config {
					deidentify_template            = "updated-deidentify-template"
					image_redact_template          = "updated-image-redact-template"
					structured_deidentify_template = "updated-structured-deidentify-template"
				}
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
	}
}
	
resource "google_bigquery_dataset" "default" {
	dataset_id                  = "tf_test_%{random_suffix}"
	friendly_name               = "terraform-test"
	description                 = "Description for the dataset created by terraform"
	location                    = "US"
	default_table_expiration_ms = 3600000
	
	labels = {
		env = "default"
	}
}
	
resource "google_bigquery_table" "default" {
	dataset_id          = google_bigquery_dataset.default.dataset_id
	table_id            = "tf_test_%{random_suffix}"
	deletion_protection = false
	
	time_partitioning {
		type = "DAY"
	}
	
	labels = {
		env = "default"
	}
	
	schema = <<EOF
		[
		{
			"name": "quantity",
			"type": "NUMERIC",
			"mode": "NULLABLE",
			"description": "The quantity"
		},
		{
			"name": "name",
			"type": "STRING",
			"mode": "NULLABLE",
			"description": "Name of the object"
		}
		]
	EOF
}
`, context)
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerJobNotificationEmails(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "actions" {
	parent       = "projects/%{project}"
	description  = "Description for the job_trigger created by terraform"
	display_name = "TerraformDisplayName"
	
	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}
	
	inspect_job {
		inspect_template_name = "sample-inspect-template"
		actions {
			job_notification_emails {}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_inspectBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "inspect" {
	parent = "projects/%{project}"
	description = "Starting description"
	display_name = "display"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "project"
						dataset_id = "dataset123"
					}
				}
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
		inspect_config {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			info_types {
				name    = "PERSON_NAME"
			}
			info_types {
				name = "LAST_NAME"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "FIRST_NAME"
			}
	
			min_likelihood     = "UNLIKELY"
			include_quote      = false
			exclude_info_types = false
			rule_set {
				info_types {
					name = "EMAIL_ADDRESS"
				}
				rules {
					exclusion_rule {
						regex {
							pattern = ".+@example.com"
						}
						matching_type = "MATCHING_TYPE_FULL_MATCH"
					}
				}
			}
			rule_set {
				info_types {
					name = "EMAIL_ADDRESS"
				}
				info_types {
					name = "DOMAIN_NAME"
				}
				info_types {
					name = "PHONE_NUMBER"
				}
				info_types {
					name = "PERSON_NAME"
				}
				info_types {
					name = "FIRST_NAME"
				}
				rules {
					exclusion_rule {
						dictionary {
							word_list {
								words = ["TEST"]
							}
						}
						matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
					}
				}
			}
	
			rule_set {
				info_types {
					name = "PERSON_NAME"
				}
				rules {
					hotword_rule {
						hotword_regex {
							pattern = "patient"
						}
						proximity {
							window_before = 50
						}
						likelihood_adjustment {
							fixed_likelihood = "VERY_LIKELY"
						}
					}
				}
			}
	
			limits {
				max_findings_per_item    = 10
				max_findings_per_request = 50
				max_findings_per_info_type {
					max_findings = "75"
					info_type {
						name = "PERSON_NAME"
					}
				}
				max_findings_per_info_type {
					max_findings = "80"
					info_type {
						name = "LAST_NAME"
					}
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_inspectUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "inspect" {
	parent = "projects/%{project}"
	description = "Starting description"
	display_name = "display"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "project"
						dataset_id = "dataset123"
					}
				}
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
		inspect_config {
			info_types {
				name    = "PERSON_NAME"
				version = "stable"
			}
			info_types {
				name = "LAST_NAME"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "FIRST_NAME"
			}
	
			min_likelihood = "UNLIKELY"
			rule_set {
				info_types {
					name = "DOMAIN_NAME"
				}
				info_types {
					name = "PHONE_NUMBER"
				}
				info_types {
					name = "PERSON_NAME"
				}
				info_types {
					name = "FIRST_NAME"
				}
				rules {
					exclusion_rule {
						dictionary {
							word_list {
								words = ["TEST"]
							}
						}
						matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
					}
				}
			}
	
			rule_set {
				info_types {
					name = "PERSON_NAME"
				}
				rules {
					hotword_rule {
						hotword_regex {
							pattern = "not-a-patient"
						}
						proximity {
							window_before = 50
						}
						likelihood_adjustment {
							fixed_likelihood = "UNLIKELY"
						}
					}
				}
			}
	
			limits {
				max_findings_per_item    = 1
				max_findings_per_request = 5
				max_findings_per_info_type {
					max_findings = "80"
					info_type {
						name = "PERSON_NAME"
					}
				}
				max_findings_per_info_type {
					max_findings = "20"
					info_type {
						name = "LAST_NAME"
					}
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_inspectCustomInfoTypes(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "inspect" {
	parent = "projects/%{project}"
	description = "Starting description"
	display_name = "display"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "project"
						dataset_id = "dataset123"
					}
				}
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
		inspect_config {
			custom_info_types {
                info_type {
                    name = "MY_CUSTOM_TYPE"
                }
    
                likelihood = "UNLIKELY"
    
                regex {
                    pattern = "test*"
                }
            }
			
			info_types {
				name = "EMAIL_ADDRESS"
			}
			info_types {
				name    = "PERSON_NAME"
				version = "latest"
			}
			info_types {
				name = "LAST_NAME"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "FIRST_NAME"
			}
	
			min_likelihood = "UNLIKELY"
			rule_set {
				info_types {
					name = "EMAIL_ADDRESS"
				}
				rules {
					exclusion_rule {
						regex {
							pattern = ".+@example.com"
						}
						matching_type = "MATCHING_TYPE_FULL_MATCH"
					}
				}
			}
			rule_set {
				info_types {
					name = "EMAIL_ADDRESS"
				}
				info_types {
					name = "DOMAIN_NAME"
				}
				info_types {
					name = "PHONE_NUMBER"
				}
				info_types {
					name = "PERSON_NAME"
				}
				info_types {
					name = "FIRST_NAME"
				}
				rules {
					exclusion_rule {
						dictionary {
							word_list {
								words = ["TEST"]
							}
						}
						matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
					}
				}
			}
	
			rule_set {
				info_types {
					name = "PERSON_NAME"
				}
				rules {
					hotword_rule {
						hotword_regex {
							pattern = "patient"
						}
						proximity {
							window_before = 50
						}
						likelihood_adjustment {
							fixed_likelihood = "VERY_LIKELY"
						}
					}
				}
			}
	
			limits {
				max_findings_per_item    = 10
				max_findings_per_request = 50
				max_findings_per_info_type {
					max_findings = "75"
					info_type {
						name = "PERSON_NAME"
					}
				}
				max_findings_per_info_type {
					max_findings = "80"
					info_type {
						name = "LAST_NAME"
					}
				}
			}
		}
	}
}
`, context)
}
