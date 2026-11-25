package bigquery_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	bigquery_tpg "github.com/hashicorp/terraform-provider-google/google/services/bigquery"
	_ "github.com/hashicorp/terraform-provider-google/google/services/bigqueryreservation"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/tags"
	"google.golang.org/api/bigquery/v2"
	"regexp"
	"strings"
	"testing"
)

func TestAccBigQueryDataset_basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_withoutLabels(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "labels.%"),
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "effective_labels.%"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDataset(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.default_table_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.default_table_expiration_ms", "3600000"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The labels field in the state is decided by the configuration.
				// During importing, the configuration is unavailable, so the labels field in the state after importing is empty.
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.env", "bar"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.default_table_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.env", "bar"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.default_table_expiration_ms", "7200000"),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated2(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_withoutLabels(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "labels.%"),
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "effective_labels.%"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryDataset_withComputedLabels(t *testing.T) {
	// Skip it in VCR test because of the randomness of uuid in "labels" field
	// which causes the replaying mode after recording mode failing in VCR test
	acctest.SkipIfVcr(t)
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.default_table_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.%", "3"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.default_table_expiration_ms", "3600000"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The labels field in the state is decided by the configuration.
				// During importing, the configuration is unavailable, so the labels field in the state after importing is empty.
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated_withComputedLabels(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_withProvider5(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.75.0", // a version that doesn't separate user defined labels and system labels
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccBigQueryDataset_withoutLabelsV4(datasetID),
				ExternalProviders: oldVersion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "labels.%"),
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "effective_labels.%"),
				),
			},
			{
				Config:                   testAccBigQueryDataset(datasetID),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.default_table_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.default_table_expiration_ms", "3600000"),
				),
			},
		},
	})
}

func TestAccBigQueryDataset_withOutOfBandLabels(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset(datasetID),
				Check:  addOutOfBandLabels(t, datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_contents_on_destroy", "labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_contents_on_destroy", "labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated_withOutOfBandLabels(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_contents_on_destroy", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_datasetWithContents(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetDeleteContents(datasetID),
				Check:  testAccAddTable(t, datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_dataset.contents_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_contents_on_destroy", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_access(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_access_%s", acctest.RandString(t, 10))
	otherDatasetID := fmt.Sprintf("tf_test_other_%s", acctest.RandString(t, 10))
	otherTableID := fmt.Sprintf("tf_test_other_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetWithOneAccess(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetWithThreeAccess(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetWithOneAccess(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetWithConditionAccess(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetWithViewAccess(datasetID, otherDatasetID, otherTableID),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_accessMixedCase_userByEmail(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_access_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_accessMixedCase(datasetID, "user_by_email", "alicE@google.COM"),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_accessMixedCase_groupByEmail(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_access_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_accessMixedCase(datasetID, "group_by_email", "MAGIC-MODULES@gOOgle.com"),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_regionalLocation(t *testing.T) {
	t.Parallel()

	datasetID1 := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryRegionalDataset(datasetID1, "asia-south1"),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_cmek(t *testing.T) {
	t.Parallel()

	bootstrapped := kms.BootstrapKMSKeyInLocation(t, "us")
	pid := envvar.GetTestProjectFromEnv()
	datasetID1 := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_cmek(pid, datasetID1, bootstrapped.CryptoKey.Name),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryDataset_storageBillModel(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetStorageBillingModel(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_invalidCharacterInID(t *testing.T) {
	t.Parallel()
	// Not an acceptance test.
	acctest.SkipIfVcr(t)

	datasetID := fmt.Sprintf("tf_test_%s-with-hyphens", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigQueryDataset(datasetID),
				ExpectError: regexp.MustCompile("must contain only letters.+numbers.+or underscores.+"),
			},
		},
	})
}

func TestAccBigQueryDataset_invalidLongID(t *testing.T) {
	t.Parallel()
	// Not an acceptance test.
	acctest.SkipIfVcr(t)

	datasetSuffix := acctest.RandString(t, 10)
	datasetID := fmt.Sprintf("tf_test_%s", strings.Repeat(datasetSuffix, 200))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigQueryDataset(datasetID),
				ExpectError: regexp.MustCompile(".+cannot be greater than 1,024 characters"),
			},
		},
	})
}

