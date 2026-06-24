package bigquerydatatransfer_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/bigquery"
	"github.com/hashicorp/terraform-provider-google/google/services/bigquerydatatransfer"
	_ "github.com/hashicorp/terraform-provider-google/google/services/kms"
	_ "github.com/hashicorp/terraform-provider-google/google/services/pubsub"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"strings"
	"testing"
	"time"
)

func TestBigqueryDataTransferConfig_resourceBigqueryDTCParamsCustomDiffFuncForceNewWhenGoogleCloudStorage(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		before   map[string]interface{}
		after    map[string]interface{}
		forcenew bool
	}{
		"changing_data_path_template": {
			before: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp-new/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			forcenew: true,
		},
		"changing_destination_table_name_template": {
			before: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-new",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			forcenew: true,
		},
		"changing_non_force_new_fields": {
			before: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 1000,
					"write_disposition":               "APPEND",
				},
			},
			forcenew: false,
		},
		"changing_destination_table_name_template_for_different_data_source_id": {
			before: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"destination_table_name_template": "table-old",
					"query":                           "SELECT 1 AS a",
					"write_disposition":               "WRITE_APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"destination_table_name_template": "table-new",
					"query":                           "SELECT 1 AS a",
					"write_disposition":               "WRITE_APPEND",
				},
			},
			forcenew: false,
		},
		"changing_data_path_template_for_different_data_source_id": {
			before: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"data_path_template": "gs://bq-bucket/*.json",
					"query":              "SELECT 1 AS a",
					"write_disposition":  "WRITE_APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"data_path_template": "gs://bq-bucket-new/*.json",
					"query":              "SELECT 1 AS a",
					"write_disposition":  "WRITE_APPEND",
				},
			},
			forcenew: false,
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				"params":         tc.before["params"],
				"data_source_id": tc.before["data_source_id"],
			},
			After: map[string]interface{}{
				"params":         tc.after["params"],
				"data_source_id": tc.after["data_source_id"],
			},
		}
		err := bigquerydatatransfer.ParamsCustomizeDiffFunc(d)
		if err != nil {
			t.Errorf("failed, expected no error but received - %s for the condition %s", err, tn)
		}
		if d.IsForceNew != tc.forcenew {
			t.Errorf("ForceNew not setup correctly for the condition-'%s', expected:%v; actual:%v", tn, tc.forcenew, d.IsForceNew)
		}
	}
}

func TestBigqueryDataTransferConfig_resourceBigqueryDTCParamsCustomDiffFuncForceNewWhenAmazonS3(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		before   map[string]interface{}
		after    map[string]interface{}
		forcenew bool
	}{
		"changing_data_path": {
			before: map[string]interface{}{
				"data_source_id": "amazon_s3",
				"params": map[string]interface{}{
					"data_path":                       "s3://s3-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "WRITE_APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "amazon_s3",
				"params": map[string]interface{}{
					"data_path":                       "s3://s3-bucket-temp-new/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "WRITE_APPEND",
				},
			},
			forcenew: true,
		},
		"changing_destination_table_name_template": {
			before: map[string]interface{}{
				"data_source_id": "amazon_s3",
				"params": map[string]interface{}{
					"data_path":                       "s3://s3-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "WRITE_APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "amazon_s3",
				"params": map[string]interface{}{
					"data_path":                       "s3://s3-bucket-temp/*.json",
					"destination_table_name_template": "table-new",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "WRITE_APPEND",
				},
			},
			forcenew: true,
		},
		"changing_non_force_new_fields": {
			before: map[string]interface{}{
				"data_source_id": "amazon_s3",
				"params": map[string]interface{}{
					"data_path":                       "s3://s3-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "WRITE_APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "amazon_s3",
				"params": map[string]interface{}{
					"data_path":                       "s3://s3-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 1000,
					"write_disposition":               "APPEND",
				},
			},
			forcenew: false,
		},
		"changing_destination_table_name_template_for_different_data_source_id": {
			before: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"destination_table_name_template": "table-old",
					"query":                           "SELECT 1 AS a",
					"write_disposition":               "WRITE_APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"destination_table_name_template": "table-new",
					"query":                           "SELECT 1 AS a",
					"write_disposition":               "WRITE_APPEND",
				},
			},
			forcenew: false,
		},
		"changing_data_path_template_for_different_data_source_id": {
			before: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"data_path":         "s3://s3-bucket-temp/*.json",
					"query":             "SELECT 1 AS a",
					"write_disposition": "WRITE_APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"data_path":         "s3://s3-bucket-temp-new/*.json",
					"query":             "SELECT 1 AS a",
					"write_disposition": "WRITE_APPEND",
				},
			},
			forcenew: false,
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				"params":         tc.before["params"],
				"data_source_id": tc.before["data_source_id"],
			},
			After: map[string]interface{}{
				"params":         tc.after["params"],
				"data_source_id": tc.after["data_source_id"],
			},
		}
		err := bigquerydatatransfer.ParamsCustomizeDiffFunc(d)
		if err != nil {
			t.Errorf("failed, expected no error but received - %s for the condition %s", err, tn)
		}
		if d.IsForceNew != tc.forcenew {
			t.Errorf("ForceNew not setup correctly for the condition-'%s', expected:%v; actual:%v", tn, tc.forcenew, d.IsForceNew)
		}
	}
}

