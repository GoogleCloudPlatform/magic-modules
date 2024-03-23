package bigtable

import (
	"reflect"
	"strings"
	"testing"

	"cloud.google.com/go/bigtable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitBigtable_flattenSubsetViewConfig(t *testing.T) {
	cases := map[string]struct {
		avt       bigtable.AuthorizedViewTypeConf
		want      []map[string]interface{}
		wantError string
	}{
		"not subset view": {
			avt: bigtable.AuthorizedViewTypeConf{
				AuthorizedViewType: bigtable.AuthorizedViewTypeUnspecified,
			},
			want:      nil,
			wantError: "not a subset view",
		},
		"empty subset view": {
			avt: bigtable.AuthorizedViewTypeConf{
				AuthorizedViewType: bigtable.AuthorizedViewTypeSubsetView,
				SubsetView:         &bigtable.SubsetViewConf{},
			},
			want: []map[string]interface{}{
				map[string]interface{}{},
			},
			wantError: "",
		},
		"subset view with row prefixes only": {
			avt: bigtable.AuthorizedViewTypeConf{
				AuthorizedViewType: bigtable.AuthorizedViewTypeSubsetView,
				SubsetView: &bigtable.SubsetViewConf{
					RowPrefixes: [][]byte{[]byte("row1"), []byte("row2")},
				},
			},
			want: []map[string]interface{}{
				map[string]interface{}{
					"row_prefixes": []string{"cm93MQ==", "cm93Mg=="},
				},
			},
			wantError: "",
		},
		"subset view with family subsets only": {
			avt: bigtable.AuthorizedViewTypeConf{
				AuthorizedViewType: bigtable.AuthorizedViewTypeSubsetView,
				SubsetView: &bigtable.SubsetViewConf{
					FamilySubsets: map[string]bigtable.FamilySubset{
						"fam1": {
							QualifierPrefixes: [][]byte{[]byte("col")},
						},
						"fam2": {
							Qualifiers: [][]byte{[]byte("col1"), []byte("col2")},
						},
					},
				},
			},
			want: []map[string]interface{}{
				map[string]interface{}{
					"family_subsets": []map[string]interface{}{
						map[string]interface{}{
							"family_name":        "fam1",
							"qualifier_prefixes": []string{"Y29s"},
						}, map[string]interface{}{
							"family_name": "fam2",
							"qualifiers":  []string{"Y29sMQ==", "Y29sMg=="},
						},
					},
				},
			},
			wantError: "",
		},
		"subset view with qualifiers only": {
			avt: bigtable.AuthorizedViewTypeConf{
				AuthorizedViewType: bigtable.AuthorizedViewTypeSubsetView,
				SubsetView: &bigtable.SubsetViewConf{
					FamilySubsets: map[string]bigtable.FamilySubset{
						"fam": {
							Qualifiers: [][]byte{[]byte("col")},
						},
					},
				},
			},
			want: []map[string]interface{}{
				map[string]interface{}{
					"family_subsets": []map[string]interface{}{
						map[string]interface{}{
							"family_name": "fam",
							"qualifiers":  []string{"Y29s"},
						},
					},
				},
			},
			wantError: "",
		},
		"subset view with qualifier prefixes only": {
			avt: bigtable.AuthorizedViewTypeConf{
				AuthorizedViewType: bigtable.AuthorizedViewTypeSubsetView,
				SubsetView: &bigtable.SubsetViewConf{
					FamilySubsets: map[string]bigtable.FamilySubset{
						"fam": {
							QualifierPrefixes: [][]byte{[]byte("col")},
						},
					},
				},
			},
			want: []map[string]interface{}{
				map[string]interface{}{
					"family_subsets": []map[string]interface{}{
						map[string]interface{}{
							"family_name":        "fam",
							"qualifier_prefixes": []string{"Y29s"},
						},
					},
				},
			},
			wantError: "",
		},
		"subset view with empty arrays": {
			avt: bigtable.AuthorizedViewTypeConf{
				AuthorizedViewType: bigtable.AuthorizedViewTypeSubsetView,
				SubsetView: &bigtable.SubsetViewConf{
					RowPrefixes:   [][]byte{},
					FamilySubsets: map[string]bigtable.FamilySubset{},
				},
			},
			want: []map[string]interface{}{
				map[string]interface{}{},
			},
			wantError: "",
		},
	}

	for tn, tc := range cases {
		got, gotErr := flattenSubsetViewConfig(tc.avt)
		if (gotErr != nil && tc.wantError == "") ||
			(gotErr == nil && tc.wantError != "") ||
			(gotErr != nil && !strings.Contains(gotErr.Error(), tc.wantError)) {
			t.Errorf("bad error: %s, got %q, want %q", tn, gotErr, tc.wantError)
		}
		if tc.want != nil && !reflect.DeepEqual(got, tc.want) {
			t.Errorf("bad: %s, got %q, want %q", tn, got, tc.want)
		}
	}
}

