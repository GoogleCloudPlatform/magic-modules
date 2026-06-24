package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigqueryDatasetIamMemberListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	dataset := "tf_test_dataset_iam_" + acctest.RandString(t, 10)
	account := "tf-test-bq-iam-" + acctest.RandString(t, 10)
	role := "roles/editor"
	member := "serviceAccount:" + envvar.ServiceAccountCanonicalEmail(account)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatasetIamMember(dataset, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_member.member", "project", project),
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_member.member", "dataset_id", dataset),
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_member.member", "role", role),
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_member.member", "member", member),
				),
			},
			{
				Query:  true,
				Config: testAccBigqueryDatasetIamMemberListQueryWithFilters(dataset, role, member),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("google_bigquery_dataset_iam_member.test", 1),
					querycheck.ExpectIdentity("google_bigquery_dataset_iam_member.test", map[string]knownvalue.Check{
						"project":         knownvalue.StringExact(project),
						"dataset_id":      knownvalue.StringExact(dataset),
						"role":            knownvalue.StringExact(role),
						"member":          knownvalue.StringExact(member),
						"condition_title": knownvalue.Null(),
					}),
				},
			},
		},
	})
}

func testAccBigqueryDatasetIamMemberListQueryWithFilters(datasetID, role, member string) string {
	return fmt.Sprintf(`
list "google_bigquery_dataset_iam_member" "test" {
  provider = google
  include_resource = true

  config {
    dataset_id = %q
    role       = %q
    member     = %q
  }
}
`, datasetID, role, member)
}
