package pubsub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/pubsub"
)

func TestAccDataSourceGooglePubsubSchema_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSchemaDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePubsubSchema_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_pubsub_schema.foo", "google_pubsub_schema.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGooglePubsubSchema_optionalProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSchemaDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePubsubSchema_optionalProject(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_pubsub_schema.foo", "google_pubsub_schema.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGooglePubsubSchema_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_schema" "foo" {
  name       = "tf-test-pubsub-schema-%{random_suffix}"
  type       = "PROTOCOL_BUFFER"
  definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\n}"
}

data "google_pubsub_schema" "foo" {
  name    = google_pubsub_schema.foo.name
  project = google_pubsub_schema.foo.project
}
`, context)
}

func testAccDataSourceGooglePubsubSchema_optionalProject(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_schema" "foo" {
  name       = "tf-test-pubsub-schema-%{random_suffix}"
  type       = "PROTOCOL_BUFFER"
  definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\n}"
}

data "google_pubsub_schema" "foo" {
  name = google_pubsub_schema.foo.name
}
`, context)
}