func TestAccBigQueryDataset_bigqueryDatasetResourceTags_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_bigqueryDatasetResourceTags_basic(context),
			},
			{
				ResourceName:            "google_bigquery_dataset.dataset",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_bigqueryDatasetResourceTags_update(context),
			},
			{
				ResourceName:            "google_bigquery_dataset.dataset",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_bigqueryDatasetExternalReferenceAws(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_bigqueryDatasetExternalReferenceAws(context),
			},
			{
				ResourceName:            "google_bigquery_dataset.dataset",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "etag", "last_modified_time"},
			},
		},
	})
}

func TestAccBigQueryDataset_externalCatalogDatasetOptions_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_externalCatalogDatasetOptions_basic(context),
			},
			{
				ResourceName:            "google_bigquery_dataset.dataset",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_externalCatalogDatasetOptions_update(context),
			},
			{
				ResourceName:            "google_bigquery_dataset.dataset",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_collationUpdate(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_defaultCollation(datasetID, "und:ci"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "default_collation", "und:ci"),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_defaultCollation(datasetID, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "default_collation", ""),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_defaultRoundingModeUpdate(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_defaultRoundingMode(datasetID, "ROUND_HALF_EVEN"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "default_rounding_mode", "ROUND_HALF_EVEN"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The labels and terraform_labels fields in the state are decided by the configuration
				// (including provider default attribution labels) and are not populated in state on import.
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_defaultRoundingMode(datasetID, "ROUND_HALF_AWAY_FROM_ZERO"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "default_rounding_mode", "ROUND_HALF_AWAY_FROM_ZERO"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The labels and terraform_labels fields in the state are decided by the configuration
				// (including provider default attribution labels) and are not populated in state on import.
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccAddTable(t *testing.T, datasetID string, tableID string) resource.TestCheckFunc {
	// Not actually a check, but adds a table independently of terraform
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		table := &bigquery.Table{
			TableReference: &bigquery.TableReference{
				DatasetId: datasetID,
				TableId:   tableID,
				ProjectId: config.Project,
			},
		}
		_, err := bigquery_tpg.NewClient(config, config.UserAgent).Tables.Insert(config.Project, datasetID, table).Do()
		if err != nil {
			return fmt.Errorf("Could not create table")
		}
		return nil
	}
}

func addOutOfBandLabels(t *testing.T, datasetID string) resource.TestCheckFunc {
	// Not actually a check, but adds labels independently of terraform
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		dataset, err := bigquery_tpg.NewClient(config, config.UserAgent).Datasets.Get(config.Project, datasetID).Do()
		if err != nil {
			return fmt.Errorf("Could not get dataset with ID %s", datasetID)
		}

		dataset.Labels["outband_key"] = "test"
		_, err = bigquery_tpg.NewClient(config, config.UserAgent).Datasets.Patch(config.Project, datasetID, dataset).Do()
		if err != nil {
			return fmt.Errorf("Could not update labele for the dataset")
		}
		return nil
	}
}

func testAccBigQueryDataset_withoutLabels(datasetID string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000
}
`, datasetID)
}

func testAccBigQueryDataset_withoutLabelsV4(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000
}
`, datasetID)
}

func testAccBigQueryDataset(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetUpdated(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "bar"
  description                     = "This is a bar description"
  location                        = "EU"
  default_partition_expiration_ms = 7200000
  default_table_expiration_ms     = 7200000

  labels = {
    env                         = "bar"
    default_table_expiration_ms = 7200000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetUpdated_withOutOfBandLabels(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "bar"
  description                     = "This is a bar description"
  location                        = "EU"
  default_partition_expiration_ms = 7200000
  default_table_expiration_ms     = 7200000

  labels = {
    env                         = "bar"
    default_table_expiration_ms = 7200000
	outband_key                 = "test-update"
  }
}
`, datasetID)
}

func testAccBigQueryDatasetUpdated2(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  # friendly_name                   = "bar"
  description                     = "This is a bar description"
  location                        = "EU"
  default_partition_expiration_ms = 7200000
  default_table_expiration_ms     = 7200000

  labels = {
    env                         = "bar"
    default_table_expiration_ms = 7200000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetUpdated_withComputedLabels(datasetID string) string {
	return fmt.Sprintf(`
resource "random_uuid" "test" {
}

resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  # friendly_name                   = "bar"
  description                     = "This is a bar description"
  location                        = "EU"
  default_partition_expiration_ms = 7200000
  default_table_expiration_ms     = 7200000

  labels = {
    env                         = "${random_uuid.test.result}"
    default_table_expiration_ms = 7200000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetDeleteContents(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "contents_test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000
  delete_contents_on_destroy      = true

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryRegionalDataset(datasetID string, location string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                  = "%s"
  friendly_name               = "foo"
  description                 = "This is a foo description"
  location                    = "%s"
  default_table_expiration_ms = 3600000

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID, location)
}

func testAccBigQueryDatasetWithOneAccess(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetWithConditionAccess(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }

  access {
    role          = "READER"
    user_by_email = "Joe@example.com"
    condition {
        title       = "test-condition"
        description = "Request after midnight of 2019-12-31"
        expression  = "request.time > timestamp(\"2020-01-01T00:00:00Z\")"
        location    = "any.file.anywhere"
    }
  }

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetWithThreeAccess(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }
  access {
    role   = "READER"
    domain = "hashicorp.com"
  }
  access {
    role       = "READER"
    iam_member = "allUsers"
  }

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetWithViewAccess(datasetID, otherDatasetID, otherTableID string) string {
	// Note that we have to add a non-view access to prevent BQ from creating 4 default
	// access entries.
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "other_dataset" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "table_with_view" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.other_dataset.dataset_id

  time_partitioning {
    type = "DAY"
  }

  view {
    query          = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"
    use_legacy_sql = true
  }
}

resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }
  access {
    view {
      project_id = google_bigquery_dataset.other_dataset.project
      dataset_id = google_bigquery_dataset.other_dataset.dataset_id
      table_id   = google_bigquery_table.table_with_view.table_id
    }
  }

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, otherDatasetID, otherTableID, datasetID)
}

func testAccBigQueryDataset_accessMixedCase(datasetID, accessType, email string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role   = "OWNER"
    %s = "%s"
  }
}
`, datasetID, accessType, email)
}

func testAccBigQueryDataset_cmek(pid, datasetID, kmsKey string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_kms_crypto_key_iam_member" "kms-member" {
  crypto_key_id = "%s"
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "test" {
  dataset_id                  = "%s"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000

  default_encryption_configuration {
    kms_key_name = "%s"
  }

  depends_on = [google_kms_crypto_key_iam_member.kms-member]
}
`, pid, kmsKey, datasetID, kmsKey)
}

