package resolvers

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestConvert_iamBinding(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Error initializing logger %s", err)
	}
	f := "iamBinding.tfplan.json"
	jsonPlan, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("Error parsing %s: %s", f, err)
	}

	idToResourceChangeMap := NewIamAdvancedResolver(logger).Resolve(jsonPlan)

	assert.Equal(t, 1, len(idToResourceChangeMap), "Expected map size is 1")
	assert.Equal(t, 2, len(idToResourceChangeMap["instance_name/google_compute_instance.tgc-iam.name/project/terraform-dev-zhenhuali/zone/us-central1-a/"]), "Expected iam list to be size 2")
	assert.Equal(t, 0, len(idToResourceChangeMap["google_compute_instance_iam_member.foo1"]), "Expected this key to return null")
}

func TestResolveParentsBasic(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Error initializing logger %s", err)
	}
	f := "compute_network.tfplan.json"
	jsonPlan, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("Error parsing %s: %s", f, err)
	}

	parentToChildMap := NewParentResourceResolver(logger).Resolve(jsonPlan)
	assert.Equal(t, "google_vmwareengine_network_peering.peering", parentToChildMap["google_compute_network.peered_network"][0])
	assert.Equal(t, "google_vmwareengine_network_peering.peering", parentToChildMap["google_vmwareengine_network.vmware_network"][0])
}

func TestResolveParentsNested(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Error initializing logger %s", err)
	}
	f := "compute_disk_nestedId.tfplan.json"
	jsonPlan, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("Error parsing %s: %s", f, err)
	}

	parentToChildMap := NewParentResourceResolver(logger).Resolve(jsonPlan)
	assert.Equal(t, "google_compute_disk.secondary", parentToChildMap["google_compute_disk.primary"][0])
}

func TestSortTraversalOrder(t *testing.T) {
	tests := []struct {
		name    string
		graph   map[string][]string
		want    map[int][]string
		wantErr bool
		errDesc string // Optional: check for specific error content
	}{
		{
			name:  "empty_graph",
			graph: map[string][]string{},
			want:  map[int][]string{},
		},
		{
			name:  "single_node_no_edges",
			graph: map[string][]string{"a": {}},
			want:  map[int][]string{0: {"a"}},
		},
		{
			name:  "simple_dag",
			graph: map[string][]string{"a": {"b"}},
			want:  map[int][]string{0: {"a"}, 1: {"b"}},
		},
		{
			name:  "linear_chain",
			graph: map[string][]string{"a": {"b"}, "b": {"c"}},
			want:  map[int][]string{0: {"a"}, 1: {"b"}, 2: {"c"}},
		},
		{
			name:  "fan_out",
			graph: map[string][]string{"a": {"b", "c"}},
			want:  map[int][]string{0: {"a"}, 1: {"b", "c"}},
		},
		{
			name:  "fan_in_diamond",
			graph: map[string][]string{"a": {"b", "c"}, "b": {"d"}, "c": {"d"}},
			want:  map[int][]string{0: {"a"}, 1: {"b", "c"}, 2: {"d"}},
		},
		{
			name:  "disconnected_components",
			graph: map[string][]string{"a": {"b"}, "c": {"d"}},
			want:  map[int][]string{0: {"a", "c"}, 1: {"b", "d"}},
		},
		{
			name:    "simple_cycle",
			graph:   map[string][]string{"a": {"b"}, "b": {"a"}},
			wantErr: true,
			errDesc: "cycle detected",
		},
		{
			name:    "self_loop",
			graph:   map[string][]string{"a": {"a"}},
			wantErr: true,
			errDesc: "cycle detected",
		},
		{
			name:    "cycle_with_entry",
			graph:   map[string][]string{"a": {"b"}, "b": {"c"}, "c": {"b"}},
			wantErr: true,
			errDesc: "cycle detected",
		},
		{
			name:    "cycle_with_entry_and_exit",
			graph:   map[string][]string{"a": {"b"}, "b": {"c"}, "c": {"b", "d"}},
			wantErr: true,
			errDesc: "cycle detected",
		},
		{
			name:  "multiple_paths_different_lengths",
			graph: map[string][]string{"a": {"b", "c"}, "b": {"c"}},
			want:  map[int][]string{0: {"a"}, 1: {"b"}, 2: {"c"}},
		},
		{
			name:  "complex_levels",
			graph: map[string][]string{"a": {"b"}, "d": {"c"}, "b": {"c"}},
			want:  map[int][]string{0: {"a", "d"}, 1: {"b"}, 2: {"c"}},
		},
		{
			name:  "node_only_as_destination",
			graph: map[string][]string{"a": {"b"}, "c": {"b"}},
			want:  map[int][]string{0: {"a", "c"}, 1: {"b"}},
		},
		{
			name: "long_path_determines_level",
			graph: map[string][]string{
				"a": {"b", "x"},
				"b": {"c"},
				"c": {"d"},
				"x": {"d"},
			},
			want: map[int][]string{
				0: {"a"},
				1: {"b", "x"},
				2: {"c"},
				3: {"d"}, // d is level 3 because a->b->c->d
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := sortTraversalOrder(tc.graph)

			if tc.wantErr {
				if err == nil {
					t.Errorf("sortTraversalOrder(%v) succeeded unexpectedly, want error", tc.graph)
				} else if tc.errDesc != "" && !strings.Contains(err.Error(), tc.errDesc) {
					t.Errorf("sortTraversalOrder(%v) returned error %q, want error containing %q", tc.graph, err, tc.errDesc)
				}
			} else {
				if err != nil {
					t.Fatalf("sortTraversalOrder(%v) failed unexpectedly: %v", tc.graph, err)
				}
				// Use cmp.Diff for map comparison. It handles order differences in slices.
				if diff := cmp.Diff(tc.want, got); diff != "" {
					// To make the output more readable in case of empty maps
					if len(tc.want) == 0 && len(got) == 0 {
						return // Treat as equal
					}
					t.Errorf("sortTraversalOrder(%v) returned unexpected diff (-want +got):\n%s", tc.graph, diff)
				}
			}
		})
	}
}