// The BigQuery Data Transfer Service agent needs a few project-level roles for
// these tests. We bootstrap them once here rather than provisioning IAM in each
// test config: managing the same shared bindings from parallel tests races and
// can also fail due to IAM propagation delays.
// See https://googlecloudplatform.github.io/magic-modules/test/test/#iam-resources
func TestAccBigqueryDataTransferConfig(t *testing.T) {
	resourcemanager.BootstrapIamMembers(t, []resourcemanager.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com",
			Role:   "roles/iam.serviceAccountTokenCreator",
		},
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com",
			Role:   "roles/pubsub.subscriber",
		},
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com",
			Role:   "roles/serviceusage.serviceUsageConsumer",
		},
	})

	testCases := map[string]func(t *testing.T){
		"basic":                            testAccBigqueryDataTransferConfig_scheduledQuery_basic,
		"update":                           testAccBigqueryDataTransferConfig_scheduledQuery_update,
		"service_account":                  testAccBigqueryDataTransferConfig_scheduledQuery_with_service_account,
		"no_destintation":                  testAccBigqueryDataTransferConfig_scheduledQuery_no_destination,
		"booleanParam":                     testAccBigqueryDataTransferConfig_copy_booleanParam,
		"update_params":                    testAccBigqueryDataTransferConfig_force_new_update_params,
		"update_service_account":           testAccBigqueryDataTransferConfig_scheduledQuery_update_service_account,
		"disable_auto_scheduling":          testAccBigqueryDataTransferConfig_disableAutoScheduling,
		"schedule_options_v2_event_driven": testAccBigqueryDataTransferConfig_scheduleOptionsV2_eventDriven,
		// Multiple connector.authentication.* fields have been deprecated and return 400 errors
		// "salesforce":             testAccBigqueryDataTransferConfig_salesforce_basic,
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

func testAccBigqueryDataTransferConfig_scheduledQuery_basic(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)
	// Use a static-but-future time so the test records/replays deterministically
	// under VCR without needing SkipIfVcr.
	base := time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 10)
	start_time := base.Format(time.RFC3339)
	end_time := base.AddDate(0, 1, 0).Format(time.RFC3339)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix, "third", start_time, end_time, "y"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_scheduleOptionsV2_eventDriven(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduleOptionsV2EventDriven(random_suffix, "subscription"),
			},
			{
				// Switch the Pub/Sub subscription and confirm the event-driven
				// schedule is updated in place rather than recreated.
				Config: testAccBigqueryDataTransferConfig_scheduleOptionsV2EventDriven(random_suffix, "subscription2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_data_transfer_config.event_driven_config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.event_driven_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_disableAutoScheduling(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_disableAutoSchedulingConfig(random_suffix),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_scheduledQuery_update(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)
	// Use static-but-future times so the test records/replays deterministically
	// under VCR without needing SkipIfVcr.
	base := time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 10)
	first_start_time := base.Format(time.RFC3339)
	first_end_time := base.AddDate(0, 1, 0).Format(time.RFC3339)
	second_start_time := base.AddDate(0, 0, 1).Format(time.RFC3339)
	second_end_time := base.AddDate(0, 2, 0).Format(time.RFC3339)
	random_suffix2 := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix, "first", first_start_time, first_end_time, "y"),
			},
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix, "second", second_start_time, second_end_time, "z"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix2, "second", second_start_time, second_end_time, "z"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_CMEK(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_CMEK_basic(random_suffix),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_scheduledQuery_no_destination(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)
	// Use a static-but-future time so the test records/replays deterministically
	// under VCR without needing SkipIfVcr.
	base := time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 10)
	start_time := base.Format(time.RFC3339)
	end_time := base.AddDate(0, 1, 0).Format(time.RFC3339)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQueryNoDestination(random_suffix, "third", start_time, end_time, "y"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_scheduledQuery_with_service_account(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery_service_account(random_suffix),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "service_account_name"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_copy_booleanParam(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_booleanParam(random_suffix),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.copy_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_force_new_update_params(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_update_params_force_new(random_suffix, "old", "old"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.update_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryDataTransferConfig_update_params_force_new(random_suffix, "new", "old"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.update_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryDataTransferConfig_update_params_force_new(random_suffix, "new", "new"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.update_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccCheckBigqueryDataTransferConfigDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigquery_data_transfer_config" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(bigquerydatatransfer.Product, config)+"{{name}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("BigqueryDataTransferConfig still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccBigqueryDataTransferConfig_scheduledQuery_update_service_account(t *testing.T) {
	random_suffix1 := acctest.RandString(t, 10)
	random_suffix2 := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery_updateServiceAccount(random_suffix1, random_suffix1),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "service_account_name"},
			},
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery_updateServiceAccount(random_suffix1, random_suffix2),
				Check:  testAccCheckDataTransferServiceAccountNamePrefix("google_bigquery_data_transfer_config.query_config", random_suffix2),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "service_account_name"},
			},
		},
	})
}

// Check if transfer config service account name starts with given prefix
func testAccCheckDataTransferServiceAccountNamePrefix(resourceName string, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if !strings.HasPrefix(rs.Primary.Attributes["service_account_name"], "bqwriter"+prefix) {
			return fmt.Errorf("Transfer config service account not updated")
		}

		return nil
	}
}