func testAccBigQueryDatasetStorageBillingModel(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000
  storage_billing_model           = "PHYSICAL"

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryDataset_bigqueryDatasetResourceTags_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_tags_tag_key" "tag_key1" {
  parent     = data.google_project.project.id
  short_name = "tf_test_tag_key1%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value1" {
  parent = google_tags_tag_key.tag_key1.id
  short_name = "tf_test_tag_value1%{random_suffix}"
}

resource "google_tags_tag_key" "tag_key2" {
  parent     = data.google_project.project.id
  short_name = "tf_test_tag_key2%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value2" {
  parent     = google_tags_tag_key.tag_key2.id
  short_name = "tf_test_tag_value2%{random_suffix}"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id                  = "tf_test_dataset%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"

  resource_tags = {
    (google_tags_tag_key.tag_key1.namespaced_name) = google_tags_tag_value.tag_value1.short_name
    (google_tags_tag_key.tag_key2.namespaced_name) = google_tags_tag_value.tag_value2.short_name
  }
}
`, context)
}

func testAccBigQueryDataset_bigqueryDatasetResourceTags_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_tags_tag_key" "tag_key1" {
  parent     = data.google_project.project.id
  short_name = "tf_test_tag_key1%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value1" {
  parent     = google_tags_tag_key.tag_key1.id
  short_name = "tf_test_tag_value1%{random_suffix}"
}

resource "google_tags_tag_key" "tag_key2" {
  parent     = data.google_project.project.id
  short_name = "tf_test_tag_key2%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value2" {
  parent     = google_tags_tag_key.tag_key2.id
  short_name = "tf_test_tag_value2%{random_suffix}"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id                  = "tf_test_dataset%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"

  resource_tags = {
  }
}
`, context)
}

