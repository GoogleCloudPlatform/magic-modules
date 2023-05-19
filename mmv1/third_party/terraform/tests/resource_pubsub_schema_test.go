package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccPubsubSchema_basic(t *testing.T) {
	t.Parallel()

	schema := fmt.Sprintf("tf-test-schema-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSchema_basic(schema),
			},
		},
	})
}

func TestAccPubsubSchema_multipleRevisions(t *testing.T) {
	t.Parallel()

	schema := fmt.Sprintf("tf-test-schema-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSchema_multipleRevisions(schema),
			},
		},
	})
}

func TestAccPubsubSchema_update(t *testing.T) {
	t.Parallel()

	schema := fmt.Sprintf("tf-test-schema-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			// Have a single revision.
			{
				Config: testAccPubsubSchema_basic(schema),
			},
			// Add a second revision.
			{
				Config: testAccPubsubSchema_multipleRevisions(schema),
			},
			// Swap the revisions, which will delete and recreate both of them.
			{
				Config: testAccPubsubSchema_swappedRevisions(schema),
			},
			// Remove the first revision.
			{
				Config: testAccPubsubSchema_basic(schema),
			},
		},
	})
}

func TestAccPubsubSchema_updateMaxRevisions(t *testing.T) {
	t.Parallel()

	schema := fmt.Sprintf("tf-test-schema-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			// Have 20 revisions.
			{
				Config: testAccPubsubSchema_maxRevisions(schema),
			},
			// Change the first revision, resulting in deleting and creating all revisions.
			{
				Config: testAccPubsubSchema_maxRevisionsUpdate(schema),
			},
		},
	})
}

func TestAccPubsubSchema_tooManyRevisions(t *testing.T) {
	t.Parallel()

	schema := fmt.Sprintf("tf-test-schema-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			// Have 20 revisions.
			{
				Config: testAccPubsubSchema_maxRevisions(schema),
			},
			// Exceed the allowed number of revisions
			{
				Config: testAccPubsubSchema_tooManyRevisions(schema),
				ExpectError: regexp.MustCompile("Cannot have more than 20 Schema revisions."),
			},
		},
	})
}

func TestAccPubsubSchema_tooFewRevisions(t *testing.T) {
	t.Parallel()

	schema := fmt.Sprintf("tf-test-schema-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			// Have 20 revisions.
			{
				Config: testAccPubsubSchema_maxRevisions(schema),
			},
			// Getting to zero revisions
			{
				Config: testAccPubsubSchema_zeroRevisions(schema),
				ExpectError: regexp.MustCompile("Must have at least one Schema revision."),
			},
		},
	})
}

func TestAccPubsubSchema_withDefinition(t *testing.T) {
	t.Parallel()

	schema := fmt.Sprintf("tf-test-schema-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSchema_withDefinition(schema),
			},
			{
				ResourceName:            "google_pubsub_schema.foo",
				ImportStateId:           schema,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func TestAccPubsubSchema_withDefinitionAndRevision(t *testing.T) {
	t.Parallel()

	schema := fmt.Sprintf("tf-test-schema-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSchema_withDefinitionAndRevision(schema),
			},
			{
				ResourceName:            "google_pubsub_schema.foo",
				ImportStateId:           schema,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func testAccPubsubSchema_basic(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
	}
`, schema)
}

func testAccPubsubSchema_multipleRevisions(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
	}
`, schema)
}

func testAccPubsubSchema_swappedRevisions(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
	}
`, schema)
}

func testAccPubsubSchema_withDefinition(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
  	definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
	}
`, schema)
}

func testAccPubsubSchema_withDefinitionAndRevision(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
  	definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"

		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
	}
`, schema)
}

func testAccPubsubSchema_maxRevisions(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
	}
`, schema)
}

func testAccPubsubSchema_maxRevisionsUpdate(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
	}
`, schema)
}

func testAccPubsubSchema_tooManyRevisions(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\nstring timestamp_response = 4;\n}"
		}
		revision {
			definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
		}
	}
`, schema)
}

func testAccPubsubSchema_zeroRevisions(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
	}
`, schema)
}
