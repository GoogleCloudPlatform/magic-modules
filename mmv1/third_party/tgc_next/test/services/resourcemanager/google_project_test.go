package resourcemanager_test

import (
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/test"
)

func TestAccProject_labels(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"billing_account",
			"auto_create_network",
			"deletion_policy",
			"tags",
		},
	)
}

func TestAccProject_parentFolder(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"billing_account",
			"auto_create_network",
			"deletion_policy",
			"tags",
		},
	)
}

func TestAccProject_abandon(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"billing_account",
			"auto_create_network",
			"deletion_policy",
			"tags",
		},
	)
}

func TestAccProject_create(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"billing_account",
			"auto_create_network",
			"deletion_policy",
			"tags",
		},
	)
}

func TestAccProject_deleteDefaultNetwork(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"billing_account",
			"auto_create_network",
			"deletion_policy",
			"tags",
		},
	)
}

func TestAccProject_billing(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"billing_account",
			"auto_create_network",
			"deletion_policy",
			"tags",
		},
	)
}

func TestAccProject_migrateParent(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"billing_account",
			"auto_create_network",
			"deletion_policy",
			"tags",
		},
	)
}

func TestAccProject_noAllowDestroy(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"billing_account",
			"auto_create_network",
			"deletion_policy",
			"tags",
		},
	)
}

// func TestAccProject_tags(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		[]string{
// 			"billing_account",
// 			"auto_create_network",
// 			"deletion_policy",
// 			"tags",
// 		},
// 	)
// }