func testAccBigQueryDataset_externalCatalogDatasetOptions_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id    = "tf_test_dataset%{random_suffix}"
  friendly_name = "test"
  description   = "This is a test description"
  location      = "US"

  external_catalog_dataset_options {
    parameters = {
      "dataset_owner" = "dataset_owner"
    }
    default_storage_location_uri = "gs://test_dataset/tables"
  }
}
`, context)
}

func testAccBigQueryDataset_externalCatalogDatasetOptions_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id    = "tf_test_dataset%{random_suffix}"
  friendly_name = "test"
  description   = "This is a test description"
  location      = "US"

  external_catalog_dataset_options {
    parameters = {
      "new_dataset_owner" = "new_dataset_owner"
    }
    default_storage_location_uri = "gs://new_test_dataset/new_tables"
  }
}
`, context)
}

func testAccBigQueryDataset_bigqueryDatasetExternalReferenceAws(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id                  = "tf_test_dataset%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "aws-us-east-1"

  external_dataset_reference {
    external_source = "aws-glue://arn:aws:glue:us-east-1:772042918353:database/db_other_formats_external"
    connection      = "projects/bigquerytestdefault/locations/aws-us-east-1/connections/external_test-connection"
  }
}
`, context)
}

func testAccBigQueryDataset_defaultCollation(datasetID, collation string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id        = "%s"
  friendly_name     = "foo"
  description       = "This is a foo description"
  location          = "US"
  default_collation = "%s"
}
`, datasetID, collation)
}