func testAccBigqueryDataTransferConfig_salesforce_basic(t *testing.T) {
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_salesforce(randomSuffix),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.salesforce_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix2, schedule, start_time, end_time, letter string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "my_dataset" {
  dataset_id    = "my_dataset%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_pubsub_topic" "my_topic" {
  name = "tf-test-my-topic-%s"
}

resource "google_bigquery_table" "my_table" {
  deletion_protection = false

  dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  table_id   = "my_table"
  schema     = <<EOF
  [
    { "name": "name", "type": "STRING" },
    { "name": "x", "type": "INTEGER" }
  ]
  EOF
}

resource "google_bigquery_data_transfer_config" "query_config" {
  display_name           = "my-query-%s"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule               = "%s sunday of quarter 00:00"
  schedule_options {
    disable_auto_scheduling = false
    start_time              = "%s"
    end_time                = "%s"
  }
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  notification_pubsub_topic = google_pubsub_topic.my_topic.id
  email_preferences {
    enable_failure_email = true
  }
  params = {
    destination_table_name_template = google_bigquery_table.my_table.table_id
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT name FROM tabl WHERE x = '%s'"
  }
}
`, random_suffix, random_suffix, random_suffix2, schedule, start_time, end_time, letter)
}

func testAccBigqueryDataTransferConfig_scheduleOptionsV2EventDriven(random_suffix, subscription string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name = "tf-test-dts-topic-%s"
}

resource "google_pubsub_subscription" "subscription" {
  name  = "tf-test-dts-subscription-%s"
  topic = google_pubsub_topic.topic.id
}

resource "google_pubsub_subscription" "subscription2" {
  name  = "tf-test-dts-subscription2-%s"
  topic = google_pubsub_topic.topic.id
}

resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-dts-bucket-%s"
  location                    = "US"
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_bigquery_dataset" "my_dataset" {
  dataset_id    = "my_dataset%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "US"
}

resource "google_bigquery_table" "my_table" {
  deletion_protection = false

  dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  table_id   = "my_table"
  schema     = <<EOF
  [
    { "name": "name", "type": "STRING" },
    { "name": "x", "type": "INTEGER" }
  ]
  EOF
}

resource "google_bigquery_data_transfer_config" "event_driven_config" {
  display_name           = "my-event-driven-%s"
  location               = google_bigquery_dataset.my_dataset.location
  data_source_id         = "google_cloud_storage"
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id

  schedule_options_v2 {
    event_driven_schedule {
      pubsub_subscription = google_pubsub_subscription.%s.id
    }
  }

  params = {
    data_path_template              = "${google_storage_bucket.bucket.url}/*.json"
    destination_table_name_template = google_bigquery_table.my_table.table_id
    file_format                     = "JSON"
    write_disposition               = "APPEND"
  }
}
`, random_suffix, random_suffix, random_suffix, random_suffix, random_suffix, random_suffix, subscription)
}

