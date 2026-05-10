package monitoring_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccMonitoringDashboard_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringDashboard_basic(),
			},
			{
				ResourceName:      "google_monitoring_dashboard.dashboard",
				ImportState:       true,
				ImportStateVerify: true,
				// Default import format uses the ID, which contains the project #
				// Testing import formats with the project name don't work because we set
				// the ID on import to what the user specified, which won't match the ID
				// from the apply
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccMonitoringDashboard_gridLayout(t *testing.T) {
	// TODO: Fix requires a breaking change https://github.com/hashicorp/terraform-provider-google/issues/9976
	t.Skip()
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringDashboard_gridLayout(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccMonitoringDashboard_rowLayout(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringDashboard_rowLayout(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccMonitoringDashboard_update(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringDashboard_rowLayout(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccMonitoringDashboard_basic(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccMonitoringDashboard_gridLayout(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccMonitoringDashboard_gridLayoutUpdate(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func testAccCheckMonitoringDashboardDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_monitoring_dashboard" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{MonitoringBasePath}}v1/{{name}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:               config,
				Method:               "GET",
				RawURL:               url,
				UserAgent:            config.UserAgent,
				ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
			})
			if err == nil {
				return fmt.Errorf("MonitoringDashboard still exists at %s", url)
			}
		}

		return nil
	}
}

func TestAccMonitoringDashboard_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "mon-dash-tagkey", map[string]interface{}{})
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "mon-dash-tagvalue", tagKey),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringDashboard_tags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_monitoring_dashboard.dashboard", "tags.%"),
					testAccCheckMonitoringDashboardHasTagBindings(t, "google_monitoring_dashboard.dashboard"),
				),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "tags"},
			},
		},
	})
}

func testAccMonitoringDashboard_basic() string {
	return fmt.Sprintf(`
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Demo Dashboard",
  "gridLayout": {
    "widgets": [
      {
        "blank": {}
      }
    ]
  }
}

EOF
}
`)
}

func testAccMonitoringDashboard_gridLayout() string {
	return fmt.Sprintf(`
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Grid Layout Example",
  "gridLayout": {
    "columns": "2",
    "widgets": [
      {
        "title": "Widget 1",
        "xyChart": {
          "dataSets": [{
            "timeSeriesQuery": {
              "timeSeriesFilter": {
                "filter": "metric.type=\"agent.googleapis.com/nginx/connections/accepted_count\"",
                "aggregation": {
                  "perSeriesAligner": "ALIGN_RATE"
                }
              },
              "unitOverride": "1"
            },
            "plotType": "LINE"
          }],
          "timeshiftDuration": "0s",
          "yAxis": {
            "label": "y1Axis",
            "scale": "LINEAR"
          }
        }
      },
      {
        "text": {
          "content": "Widget 2",
          "format": "MARKDOWN"
        }
      },
      {
        "title": "Widget 3",
        "xyChart": {
          "dataSets": [{
            "timeSeriesQuery": {
              "timeSeriesFilter": {
                "filter": "metric.type=\"agent.googleapis.com/nginx/connections/accepted_count\"",
                "aggregation": {
                  "perSeriesAligner": "ALIGN_RATE"
                }
              },
              "unitOverride": "1"
            },
            "plotType": "STACKED_BAR"
          }],
          "timeshiftDuration": "0s",
          "yAxis": {
            "label": "y1Axis",
            "scale": "LINEAR"
          }
        }
      }
    ]
  }
}

EOF
}
`)
}

func testAccMonitoringDashboard_gridLayoutUpdate() string {
	return fmt.Sprintf(`
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Grid Layout Example",
  "gridLayout": {
    "columns": "2",
    "widgets": [
      {
        "title": "Widget 1",
        "xyChart": {
          "dataSets": [{
            "timeSeriesQuery": {
              "timeSeriesFilter": {
                "filter": "metric.type=\"agent.googleapis.com/nginx/connections/accepted_count\"",
                "aggregation": {
                  "perSeriesAligner": "ALIGN_RATE"
                }
              },
              "unitOverride": "1"
            },
            "plotType": "LINE"
          }],
          "timeshiftDuration": "0s",
          "yAxis": {
            "label": "y1Axis",
            "scale": "LINEAR"
          }
        }
      },
      {
        "text": {
          "content": "Widget 2",
          "format": "MARKDOWN"
        }
      }
    ]
  }
}

EOF
}
`)
}

func testAccMonitoringDashboard_rowLayout() string {
	return fmt.Sprintf(`
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Row Layout Example",
  "rowLayout": {
    "rows": [
      {
        "weight": "1",
        "widgets": [
          {
            "text": {
              "content": "Widget 1",
              "format": "MARKDOWN"
            }
          },
          {
            "text": {
              "content": "Widget 3",
              "format": "MARKDOWN"
            }
          },
          {
            "text": {
              "content": "Widget 2",
              "format": "MARKDOWN"
            }
          }
        ]
      }
    ]
  }
}

EOF
}
`)
}

func testAccMonitoringDashboard_tags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Dashboard with Tags %{random_suffix}",
  "gridLayout": {
    "widgets": [
      {
        "blank": {}
      }
    ]
  }
}
EOF
  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

func testAccCheckMonitoringDashboardHasTagBindings(t *testing.T, resourceName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		config := acctest.GoogleProviderConfig(t)
		var configuredTagValueNamespacedName string
		for key, val := range rs.Primary.Attributes {
			if strings.HasPrefix(key, "tags.") && key != "tags.%" {
				tagKeyNamespacedName := strings.TrimPrefix(key, "tags.")
				tagValueShortName := val
				if tagValueShortName != "" {
					configuredTagValueNamespacedName = fmt.Sprintf("%s/%s", tagKeyNamespacedName, tagValueShortName)
					break
				}
			}
		}

		if configuredTagValueNamespacedName == "" {
			return fmt.Errorf("could not find a configured tag value in the state for resource %s", rs.Primary.ID)
		}

		parentURL := fmt.Sprintf("//monitoring.googleapis.com/%s", rs.Primary.ID)
		listBindingsURL := fmt.Sprintf("https://cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", url.QueryEscape(parentURL))

		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    listBindingsURL,
			UserAgent: config.UserAgent,
		})

		if err != nil {
			t.Logf("[EXPECTED FAILURE] Error calling TagBindings API: %v", err)
			return nil
		}

		return fmt.Errorf("TagBindings API call succeeded for unsupported resource type")
	}
}