func testAccBigQueryDataset_defaultRoundingMode(datasetID, roundingMode string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id            = "%s"
  friendly_name         = "foo"
  description           = "This is a foo description"
  location              = "US"
  default_rounding_mode = "%s"
}
`, datasetID, roundingMode)
}

func TestAccBigQueryDataset_replicasLifecycle(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_crr_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_replicasNone(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "0"),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_replicasTwo(datasetID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "2"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 2, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_replicasThree(datasetID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "3"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 3, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_replicasPrimaryOnly(datasetID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "1"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 1, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_replicasCmek(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()
	datasetID := fmt.Sprintf("tf_test_crr_cmek_%s", acctest.RandString(t, 10))
	kmsCentral1 := kms.BootstrapKMSKeyInLocation(t, "us-central1")
	kmsEast4 := kms.BootstrapKMSKeyInLocation(t, "us-east4")
	kmsEast4Alt := kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-east4", "tf-bootstrap-bq-crr-us-east4-alt")
	kmsWest1 := kms.BootstrapKMSKeyInLocation(t, "us-west1")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create CMEK dataset without replicas first so the primary replica
				// inherits default_encryption_configuration on update (BigQuery requires
				// replica_options on all replicas including primary during datasets.insert).
				Config: testAccBigQueryDataset_replicasCmekNone(pid, datasetID, kmsCentral1.CryptoKey.Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "0"),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_replicasCmekTwo(pid, datasetID, kmsCentral1.CryptoKey.Name, kmsEast4.CryptoKey.Name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "2"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 2, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "replicas"},
			},
			{
				Config: testAccBigQueryDataset_replicasCmekThree(pid, datasetID, kmsCentral1.CryptoKey.Name, kmsEast4.CryptoKey.Name, kmsWest1.CryptoKey.Name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "3"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 3, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "replicas"},
			},
			{
				Config: testAccBigQueryDataset_replicasCmekPrimaryOnly(pid, datasetID, kmsCentral1.CryptoKey.Name, kmsEast4.CryptoKey.Name, kmsWest1.CryptoKey.Name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "1"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 1, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				// Populate write-only replica_options on the existing primary replica ("" -> key)
				// and re-add the secondary replica in-place without forcing replacement.
				Config: testAccBigQueryDataset_replicasCmekTwoWithPrimaryKey(pid, datasetID, kmsCentral1.CryptoKey.Name, kmsEast4.CryptoKey.Name, kmsEast4Alt.CryptoKey.Name, kmsWest1.CryptoKey.Name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "2"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 2, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "replicas"},
			},
			{
				// Mutating a non-empty default_kms_key_name on an existing replica (2 -> 2)
				// forces replacement and exercises single-step CMEK dataset creation with replicas.
				Config: testAccBigQueryDataset_replicasCmekTwoKeyReplaced(pid, datasetID, kmsCentral1.CryptoKey.Name, kmsEast4.CryptoKey.Name, kmsEast4Alt.CryptoKey.Name, kmsWest1.CryptoKey.Name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "2"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 2, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "replicas"},
			},
		},
	})
}

func TestAccBigQueryDataset_failoverReservation(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()
	suffix := acctest.RandString(t, 10)
	datasetID := fmt.Sprintf("tf_test_dr_%s", suffix)
	resName1 := fmt.Sprintf("tf-test-dr-res1-%s", suffix)
	resName2 := fmt.Sprintf("tf-test-dr-res2-%s", suffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_failoverReservationRes1(datasetID, resName1, resName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "failover_reservation_id", fmt.Sprintf("%s.%s", pid, resName1)),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "2"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 2, true),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_failoverReservationRes2(datasetID, resName1, resName2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "failover_reservation_id", fmt.Sprintf("%s.%s", pid, resName2)),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "2"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 2, true),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				// Simultaneously disassociate failover_reservation_id and drop the secondary replica
				// (exercises pre_update disassociation before PUT replica drop and GetRawConfig fallback).
				Config: testAccBigQueryDataset_failoverReservationDisassociatedPrimaryOnly(datasetID, resName1, resName2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "failover_reservation_id", ""),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "1"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 1, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				// Simultaneously add back the secondary replica and associate failover_reservation_id
				// (exercises post_update replica backfill wait before PATCH failover association).
				Config: testAccBigQueryDataset_failoverReservationRes1(datasetID, resName1, resName2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "failover_reservation_id", fmt.Sprintf("%s.%s", pid, resName1)),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "2"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 2, true),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_failoverReservationDisassociated(datasetID, resName1, resName2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigquery_dataset.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "failover_reservation_id", ""),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "replicas.#", "2"),
					testAccCheckBigQueryDatasetReplicasMetadata("google_bigquery_dataset.test", 2, false),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccCheckBigQueryDatasetReplicasMetadata(resourceName string, expectedCount int, expectFailoverURI bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}
		attrs := rs.Primary.Attributes
		var primaryCount int
		var secondaryCount int
		seenIndices := map[string]bool{}
		for k := range attrs {
			if !strings.HasPrefix(k, "replicas.") {
				continue
			}
			parts := strings.Split(k, ".")
			if len(parts) < 3 || parts[1] == "#" {
				continue
			}
			seenIndices[parts[1]] = true
		}
		if len(seenIndices) != expectedCount {
			return fmt.Errorf("expected %d replicas in state for %s, got %d", expectedCount, resourceName, len(seenIndices))
		}
		for idx := range seenIndices {
			prefix := "replicas." + idx + "."
			state := attrs[prefix+"replica_state"]
			if attrs[prefix+"creation_time"] == "" {
				return fmt.Errorf("expected non-empty creation_time at %s", prefix)
			}
			if attrs[prefix+"backfill_completion_time"] == "" {
				return fmt.Errorf("expected non-empty backfill_completion_time at %s (state=%q)", prefix, state)
			}
			switch state {
			case "PRIMARY":
				primaryCount++
				if attrs[prefix+"primary_assignment_time"] == "" {
					return fmt.Errorf("expected non-empty primary_assignment_time on PRIMARY replica at %s", prefix)
				}
				if attrs[prefix+"primary_assignment_completion_time"] == "" {
					return fmt.Errorf("expected non-empty primary_assignment_completion_time on PRIMARY replica at %s", prefix)
				}
			case "SECONDARY":
				secondaryCount++
			default:
				return fmt.Errorf("unexpected replica_state %q at %s", state, prefix)
			}
			if expectFailoverURI && attrs[prefix+"failover_reservation_uri"] == "" {
				return fmt.Errorf("expected non-empty failover_reservation_uri at %s", prefix)
			}
			if !expectFailoverURI && attrs[prefix+"failover_reservation_uri"] != "" {
				return fmt.Errorf("expected empty failover_reservation_uri at %s, got %q", prefix, attrs[prefix+"failover_reservation_uri"])
			}
		}
		if primaryCount != 1 || secondaryCount != expectedCount-1 {
			return fmt.Errorf("expected 1 PRIMARY and %d SECONDARY replicas, got %d PRIMARY and %d SECONDARY", expectedCount-1, primaryCount, secondaryCount)
		}
		return nil
	}
}

func testAccBigQueryDataset_replicasNone(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR lifecycle test"
  location    = "us-central1"
}
`, datasetID)
}

