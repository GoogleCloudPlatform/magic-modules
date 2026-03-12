package chronicle_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// TestAccChronicleNativeDashboard_chronicleNativedashboardUpdateExample tests updating a Native Dashboard and the layout/filter association of an existing Chart within it.
func TestAccChronicleNativeDashboard_chronicleNativedashboardUpdateExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"chronicle_id":       envvar.GetTestChronicleInstanceIdFromEnv(t),
		"random_suffix":      acctest.RandString(t, 10),
		"chart_display_name": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckChronicleNativeDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccChronicleNativeDashboard_chronicleNativedashboardBasicExample_basic(context),
			},
			{
				Config:                  testAccChronicleNativeDashboard_chronicleNativedashboardBasicExample_update(context),
				ImportState:             true,
				ImportStateVerify:       true,
				ResourceName:            "google_chronicle_native_dashboard.my_dashboard",
				ImportStateVerifyIgnore: []string{"last_viewed_time"},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_chronicle_native_dashboard.my_dashboard", plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

// testAccChronicleNativeDashboard_chronicleNativedashboardBasicExample_basic
func testAccChronicleNativeDashboard_chronicleNativedashboardBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_native_dashboard" "my_dashboard" {
  location     = "us"
  instance     = "%{chronicle_id}"
  display_name = "Initial Dashboard Name-%{random_suffix}"
  description  = "Description for update test - initial."
  access       = "DASHBOARD_PRIVATE"
  type         = "CUSTOM"

  # Flattened: 'filters' is now top-level
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
  location         = "us"
  instance         = "%{chronicle_id}"
  native_dashboard = google_chronicle_native_dashboard.my_dashboard.name

  chart_layout {
    span_x  = 42
    span_y  = 27
    start_x = 0
    start_y = 0
  }

  dashboard_chart {
    display_name = "Chart Name-%{chart_display_name}"
    description  = "Managed by Terraform AccTest"
    tile_type    = "TILE_TYPE_VISUALIZATION"
    chart_datasource {
      data_sources = ["INGESTION_METRICS"]
    }
    visualization {
      column_defs {
        field  = "event_count"
        header = "event_count"
      }
      x_axes { axis_type = "VALUE" }
      y_axes { axis_type = "VALUE" }
      table_config { enable_text_wrap = false }
      legends { legend_orient = "HORIZONTAL" }
      series_column = []
      grouping_type = "Off"
      threshold_coloring_enabled = false
    }
  }

  dashboard_query {
    query = "ingestion.component = \"Ingestion API\" \ningestion.log_type != \"\" \noutcome: $event_count = math.round(sum(ingestion.log_count)/(1000*1000), 2)\n"
    input {
      relative_time {
        time_unit      = "SECOND"
        start_time_val = "1"
      }
    }
  }
}
`, context)
}

// testAccChronicleNativeDashboard_chronicleNativedashboardBasicExample_update
func testAccChronicleNativeDashboard_chronicleNativedashboardBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_native_dashboard" "my_dashboard" {
  location     = "us"
  instance     = "%{chronicle_id}"
  display_name = "Updated Dashboard Name-%{random_suffix}"
  description  = "Updated Description for update test."
  access       = "DASHBOARD_PUBLIC"
  type         = "CUSTOM"

  # DashboardUserData
  is_pinned = true


  # Flattened: 'filters' and 'charts' are top-level
  filters {
    id                                    = "GlobalTimeFilter"
    display_name                          = "Updated Global Time Filter"
    data_source                           = "GLOBAL"
    is_standard_time_range_filter         = true
    is_standard_time_range_filter_enabled = true
    filter_operator_and_field_values {
      filter_operator = "PAST"
      field_values    = ["7", "DAY"]
    }
  }

  filters {
    id           = "HostnameFilter"
    display_name = "Hostname Filter"
    data_source  = "UDM"
    field_path   = "principal.hostname"
    filter_operator_and_field_values {
      filter_operator = "EQUAL"
      field_values    = ["test-host-1"]
    }
  }

  charts {
    dashboard_chart = google_chronicle_dashboard_chart.my_chart.name
    chart_layout {
      span_x  = 8
      span_y  = 6
      start_x = 0
      start_y = 0
    }
    filters_ids = ["GlobalTimeFilter", "HostnameFilter"]
  }
}

# Reactivate the chart resource so Terraform can manage its destruction
`, context)
}
