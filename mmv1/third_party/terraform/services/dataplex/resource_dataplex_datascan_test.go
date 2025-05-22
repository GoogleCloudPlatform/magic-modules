package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataplexDatascanDataplexDatascanFullQuality_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexDatascanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexDatascanDataplexDatascanFullQuality_full(context),
			},
			{
				ResourceName:            "google_dataplex_datascan.full_quality",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_scan_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccDataplexDatascanDataplexDatascanFullQuality_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataplex_datascan.full_quality", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_dataplex_datascan.full_quality",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_scan_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexDatascanDataplexDatascanFullQuality_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_datascan" "full_quality" {
  location = "us-central1"
  display_name = "Full Datascan Quality"
  data_scan_id = "tf-test-dataquality-full%{random_suffix}"
  description = "Example resource - Full Datascan Quality"
  labels = {
    author = "billing"
  }

  data {
    resource = "//bigquery.googleapis.com/projects/bigquery-public-data/datasets/austin_bikeshare/tables/bikeshare_stations"
  }

  execution_spec {
    trigger {
      schedule {
        cron = "TZ=America/New_York 1 1 * * *"
      }
    }
    field = "modified_date"
  }

  data_quality_spec {
    sampling_percent = 5
    row_filter = "station_id > 1000"
    post_scan_actions {
      notification_report {
        recipients {
          emails = ["jane.doe@example.com"]
        }
        score_threshold_trigger {
          score_threshold = 86
        }
      }
    }
    
    rules {
      column = "address"
      dimension = "VALIDITY"
      threshold = 0.99
      non_null_expectation {}
    }

    rules {
      column = "council_district"
      dimension = "VALIDITY"
      ignore_null = true
      threshold = 0.9
      range_expectation {
        min_value = 1
        max_value = 10
        strict_min_enabled = true
        strict_max_enabled = false
      }
    }

    rules {
      column = "power_type"
      dimension = "VALIDITY"
      ignore_null = false
      regex_expectation {
        regex = ".*solar.*"
      }
    }

    rules {
      column = "property_type"
      dimension = "VALIDITY"
      ignore_null = false
      set_expectation {
        values = ["sidewalk", "parkland"]
      }
    }


    rules {
      column = "address"
      dimension = "UNIQUENESS"
      uniqueness_expectation {}
    }

    rules {
      column = "number_of_docks"
      dimension = "VALIDITY"
      statistic_range_expectation {
        statistic = "MEAN"
        min_value = 5
        max_value = 15
        strict_min_enabled = true
        strict_max_enabled = true
      }
    }

    rules {
      column = "footprint_length"
      dimension = "VALIDITY"
      row_condition_expectation {
        sql_expression = "footprint_length > 0 AND footprint_length <= 10"
      }
    }

    rules {
      dimension = "VALIDITY"
      table_condition_expectation {
        sql_expression = "COUNT(*) > 0"
      }
    }

    rules {
      dimension = "VALIDITY"
      sql_assertion {
        sql_statement = "select * from bigquery-public-data.austin_bikeshare.bikeshare_stations where station_id is null"
      }
    }
  }


  project = "%{project_name}"
}
`, context)
}

func testAccDataplexDatascanDataplexDatascanFullQuality_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_datascan" "full_quality" {
  location = "us-central1"
  display_name = "Full Datascan Quality"
  data_scan_id = "tf-test-dataquality-full%{random_suffix}"
  description = "Example resource - Full Datascan Quality"
  labels = {
    author = "billing"
  }

  data {
    resource = "//bigquery.googleapis.com/projects/bigquery-public-data/datasets/austin_bikeshare/tables/bikeshare_stations"
  }

  execution_spec {
    trigger {
      schedule {
        cron = "TZ=America/New_York 1 1 * * *"
      }
    }
    field = "modified_date"
  }

  data_quality_spec {
    sampling_percent = 5
    row_filter = "station_id > 1000"
    catalog_publishing_enabled = true
    post_scan_actions {
      notification_report {
        recipients {
          emails = ["jane.doe@example.com"]
        }
        score_threshold_trigger {
          score_threshold = 86
        }
      }
    }
    
    rules {
      column = "address"
      dimension = "VALIDITY"
      threshold = 0.99
      non_null_expectation {}
    }

    rules {
      column = "council_district"
      dimension = "VALIDITY"
      ignore_null = true
      threshold = 0.9
      range_expectation {
        min_value = 1
        max_value = 10
        strict_min_enabled = true
        strict_max_enabled = false
      }
    }

    rules {
      column = "power_type"
      dimension = "VALIDITY"
      ignore_null = false
      regex_expectation {
        regex = ".*solar.*"
      }
    }

    rules {
      column = "property_type"
      dimension = "VALIDITY"
      ignore_null = false
      set_expectation {
        values = ["sidewalk", "parkland"]
      }
    }


    rules {
      column = "address"
      dimension = "UNIQUENESS"
      uniqueness_expectation {}
    }

    rules {
      column = "number_of_docks"
      dimension = "VALIDITY"
      statistic_range_expectation {
        statistic = "MEAN"
        min_value = 5
        max_value = 15
        strict_min_enabled = true
        strict_max_enabled = true
      }
    }

    rules {
      column = "footprint_length"
      dimension = "VALIDITY"
      row_condition_expectation {
        sql_expression = "footprint_length > 0 AND footprint_length <= 10"
      }
    }

    rules {
      dimension = "VALIDITY"
      table_condition_expectation {
        sql_expression = "COUNT(*) > 0"
      }
    }

    rules {
      dimension = "VALIDITY"
      sql_assertion {
        sql_statement = "select * from bigquery-public-data.austin_bikeshare.bikeshare_stations where station_id is null"
      }
    }
  }


  project = "%{project_name}"
}
`, context)
}
