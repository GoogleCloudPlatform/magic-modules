package chronicle_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccChronicleDashboardChart_chronicleDashboardchartFullExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"chronicle_id":  envvar.GetTestChronicleInstanceIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckChronicleDashboardChartDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccChronicleDashboardChart_chronicleDashboardchartFullExample_full(context),
			},
			{
				ResourceName:            "google_chronicle_dashboard_chart.my_chart",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"chart_layout", "dashboardUserData.0.lastViewedTime", "dashboard_query", "instance", "location", "native_dashboard"},
			},
			{
				Config: testAccChronicleDashboardChart_chronicleDashboardchartFullExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_chronicle_dashboard_chart.my_chart", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_chronicle_dashboard_chart.my_chart",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"chart_layout", "dashboardUserData.0.lastViewedTime", "dashboard_query", "instance", "location", "native_dashboard"},
			},
		},
	})
}

func testAccChronicleDashboardChart_chronicleDashboardchartFullExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
# A Native Dashboard is required to create a Dashboard Chart.
resource "google_chronicle_native_dashboard" "my_dashboard" {
  location     = "us" 
  instance     = "%{chronicle_id}"
  display_name = "tf-test-dashboard-1%{random_suffix}"
  description  = "tf-test-dashboard-description%{random_suffix}"
  access       = "DASHBOARD_PRIVATE"
  type         = "CUSTOM"

  # REFINED: The 'definition' wrapper is removed because flatten_object is true.
  filters {
    id                                    = "GlobalTimeFilter"
    display_name                          = "Global Time Filter"
    data_source                           = "GLOBAL"
    is_standard_time_range_filter         = true
    is_standard_time_range_filter_enabled = true
    filter_operator_and_field_values {
      filter_operator = "PAST"
      field_values    = ["1", "DAY"]
    }
  }
}
resource "google_chronicle_dashboard_chart" "my_chart" {
  location         = google_chronicle_native_dashboard.my_dashboard.location
  instance         = google_chronicle_native_dashboard.my_dashboard.instance
  native_dashboard = google_chronicle_native_dashboard.my_dashboard.name

  # Chart layout from JSON for "Data Source Health Overview"
  chart_layout {
    start_x = 0
    span_x  = 50
    start_y = 12
    span_y  = 18
  }

  # Filters applied to this chart, from JSON
  # filters_ids = ["GlobalTimeFilter"]

  dashboard_chart {
    display_name = "Data Source Health Overview" # From JSON
    description  = "Health of data sources over time"  # From JSON
    tile_type    = "TILE_TYPE_VISUALIZATION"     # From JSON

    chart_datasource {
      # When defining the query directly in dashboard_query, only data_sources is needed here.
      data_sources = ["INGESTION_METRICS"] # From JSON
    }

    # Visualization configuration from JSON
    visualization {
      series {
        series_type = "LINE"
        encode {
          x = "timestamp"
          y = "total_count"
        }
        data_label {}
      }

      x_axes {
        axis_type   = "CATEGORY"
        display_name = "Date"
      }

      y_axes {
        axis_type   = "CATEGORY"
        display_name = "Sources"
      }

      legends {
        top           = 12
        legend_orient = "HORIZONTAL"
      }

      series_column = ["health_status"]
      grouping_type = "Grouped"
      # Note: column_defs, table_config, and threshold_coloring_enabled are omitted
      # to match the provided JSON's visualization block.
    }

    # Drill down configuration from JSON
    drill_down_config {
      left_drill_downs {
        id          = "D89B834D-977A-4E0C-83B0-12AB1D05E76B"
        display_name = "Link to the google"
        custom_settings {
          new_tab = true
          external_link {
            link = "www.google.com"
          }
        }
      }
    }
  }

  # The Malachite Query Language (MQL) query content.
  # This block is used for "AddChart" to define the query within the chart resource.
  dashboard_query {
    query = <<EOT
$component = ingestion.component
$collector_id = ingestion.collector_id
$feed_id = ingestion.feed_id
$log_type = ingestion.log_type
$timestamp = timestamp.get_timestamp(ingestion.end_time, "%Y-%m-%d")
$health_status = ingestion.health_status
ingestion.log_type != ""
ingestion.component != ""
ingestion.component != "Normalizer"
ingestion.collector_id != "bbbb1111-bbbb-1111-bbbb-1111bbbb1111" // Risk Entity Collector ID
ingestion.collector_id != "01010101-1111-0101-1111-000111110000" // Prbr
ingestion.collector_id != "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa" // Out-of-Band
ingestion.collector_id != "aaaa2222-aaaa-2222-aaaa-2222aaaa2222" // HTTPS Push
ingestion.collector_id != "aaaa4444-aaaa-4444-aaaa-4444aaaa4444" // Event Hub

//Filtering GCP Native Ingestion as they will be captured individually
not (ingestion.component = "Ingestion API" and ingestion.collector_id = "aaaa3333-aaaa-3333-aaaa-3333aaaa3333") and
//Filtering Workspace Native Ingestion as they will be captured individually
not (ingestion.component = "Ingestion API" and ingestion.collector_id = "dddddddd-dddd-dddd-dddd-dddddddddddd")
($feed_id != "" or $collector_id != "")
$health_status != ""
match:
  $health_status,
  $timestamp
outcome:
  $feeds = count_distinct($feed_id) - 1
  $collectors = count_distinct(if($collector_id != "", strings.concat($collector_id, $log_type), "")) - 1
  $total_count = $collectors + $feeds
order:
  $timestamp asc
limit:
  10000
EOT
    input {
      # Matches the relativeTime from the original HCL.
      relative_time {
        time_unit      = "SECOND"
        start_time_val = "1"
      }
    }
  }
}
`, context)
}

func testAccChronicleDashboardChart_chronicleDashboardchartFullExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
# A Native Dashboard is required to create a Dashboard Chart.
resource "google_chronicle_native_dashboard" "my_dashboard" {
  location     = "us" 
  instance     = "%{chronicle_id}"
  display_name = "tf-test-dashboard-1%{random_suffix}"
  description  = "tf-test-dashboard-description%{random_suffix}"
  access       = "DASHBOARD_PRIVATE"
  type         = "CUSTOM"

  # REFINED: The 'definition' wrapper is removed because flatten_object is true.
  filters {
    id                                    = "GlobalTimeFilter"
    display_name                          = "Global Time Filter"
    data_source                           = "GLOBAL"
    is_standard_time_range_filter         = true
    is_standard_time_range_filter_enabled = true
    filter_operator_and_field_values {
      filter_operator = "PAST"
      field_values    = ["1", "DAY"]
    }
  }
}

resource "google_chronicle_dashboard_chart" "my_chart" {
  location         = google_chronicle_native_dashboard.my_dashboard.location
  instance         = google_chronicle_native_dashboard.my_dashboard.instance
  native_dashboard = google_chronicle_native_dashboard.my_dashboard.name

  # Chart layout from JSON for "Data Source Health Overview"
  chart_layout {
    start_x = 0
    span_x  = 50
    start_y = 12
    span_y  = 18
  }

  # Filters applied to this chart, from JSON
  # filters_ids = ["GlobalTimeFilter"]

  dashboard_chart {
    display_name = "updated_chart" # From JSON
    description  = "updated_description"  # From JSON
    tile_type    = "TILE_TYPE_VISUALIZATION"     # From JSON

    chart_datasource {
      # When defining the query directly in dashboard_query, only data_sources is needed here.
      data_sources = ["INGESTION_METRICS"] # From JSON
    }

    # Visualization configuration from JSON
    visualization {
      series {
        series_type = "LINE"
        encode {
          x = "updated_x_axis"
          y = "updated_y_axis"
        }
        data_label {}
      }

      x_axes {
        axis_type   = "CATEGORY"
        display_name = "Date_updated"
      }

      y_axes {
        axis_type   = "CATEGORY"
        display_name = "Sources_updated"
      }

      legends {
        top           = 15
        legend_orient = "VERTICAL"
      }

      series_column = ["health_status"]
      grouping_type = "Grouped"
      # Note: column_defs, table_config, and threshold_coloring_enabled are omitted
      # to match the provided JSON's visualization block.
    }

    # Drill down configuration from JSON
    drill_down_config {
      left_drill_downs {
        id          = "D89B834D-977A-4E0C-83B0-12AB1D05E76B"
        display_name = "Link to the google updated"
        custom_settings {
          new_tab = true
          external_link {
            link = "www.google3.com"
          }
        }
      }
    }
  }

  # The Malachite Query Language (MQL) query content.
  # This block is used for "AddChart" to define the query within the chart resource.
  dashboard_query {
    query = <<EOT
$component = ingestion.component
$collector_id = ingestion.collector_id
$log_type = ingestion.log_type
$timestamp = timestamp.get_timestamp(ingestion.end_time, "%Y-%m-%d")
$health_status = ingestion.health_status
ingestion.log_type != ""
ingestion.component != ""
ingestion.component != "Normalizer"
ingestion.collector_id != "bbbb1111-bbbb-1111-bbbb-1111bbbb1111" // Risk Entity Collector ID
ingestion.collector_id != "01010101-1111-0101-1111-000111110000" // Prbr
ingestion.collector_id != "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa" // Out-of-Band
ingestion.collector_id != "aaaa2222-aaaa-2222-aaaa-2222aaaa2222" // HTTPS Push
ingestion.collector_id != "aaaa4444-aaaa-4444-aaaa-4444aaaa4444" // Event Hub

//Filtering GCP Native Ingestion as they will be captured individually
not (ingestion.component = "Ingestion API" and ingestion.collector_id = "aaaa3333-aaaa-3333-aaaa-3333aaaa3333") and
//Filtering Workspace Native Ingestion as they will be captured individually
not (ingestion.component = "Ingestion API" and ingestion.collector_id = "dddddddd-dddd-dddd-dddd-dddddddddddd")
($feed_id != "" or $collector_id != "")
$health_status != ""
match:
  $health_status,
  $timestamp
outcome:
  $feeds = count_distinct($feed_id) - 1
  $collectors = count_distinct(if($collector_id != "", strings.concat($collector_id, $log_type), "")) - 1
  $total_count = $collectors + $feeds
order:
  $timestamp asc
limit:
  10000
EOT
    input {
      # Matches the relativeTime from the original HCL.
      relative_time {
        time_unit      = "SECOND"
        start_time_val = "2"
      }
    }
  }
}
`, context)
}
