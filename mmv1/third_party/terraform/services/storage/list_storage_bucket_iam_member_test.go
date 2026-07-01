package storage_test

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

func TestAccStorageBucketIamMemberListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	bucket := "tf-test-bucket-iam-" + acctest.RandString(t, 10)
	account := "tf-test-storage-iam-" + acctest.RandString(t, 10)
	role := "roles/storage.objectViewer"
	member := "serviceAccount:" + envvar.ServiceAccountCanonicalEmail(account)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketIamMember(bucket, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_storage_bucket_iam_member.member", "bucket", "b/"+bucket),
					resource.TestCheckResourceAttr("google_storage_bucket_iam_member.member", "role", role),
					resource.TestCheckResourceAttr("google_storage_bucket_iam_member.member", "member", member),
				),
			},
			{
				Query:  true,
				Config: testAccStorageBucketIamMemberListQueryWithFilters(bucket, role, member),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("google_storage_bucket_iam_member.test", 1),
					querycheck.ExpectIdentity("google_storage_bucket_iam_member.test", map[string]knownvalue.Check{
						"bucket":          knownvalue.StringExact("b/" + bucket),
						"role":            knownvalue.StringExact(role),
						"member":          knownvalue.StringExact(member),
						"condition_title": knownvalue.Null(),
					}),
				},
			},
		},
	})
}

func testAccStorageBucketIamMember(bucket, account, role string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Storage Bucket IAM Testing Account"
}

resource "google_storage_bucket_iam_member" "member" {
  bucket = google_storage_bucket.bucket.name
  role   = "%s"
  member = "serviceAccount:${google_service_account.test-account.email}"
}
`, bucket, account, role)
}

func testAccStorageBucketIamMemberListQueryWithFilters(bucket, role, member string) string {
	return fmt.Sprintf(`
list "google_storage_bucket_iam_member" "test" {
  provider = google
  include_resource = true

  config {
    bucket = %q
    role   = %q
    member = %q
  }
}
`, bucket, role, member)
}