func testAccBigqueryDataTransferConfig_disableAutoSchedulingConfig(random_suffix string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "my_dataset" {
  dataset_id    = "my_dataset%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_table" "my_table" {
  deletion_protection = false

  dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  table_id   = "my_table"
  schema     = <<EOF
  [
    { "name": "name", "type": "STRING" },
    { "name": "x", "type": "INTEGER" }
  ]
  EOF
}

resource "google_bigquery_data_transfer_config" "query_config" {
  display_name           = "my-query-%s"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule_options {
    disable_auto_scheduling = true
  }
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  params = {
    destination_table_name_template = google_bigquery_table.my_table.table_id
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT name FROM tabl WHERE x = 'y'"
  }
}
`, random_suffix, random_suffix)
}

func testAccBigqueryDataTransferConfig_scheduledQuery_service_account(random_suffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "bqwriter" {
  account_id = "bqwriter%s"
}

resource "google_project_iam_member" "data_editor" {
  project = data.google_project.project.project_id

  role   = "roles/bigquery.dataEditor"
  member = "serviceAccount:${google_service_account.bqwriter.email}"
}

resource "google_bigquery_dataset" "my_dataset" {
  dataset_id    = "my_dataset%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_table" "my_table" {
  deletion_protection = false

  dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  table_id   = "my_table"
}

resource "google_bigquery_data_transfer_config" "query_config" {
  depends_on = [google_project_iam_member.data_editor]

  display_name           = "my-query-%s"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule               = "every day 00:00"
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  service_account_name   = google_service_account.bqwriter.email
  params = {
    destination_table_name_template = google_bigquery_table.my_table.table_id
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT 1 AS a"
  }
}
`, random_suffix, random_suffix, random_suffix)
}

func testAccBigqueryDataTransferConfig_scheduledQueryNoDestination(random_suffix, schedule, start_time, end_time, letter string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "my_topic" {
  name = "tf-test-my-topic-%s"
}

resource "google_bigquery_dataset" "my_dataset" {
  dataset_id    = "my_dataset%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_table" "my_table" {
  deletion_protection = false

  dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  table_id   = "my_table"
  schema     = <<EOF
  [
    { "name": "name", "type": "STRING" },
    { "name": "x", "type": "INTEGER" }
  ]
  EOF
}

resource "google_bigquery_data_transfer_config" "query_config" {
  display_name           = "my-query-%s"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule               = "%s sunday of quarter 00:00"
  schedule_options {
    disable_auto_scheduling = false
    start_time              = "%s"
    end_time                = "%s"
  }
  notification_pubsub_topic = google_pubsub_topic.my_topic.id
  email_preferences {
    enable_failure_email = true
  }
  params = {
    destination_table_name_template = google_bigquery_table.my_table.table_id
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT name FROM tabl WHERE x = '%s'"
  }
}
`, random_suffix, random_suffix, random_suffix, schedule, start_time, end_time, letter)
}

func testAccBigqueryDataTransferConfig_booleanParam(random_suffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_bigquery_dataset" "source_dataset" {
  dataset_id    = "source_%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_dataset" "destination_dataset" {
  dataset_id    = "destination_%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_data_transfer_config" "copy_config" {
  location = "asia-northeast1"

  display_name           = "Copy test %s"
  data_source_id         = "cross_region_copy"
  destination_dataset_id = google_bigquery_dataset.destination_dataset.dataset_id
  params = {
    overwrite_destination_table = "true"
    source_dataset_id           = google_bigquery_dataset.source_dataset.dataset_id
    source_project_id           = data.google_project.project.project_id
  }
}
`, random_suffix, random_suffix, random_suffix)
}

func testAccBigqueryDataTransferConfig_CMEK_basic(random_suffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
}

resource "google_kms_key_ring" "example_keyring" {
  name     = "keyring-test-%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "example_crypto_key" {
  name = "crypto-key-%s"
  key_ring = google_kms_key_ring.example_keyring.id
  purpose = "ENCRYPT_DECRYPT"
}

resource "google_service_account" "bqwriter%s" {
  account_id = "bqwriter%s"
}

resource "google_project_iam_member" "data_editor" {
  project = data.google_project.project.project_id

  role   = "roles/bigquery.dataEditor"
  member = "serviceAccount:${google_service_account.bqwriter%s.email}"
}

data "google_iam_policy" "owner" {
  binding {
    role = "roles/bigquery.dataOwner"

    members = [
      "serviceAccount:${google_service_account.bqwriter%s.email}",
    ]
  }
}