func testAccBigQueryDataset_replicasTwo(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR lifecycle test"
  location    = "us-central1"

  replicas {
    id       = "us-central1-primary"
    location = "us-central1"
  }

  replicas {
    id       = "us-east4-secondary"
    location = "us-east4"
  }
}
`, datasetID)
}

func testAccBigQueryDataset_replicasThree(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR lifecycle test"
  location    = "us-central1"

  replicas {
    id       = "us-central1-primary"
    location = "us-central1"
  }

  replicas {
    id       = "us-east4-secondary"
    location = "us-east4"
  }

  replicas {
    id       = "us-west1-secondary"
    location = "us-west1"
  }
}
`, datasetID)
}

func testAccBigQueryDataset_replicasPrimaryOnly(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR lifecycle test"
  location    = "us-central1"

  replicas {
    id       = "us-central1-primary"
    location = "us-central1"
  }
}
`, datasetID)
}

func testAccBigQueryDataset_replicasCmekNone(pid, datasetID, kmsCentral1 string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_kms_crypto_key_iam_member" "kms_central1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR CMEK test"
  location    = "us-central1"

  default_encryption_configuration {
    kms_key_name = "%s"
  }

  depends_on = [
    google_kms_crypto_key_iam_member.kms_central1,
  ]
}
`, pid, kmsCentral1, datasetID, kmsCentral1)
}

func testAccBigQueryDataset_replicasCmekTwo(pid, datasetID, kmsCentral1, kmsEast4 string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_kms_crypto_key_iam_member" "kms_central1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_east4" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR CMEK test"
  location    = "us-central1"

  default_encryption_configuration {
    kms_key_name = "%s"
  }

  replicas {
    id       = "us-central1"
    location = "us-central1"
  }

  replicas {
    id       = "us-east4"
    location = "us-east4"
    replica_options {
      default_kms_key_name = "%s"
    }
  }

  depends_on = [
    google_kms_crypto_key_iam_member.kms_central1,
    google_kms_crypto_key_iam_member.kms_east4,
  ]
}
`, pid, kmsCentral1, kmsEast4, datasetID, kmsCentral1, kmsEast4)
}

func testAccBigQueryDataset_replicasCmekThree(pid, datasetID, kmsCentral1, kmsEast4, kmsWest1 string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_kms_crypto_key_iam_member" "kms_central1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_east4" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_west1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR CMEK test"
  location    = "us-central1"

  default_encryption_configuration {
    kms_key_name = "%s"
  }

  replicas {
    id       = "us-central1"
    location = "us-central1"
  }

  replicas {
    id       = "us-east4"
    location = "us-east4"
    replica_options {
      default_kms_key_name = "%s"
    }
  }

  replicas {
    id       = "us-west1"
    location = "us-west1"
    replica_options {
      default_kms_key_name = "%s"
    }
  }

  depends_on = [
    google_kms_crypto_key_iam_member.kms_central1,
    google_kms_crypto_key_iam_member.kms_east4,
    google_kms_crypto_key_iam_member.kms_west1,
  ]
}
`, pid, kmsCentral1, kmsEast4, kmsWest1, datasetID, kmsCentral1, kmsEast4, kmsWest1)
}

