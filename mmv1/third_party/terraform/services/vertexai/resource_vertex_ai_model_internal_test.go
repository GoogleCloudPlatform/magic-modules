package vertexai

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestVertexAIModelVersionAliasChanges(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		current []interface{}
		desired []interface{}
		want    []interface{}
	}{
		{"create", []interface{}{"default"}, []interface{}{"alias"}, []interface{}{"alias"}},
		{"add", []interface{}{"default", "alias"}, []interface{}{"alias", "new"}, []interface{}{"new"}},
		{"replace", []interface{}{"default", "old"}, []interface{}{"new"}, []interface{}{"new", "-old"}},
		{"clear", []interface{}{"default", "alias"}, nil, []interface{}{"-alias"}},
		{"preserve default", []interface{}{"default"}, nil, nil},
		{"explicit default", []interface{}{"default", "alias"}, []interface{}{"default", "alias"}, nil},
		{"assign default", []interface{}{"alias"}, []interface{}{"default", "alias"}, []interface{}{"default"}},
		{"reorder", []interface{}{"a", "b"}, []interface{}{"b", "a"}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := vertexAiModelVersionAliasChanges(schema.NewSet(schema.HashString, tc.current), schema.NewSet(schema.HashString, tc.desired))
			if !schema.NewSet(schema.HashString, got).Equal(schema.NewSet(schema.HashString, tc.want)) {
				t.Errorf("alias changes = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestVertexAIModelFilterVersionAliases(t *testing.T) {
	t.Parallel()
	aliases := []interface{}{"default", "alias", "drifted"}
	for _, tc := range []struct {
		name       string
		configured cty.Value
		want       []interface{}
	}{
		{"unset", cty.NullVal(cty.Set(cty.String)), aliases},
		{"unknown", cty.UnknownVal(cty.Set(cty.String)), aliases},
		{"partially unknown", cty.SetVal([]cty.Value{cty.UnknownVal(cty.String)}), aliases},
		{"explicit default", cty.SetVal([]cty.Value{cty.StringVal("default")}), aliases},
		{"implicit default", cty.SetVal([]cty.Value{cty.StringVal("alias")}), []interface{}{"alias", "drifted"}},
		{"empty", cty.SetValEmpty(cty.String), []interface{}{"alias", "drifted"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := vertexAiModelFilterVersionAliases(aliases, tc.configured)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("aliases = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestVertexAIModelFlattenVersionAliasesRefresh(t *testing.T) {
	t.Parallel()
	aliases := []interface{}{"default", "alias", "drifted"}
	for _, tc := range []struct {
		name  string
		state []interface{}
		want  []interface{}
	}{
		{"import", nil, aliases},
		{"API default", []interface{}{"default"}, aliases},
		{"implicit default", []interface{}{"alias"}, []interface{}{"alias", "drifted"}},
		{"empty", []interface{}{}, []interface{}{"alias", "drifted"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			attributes := map[string]string{}
			if tc.state != nil {
				attributes["version_aliases.#"] = fmt.Sprint(len(tc.state))
				for _, alias := range tc.state {
					attributes[fmt.Sprintf("version_aliases.%d", schema.HashString(alias))] = alias.(string)
				}
			}
			d := ResourceVertexAIModel().Data(&terraform.InstanceState{ID: "test-model", Attributes: attributes})
			got := flattenVertexAIModelVersionAliases(aliases, d, nil)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("aliases = %v, want %v", got, tc.want)
			}
		})
	}
}
