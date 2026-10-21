package publicca

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Exercise the generated Open method, including URL defaults, API field aliases,
// response transforms, and diagnostics. Synthetic key material never reaches GCP.
func TestEphemeralPublicCAExternalAccountKeyOpen(t *testing.T) {
	ctx := context.Background()
	macKey := base64.RawURLEncoding.EncodeToString([]byte{0xfb})
	encoded := base64.URLEncoding.EncodeToString([]byte(macKey))
	for _, tc := range []struct {
		name      string
		key       interface{}
		wantError bool
	}{
		{"decoded values", encoded, false},
		{"malformed base64", "invalid-secret-!", true},
		{"wrong API type", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" || r.URL.Path != "/projects/test-project/locations/global/externalAccountKeys" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"name":  "projects/test-project/locations/global/externalAccountKeys/test-key",
					"keyId": "test-key", "b64MacKey": tc.key,
				})
			}))
			defer server.Close()
			p := GoogleEphemeralPublicCAExternalAccountKey()
			var schemaResp ephemeral.SchemaResponse
			p.Schema(ctx, ephemeral.SchemaRequest{}, &schemaResp)
			var configured ephemeral.ConfigureResponse
			p.(ephemeral.EphemeralResourceWithConfigure).Configure(ctx, ephemeral.ConfigureRequest{ProviderData: &transport_tpg.Config{
				Project: "test-project", Client: server.Client(),
				CustomEndpoints: map[string]string{"public_ca_custom_endpoint": server.URL + "/"},
			}}, &configured)
			if configured.Diagnostics.HasError() {
				t.Fatal(configured.Diagnostics)
			}
			attrs := make(map[string]tftypes.Value)
			for name, a := range schemaResp.Schema.Attributes {
				attrs[name] = tftypes.NewValue(a.GetType().TerraformType(ctx), nil)
			}
			raw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), attrs)
			resp := ephemeral.OpenResponse{Result: tfsdk.EphemeralResultData{Schema: schemaResp.Schema, Raw: raw}}
			p.Open(ctx, ephemeral.OpenRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: raw}}, &resp)
			if resp.Diagnostics.HasError() != tc.wantError {
				t.Fatalf("diagnostics = %v, want error = %v", resp.Diagnostics, tc.wantError)
			}
			if tc.wantError {
				if strings.Contains(resp.Diagnostics.Errors()[0].Detail(), "invalid-secret") {
					t.Fatal("secret exposed in diagnostic")
				}
				return
			}
			var result ephemeralPublicCAExternalAccountKeyModel
			if diags := resp.Result.Get(ctx, &result); diags.HasError() {
				t.Fatal(diags)
			}
			if result.MacKey.ValueString() != macKey {
				t.Fatal("mac_key was not decoded exactly once")
			}
			if result.B64urlMacKey.ValueString() != base64.RawURLEncoding.EncodeToString([]byte(macKey)) {
				t.Fatal("b64url_mac_key has incorrect encoding")
			}
			if result.B64MacKey.ValueString() != encoded {
				t.Fatal("raw API value changed")
			}
		})
	}
}