func TestUnitBigtable_generateSubsetViewConfig(t *testing.T) {
	cases := map[string]struct {
		sv        []interface{}
		want      *bigtable.SubsetViewConf
		orWant    *bigtable.SubsetViewConf
		wantError string
	}{
		"empty subset view list": {
			sv:        []interface{}{},
			want:      nil,
			orWant:    nil,
			wantError: "empty subset_view list",
		},
		"subset view list with wrong type element": {
			sv: []interface{}{
				"random-string",
			},
			want:      nil,
			orWant:    nil,
			wantError: "element in subset_view list has wrong type",
		},
		"subset view list with nil element": {
			sv: []interface{}{
				nil,
			},
			want:      &bigtable.SubsetViewConf{},
			orWant:    nil,
			wantError: "",
		},
		"subset view list with empty element": {
			sv: []interface{}{
				map[string]interface{}{},
			},
			want:      &bigtable.SubsetViewConf{},
			orWant:    nil,
			wantError: "",
		},
		"subset view list with empty lists": {
			sv: []interface{}{
				map[string]interface{}{
					"row_prefixes":   schema.NewSet(schema.HashString, []interface{}{}),
					"family_subsets": schema.NewSet(schema.HashResource(familySubsetSchemaElem), []interface{}{}),
				},
			},
			want:      &bigtable.SubsetViewConf{},
			orWant:    nil,
			wantError: "",
		},
		"subset view list with row prefixes only": {
			sv: []interface{}{
				map[string]interface{}{
					"row_prefixes": schema.NewSet(schema.HashString, []interface{}{"cm93MQ==", "cm93Mg=="}),
				},
			},
			want: &bigtable.SubsetViewConf{
				RowPrefixes: [][]byte{[]byte("row1"), []byte("row2")},
			},
			orWant: &bigtable.SubsetViewConf{
				RowPrefixes: [][]byte{[]byte("row2"), []byte("row1")},
			},
			wantError: "",
		},
		"subset view list with invalid row prefixes encoding": {
			sv: []interface{}{
				map[string]interface{}{
					"row_prefixes": schema.NewSet(schema.HashString, []interface{}{"#"}),
				},
			},
			want:      nil,
			orWant:    nil,
			wantError: "illegal base64 data",
		},
		"subset view list with empty row prefixes element": {
			sv: []interface{}{
				map[string]interface{}{
					"row_prefixes": schema.NewSet(schema.HashString, []interface{}{""}),
				},
			},
			want: &bigtable.SubsetViewConf{
				RowPrefixes: [][]byte{[]byte("")},
			},
			orWant:    nil,
			wantError: "",
		},
		"subset view list with family subsets only": {
			sv: []interface{}{
				map[string]interface{}{
					"family_subsets": schema.NewSet(schema.HashResource(familySubsetSchemaElem), []interface{}{
						map[string]interface{}{
							"family_name":        "fam1",
							"qualifier_prefixes": schema.NewSet(schema.HashString, []interface{}{"Y29s"}),
						}, map[string]interface{}{
							"family_name": "fam2",
							"qualifiers":  schema.NewSet(schema.HashString, []interface{}{"Y29sMQ==", "Y29sMg=="}),
						},
					}),
				},
			},
			want: &bigtable.SubsetViewConf{
				FamilySubsets: map[string]bigtable.FamilySubset{
					"fam1": {
						QualifierPrefixes: [][]byte{[]byte("col")},
					},
					"fam2": {
						Qualifiers: [][]byte{[]byte("col1"), []byte("col2")},
					},
				},
			},
			orWant: &bigtable.SubsetViewConf{
				FamilySubsets: map[string]bigtable.FamilySubset{
					"fam1": {
						QualifierPrefixes: [][]byte{[]byte("col")},
					},
					"fam2": {
						Qualifiers: [][]byte{[]byte("col2"), []byte("col1")},
					},
				},
			},
			wantError: "",
		},
		"subset view list with qualifiers only": {
			sv: []interface{}{
				map[string]interface{}{
					"family_subsets": schema.NewSet(schema.HashResource(familySubsetSchemaElem), []interface{}{
						map[string]interface{}{
							"family_name": "fam",
							"qualifiers":  schema.NewSet(schema.HashString, []interface{}{"Y29sMQ==", "Y29sMg=="}),
						},
					}),
				},
			},
			want: &bigtable.SubsetViewConf{
				FamilySubsets: map[string]bigtable.FamilySubset{
					"fam": {
						Qualifiers: [][]byte{[]byte("col1"), []byte("col2")},
					},
				},
			},
			orWant: &bigtable.SubsetViewConf{
				FamilySubsets: map[string]bigtable.FamilySubset{
					"fam": {
						Qualifiers: [][]byte{[]byte("col2"), []byte("col1")},
					},
				},
			},
			wantError: "",
		},
		"subset view list with qualifier prefixes only": {
			sv: []interface{}{
				map[string]interface{}{
					"family_subsets": schema.NewSet(schema.HashResource(familySubsetSchemaElem), []interface{}{
						map[string]interface{}{
							"family_name":        "fam",
							"qualifier_prefixes": schema.NewSet(schema.HashString, []interface{}{"Y29s"}),
						},
					}),
				},
			},
			want: &bigtable.SubsetViewConf{
				FamilySubsets: map[string]bigtable.FamilySubset{
					"fam": {
						QualifierPrefixes: [][]byte{[]byte("col")},
					},
				},
			},
			orWant:    nil,
			wantError: "",
		},
		"subset view list with invalid qualifiers encoding": {
			sv: []interface{}{
				map[string]interface{}{
					"family_subsets": schema.NewSet(schema.HashResource(familySubsetSchemaElem), []interface{}{
						map[string]interface{}{
							"family_name": "fam",
							"qualifiers":  schema.NewSet(schema.HashString, []interface{}{"#"}),
						},
					}),
				},
			},
			want:      nil,
			orWant:    nil,
			wantError: "illegal base64 data",
		},
		"subset view list with invalid qualifier prefixes encoding": {
			sv: []interface{}{
				map[string]interface{}{
					"family_subsets": schema.NewSet(schema.HashResource(familySubsetSchemaElem), []interface{}{
						map[string]interface{}{
							"family_name":        "fam",
							"qualifier_prefixes": schema.NewSet(schema.HashString, []interface{}{"#"}),
						},
					}),
				},
			},
			want:      nil,
			orWant:    nil,
			wantError: "illegal base64 data",
		},
	}

	for tn, tc := range cases {
		got, gotErr := generateSubsetViewConfig(tc.sv)
		if (gotErr != nil && tc.wantError == "") ||
			(gotErr == nil && tc.wantError != "") ||
			(gotErr != nil && !strings.Contains(gotErr.Error(), tc.wantError)) {
			t.Errorf("bad error: %s, got %q, want %q", tn, gotErr, tc.wantError)
		}
		if tc.want != nil && !(reflect.DeepEqual(got, tc.want) || reflect.DeepEqual(got, tc.orWant)) {
			t.Errorf("bad: %s, got %q, want %q", tn, got, tc.want)
		}
	}
}
