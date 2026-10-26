package vertexai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
	"github.com/hashicorp/terraform-provider-google/google/services/vertexai"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func createEvaluationItem(t *testing.T, config *transport_tpg.Config, project, region, itemDisplayName string) string {
	rawUrl := fmt.Sprintf("%sprojects/%s/locations/%s/evaluationItems",
		transport_tpg.BaseUrl(vertexai.Product, config), project, region)
	url := strings.ReplaceAll(rawUrl, "{{region}}", region)

	body := map[string]interface{}{
		"displayName":        itemDisplayName,
		"evaluationItemType": "REQUEST",
		"evaluationRequest": map[string]interface{}{
			"prompt": map[string]interface{}{
				"text": "Test Prompt for " + itemDisplayName,
			},
		},
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    url,
		UserAgent: config.UserAgent,
		Headers:   headers,
		Body:      body,
	})
	if err != nil {
		t.Fatalf("Failed to create bootstrap evaluationItem %s: %s", itemDisplayName, err)
	}

	var parsed map[string]interface{}
	resBytes, err := json.Marshal(res)
	if err == nil {
		_ = json.Unmarshal(resBytes, &parsed)
	}

	fullName, ok := parsed["name"].(string)
	if !ok || fullName == "" {
		t.Fatalf("Created evaluationItem did not return a valid resource name: %v", res)
	}
	return fullName
}

func deleteEvaluationItem(config *transport_tpg.Config, project, region, fullName string) {
	if fullName == "" {
		return
	}
	rawUrl := fmt.Sprintf("%s%s", transport_tpg.BaseUrl(vertexai.Product, config), fullName)
	url := strings.ReplaceAll(rawUrl, "{{region}}", region)

	headers := make(http.Header)
	_, _ = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   project,
		RawURL:    url,
		UserAgent: config.UserAgent,
		Headers:   headers,
	})
}

func TestAccVertexAIEvaluationSet_basic(t *testing.T) {
	t.Parallel()

	acctest.AccTestPreCheck(t)
	config := acctest.GoogleProviderConfig(t)

	randomSuffix := acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()
	region := "us-central1"
	itemName1 := "tf-test-eval-item-1-" + randomSuffix
	itemName2 := "tf-test-eval-item-2-" + randomSuffix

	itemPath1 := createEvaluationItem(t, config, project, region, itemName1)
	itemPath2 := createEvaluationItem(t, config, project, region, itemName2)
	t.Cleanup(func() {
		deleteEvaluationItem(config, project, region, itemPath1)
		deleteEvaluationItem(config, project, region, itemPath2)
	})

	context := map[string]interface{}{
		"project":       project,
		"region":        region,
		"display_name":  "tf-test-eval-set-" + randomSuffix,
		"random_suffix": randomSuffix,
		"item_1":        itemPath1,
		"item_2":        itemPath2,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEvaluationSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEvaluationSet_step1(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "region", region),
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "display_name", "tf-test-eval-set-"+randomSuffix),
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "evaluation_items.#", "1"),
					resource.TestCheckResourceAttrSet("google_vertex_ai_evaluation_set.example", "name"),
					resource.TestCheckResourceAttrSet("google_vertex_ai_evaluation_set.example", "create_time"),
					resource.TestCheckResourceAttrSet("google_vertex_ai_evaluation_set.example", "update_time"),
				),
			},
			{
				ResourceName:            "google_vertex_ai_evaluation_set.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccVertexAIEvaluationSet_step2(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "region", region),
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "display_name", "tf-test-eval-set-"+randomSuffix),
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "evaluation_items.#", "2"),
				),
			},
			{
				ResourceName:            "google_vertex_ai_evaluation_set.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func TestAccVertexAIEvaluationSet_withEncryption(t *testing.T) {
	t.Parallel()

	acctest.AccTestPreCheck(t)
	config := acctest.GoogleProviderConfig(t)

	randomSuffix := acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()
	region := "us-central1"
	itemName := "tf-test-eval-item-enc-" + randomSuffix

	itemPath := createEvaluationItem(t, config, project, region, itemName)
	t.Cleanup(func() {
		deleteEvaluationItem(config, project, region, itemPath)
	})

	bootstrapped := kms.BootstrapKMSKeyInLocation(t, "us-central1")

	context := map[string]interface{}{
		"project":       project,
		"region":        region,
		"display_name":  "tf-test-eval-set-enc-" + randomSuffix,
		"random_suffix": randomSuffix,
		"item":          itemPath,
		"kms_key_name":  bootstrapped.CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEvaluationSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEvaluationSet_withEncryption(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "region", region),
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "display_name", "tf-test-eval-set-enc-"+randomSuffix),
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "evaluation_items.#", "1"),
					resource.TestCheckResourceAttr("google_vertex_ai_evaluation_set.example", "encryption_spec.0.kms_key_name", bootstrapped.CryptoKey.Name),
					resource.TestCheckResourceAttrSet("google_vertex_ai_evaluation_set.example", "name"),
					resource.TestCheckResourceAttrSet("google_vertex_ai_evaluation_set.example", "create_time"),
					resource.TestCheckResourceAttrSet("google_vertex_ai_evaluation_set.example", "update_time"),
				),
			},
			{
				ResourceName:            "google_vertex_ai_evaluation_set.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccVertexAIEvaluationSet_step1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_evaluation_set" "example" {
  region           = "%{region}"
  display_name     = "%{display_name}"
  evaluation_items = ["%{item_1}"]
}
`, context)
}

func testAccVertexAIEvaluationSet_step2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_evaluation_set" "example" {
  region           = "%{region}"
  display_name     = "%{display_name}"
  evaluation_items = ["%{item_1}", "%{item_2}"]
}
`, context)
}

func testAccVertexAIEvaluationSet_withEncryption(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_evaluation_set" "example" {
  region           = "%{region}"
  display_name     = "%{display_name}"
  evaluation_items = ["%{item}"]

  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }
}
`, context)
}

func testAccCheckVertexAIEvaluationSetDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vertex_ai_evaluation_set" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(vertexai.Product, config)+"projects/{{project}}/locations/{{region}}/evaluationSets/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("VertexAIEvaluationSet still exists at %s", url)
			}
		}

		return nil
	}
}
