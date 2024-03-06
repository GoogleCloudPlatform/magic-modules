package ancestrymanager

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	resources "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/tfdata"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"

	"github.com/google/go-cmp/cmp"
	provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
	"go.uber.org/zap"
	crmv1 "google.golang.org/api/cloudresourcemanager/v1"
	crmv3 "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/storage/v1"
)

func TestGetAncestors(t *testing.T) {
	ownerProject := "foo"
	ownerAncestryPath := "organization/qux/folder/bar/project/foo"
	anotherProject := "foo2"

	// Setup a simple test server to mock the response of resource manager.
	v3Responses := map[string]*crmv3.Project{
		"folders/bar":        {Name: "folders/bar", Parent: "organizations/qux"},
		"organizations/qux":  {Name: "organizations/qux", Parent: ""},
		"folders/bar2":       {Name: "folders/bar2", Parent: "organizations/qux2"},
		"organizations/qux2": {Name: "organizations/qux2", Parent: ""},
	}
	v1Responses := map[string][]*crmv1.Ancestor{
		ownerProject: {
			{ResourceId: &crmv1.ResourceId{Id: "foo", Type: "project"}},
			{ResourceId: &crmv1.ResourceId{Id: "bar", Type: "folder"}},
			{ResourceId: &crmv1.ResourceId{Id: "qux", Type: "organization"}},
		},
		"12345": {
			{ResourceId: &crmv1.ResourceId{Id: "foo", Type: "project"}},
			{ResourceId: &crmv1.ResourceId{Id: "bar", Type: "folder"}},
			{ResourceId: &crmv1.ResourceId{Id: "qux", Type: "organization"}},
		},
		anotherProject: {
			{ResourceId: &crmv1.ResourceId{Id: "foo2", Type: "project"}},
			{ResourceId: &crmv1.ResourceId{Id: "bar2", Type: "folder"}},
			{ResourceId: &crmv1.ResourceId{Id: "qux2", Type: "organization"}},
		},
	}

	ts := newTestServer(t, v1Responses, v3Responses)
	defer ts.Close()
	mockV1Client, err := crmv1.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatal(err)
	}
	mockV3Client, err := crmv3.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatal(err)
	}

	entries := map[string]string{
		ownerProject: ownerAncestryPath,
	}

	p := provider.Provider()

	// offline return errors when the cache cannot cover the request.
	// online return errors when neither cache and mock server cannot cover the request.
	cases := []struct {
		name             string
		data             tpgresource.TerraformResourceData
		asset            *resources.Asset
		cfg              *transport_tpg.Config
		want             []string
		wantParent       string
		wantOnlineError  bool
		wantOfflineError bool
	}{
		{
			name: "owner project - project id",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": ownerProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:       []string{"projects/foo", "folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			// google_project does not expect `project` attribute
			name: "owner project - project",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project": ownerProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:       []string{"organizations/unknown"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/unknown",
		},
		{
			// google_project_iam expect `project` attribute
			name: "owner project - project",
			data: tfdata.NewFakeResourceData(
				"google_project_iam_member",
				p.ResourcesMap["google_project_iam_member"].Schema,
				map[string]interface{}{
					"project": ownerProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:       []string{"projects/foo", "folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "owner project - project number",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"number": "12345",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:             []string{"projects/foo", "folders/bar", "organizations/qux"},
			wantOfflineError: true,
			wantParent:       "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "owner project - project from config",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			cfg: &transport_tpg.Config{
				Project: ownerProject,
			},
			want:       []string{"projects/foo", "folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "another project",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": anotherProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:             []string{"projects/foo2", "folders/bar2", "organizations/qux2"},
			wantOfflineError: true,
			wantParent:       "//cloudresourcemanager.googleapis.com/folders/bar2",
		},
		{
			name: "owner folder",
			data: tfdata.NewFakeResourceData(
				"google_folder_iam_policy",
				p.ResourcesMap["google_folder_iam_policy"].Schema,
				map[string]interface{}{
					"folder": "bar",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:       []string{"folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "owner folder with prefix",
			data: tfdata.NewFakeResourceData(
				"google_folder_iam_policy",
				p.ResourcesMap["google_folder_iam_policy"].Schema,
				map[string]interface{}{
					"folder": "folders/bar",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:       []string{"folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "another folder online",
			data: tfdata.NewFakeResourceData(
				"google_folder_iam_policy",
				p.ResourcesMap["google_folder_iam_policy"].Schema,
				map[string]interface{}{
					"folder": "bar2",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:             []string{"folders/bar2", "organizations/qux2"},
			wantOfflineError: true,
			wantParent:       "//cloudresourcemanager.googleapis.com/organizations/qux2",
		},
		{
			// Not supporting folder create resource yet.
			name: "not exist folder online",
			data: tfdata.NewFakeResourceData(
				"google_folder_iam_policy",
				p.ResourcesMap["google_folder_iam_policy"].Schema,
				map[string]interface{}{
					"folder": "notexist",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			wantOfflineError: true,
			wantOnlineError:  true,
		},
		{
			name: "owner org",
			data: tfdata.NewFakeResourceData(
				"google_organization_iam_policy",
				p.ResourcesMap["google_organization_iam_policy"].Schema,
				map[string]interface{}{
					"org_id": "qux",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Organization",
			},
			want:       []string{"organizations/qux"},
			wantParent: "",
		},
		{
			// organization do not have ancestors except itself
			// hence offline also pass.
			name: "another org",
			data: tfdata.NewFakeResourceData(
				"google_organization_iam_policy",
				p.ResourcesMap["google_organization_iam_policy"].Schema,
				map[string]interface{}{
					"org_id": "qux2",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Organization",
			},
			want:       []string{"organizations/qux2"},
			wantParent: "",
		},
		{
			name: "other resource with owner project",
			data: tfdata.NewFakeResourceData(
				"google_compute_disk",
				p.ResourcesMap["google_compute_disk"].Schema,
				map[string]interface{}{
					"project": ownerProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Disk",
			},
			want:       []string{"projects/foo", "folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/projects/foo",
		},
		{
			name: "other resource online with another project",
			data: tfdata.NewFakeResourceData(
				"google_compute_disk",
				p.ResourcesMap["google_compute_disk"].Schema,
				map[string]interface{}{
					"project": anotherProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Disk",
			},
			want:             []string{"projects/foo2", "folders/bar2", "organizations/qux2"},
			wantOfflineError: true,
			wantParent:       "//cloudresourcemanager.googleapis.com/projects/foo2",
		},
		{
			name: "custom role with org",
			data: tfdata.NewFakeResourceData(
				"google_organization_iam_custom_role",
				p.ResourcesMap["google_organization_iam_custom_role"].Schema,
				map[string]interface{}{
					"org_id": "qux",
				},
			),
			asset: &resources.Asset{
				Type: "iam.googleapis.com/Role",
			},
			want:       []string{"organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "custom role with project",
			data: tfdata.NewFakeResourceData(
				"google_project_iam_custom_role",
				p.ResourcesMap["google_project_iam_custom_role"].Schema,
				map[string]interface{}{
					"project": "foo",
				},
			),
			asset: &resources.Asset{
				Type: "iam.googleapis.com/Role",
			},
			want:       []string{"projects/foo", "folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/projects/foo",
		},
		{
			name: "custom role with empty project ID",
			data: tfdata.NewFakeResourceData(
				"google_project_iam_custom_role",
				p.ResourcesMap["google_project_iam_custom_role"].Schema,
				map[string]interface{}{
					"project": "",
				},
			),
			asset: &resources.Asset{
				Type: "iam.googleapis.com/Role",
			},
			want:       []string{"organizations/unknown"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/unknown",
		},
		{
			name: "new project in folder",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"folder_id":  "bar",
					"project_id": "new-project",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:       []string{"projects/new-project", "folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "new project in organization",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"org_id":     "qux",
					"project_id": "new-project",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:       []string{"projects/new-project", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			// for new projects, if it cannot find ancestors in online mode,
			// it just returns 403 error.
			// offline will fail because no cloud resource manager.
			name: "new project without org_id or folder_id",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": "new-project",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			wantOnlineError:  true,
			wantOfflineError: true,
		},
		{
			name: "Org policy v2 on Project",
			data: tfdata.NewFakeResourceData(
				"google_org_policy_policy",
				p.ResourcesMap["google_org_policy_policy"].Schema,
				map[string]interface{}{
					"parent": "projects/foo",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:       []string{"projects/foo", "folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "Org policy v2 on Folder",
			data: tfdata.NewFakeResourceData(
				"google_org_policy_policy",
				p.ResourcesMap["google_org_policy_policy"].Schema,
				map[string]interface{}{
					"parent": "folders/bar",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:       []string{"folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "Org policy v2 on Organization",
			data: tfdata.NewFakeResourceData(
				"google_org_policy_policy",
				p.ResourcesMap["google_org_policy_policy"].Schema,
				map[string]interface{}{
					"parent": "organizations/qux",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Organization",
			},
			want: []string{"organizations/qux"},
		},
		{
			name: "Google folder with organizations/ as {parent}",
			data: tfdata.NewFakeResourceData(
				"google_folder",
				p.ResourcesMap["google_folder"].Schema,
				map[string]interface{}{
					"parent": "organizations/qux",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:       []string{"organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "Google folder with folders/ as {parent}",
			data: tfdata.NewFakeResourceData(
				"google_folder",
				p.ResourcesMap["google_folder"].Schema,
				map[string]interface{}{
					"parent": "folders/bar",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:       []string{"folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "Google folder with both folder_id and parent fields present",
			data: tfdata.NewFakeResourceData(
				"google_folder",
				p.ResourcesMap["google_folder"].Schema,
				map[string]interface{}{
					"folder_id": "bar",
					"parent":    "organizations/qux",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:       []string{"folders/bar", "organizations/qux"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "Google folder with missing parent field",
			data: tfdata.NewFakeResourceData(
				"google_folder",
				p.ResourcesMap["google_folder"].Schema,
				map[string]interface{}{},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:       []string{"organizations/unknown"},
			wantParent: "//cloudresourcemanager.googleapis.com/organizations/unknown",
		},
	}
	for _, c := range cases {
		for _, offline := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s offline = %t", c.name, offline), func(t *testing.T) {
				if c.cfg == nil {
					c.cfg = &transport_tpg.Config{}
				}
				ancestryManager := &manager{
					errorLogger:   zap.NewExample(),
					ancestorCache: make(map[string][]string),
				}
				if !offline {
					ancestryManager.resourceManagerV3 = mockV3Client
					ancestryManager.resourceManagerV1 = mockV1Client
				}
				ancestryManager.initAncestryCache(entries)

				got, gotParent, err := ancestryManager.Ancestors(c.cfg, c.data, c.asset)
				if !offline {
					if c.wantOnlineError {
						if err == nil {
							t.Fatalf("onlineMgr.Ancestors(%v, %v, %v) = nil, want = err", c.cfg, c.data, c.asset)
						}
					} else {
						if err != nil {
							t.Fatalf("onlineMgr.Ancestors(%v, %v, %v) = %s, want = nil", c.cfg, c.data, c.asset, err)
						}
						if gotParent != c.wantParent {
							t.Errorf("onlineMgr.Ancestors(%v, %v, %v) parent = %s, want = %s", c.cfg, c.data, c.asset, gotParent, c.wantParent)
						}
						if diff := cmp.Diff(c.want, got); diff != "" {
							t.Errorf("onlineMgr.Ancestors(%v, %v, %v) returned unexpected diff (-want +got):\n%s", c.cfg, c.data, c.asset, diff)
						}
					}
				} else {
					if c.wantOfflineError {
						if err == nil {
							t.Fatalf("offlineMgr.Ancestors(%v, %v, %v) = nil, want = err", c.cfg, c.data, c.asset)
						}
					} else {
						if err != nil {
							t.Fatalf("offlineMgr.Ancestors(%v, %v, %v) = %s, want = nil", c.cfg, c.data, c.asset, err)
						}
						if gotParent != c.wantParent {
							t.Errorf("offlineMgr.Ancestors(%v, %v, %v) parent = %s, want = %s", c.cfg, c.data, c.asset, gotParent, c.wantParent)
						}
						if diff := cmp.Diff(c.want, got); diff != "" {
							t.Errorf("offlineMgr.Ancestors(%v, %v, %v) returned unexpected diff (-want +got):\n%s", c.cfg, c.data, c.asset, diff)
						}
					}
				}
			})
		}
	}
}

type testServer struct {
	*httptest.Server
	v1Count int
	v3Count int
}

func newTestServer(t *testing.T, v1Responses map[string][]*crmv1.Ancestor, v3Responses map[string]*crmv3.Project) *testServer {
	t.Helper()
	ts := &testServer{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v3/") {
			name := strings.TrimPrefix(r.URL.Path, "/v3/")
			resp, ok := v3Responses[name]
			if !ok {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(fmt.Sprintf("no response for request path %s", "/v3/"+name)))
				return
			}

			ts.v3Count++
			payload, err := resp.MarshalJSON()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("failed to MarshalJSON: %s", err)))
				return
			}
			w.Write(payload)
		} else if strings.HasPrefix(r.URL.Path, "/v1/") {
			re := regexp.MustCompile(`([^/]*):getAncestry`)
			path := re.FindStringSubmatch(r.URL.Path)
			if path == nil || v1Responses[path[1]] == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ts.v1Count++
			payload, err := (&crmv1.GetAncestryResponse{Ancestor: v1Responses[path[1]]}).MarshalJSON()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("failed to MarshalJSON: %s", err)))
				return
			}
			w.Write(payload)
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(fmt.Sprintf("no response for url path %s", r.URL.Path)))
			return
		}
	}))
	ts.Server = server
	return ts
}

func TestGetAncestors_Folder(t *testing.T) {
	// v1 API is used to fetch ancestry first, and compared with the resource data
	// to find out whether there is a folder ID change.
	p := provider.Provider()
	cases := []struct {
		name        string
		data        tpgresource.TerraformResourceData
		asset       *resources.Asset
		v1Responses map[string][]*crmv1.Ancestor
		v3Responses map[string]*crmv3.Project
		want        []string
		parent      string
		v3Count     int
		v1Count     int
	}{
		{
			name: "folder not changed",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": "foo",
					"folder_id":  "folders/bar",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			v3Responses: map[string]*crmv3.Project{},
			v1Responses: map[string][]*crmv1.Ancestor{
				"foo": {
					{ResourceId: &crmv1.ResourceId{Id: "foo", Type: "project"}},
					{ResourceId: &crmv1.ResourceId{Id: "bar", Type: "folder"}},
					{ResourceId: &crmv1.ResourceId{Id: "qux", Type: "organization"}},
				},
			},
			want:    []string{"projects/foo", "folders/bar", "organizations/qux"},
			parent:  "//cloudresourcemanager.googleapis.com/folders/bar",
			v1Count: 1,
		},
		{
			// project moving from organizations/qux/folders/bar to organizations/qux2/folders/bar2
			name: "project moved from a top-level folder in one org to a top-level folder in a different org",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": "foo",
					"folder_id":  "folders/bar2",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			v3Responses: map[string]*crmv3.Project{
				"folders/bar2":       {Name: "folders/bar2", Parent: "organizations/qux2"},
				"organizations/qux2": {Name: "organizations/qux2", Parent: ""},
			},
			v1Responses: map[string][]*crmv1.Ancestor{
				"foo": {
					{ResourceId: &crmv1.ResourceId{Id: "foo", Type: "project"}},
					{ResourceId: &crmv1.ResourceId{Id: "bar", Type: "folder"}},
					{ResourceId: &crmv1.ResourceId{Id: "qux", Type: "organization"}},
				},
			},
			want:    []string{"projects/foo", "folders/bar2", "organizations/qux2"},
			parent:  "//cloudresourcemanager.googleapis.com/folders/bar2",
			v3Count: 1,
			v1Count: 1,
		},
		{
			// project moving from folders/bar2/folders/bar to folders/bar2
			name: "project moved from child folder to parent folder",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": "foo",
					"folder_id":  "folders/bar2",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			v3Responses: map[string]*crmv3.Project{},
			v1Responses: map[string][]*crmv1.Ancestor{
				"foo": {
					{ResourceId: &crmv1.ResourceId{Id: "foo", Type: "project"}},
					{ResourceId: &crmv1.ResourceId{Id: "bar", Type: "folder"}},
					{ResourceId: &crmv1.ResourceId{Id: "bar2", Type: "folder"}},
					{ResourceId: &crmv1.ResourceId{Id: "qux", Type: "organization"}},
				},
			},
			want:    []string{"projects/foo", "folders/bar2", "organizations/qux"},
			parent:  "//cloudresourcemanager.googleapis.com/folders/bar2",
			v1Count: 1,
		},
		{
			// project moving from folders/bar2 to folders/bar2/folders/bar
			name: "project moved from parent folder to child folder",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": "foo",
					"folder_id":  "folders/bar",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			v3Responses: map[string]*crmv3.Project{
				"folders/bar": {Name: "folders/bar", Parent: "folders/bar2"},
			},
			v1Responses: map[string][]*crmv1.Ancestor{
				"foo": {
					{ResourceId: &crmv1.ResourceId{Id: "foo", Type: "project"}},
					{ResourceId: &crmv1.ResourceId{Id: "bar2", Type: "folder"}},
					{ResourceId: &crmv1.ResourceId{Id: "qux", Type: "organization"}},
				},
			},
			want:    []string{"projects/foo", "folders/bar", "folders/bar2", "organizations/qux"},
			parent:  "//cloudresourcemanager.googleapis.com/folders/bar",
			v1Count: 1,
			v3Count: 1,
		},
		{
			name: "folder ID is empty string and no ancestor",
			data: tfdata.NewFakeResourceData(
				"google_folder_iam_member",
				p.ResourcesMap["google_folder_iam_member"].Schema,
				map[string]interface{}{
					"folder": "",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			v3Responses: map[string]*crmv3.Project{},
			v1Responses: map[string][]*crmv1.Ancestor{},
			want:        []string{"organizations/unknown"},
			parent:      "//cloudresourcemanager.googleapis.com/organizations/unknown",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ts := newTestServer(t, c.v1Responses, c.v3Responses)
			defer ts.Close()

			mockV1Client, err := crmv1.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}
			mockV3Client, err := crmv3.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}

			cfg := &transport_tpg.Config{
				Project: "foo",
			}
			ancestryManager := &manager{
				errorLogger:       zap.NewExample(),
				ancestorCache:     make(map[string][]string),
				resourceManagerV3: mockV3Client,
				resourceManagerV1: mockV1Client,
			}
			// empty cache
			ancestryManager.initAncestryCache(map[string]string{})

			got, parent, err := ancestryManager.Ancestors(cfg, c.data, c.asset)
			if err != nil {
				t.Fatalf("Ancestors(%v, %v, %v) = %s, want = nil", cfg, c.data, c.asset, err)
			}
			if parent != c.parent {
				t.Errorf("Ancestors(%v, %v, %v) parent = %s, want = %s", cfg, c.data, c.asset, parent, c.parent)
			}
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("Ancestors(%v, %v, %v) returned unexpected diff (-want +got):\n%s", cfg, c.data, c.asset, diff)
			}
			if ts.v3Count != c.v3Count {
				t.Errorf("Ancestors(%v, %v, %v) v3 API called = %d, want = %d", cfg, c.data, c.asset, ts.v3Count, c.v3Count)
			}
			if ts.v1Count != c.v1Count {
				t.Errorf("Ancestors(%v, %v, %v) v1 API called = %d, want = %d", cfg, c.data, c.asset, ts.v1Count, c.v1Count)
			}
		})
	}
}

func TestGetAncestorsWithCache(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		cache       map[string][]string
		v3Responses map[string]*crmv3.Project
		v1Responses map[string][]*crmv1.Ancestor
		want        []string
		wantCache   map[string][]string
	}{
		{
			name:        "empty cache",
			input:       "projects/abc",
			cache:       make(map[string][]string),
			v3Responses: map[string]*crmv3.Project{},
			v1Responses: map[string][]*crmv1.Ancestor{
				"abc": {
					{ResourceId: &crmv1.ResourceId{Id: "abc", Type: "project"}},
					{ResourceId: &crmv1.ResourceId{Id: "456", Type: "folder"}},
					{ResourceId: &crmv1.ResourceId{Id: "789", Type: "folder"}},
					{ResourceId: &crmv1.ResourceId{Id: "321", Type: "organization"}},
				},
			},
			want: []string{"projects/abc", "folders/456", "folders/789", "organizations/321"},
			wantCache: map[string][]string{
				"projects/abc":      {"projects/abc", "folders/456", "folders/789", "organizations/321"},
				"folders/456":       {"folders/456", "folders/789", "organizations/321"},
				"folders/789":       {"folders/789", "organizations/321"},
				"organizations/321": {"organizations/321"},
			},
		},
		{
			name:  "partial cache",
			input: "projects/abc",
			cache: map[string][]string{
				"folders/789":       {"folders/789", "organizations/321"},
				"organizations/321": {"organizations/321"},
			},
			v3Responses: map[string]*crmv3.Project{},
			v1Responses: map[string][]*crmv1.Ancestor{
				"abc": {
					{ResourceId: &crmv1.ResourceId{Id: "abc", Type: "project"}},
					{ResourceId: &crmv1.ResourceId{Id: "456", Type: "folder"}},
					{ResourceId: &crmv1.ResourceId{Id: "789", Type: "folder"}},
					{ResourceId: &crmv1.ResourceId{Id: "321", Type: "organization"}},
				},
			},
			want: []string{"projects/abc", "folders/456", "folders/789", "organizations/321"},
			wantCache: map[string][]string{
				"projects/abc":      {"projects/abc", "folders/456", "folders/789", "organizations/321"},
				"folders/456":       {"folders/456", "folders/789", "organizations/321"},
				"folders/789":       {"folders/789", "organizations/321"},
				"organizations/321": {"organizations/321"},
			},
		},
		{
			name:  "all response from cache",
			input: "projects/abc",
			cache: map[string][]string{
				"projects/abc": {"projects/123", "folders/456", "folders/789", "organizations/321"},
			},
			v3Responses: map[string]*crmv3.Project{},
			v1Responses: map[string][]*crmv1.Ancestor{},
			want:        []string{"projects/123", "folders/456", "folders/789", "organizations/321"},
			wantCache: map[string][]string{
				"projects/abc":      {"projects/123", "folders/456", "folders/789", "organizations/321"},
				"projects/123":      {"projects/123", "folders/456", "folders/789", "organizations/321"},
				"folders/456":       {"folders/456", "folders/789", "organizations/321"},
				"folders/789":       {"folders/789", "organizations/321"},
				"organizations/321": {"organizations/321"},
			},
		},
		{
			name:        "organization",
			input:       "organizations/321",
			cache:       map[string][]string{},
			v3Responses: map[string]*crmv3.Project{},
			v1Responses: map[string][]*crmv1.Ancestor{},
			want:        []string{"organizations/321"},
			wantCache: map[string][]string{
				"organizations/321": {"organizations/321"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := newTestServer(t, test.v1Responses, test.v3Responses)
			defer ts.Close()
			mockV1Client, err := crmv1.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}
			mockV3Client, err := crmv3.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}
			m := &manager{
				errorLogger:       zap.NewExample(),
				ancestorCache:     test.cache,
				resourceManagerV3: mockV3Client,
				resourceManagerV1: mockV1Client,
			}

			got, err := m.getAncestorsWithCache(test.input)
			if err != nil {
				t.Fatalf("getAncestorsWithCache(%s) = %s, want = nil", test.input, err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("getAncestorsWithCache(%v) returned unexpected diff (-want +got):\n%s", test.input, diff)
			}
			if diff := cmp.Diff(test.wantCache, m.ancestorCache); diff != "" {
				t.Errorf("getAncestorsWithCache(%v) cache returned unexpected diff (-want +got):\n%s", test.input, diff)
			}
		})
	}
}

func TestGetAncestorsWithCache_Fail(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		cache       map[string][]string
		v3Responses map[string]*crmv3.Project
		v1Responses map[string][]*crmv1.Ancestor
		wantErr     string
	}{
		{
			name:  "no parent response",
			input: "projects/abc",
			cache: make(map[string][]string),
			v3Responses: map[string]*crmv3.Project{
				"projects/abc": {Name: "projects/123", ProjectId: "projects/abc", Parent: "folders/not-exist"},
			},
			wantErr: "no response",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := newTestServer(t, test.v1Responses, test.v3Responses)
			defer ts.Close()
			mockV1Client, err := crmv1.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}
			mockV3Client, err := crmv3.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}
			m := &manager{
				errorLogger:       zap.NewExample(),
				ancestorCache:     test.cache,
				resourceManagerV3: mockV3Client,
				resourceManagerV1: mockV1Client,
			}

			_, err = m.getAncestorsWithCache(test.input)
			if err == nil {
				t.Fatalf("getAncestorsWithCache(%s) = nil, want = %s", test.input, test.wantErr)
			}
		})
	}
}

func TestParseAncestryPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []string
	}{
		{
			name: "all kinds of resource",
			path: "organizations/123/folders/456/projects/789",
			want: []string{"projects/789", "folders/456", "organizations/123"},
		},
		{
			name: "multiple folders",
			path: "organizations/123/folders/456/folders/789",
			want: []string{"folders/789", "folders/456", "organizations/123"},
		},
		{
			name: "normalize resource name",
			path: "organization/123/folder/456/project/789",
			want: []string{"projects/789", "folders/456", "organizations/123"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseAncestryPath(test.path)
			if err != nil {
				t.Fatalf("parseAncestryPath(%s) = %s, want = nil", test.path, err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("parseAncestryPath(%v) returned unexpected diff (-want +got):\n%s", test.path, diff)
			}
		})
	}
}

func TestParseAncestryPath_Fail(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{
			name:    "malform with single word",
			path:    "organizations",
			wantErr: "unexpected format",
		},
		{
			name:    "malform",
			path:    "organizations/123/folders",
			wantErr: "unexpected format",
		},
		{
			name:    "invalid keyword",
			path:    "org/123/folders/123",
			wantErr: "invalid ancestry path",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseAncestryPath(test.path)
			if err == nil {
				t.Fatalf("parseAncestryPath(%s) = nil, want = %s", test.path, test.wantErr)
			}
		})
	}
}

func TestInitAncestryCache(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string]string
		want    map[string][]string
	}{
		{
			name: "empty ancestry",
			entries: map[string]string{
				"test-proj": "",
			},
			want: map[string][]string{},
		},
		{
			name: "empty key",
			entries: map[string]string{
				"": "organizations/123/folders/345",
			},
			want: map[string][]string{},
		},
		{
			name: "default key to project",
			entries: map[string]string{
				"test-proj": "organizations/123/folders/345",
			},
			want: map[string][]string{
				"projects/test-proj": {"projects/test-proj", "folders/345", "organizations/123"},
				"folders/345":        {"folders/345", "organizations/123"},
				"organizations/123":  {"organizations/123"},
			},
		},
		{
			name: "key has prefix folders/",
			entries: map[string]string{
				"folders/345": "organizations/123",
			},
			want: map[string][]string{
				"folders/345":       {"folders/345", "organizations/123"},
				"organizations/123": {"organizations/123"},
			},
		},
		{
			name: "key has prefix projects/",
			entries: map[string]string{
				"projects/test-proj": "organizations/123",
			},
			want: map[string][]string{
				"projects/test-proj": {"projects/test-proj", "organizations/123"},
				"organizations/123":  {"organizations/123"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := &manager{
				ancestorCache: make(map[string][]string),
			}
			err := m.initAncestryCache(test.entries)
			if err != nil {
				t.Fatalf("initAncestryCache(%v) = %s, want = nil", test.entries, err)
			}
			if diff := cmp.Diff(test.want, m.ancestorCache); diff != "" {
				t.Errorf("initAncestryCache(%v) returned unexpected diff (-want +got):\n%s", test.entries, diff)
			}
		})
	}
}

func TestInitAncestryCache_Fail(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string]string
	}{
		{
			name: "typo",
			entries: map[string]string{
				"foldres/def": "organizations/123",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := &manager{
				ancestorCache: make(map[string][]string),
			}
			err := m.initAncestryCache(test.entries)
			if err == nil {
				t.Fatalf("initAncestryCache(%v) = nil, want = err", test.entries)
			}
		})
	}
}

func TestParseAncestryKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "not contain /",
			key:  "proj",
			want: "projects/proj",
		},
		{
			name: "contain projects/",
			key:  "projects/1",
			want: "projects/1",
		},
		{
			name: "contain folders/",
			key:  "folders/1",
			want: "folders/1",
		},
		{
			name: "contain organizations/",
			key:  "organizations/1",
			want: "organizations/1",
		},
		{
			name: "contain project/",
			key:  "project/1",
			want: "projects/1",
		},
		{
			name: "contain folder/",
			key:  "folder/1",
			want: "folders/1",
		},
		{
			name: "contain organization/",
			key:  "organization/1",
			want: "organizations/1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseAncestryKey(test.key)
			if err != nil {
				t.Fatalf("parseAncestryKey(%v) = %v, want = nil", test.key, err)
			}
			if got != test.want {
				t.Errorf("parseAncestryKey(%v) = %v, want = %v", test.key, got, test.want)
			}
		})
	}
}

func TestParseAncestryKey_Fail(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "invalid spell",
			key:  "org/1",
		},
		{
			name: "multiple /",
			key:  "folders/123/folders/456",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseAncestryKey(test.key)
			if err == nil {
				t.Fatalf("parseAncestryKey(%v) = %v, want error", test.key, got)
			}
		})
	}
}

func TestUnknownProject(t *testing.T) {
	p := provider.Provider()
	cases := []struct {
		name        string
		data        tpgresource.TerraformResourceData
		asset       *resources.Asset
		v1Responses map[string][]*crmv1.Ancestor
		v3Responses map[string]*crmv3.Project
		want        []string
		wantParent  string
		v3Count     int
		v1Count     int
	}{
		{
			name: "project ID is empty string",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": "",
					"folder_id":  "folders/bar",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			v3Responses: map[string]*crmv3.Project{
				"folders/bar": {Name: "folders/bar", Parent: "organizations/qux"},
			},
			v1Responses: map[string][]*crmv1.Ancestor{},
			want:        []string{"folders/bar", "organizations/qux"},
			wantParent:  "//cloudresourcemanager.googleapis.com/folders/bar",
			v3Count:     1,
		},
		{
			name: "project ID is empty string and no ancestor",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": "",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			v3Responses: map[string]*crmv3.Project{},
			v1Responses: map[string][]*crmv1.Ancestor{},
			want:        []string{"organizations/unknown"},
			wantParent:  "//cloudresourcemanager.googleapis.com/organizations/unknown",
		},
		{
			name: "project ID not exist for bucket",
			data: tfdata.NewFakeResourceData(
				"google_storage_bucket",
				p.ResourcesMap["google_storage_bucket"].Schema,
				map[string]interface{}{},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Bucket",
			},
			v3Responses: map[string]*crmv3.Project{},
			v1Responses: map[string][]*crmv1.Ancestor{},
			want:        []string{"organizations/unknown"},
			wantParent:  "//cloudresourcemanager.googleapis.com/organizations/unknown",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ts := newTestServer(t, c.v1Responses, c.v3Responses)
			defer ts.Close()

			mockV1Client, err := crmv1.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}
			mockV3Client, err := crmv3.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}

			cfg := &transport_tpg.Config{}
			ancestryManager := &manager{
				errorLogger:       zap.NewExample(),
				ancestorCache:     make(map[string][]string),
				resourceManagerV3: mockV3Client,
				resourceManagerV1: mockV1Client,
			}
			// empty cache
			ancestryManager.initAncestryCache(map[string]string{})

			got, gotParent, err := ancestryManager.Ancestors(cfg, c.data, c.asset)
			if err != nil {
				t.Fatalf("Ancestors(%v, %v, %v) = %s, want = nil", cfg, c.data, c.asset, err)
			}
			if gotParent != c.wantParent {
				t.Errorf("Ancestors(%v, %v, %v) parent = %s, want = %s", cfg, c.data, c.asset, gotParent, c.wantParent)
			}
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("Ancestors(%v, %v, %v) returned unexpected diff (-want +got):\n%s", cfg, c.data, c.asset, diff)
			}
			if ts.v3Count != c.v3Count {
				t.Errorf("Ancestors(%v, %v, %v) v3 API called = %d, want = %d", cfg, c.data, c.asset, ts.v3Count, c.v3Count)
			}
			if ts.v1Count != c.v1Count {
				t.Errorf("Ancestors(%v, %v, %v) v1 API called = %d, want = %d", cfg, c.data, c.asset, ts.v1Count, c.v1Count)
			}
		})
	}
}

func TestGetProjectFromResource(t *testing.T) {
	p := provider.Provider()
	cases := []struct {
		name   string
		asset  *resources.Asset
		config *transport_tpg.Config
		d      tpgresource.TerraformResourceData
		resp   *storage.Bucket
		want   string
	}{
		{
			name: "bucket - from cai resource",
			config: &transport_tpg.Config{
				Project: "test-project",
			},
			d: tfdata.NewFakeResourceData(
				"google_storage_bucket_iam_member",
				p.ResourcesMap["google_storage_bucket_iam_member"].Schema,
				map[string]interface{}{
					"bucket": "bucket-name",
				},
			),
			asset: &resources.Asset{
				Type: "storage.googleapis.com/Bucket",
				Resource: &resources.AssetResource{
					Data: map[string]interface{}{
						"project": "resource-project",
					},
				},
			},
			resp: &storage.Bucket{
				ProjectNumber: 123,
			},
			want: "resource-project",
		},
		{
			name:   "bucket - from storage API",
			config: &transport_tpg.Config{Project: "test-project"},
			d: tfdata.NewFakeResourceData(
				"google_storage_bucket_iam_member",
				p.ResourcesMap["google_storage_bucket_iam_member"].Schema,
				map[string]interface{}{
					"bucket": "bucket-name",
				},
			),
			asset: &resources.Asset{
				Type: "storage.googleapis.com/Bucket",
			},
			resp: &storage.Bucket{
				ProjectNumber: 123,
			},
			want: "123",
		},
		{
			name:   "bucket - from provider config",
			config: &transport_tpg.Config{Project: "test-project"},
			d: tfdata.NewFakeResourceData(
				"google_storage_bucket_iam_member",
				p.ResourcesMap["google_storage_bucket_iam_member"].Schema,
				map[string]interface{}{
					"bucket": "bucket-name",
				},
			),
			asset: &resources.Asset{
				Type: "storage.googleapis.com/Bucket",
			},
			resp: nil,
			want: "test-project",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if c.resp == nil {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("no response"))
					return
				}

				payload, err := c.resp.MarshalJSON()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(fmt.Sprintf("failed to MarshalJSON: %s", err)))
					return
				}
				w.Write(payload)
			}))
			defer ts.Close()

			mockClient, err := storage.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}

			ancestryManager := &manager{
				errorLogger:   zap.NewExample(),
				storageClient: mockClient,
			}
			got, err := ancestryManager.getProjectFromResource(c.d, c.config, c.asset)
			if err != nil {
				t.Fatalf("getProjectFromResource() = %s, want = nil", err)
			}
			if got != c.want {
				t.Fatalf("getProjectFromResource() = %s, want = %s", got, c.want)
			}
		})
	}
}

