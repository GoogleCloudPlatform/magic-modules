package sql

import (
	"context"
	"fmt"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSqlSslCert_readParentNotFound(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name         string
		instanceCode int
		wantError    bool
	}{{"missing instance", 404, false}, {"forbidden instance", 403, true}, {"live instance forbidden cert", 200, true}} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/sql/v1beta4/projects/test-project/instances/test-instance" {
					if tc.instanceCode == 200 {
						fmt.Fprint(w, `{"name":"test-instance","settings":{"activationPolicy":"ALWAYS"}}`)
						return
					}
					w.WriteHeader(tc.instanceCode)
					fmt.Fprintf(w, `{"error":{"code":%d,"message":"instance lookup failed"}}`, tc.instanceCode)
					return
				}
				if r.URL.Path == "/sql/v1beta4/projects/test-project/instances/test-instance/sslCerts/fingerprint" {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprint(w, `{"error":{"code":403,"message":"not authorized"}}`)
					return
				}
				t.Errorf("unexpected request: %s", r.URL.Path)
				w.WriteHeader(http.StatusForbidden)
			}))
			defer server.Close()
			cfg := &transport_tpg.Config{Context: context.Background(), Client: server.Client(), Project: "test-project", CustomEndpoints: map[string]string{"sql_custom_endpoint": server.URL + "/sql/v1beta4/"}}
			res := ResourceSqlSslCert()
			d := res.Data(nil)
			const id = "projects/test-project/instances/test-instance/child"
			d.SetId(id)
			for key, value := range map[string]interface{}{"project": "test-project", "instance": "test-instance", "common_name": "child", "sha1_fingerprint": "fingerprint"} {
				if err := d.Set(key, value); err != nil {
					t.Fatal(err)
				}
			}
			err := res.Read(d, cfg)
			if (err != nil) != tc.wantError {
				t.Fatalf("Read error = %v, wantError = %v", err, tc.wantError)
			}
			wantID := ""
			if tc.wantError {
				wantID = id
			}
			if d.Id() != wantID {
				t.Fatalf("ID = %q, want %q", d.Id(), wantID)
			}
		})
	}
}
