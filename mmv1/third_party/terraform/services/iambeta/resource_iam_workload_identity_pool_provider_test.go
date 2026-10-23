package iambeta_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/iambeta"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIAMBetaWorkloadIdentityPoolProvider_aws(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolProviderDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProvider_aws_full(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProvider_aws_enabled(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProvider_aws_basic(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIAMBetaWorkloadIdentityPoolProvider_oidc(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolProviderDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProvider_oidc_full(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProvider_oidc_update(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProvider_oidc_basic(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIAMBetaWorkloadIdentityPoolProvider_aws_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "my_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.my_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-my-provider-%{random_suffix}"
  display_name                       = "Name of provider"
  description                        = "AWS identity pool provider for automated test"
  disabled                           = true
  attribute_condition                = "attribute.aws_role==\"arn:aws:sts::999999999999:assumed-role/stack-eu-central-1-lambdaRole\""
  attribute_mapping                  = {
    "google.subject"        = "assertion.arn"
    "attribute.aws_account" = "assertion.account"
    "attribute.environment" = "assertion.arn.contains(\":instance-profile/Production\") ? \"prod\" : \"test\""
  }
  aws {
    account_id = "999999999999"
  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPoolProvider_aws_enabled(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "my_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.my_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-my-provider-%{random_suffix}"
  display_name                       = "Name of provider"
  description                        = "AWS identity pool provider for automated test"
  disabled                           = false
  attribute_condition                = "attribute.aws_role==\"arn:aws:sts::999999999999:assumed-role/stack-eu-central-1-lambdaRole\""
  attribute_mapping                  = {
    "google.subject"        = "assertion.arn"
    "attribute.aws_account" = "assertion.account"
    "attribute.environment" = "assertion.arn.contains(\":instance-profile/Production\") ? \"prod\" : \"test\""
  }
  aws {
    account_id = "999999999999"
  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPoolProvider_oidc_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "my_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.my_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-my-provider-%{random_suffix}"
  display_name                       = "Name of provider"
  description                        = "OIDC identity pool provider for automated test"
  disabled                           = true
  attribute_condition                = "\"e968c2ef-047c-498d-8d79-16ca1b61e77e\" in assertion.groups"
  attribute_mapping                  = {
    "google.subject"                  = "\"azure::\" + assertion.tid + \"::\" + assertion.sub"
    "attribute.tid"                   = "assertion.tid"
    "attribute.managed_identity_name" = <<EOT
      {
        "8bb39bdb-1cc5-4447-b7db-a19e920eb111":"workload1",
        "55d36609-9bcf-48e0-a366-a3cf19027d2a":"workload2"
      }[assertion.oid]
EOT
  }
  oidc {
    allowed_audiences = ["https://example.com/gcp-oidc-federation", "example.com/gcp-oidc-federation"]
    issuer_uri        = "https://sts.windows.net/azure-tenant-id-full"
    jwks_json         = "{\"keys\":[{\"kty\":\"RSA\",\"alg\":\"RS256\",\"kid\":\"sif0AR-F6MuvksAyAOv-Pds08Bcf2eUMlxE30NofddA\",\"use\":\"sig\",\"e\":\"AQAB\",\"n\":\"ylH1Chl1tpfti3lh51E1g5dPogzXDaQseqjsefGLknaNl5W6Wd4frBhHyE2t41Q5zgz_Ll0-NvWm0FlaG6brhrN9QZu6sJP1bM8WPfJVPgXOanxi7d7TXCkeNubGeiLTf5R3UXtS9Lm_guemU7MxDjDTelxnlgGCihOVTcL526suNJUdfXtpwUsvdU6_ZnAp9IpsuYjCtwPm9hPumlcZGMbxstdh07O4y4O90cVQClJOKSGQjAUCKJWXIQ0cqffGS_HuS_725CPzQ85SzYZzaNpgfhAER7kx_9P16ARM3BJz0PI5fe2hECE61J4GYU_BY43sxDfs7HyJpEXKLU9eWw\"}]}"
  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPoolProvider_oidc_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "my_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.my_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-my-provider-%{random_suffix}"
  display_name                       = "Name of provider"
  description                        = "OIDC identity pool provider for automated test"
  disabled                           = true
  attribute_condition                = "\"e968c2ef-047c-498d-8d79-16ca1b61e77e\" in assertion.groups"
  attribute_mapping                  = {
    "google.subject"                  = "\"azure::\" + assertion.tid + \"::\" + assertion.sub"
    "attribute.tid"                   = "assertion.tid"
    "attribute.managed_identity_name" = <<EOT
      {
        "8bb39bdb-1cc5-4447-b7db-a19e920eb111":"workload1",
        "55d36609-9bcf-48e0-a366-a3cf19027d2a":"workload2"
      }[assertion.oid]
EOT
  }
  oidc {
    allowed_audiences = ["https://example.com/gcp-oidc-federation-update", "example.com/gcp-oidc-federation-update"]
    issuer_uri        = "https://sts.windows.net/azure-tenant-id-update"
    jwks_json         = "{\"keys\":[{\"kty\":\"RSA\",\"alg\":\"RS256\",\"kid\":\"F6MuvksAyAOv-Pds08Bcf2eUMlxE30NofddA\",\"use\":\"sig\",\"e\":\"AQAB\",\"n\":\"ylH1Chl1tpfti3lh51E1g5dPogzXDaQseqjsefGLknaNl5W6Wd4frBhHyE2t41Q5zgz_Ll0-NvWm0FlaG6brhrN9QZu6sJP1bM8WPfJVPgXOanxi7d7TXCkeNubGeiLTf5R3UXtS9Lm_guemU7MxDjDTelxnlgGCihOVTcL526suNJUdfXtpwUsvdU6_ZnAp9IpsuYjCtwPm9hPumlcZGMbxstdh07O4y4O90cVQClJOKSGQjAUCKJWXIQ0cqffGS_HuS_725CPzQ85SzYZzaNpgfhAER7kx_9P16ARM3BJz0PI5fe2hECE61J4GYU_BY43sxDfs7HyJpEXKLU9eWw\"}]}"

  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPoolProvider_aws_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "my_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.my_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-my-provider-%{random_suffix}"
  aws {
    account_id = "999999999999"
  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPoolProvider_oidc_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "my_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.my_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-my-provider-%{random_suffix}"
  attribute_mapping                  = {
	"google.subject"                  = "assertion.sub"
  }
  oidc {
    allowed_audiences = ["https://example.com/gcp-oidc-federation", "example.com/gcp-oidc-federation"]
    issuer_uri        = "https://sts.windows.net/azure-tenant-id-full"
  }
}
`, context)
}

func TestAccIAMBetaWorkloadIdentityPoolProvider_x509(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolProviderDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProvider_x509_full(context),
			},
			{
				ResourceName:            "google_iam_workload_identity_pool_provider.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workload_identity_pool_id", "workload_identity_pool_provider_id"},
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProvider_x509_update(context),
			},
			{
				ResourceName:            "google_iam_workload_identity_pool_provider.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workload_identity_pool_id", "workload_identity_pool_provider_id"},
			},
		},
	})
}

func testAccIAMBetaWorkloadIdentityPoolProvider_x509_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "pool" {
  workload_identity_pool_id = "tf-test-example-pool%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "example" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-example-prvdr%{random_suffix}"
  display_name                       = "Name of provider"
  description                        = "X.509 identity pool provider for automated test"
  disabled                           = true
  attribute_mapping                  = {
    "google.subject"        = "assertion.subject.dn.cn"
  }
  x509 {
    trust_store {
        trust_anchors {
            pem_certificate = file("test-fixtures/trust_anchor.pem")
        }
        intermediate_cas {
            pem_certificate = file("test-fixtures/intermediate_ca.pem")
        }
    }
  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPoolProvider_x509_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "pool" {
  workload_identity_pool_id = "tf-test-example-pool%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "example" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-example-prvdr%{random_suffix}"
  display_name                       = "Name of provider"
  description                        = "X.509 identity pool provider for automated test"
  disabled                           = true
  attribute_mapping                  = {
    "google.subject"        = "assertion.subject.dn.cn"
  }
  x509 {
    trust_store {
        trust_anchors {
            pem_certificate = file("test-fixtures/trust_anchor_updated.pem")
        }
        trust_anchors {
            pem_certificate = file("test-fixtures/intermediate_ca.pem")
        }
    }
  }
}
`, context)
}

func TestAccIAMBetaWorkloadIdentityPoolProvider_undeleteOnCreate(t *testing.T) {
	t.Parallel()

	context := map[string]any{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMBetaWorkloadIdentityPoolProviderDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProviderUndeleteInitial(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProviderUndeletePoolOnly(context),
			},
			{
				Config: testAccIAMBetaWorkloadIdentityPoolProviderUndeleteRecreated(context),
			},
			{
				ResourceName:      "google_iam_workload_identity_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIAMBetaWorkloadIdentityPoolProviderUndeleteInitial(context map[string]any) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "my_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.my_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-my-provider-%{random_suffix}"
  display_name                       = "Initial provider name"
  description                        = "Initial OIDC provider"
  disabled                           = true
  attribute_mapping                  = {
    "google.subject" = "assertion.sub"
  }
  oidc {
    issuer_uri = "https://sts.windows.net/initial-tenant-id"
  }
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPoolProviderUndeletePoolOnly(context map[string]any) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
}
`, context)
}

func testAccIAMBetaWorkloadIdentityPoolProviderUndeleteRecreated(context map[string]any) string {
	return acctest.Nprintf(`
resource "google_iam_workload_identity_pool" "my_pool" {
  workload_identity_pool_id = "tf-test-my-pool-%{random_suffix}"
}

resource "google_iam_workload_identity_pool_provider" "my_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.my_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "tf-test-my-provider-%{random_suffix}"
  display_name                       = "Recreated provider name"
  description                        = "Recreated OIDC provider"
  disabled                           = false
  attribute_mapping                  = {
    "google.subject" = "assertion.sub"
  }
  oidc {
    issuer_uri = "https://sts.windows.net/recreated-tenant-id"
  }
}
`, context)
}

func TestResourceIAMBetaWorkloadIdentityPoolProviderCreate_undeleteWhenDeleted(t *testing.T) {
	_ = t.Context()

	var requests []recordedRequest
	undeleted := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var body map[string]any
		if r.Body != nil {
			raw, _ := io.ReadAll(r.Body)
			if len(raw) > 0 {
				_ = json.Unmarshal(raw, &body)
			}
		}

		if strings.HasPrefix(r.URL.Path, "/v1/operations/") {
			_, _ = w.Write([]byte(`{"name": "` + strings.TrimPrefix(r.URL.Path, "/v1/") + `", "done": true}`))
			return
		}

		requests = append(requests, recordedRequest{
			Method:     r.Method,
			Path:       r.URL.Path,
			UpdateMask: r.URL.Query().Get("updateMask"),
			Body:       body,
		})

		switch {
		case r.Method == "GET" && !undeleted:
			_, _ = w.Write([]byte(`{
				"name": "projects/test-project/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
				"state": "DELETED",
				"displayName": "Old Provider Name",
				"description": "Old Provider Description",
				"disabled": true,
				"attributeMapping": {
					"google.subject": "assertion.sub"
				},
				"oidc": {
					"issuerUri": "https://old.issuer.example.com"
				}
			}`))
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, ":undelete"):
			undeleted = true
			_, _ = w.Write([]byte(`{"name": "operations/undelete-provider-op", "done": true}`))
		case r.Method == "PATCH":
			_, _ = w.Write([]byte(`{"name": "operations/patch-provider-op", "done": true}`))
		case r.Method == "GET" && undeleted:
			_, _ = w.Write([]byte(`{
				"name": "projects/test-project/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
				"state": "ACTIVE",
				"displayName": "Recreated Provider Name",
				"description": "Recreated Provider Description",
				"disabled": false,
				"attributeMapping": {
					"google.subject": "assertion.sub"
				},
				"attributeCondition": "assertion.sub != ''",
				"oidc": {
					"issuerUri": "https://new.issuer.example.com",
					"allowedAudiences": ["aud1"]
				}
			}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "unexpected request"}`))
		}
	}))
	defer ts.Close()

	res := iambeta.ResourceIAMBetaWorkloadIdentityPoolProvider()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]any{
		"project":                            "test-project",
		"workload_identity_pool_id":          "my-pool",
		"workload_identity_pool_provider_id": "my-provider",
		"display_name":                       "Recreated Provider Name",
		"description":                        "Recreated Provider Description",
		"disabled":                           false,
		"attribute_mapping": map[string]any{
			"google.subject": "assertion.sub",
		},
		"attribute_condition": "assertion.sub != ''",
		"oidc": []any{
			map[string]any{
				"issuer_uri":        "https://new.issuer.example.com",
				"allowed_audiences": []any{"aud1"},
			},
		},
	})

	config := &transport_tpg.Config{
		Client:       ts.Client(),
		Project:      "test-project",
		UserAgent:    "test-agent",
		PollInterval: 1 * time.Millisecond,
		CustomEndpoints: map[string]string{
			"iam_beta_custom_endpoint": ts.URL + "/v1/",
		},
	}

	if err := res.Create(d, config); err != nil {
		t.Fatalf("ResourceIAMBetaWorkloadIdentityPoolProvider().Create(%q) returned unexpected error: %v", "my-provider", err)
	}

	wantRequests := []recordedRequest{
		{
			Method: "GET",
			Path:   "/v1/projects/test-project/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
		},
		{
			Method: "POST",
			Path:   "/v1/projects/test-project/locations/global/workloadIdentityPools/my-pool/providers/my-provider:undelete",
			Body:   map[string]any{},
		},
		{
			Method:     "PATCH",
			Path:       "/v1/projects/test-project/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
			UpdateMask: "displayName,description,disabled,attributeMapping,attributeCondition,oidc.allowedAudiences,oidc.issuerUri,oidc.jwksJson",
			Body: map[string]any{
				"displayName":        "Recreated Provider Name",
				"description":        "Recreated Provider Description",
				"disabled":           false,
				"attributeCondition": "assertion.sub != ''",
				"attributeMapping": map[string]any{
					"google.subject": "assertion.sub",
				},
				"oidc": map[string]any{
					"allowedAudiences": []any{"aud1"},
					"issuerUri":        "https://new.issuer.example.com",
				},
			},
		},
		{
			Method: "GET",
			Path:   "/v1/projects/test-project/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
		},
	}

	if diff := cmp.Diff(wantRequests, requests); diff != "" {
		t.Errorf("ResourceIAMBetaWorkloadIdentityPoolProvider().Create(%q) HTTP requests mismatch (-want +got):\n%s", "my-provider", diff)
	}
}