func TestGetAncestorsRetry(t *testing.T) {
	v3Called := 0
	v1Called := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload []byte
		var err error
		if strings.HasPrefix(r.URL.Path, "/v3/") {
			v3Called++
			if v3Called == 1 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				resp := &crmv3.Project{
					Name:   "folders/123",
					Parent: "organizations/456",
				}
				payload, err = resp.MarshalJSON()
			}
		} else if strings.HasPrefix(r.URL.Path, "/v1/") {
			v1Called++
			if v1Called == 1 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				resp := &crmv1.GetAncestryResponse{
					Ancestor: []*crmv1.Ancestor{
						{ResourceId: &crmv1.ResourceId{Id: "abc", Type: "project"}},
						{ResourceId: &crmv1.ResourceId{Id: "888", Type: "organization"}},
					},
				}
				payload, err = resp.MarshalJSON()
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(fmt.Sprintf("no response for url path %s", r.URL.Path)))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to MarshalJSON: %s", err)))
			return
		}
		w.Write(payload)
	}))
	defer ts.Close()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "access v1 API",
			input: "projects/abc",
			want:  []string{"projects/abc", "organizations/888"},
		},
		{
			name:  "access v3 API",
			input: "folders/123",
			want:  []string{"folders/123", "organizations/456"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockV1Client, err := crmv1.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}
			mockV3Client, err := crmv3.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatal(err)
			}
			m := &manager{
				errorLogger:       zap.NewExample(),
				ancestorCache:     make(map[string][]string),
				resourceManagerV3: mockV3Client,
				resourceManagerV1: mockV1Client,
			}
			got, err := m.getAncestorsWithCache(test.input)
			if err != nil {
				t.Fatalf("getAncestorsWithCache(%s) = %s, want = nil", test.input, err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("getAncestorsWithCache(%v) returned unexpected diff (-want +got):\n%s", test.input, diff)
			}
		})
	}
}