func testAccBigQueryDataset_replicasCmekPrimaryOnly(pid, datasetID, kmsCentral1, kmsEast4, kmsWest1 string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_kms_crypto_key_iam_member" "kms_central1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_east4" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_west1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR CMEK test"
  location    = "us-central1"

  default_encryption_configuration {
    kms_key_name = "%s"
  }

  replicas {
    id       = "us-central1"
    location = "us-central1"
  }

  depends_on = [
    google_kms_crypto_key_iam_member.kms_central1,
    google_kms_crypto_key_iam_member.kms_east4,
    google_kms_crypto_key_iam_member.kms_west1,
  ]
}
`, pid, kmsCentral1, kmsEast4, kmsWest1, datasetID, kmsCentral1)
}

func testAccBigQueryDataset_replicasCmekTwoWithPrimaryKey(pid, datasetID, kmsCentral1, kmsEast4, kmsEast4Alt, kmsWest1 string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_kms_crypto_key_iam_member" "kms_central1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_east4" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_east4_alt" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_west1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR CMEK test"
  location    = "us-central1"

  default_encryption_configuration {
    kms_key_name = "%s"
  }

  replicas {
    id       = "us-central1"
    location = "us-central1"
    replica_options {
      default_kms_key_name = "%s"
    }
  }

  replicas {
    id       = "us-east4"
    location = "us-east4"
    replica_options {
      default_kms_key_name = "%s"
    }
  }

  depends_on = [
    google_kms_crypto_key_iam_member.kms_central1,
    google_kms_crypto_key_iam_member.kms_east4,
    google_kms_crypto_key_iam_member.kms_east4_alt,
    google_kms_crypto_key_iam_member.kms_west1,
  ]
}
`, pid, kmsCentral1, kmsEast4, kmsEast4Alt, kmsWest1, datasetID, kmsCentral1, kmsCentral1, kmsEast4)
}

func testAccBigQueryDataset_replicasCmekTwoKeyReplaced(pid, datasetID, kmsCentral1, kmsEast4, kmsEast4Alt, kmsWest1 string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_kms_crypto_key_iam_member" "kms_central1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_east4" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_east4_alt" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms_west1" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "CRR CMEK test"
  location    = "us-central1"

  default_encryption_configuration {
    kms_key_name = "%s"
  }

  replicas {
    id       = "us-central1"
    location = "us-central1"
    replica_options {
      default_kms_key_name = "%s"
    }
  }

  replicas {
    id       = "us-east4"
    location = "us-east4"
    replica_options {
      default_kms_key_name = "%s"
    }
  }

  depends_on = [
    google_kms_crypto_key_iam_member.kms_central1,
    google_kms_crypto_key_iam_member.kms_east4,
    google_kms_crypto_key_iam_member.kms_east4_alt,
    google_kms_crypto_key_iam_member.kms_west1,
  ]
}
`, pid, kmsCentral1, kmsEast4, kmsEast4Alt, kmsWest1, datasetID, kmsCentral1, kmsCentral1, kmsEast4Alt)
}

func testAccBigQueryDataset_failoverReservationRes1(datasetID, resName1, resName2 string) string {
	return fmt.Sprintf(`
resource "google_bigquery_reservation" "res1" {
  name               = "%s"
  location           = "us-central1"
  secondary_location = "us-east4"
  slot_capacity      = 0
  edition            = "ENTERPRISE_PLUS"
  ignore_idle_slots  = true
  concurrency        = 0
  autoscale {
    max_slots = 100
  }
}

resource "google_bigquery_reservation" "res2" {
  name               = "%s"
  location           = "us-central1"
  secondary_location = "us-east4"
  slot_capacity      = 0
  edition            = "ENTERPRISE_PLUS"
  ignore_idle_slots  = true
  concurrency        = 0
  autoscale {
    max_slots = 100
  }
}

