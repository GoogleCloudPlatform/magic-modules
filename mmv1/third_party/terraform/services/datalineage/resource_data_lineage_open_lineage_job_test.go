// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datalineage_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/datacatalog/lineage/apiv1/lineagepb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/datalineage"
	_ "github.com/hashicorp/terraform-provider-google/google/services/dataplex"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/googleapi"
)

var (
	_ = fmt.Sprintf
	_ = log.Print
	_ = strconv.Atoi
	_ = strings.Trim
	_ = time.Now
	_ = resource.TestMain
	_ = terraform.NewState
	_ = envvar.TestEnvVar
	_ = tpgresource.SetLabels
	_ = transport_tpg.Config{}
	_ = googleapi.Error{}
	_ = datalineage.Product
)

func TestAccDataLineageOpenLineageJob_dataLineageOpenLineageJobSimpleExample(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckDataLineageOpenLineageJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLineageOpenLineageJob_dataLineageOpenLineageJobSimpleExample(context),
			},
			{
				ResourceName: "google_data_lineage_open_lineage_job.simple",
				RefreshState: true,
			},
		},
	})
}

func testAccDataLineageOpenLineageJob_dataLineageOpenLineageJobSimpleExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_lineage_open_lineage_job" "simple" {
  namespace   = "example_simple_namespace"
  name        = "example_simple_name"
  description = "Nightly ETL from raw to curated"

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset_simple/source_table_1"
  }

  output {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target_simple/target_table_1"
  }
}
`, context)
}

func TestAccDataLineageOpenLineageJob_dataLineageOpenLineageJobWithFacetsExample(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckDataLineageOpenLineageJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLineageOpenLineageJob_dataLineageOpenLineageJobWithFacetsExample(context),
			},
			{
				ResourceName: "google_data_lineage_open_lineage_job.with_facets",
				RefreshState: true,
			},
		},
	})
}

func testAccDataLineageOpenLineageJob_dataLineageOpenLineageJobWithFacetsExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_lineage_open_lineage_job" "with_facets" {
  namespace   = "example_with_facets_namespace"
  name        = "example_with_facets_name"
  description = "Nightly ETL from raw to curated"

  owner {
    name = "team:data-engineering"
    type = "MAINTAINER"
  }

  input {
    namespace = "gs://example-bucket/"
    name      = "warehouse/raw_dataset_with_facets/source_table_1"

    symlink {
      namespace = "bigquery"
      name      = "my-project-name.raw_dataset_with_facets.source_table_1"
      type      = "TABLE"
    }

    catalog {
      framework = "bigquery"
      type      = "TABLE"
      name      = "my-project-name"
    }
  }

  output {
    namespace = "gs://example-bucket/"
    name      = "warehouse/target_with_facets/target_table_1"

    symlink {
      namespace = "bigquery"
      name      = "my-project-name.target_dataset_with_facets.target_table_1"
      type      = "TABLE"
    }

    catalog {
      framework = "bigquery"
      type      = "TABLE"
      name      = "my-project-name"
    }

    column_lineage {
      field {
        name = "user_id"
        input {
          namespace = "gs://example-bucket/"
          name      = "warehouse/raw_dataset_with_facets/source_table_1"
          field     = "id"
          transformation {
            type    = "DIRECT"
            subtype = "IDENTITY"
          }
        }
      }
    }
  }
}
`, context)
}

func testAccCheckDataLineageOpenLineageJobDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_data_lineage_open_lineage_job" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}
			conf := acctest.GoogleProviderConfig(t)
			ctx := context.Background()
			client, err := datalineage.LineageClientFromConfig(ctx, conf)
			if err != nil {
				return err
			}

			n := rs.Primary.Attributes["knowledge_catalog.0.process"]
			_, pErr := client.GetProcess(ctx, &lineagepb.GetProcessRequest{
				Name: n,
			})

			if pErr == nil {
				return fmt.Errorf("DataLineageOpenLineageJob still exists at %s", n)
			}
		}

		return nil
	}
}
