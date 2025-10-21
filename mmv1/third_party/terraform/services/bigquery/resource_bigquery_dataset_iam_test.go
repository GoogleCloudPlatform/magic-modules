package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigqueryDatasetIamBinding(t *testing.T) {
	t.Parallel()

	dataset := "tf_test_dataset_iam_" + acctest.RandString(t, 10)
	account := "tf-test-bq-iam-" + acctest.RandString(t, 10)
	role := "roles/bigquery.dataViewer"

	importId := fmt.Sprintf("projects/%s/datasets/%s %s",
		envvar.GetTestProjectFromEnv(), dataset, role)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigqueryDatasetIamBinding_basic(dataset, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigquery_dataset_iam_binding.binding", "role", role),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Binding update
				Config: testAccBigqueryDatasetIamBinding_update(dataset, account, role),
			},
			{
				ResourceName:      "google_bigquery_dataset_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigqueryDatasetIamMember(t *testing.T) {
	t.Parallel()

	dataset := "tf_test_dataset_iam_" + acctest.RandString(t, 10)
	account := "tf-test-bq-iam-" + acctest.RandString(t, 10)
	role := "roles/editor"

	importId := fmt.Sprintf("projects/%s/datasets/%s %s serviceAccount:%s",
		envvar.GetTestProjectFromEnv(),
		dataset,
		role,
		envvar.ServiceAccountCanonicalEmail(account))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigqueryDatasetIamMember(dataset, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigquery_dataset_iam_member.member", "role", role),
					resource.TestCheckResourceAttr(
						"google_bigquery_dataset_iam_member.member", "member", "serviceAccount:"+envvar.ServiceAccountCanonicalEmail(account)),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset_iam_member.member",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigqueryDatasetIamPolicy(t *testing.T) {
	t.Parallel()

	dataset := "tf_test_dataset_iam_" + acctest.RandString(t, 10)
	account := "tf-test-bq-iam-" + acctest.RandString(t, 10)
	role := "roles/bigquery.dataOwner"

	importId := fmt.Sprintf("projects/%s/datasets/%s",
		envvar.GetTestProjectFromEnv(), dataset)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigqueryDatasetIamPolicy(dataset, account, role),
				Check:  resource.TestCheckResourceAttrSet("data.google_bigquery_dataset_iam_policy.policy", "policy_data"),
			},
			{
				ResourceName:      "google_bigquery_dataset_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigqueryDatasetIamBindingWithIAMCondition(t *testing.T) {
	t.Parallel()

	dataset := "tf_test_dataset_iam_" + acctest.RandString(t, 10)
	account := "tf-test-bq-iam-" + acctest.RandString(t, 10)
	role := "roles/bigquery.dataViewer"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatasetIamBindingWithIAMCondition(dataset, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_binding.binding", "condition.0.title", "Expire access on 2050-12-31"),
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_binding.binding", "condition.0.description", "This condition will automatically remove access after 2050-12-31"),
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_binding.binding", "condition.0.expression", "request.time < timestamp('2050-12-31T23:59:59Z')"),
				),
			},
			{
				Config: testAccBigqueryDatasetIamBindingWithIAMCondition_update(dataset, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_bigquery_dataset_iam_binding.binding", "members.1"),
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_binding.binding", "condition.0.title", "Expire access on 2040-12-31"),
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_binding.binding", "condition.0.description", "This condition will automatically remove access after 2040-12-31"),
					resource.TestCheckResourceAttr("google_bigquery_dataset_iam_binding.binding", "condition.0.expression", "request.time < timestamp('2040-12-31T23:59:59Z')"),
				),
			},
		},
	})
}

func TestAccBigqueryDatasetIamPolicyWithIAMCondition(t *testing.T) {
	t.Parallel()

	owner := "tf-test-" + acctest.RandString(t, 10)
	dataset := "tf_test_dataset_iam_" + acctest.RandString(t, 10)
	account := "tf-test-bq-iam-" + acctest.RandString(t, 10)
	role := "roles/bigquery.dataViewer"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatasetIamPolicyWithIAMCondition(dataset, owner, account, role),
				Check:  resource.TestCheckResourceAttrSet("data.google_bigquery_dataset_iam_policy.policy", "policy_data"),
			},
		},
	})
}

func testAccBigqueryDatasetIamBinding_basic(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Bigquery Dataset Iam Testing Account"
}

resource "google_bigquery_dataset_iam_binding" "binding" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role     = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
  ]
}
`, dataset, account, account, role)
}

func testAccBigqueryDatasetIamBinding_update(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_bigquery_dataset_iam_binding" "binding" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role     = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
    "serviceAccount:${google_service_account.test-account2.email}",
  ]
}
`, dataset, account, account, role)
}

func testAccBigqueryDatasetIamMember(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_bigquery_dataset_iam_member" "member" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role     = "%s"
  member   = "serviceAccount:${google_service_account.test-account.email}"
}
`, dataset, account, role)
}

func testAccBigqueryDatasetIamPolicy(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Bigquery Dataset IAM Testing Account"
}

data "google_iam_policy" "policy" {
  binding {
    role    = "%s"
    members = ["serviceAccount:${google_service_account.test-account.email}"]
  }
}

resource "google_bigquery_dataset_iam_policy" "policy" {
  dataset_id  = google_bigquery_dataset.dataset.dataset_id
  policy_data = data.google_iam_policy.policy.policy_data
}

data "google_bigquery_dataset_iam_policy" "policy" {
  dataset_id  = google_bigquery_dataset.dataset.dataset_id
}
`, dataset, account, role)
}

var testBigqueryDatasetIam = `
resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}
`

func testAccBigqueryDatasetIamBindingWithIAMCondition(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Bigquery Dataset Iam Testing Account"
}

resource "google_bigquery_dataset_iam_binding" "binding" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role     = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
  ]

  condition {
    title       = "Expire access on 2050-12-31"
    description = "This condition will automatically remove access after 2050-12-31"
		expression  = "request.time < timestamp('2050-12-31T23:59:59Z')"
  }
}
`, dataset, account, account, role)
}

func testAccBigqueryDatasetIamBindingWithIAMCondition_update(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Bigquery Dataset Iam Testing Account"
}

resource "google_bigquery_dataset_iam_binding" "binding" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role     = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
    "serviceAccount:${google_service_account.test-account2.email}",
  ]

  condition {
    title       = "Expire access on 2040-12-31"
    description = "This condition will automatically remove access after 2040-12-31"
		expression  = "request.time < timestamp('2040-12-31T23:59:59Z')"
  }
}
`, dataset, account, account, role)
}

func testAccBigqueryDatasetIamPolicyWithIAMCondition(dataset, owner, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "owner" {
  account_id   = "%s"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Bigquery Dataset IAM Testing Account"
}

data "google_iam_policy" "policy" {
	binding {
		role    = "roles/bigquery.dataOwner"
		members = ["serviceAccount:${google_service_account.owner.email}"]
	}

  binding {
    role    = "%s"
    members = ["serviceAccount:${google_service_account.test-account.email}"]

		condition {
			title       = "Expire access on 2050-12-31"
			description = "This condition will automatically remove access after 2050-12-31"
			expression  = "request.time < timestamp('2050-12-31T23:59:59Z')"
		}
  }
}

resource "google_bigquery_dataset_iam_policy" "policy" {
  dataset_id  = google_bigquery_dataset.dataset.dataset_id
  policy_data = data.google_iam_policy.policy.policy_data
}

data "google_bigquery_dataset_iam_policy" "policy" {
  dataset_id  = google_bigquery_dataset.dataset.dataset_id
}
`, dataset, owner, account, role)
}