resource "google_bigquery_dataset" "test" {
  dataset_id              = "%s"
  description             = "Managed DR dataset test"
  location                = "us-central1"
  failover_reservation_id = "${google_bigquery_reservation.res1.project}.${google_bigquery_reservation.res1.name}"

  replicas {
    id       = "us-central1"
    location = "us-central1"
  }

  replicas {
    id       = "us-east4"
    location = "us-east4"
  }
}
`, resName1, resName2, datasetID)
}

func testAccBigQueryDataset_failoverReservationRes2(datasetID, resName1, resName2 string) string {
	return fmt.Sprintf(`
resource "google_bigquery_reservation" "res1" {
  name               = "%s"
  location           = "us-central1"
  secondary_location = "us-east4"
  slot_capacity      = 0
  edition            = "ENTERPRISE_PLUS"
  ignore_idle_slots  = true
  concurrency        = 0
  autoscale {
    max_slots = 100
  }
}

resource "google_bigquery_reservation" "res2" {
  name               = "%s"
  location           = "us-central1"
  secondary_location = "us-east4"
  slot_capacity      = 0
  edition            = "ENTERPRISE_PLUS"
  ignore_idle_slots  = true
  concurrency        = 0
  autoscale {
    max_slots = 100
  }
}

resource "google_bigquery_dataset" "test" {
  dataset_id              = "%s"
  description             = "Managed DR dataset test"
  location                = "us-central1"
  failover_reservation_id = "${google_bigquery_reservation.res2.project}.${google_bigquery_reservation.res2.name}"

  replicas {
    id       = "us-central1"
    location = "us-central1"
  }

  replicas {
    id       = "us-east4"
    location = "us-east4"
  }
}
`, resName1, resName2, datasetID)
}

func testAccBigQueryDataset_failoverReservationDisassociated(datasetID, resName1, resName2 string) string {
	return fmt.Sprintf(`
resource "google_bigquery_reservation" "res1" {
  name               = "%s"
  location           = "us-central1"
  secondary_location = "us-east4"
  slot_capacity      = 0
  edition            = "ENTERPRISE_PLUS"
  ignore_idle_slots  = true
  concurrency        = 0
  autoscale {
    max_slots = 100
  }
}

resource "google_bigquery_reservation" "res2" {
  name               = "%s"
  location           = "us-central1"
  secondary_location = "us-east4"
  slot_capacity      = 0
  edition            = "ENTERPRISE_PLUS"
  ignore_idle_slots  = true
  concurrency        = 0
  autoscale {
    max_slots = 100
  }
}

resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "Managed DR dataset test"
  location    = "us-central1"

  replicas {
    id       = "us-central1"
    location = "us-central1"
  }

  replicas {
    id       = "us-east4"
    location = "us-east4"
  }
}
`, resName1, resName2, datasetID)
}

func testAccBigQueryDataset_failoverReservationDisassociatedPrimaryOnly(datasetID, resName1, resName2 string) string {
	return fmt.Sprintf(`
resource "google_bigquery_reservation" "res1" {
  name               = "%s"
  location           = "us-central1"
  secondary_location = "us-east4"
  slot_capacity      = 0
  edition            = "ENTERPRISE_PLUS"
  ignore_idle_slots  = true
  concurrency        = 0
  autoscale {
    max_slots = 100
  }
}

resource "google_bigquery_reservation" "res2" {
  name               = "%s"
  location           = "us-central1"
  secondary_location = "us-east4"
  slot_capacity      = 0
  edition            = "ENTERPRISE_PLUS"
  ignore_idle_slots  = true
  concurrency        = 0
  autoscale {
    max_slots = 100
  }
}

resource "google_bigquery_dataset" "test" {
  dataset_id  = "%s"
  description = "Managed DR dataset test"
  location    = "us-central1"

  replicas {
    id       = "us-central1"
    location = "us-central1"
  }
}
`, resName1, resName2, datasetID)
}
