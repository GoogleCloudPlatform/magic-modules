package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestStorageBucketAclRoleEntityCustomizeDiff(t *testing.T) {
	t.Parallel()

	// Default ACLs of a bucket owned by project number 111.
	const (
		owners  = "OWNER:project-owners-111"
		editors = "OWNER:project-editors-111"
		viewers = "READER:project-viewers-111"
		user    = "OWNER:user-alice@example.com"
	)

	cases := map[string]struct {
		state             []string
		config            []string
		bucketLookupFails bool
		want              []string
		wantBucketLookups int32
		wantErr           bool
	}{
		"default ACLs of the bucket's project are kept": {
			state:             []string{owners, editors, viewers, user},
			config:            []string{user},
			want:              []string{owners, editors, viewers, user},
			wantBucketLookups: 1,
		},
		"teams of another project are removed": {
			state:             []string{owners, editors, viewers, user, "OWNER:project-owners-999", "OWNER:project-editors-999", "READER:project-viewers-999"},
			config:            []string{user},
			want:              []string{owners, editors, viewers, user},
			wantBucketLookups: 1,
		},
		"project number that only shares a prefix is removed": {
			state:             []string{owners, editors, viewers, user, "OWNER:project-editors-1111"},
			config:            []string{user},
			want:              []string{owners, editors, viewers, user},
			wantBucketLookups: 1,
		},
		"teams of another project are kept when configured": {
			state:             []string{owners, editors, viewers, user, "READER:project-viewers-999"},
			config:            []string{user, "READER:project-viewers-999"},
			want:              []string{owners, editors, viewers, user, "READER:project-viewers-999"},
			wantBucketLookups: 1,
		},
		"no bucket lookup when the config lists the default ACLs": {
			state:  []string{owners, editors, viewers, user, "READER:allUsers"},
			config: []string{owners, editors, viewers, user},
			want:   []string{owners, editors, viewers, user},
		},
		"failed bucket lookup fails the plan": {
			state:             []string{owners, editors, viewers, user},
			config:            []string{user},
			bucketLookupFails: true,
			wantBucketLookups: 1,
			wantErr:           true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var bucketLookups atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/storage/v1/b/my-bucket" {
					http.NotFound(w, r)
					return
				}
				bucketLookups.Add(1)
				if tc.bucketLookupFails {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"name": "my-bucket", "projectNumber": "111"}`)
			}))
			defer srv.Close()

			config := &transport_tpg.Config{
				Client:          srv.Client(),
				Context:         context.Background(),
				CustomEndpoints: map[string]string{Product.CustomEndpointField: srv.URL + "/storage/v1/"},
			}

			attributes := map[string]string{
				"bucket":          "my-bucket",
				"deletion_policy": "DELETE",
				"role_entity.#":   strconv.Itoa(len(tc.state)),
			}
			for i, re := range tc.state {
				attributes[fmt.Sprintf("role_entity.%d", i)] = re
			}
			state := &terraform.InstanceState{ID: "my-bucket-acl", Attributes: attributes}

			roleEntity := make([]interface{}, 0, len(tc.config))
			for _, re := range tc.config {
				roleEntity = append(roleEntity, re)
			}
			resourceConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
				"bucket":      "my-bucket",
				"role_entity": roleEntity,
			})

			// Only the role_entity customize diff is under test.
			r := &schema.Resource{
				Schema:        ResourceStorageBucketAcl().Schema,
				CustomizeDiff: resourceStorageRoleEntityCustomizeDiff,
			}
			diff, err := r.SimpleDiff(context.Background(), state, resourceConfig, config)
			if got := bucketLookups.Load(); got != tc.wantBucketLookups {
				t.Errorf("got %d bucket lookups, want %d", got, tc.wantBucketLookups)
			}
			if tc.wantErr {
				if err == nil || !strings.Contains(err.Error(), `Error reading bucket "my-bucket"`) {
					t.Fatalf("got error %v, want a bucket read error", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			planned := state.MergeDiff(diff).Attributes
			count, err := strconv.Atoi(planned["role_entity.#"])
			if err != nil {
				t.Fatalf("unexpected planned role_entity count %q: %s", planned["role_entity.#"], err)
			}
			got := make([]string, 0, count)
			for i := 0; i < count; i++ {
				got = append(got, planned[fmt.Sprintf("role_entity.%d", i)])
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("planned role_entity:\n got %v\nwant %v", got, tc.want)
			}
		})
	}
}
