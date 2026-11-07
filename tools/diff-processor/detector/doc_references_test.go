package detector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDetectMissingDocsFollowsDocReferences(t *testing.T) {
	repo := t.TempDir()
	docs := map[string]string{
		"compute_instance": "# google_compute_instance\n\n## Argument Reference\n\n* `name` - (Required) Name.\n",
		"compute_instance_from_machine_image": "# google_compute_instance_from_machine_image\n\n## Argument Reference\n\n* `source_machine_image` - (Required) Image.\n\n" +
			"In addition to these, all arguments from `google_compute_instance` are supported.\n",
		"compute_instance_from_template": "# google_compute_instance_from_template\n\n## Argument Reference\n\n* `source_instance_template` - (Required) Template.\n\n" +
			"In addition to these, all arguments from `google_compute_instance` are supported.\n",
		"container_cluster": "# google_container_cluster\n\n## Argument Reference\n\n* `node_pool` - (Optional) See google_container_node_pool.\n\n" +
			"* `node_config` - (Optional) Structure is [documented below](#nested_node_config).\n\n" +
			"<a name=\"nested_node_config\"></a>The `node_config` block supports:\n\n* `machine_type` - (Optional) Machine type.\n",
		"container_node_pool": "# google_container_node_pool\n\n## Argument Reference\n\n* `node_count` - (Optional) Count.\n\n" +
			"* `node_config` - (Optional) See google_container_cluster for exact schema.\n",
	}
	dir := filepath.Join(repo, "website", "docs", "r")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	for name, content := range docs {
		if err := os.WriteFile(filepath.Join(dir, name+".html.markdown"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	newField := diff.FieldDiff{New: &schema.Schema{Optional: true}}
	schemaDiff := diff.SchemaDiff{
		"google_compute_instance_from_machine_image": {
			Fields: map[string]diff.FieldDiff{
				"source_machine_image": newField,
				"name":                 newField,
				"undocumented_field":   newField,
			},
		},
		"google_compute_instance_from_template": {
			Fields: map[string]diff.FieldDiff{
				"source_instance_template": newField,
				"name":                     newField,
			},
		},
		"google_container_node_pool": {
			Fields: map[string]diff.FieldDiff{
				"node_count":               newField,
				"node_config":              newField,
				"node_config.machine_type": newField,
			},
		},
		"google_container_cluster": {
			Fields: map[string]diff.FieldDiff{
				"node_config":                        newField,
				"node_config.machine_type":           newField,
				"node_pool":                          newField,
				"node_pool.node_count":               newField,
				"node_pool.node_config":              newField,
				"node_pool.node_config.machine_type": newField,
			},
		},
	}

	got, err := DetectMissingDocs(schemaDiff, repo)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]MissingDocDetails{
		"google_compute_instance_from_machine_image": {
			Name:     "google_compute_instance_from_machine_image",
			FilePath: "/website/docs/r/compute_instance_from_machine_image.html.markdown",
			Fields:   []string{"undocumented_field"},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("DetectMissingDocs returned diff (-want, +got):\n%s", diff)
	}
}
