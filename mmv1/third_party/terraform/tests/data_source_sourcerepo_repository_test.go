package google_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	google "internal/terraform-provider-google"
)

func TestAccDataSourceGoogleSourceRepoRepository_basic(t *testing.T) {
	t.Parallel()

	name := "tf-repository-" + google.RandString(t, 10)

	google.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { google.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: google.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSourceRepoRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleSourceRepoRepositoryConfig(name),
				Check: resource.ComposeTestCheckFunc(
					google.CheckDataSourceStateMatchesResourceState("data.google_sourcerepo_repository.bar", "google_sourcerepo_repository.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleSourceRepoRepositoryConfig(name string) string {
	return fmt.Sprintf(`
resource "google_sourcerepo_repository" "foo" {
  name               = "%s"
}

data "google_sourcerepo_repository" "bar" {
  name = google_sourcerepo_repository.foo.name
  depends_on = [
    google_sourcerepo_repository.foo,
  ]
}
`, name)
}
