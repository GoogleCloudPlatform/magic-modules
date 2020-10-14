package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigQueryRoutine_bigQueryRoutine_Update(t *testing.T) {
	t.Parallel()

	dataset := fmt.Sprintf("tfmanualdataset%s", randString(t, 10))
	routine := fmt.Sprintf("tfmanualroutine%s", randString(t, 10))

	body := "CREATE FUNCTION Add(x FLOAT64, y FLOAT64) RETURNS FLOAT64 AS (x + y);"
	body_updated := "CREATE FUNCTION Minus(x FLOAT64, y FLOAT64) RETURNS FLOAT64 AS (x + y);"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckBigQueryRoutineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryRoutine_bigQueryRoutine_Update(dataset, routine, body),
			},
			{
				ResourceName:      "google_bigquery_routine.sproc",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryRoutine_bigQueryRoutine_Update(dataset, routine, body_updated),
			},
			{
				ResourceName:      "google_bigquery_routine.sproc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigQueryRoutine_bigQueryRoutine_Update(dataset, routine, body string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_routine" "sproc" {
  dataset_id = google_bigquery_dataset.test.dataset_id
  routine_id     = "%s"
  routine_type = "PROCEDURE"
  language = "SQL"
  definition_body = "%s"
}
`, dataset, routine, body)
}