resource "google_bigquery_dataset_iam_policy" "dataset" {
  dataset_id  = google_bigquery_dataset.my_dataset.dataset_id
  policy_data = data.google_iam_policy.owner.policy_data
}

resource "google_bigquery_data_transfer_config" "query_config" {
  depends_on = [ google_kms_crypto_key.example_crypto_key ]
  encryption_configuration {
    kms_key_name = google_kms_crypto_key.example_crypto_key.id
  }
  display_name           = "my-query-%s"
  location               = "us-central1"
  data_source_id         = "scheduled_query"
  schedule               = "first sunday of quarter 00:00"
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  params = {
    destination_table_name_template = google_bigquery_table.my_table.table_id
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT name FROM table WHERE x = 'y'"
  }
}

resource "google_bigquery_dataset" "my_dataset" {
  dataset_id    = "my_dataset_%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "us-central1"
}

resource "google_bigquery_table" "my_table" {
  deletion_protection = false

  dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  table_id   = "my_table"
  schema     = <<EOF
  [
    { "name": "name", "type": "STRING" },
    { "name": "x", "type": "INTEGER" }
  ]
  EOF
}
`, random_suffix, random_suffix, random_suffix, random_suffix, random_suffix, random_suffix, random_suffix, random_suffix)
}

func testAccBigqueryDataTransferConfig_update_params_force_new(random_suffix, path, table string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id       = "tf_test_%s"
  friendly_name    = "foo"
  description      = "bar"
  location         = "US"
}

resource "google_bigquery_table" "my_table" {
  deletion_protection = false

  dataset_id = google_bigquery_dataset.dataset.dataset_id
  table_id   = "the-table-%s-%s"
  schema     = <<EOF
  [
    { "name": "name", "type": "STRING" },
    { "name": "x", "type": "INTEGER" }
  ]
  EOF
}

resource "google_bigquery_data_transfer_config" "update_config" {
  display_name           = "tf-test-%s"
  data_source_id         = "google_cloud_storage"
  destination_dataset_id = google_bigquery_dataset.dataset.dataset_id
  location               = google_bigquery_dataset.dataset.location

  params = {
    data_path_template              = "gs://bq-bucket-%s-%s/*.json"
    destination_table_name_template = google_bigquery_table.my_table.table_id
    file_format                     = "JSON"
    max_bad_records                 = 0
    write_disposition               = "APPEND"
  }
}
`, random_suffix, random_suffix, table, random_suffix, random_suffix, path)
}

func testAccBigqueryDataTransferConfig_scheduledQuery_updateServiceAccount(random_suffix string, service_account string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "bqwriter%s" {
  account_id = "bqwriter%s"
}

resource "google_project_iam_member" "data_editor" {
  project = data.google_project.project.project_id

  role   = "roles/bigquery.dataEditor"
  member = "serviceAccount:${google_service_account.bqwriter%s.email}"
}

resource "google_bigquery_dataset" "my_dataset" {
  dataset_id    = "my_dataset%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_table" "my_table" {
  deletion_protection = false

  dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  table_id   = "my_table"
  schema     = <<EOF
  [
    { "name": "name", "type": "STRING" },
    { "name": "x", "type": "INTEGER" }
  ]
  EOF
}

resource "google_bigquery_data_transfer_config" "query_config" {
  depends_on = [google_project_iam_member.data_editor]

  display_name           = "my-query-%s"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule               = "every 15 minutes"
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  service_account_name   = google_service_account.bqwriter%s.email
  params = {
    destination_table_name_template = google_bigquery_table.my_table.table_id
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT 1 AS a"
  }
}
`, service_account, service_account, service_account, random_suffix, random_suffix, service_account)
}

func testAccBigqueryDataTransferConfig_salesforce(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id       = "tf_test_%s"
  friendly_name    = "foo"
  description      = "bar"
  location         = "US"
}

resource "google_bigquery_data_transfer_config" "salesforce_config" {
  display_name           = "tf-test-%s"
  data_source_id         = "salesforce"
  destination_dataset_id = google_bigquery_dataset.dataset.dataset_id
  location               = google_bigquery_dataset.dataset.location

  params = {
    "connector.authentication.oauth.clientId"     = ""
    "connector.authentication.oauth.clientSecret" = ""
    "connector.authentication.oauth.myDomain"     = "MyDomain"
    "assets"                                      = "[\"asset-a\",\"asset-b\"]"
  }
}
`, randomSuffix, randomSuffix)
}
