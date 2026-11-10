package cloudidentity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestCloudIdentityInboundOidcSsoProfile_Create_MPA_DoneFalse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/inboundOidcSsoProfiles") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"done": false,
				"name": "",
				"metadata": map[string]interface{}{
					"state": "awaiting-multi-party-approval",
				},
			})
			return
		}
		t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()

	config := &transport_tpg.Config{
		CustomEndpoints: map[string]string{
			"cloud_identity_custom_endpoint": srv.URL + "/",
		},
		Client: srv.Client(),
	}
	d := schema.TestResourceDataRaw(t, ResourceCloudIdentityInboundOidcSsoProfile().Schema, map[string]interface{}{
		"customer": "customers/C0123abc",
		"idp_config": []interface{}{
			map[string]interface{}{
				"issuer_uri": "https://idp.example.com",
			},
		},
	})

	err := resourceCloudIdentityInboundOidcSsoProfileCreate(d, config)
	if err == nil {
		t.Fatalf("Expected error for MPA awaiting approval, got nil")
	}
	if !strings.Contains(err.Error(), "awaiting-multi-party-approval") {
		t.Fatalf("Expected error to contain awaiting-multi-party-approval, got: %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("Expected empty ID, got %q", d.Id())
	}
}

func TestCloudIdentityInboundOidcSsoProfile_Create_DoneTrue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/inboundOidcSsoProfiles") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"done": true,
				"response": map[string]interface{}{
					"name": "inboundOidcSsoProfiles/profile-123",
				},
			})
			return
		}
		if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/inboundOidcSsoProfiles/profile-123") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"name": "inboundOidcSsoProfiles/profile-123",
				"customer": "customers/C0123abc",
			})
			return
		}
		t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()

	config := &transport_tpg.Config{
		CustomEndpoints: map[string]string{
			"cloud_identity_custom_endpoint": srv.URL + "/",
		},
		Client: srv.Client(),
	}
	d := schema.TestResourceDataRaw(t, ResourceCloudIdentityInboundOidcSsoProfile().Schema, map[string]interface{}{
		"customer": "customers/C0123abc",
		"idp_config": []interface{}{
			map[string]interface{}{
				"issuer_uri": "https://idp.example.com",
			},
		},
	})

	err := resourceCloudIdentityInboundOidcSsoProfileCreate(d, config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if d.Id() != "inboundOidcSsoProfiles/profile-123" {
		t.Fatalf("Expected ID to be inboundOidcSsoProfiles/profile-123, got %q", d.Id())
	}
	if d.Get("name").(string) != "inboundOidcSsoProfiles/profile-123" {
		t.Fatalf("Expected name to be inboundOidcSsoProfiles/profile-123, got %q", d.Get("name"))
	}
	if d.Get("inbound_oidc_sso_profile_id").(string) != "profile-123" {
		t.Fatalf("Expected inbound_oidc_sso_profile_id to be profile-123, got %q", d.Get("inbound_oidc_sso_profile_id"))
	}
}

func TestCheckCloudIdentityInboundOidcSsoProfileOperation_DoneFalse_MPA(t *testing.T) {
	res := map[string]interface{}{
		"done": false,
		"metadata": map[string]interface{}{
			"state": "awaiting-multi-party-approval",
		},
	}
	_, err := checkCloudIdentityInboundOidcSsoProfileOperation(res, "updating")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "awaiting-multi-party-approval") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestCheckCloudIdentityInboundOidcSsoProfileOperation_DoneTrue(t *testing.T) {
	res := map[string]interface{}{
		"done": true,
		"response": map[string]interface{}{
			"name": "inboundOidcSsoProfiles/profile-123",
		},
	}
	respBody, err := checkCloudIdentityInboundOidcSsoProfileOperation(res, "updating")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if respBody["name"] != "inboundOidcSsoProfiles/profile-123" {
		t.Fatalf("unexpected name")
	}
}

func TestCheckCloudIdentityInboundOidcSsoProfileOperation_Error(t *testing.T) {
	res := map[string]interface{}{
		"error": map[string]interface{}{
			"message": "some error",
		},
	}
	_, err := checkCloudIdentityInboundOidcSsoProfileOperation(res, "updating")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
