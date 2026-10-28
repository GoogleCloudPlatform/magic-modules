package google

import (
	"bytes"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type strictInner struct {
	Allowed string `yaml:"allowed,omitempty"`
}

type strictEmbedded struct {
	InlinedField string `yaml:"inlined_field,omitempty"`
}

type strictCustom struct {
	Name           string                 `yaml:"name"`
	DefaultLower   string                 // untagged: defaults to "defaultlower"
	Ignored        string                 `yaml:"-"`
	Inner          *strictInner           `yaml:"inner,omitempty"`
	Items          []strictInner          `yaml:"items,omitempty"`
	NamedItems     map[string]strictInner `yaml:"named_items,omitempty"`
	strictEmbedded `yaml:",inline"`
}

func (s *strictCustom) UnmarshalYAML(value *yaml.Node) error {
	type alias strictCustom
	return DecodeYAMLStrict(value, (*alias)(s))
}

type strictInlineMap struct {
	Name  string                 `yaml:"name"`
	Extra map[string]interface{} `yaml:",inline"`
}

func (s *strictInlineMap) UnmarshalYAML(value *yaml.Node) error {
	type alias strictInlineMap
	return DecodeYAMLStrict(value, (*alias)(s))
}

func decodeWithKnownFields(raw string, out interface{}) error {
	dec := yaml.NewDecoder(bytes.NewReader([]byte(raw)))
	dec.KnownFields(true)
	return dec.Decode(out)
}

func TestDecodeYAMLStrict(t *testing.T) {
	t.Parallel()

	t.Run("accepts valid top-level, nested, inlined, and merged YAML fields", func(t *testing.T) {
		t.Parallel()

		raw := `
base: &base
  allowed: from_anchor
name: example
defaultlower: ok
inlined_field: present
inner:
  <<: *base
items:
  - allowed: item_one
named_items:
  first:
    allowed: map_one
`
		type doc struct {
			Base   strictInner  `yaml:"base"`
			Target strictCustom `yaml:",inline"`
		}
		var got doc
		if err := decodeWithKnownFields(raw, &got.Target); err == nil {
			t.Fatalf("expected error when 'base' anchor key is not on strictCustom")
		}

		validRaw := `
name: example
defaultlower: ok
inlined_field: present
inner:
  allowed: nested_ok
items:
  - allowed: item_one
named_items:
  first:
    allowed: map_one
`
		var valid strictCustom
		if err := decodeWithKnownFields(validRaw, &valid); err != nil {
			t.Fatalf("expected valid YAML to decode cleanly, got: %v", err)
		}
		if valid.Name != "example" || valid.DefaultLower != "ok" || valid.InlinedField != "present" {
			t.Errorf("unexpected decoded values: %+v", valid)
		}
	})

	t.Run("rejects unknown top-level field on custom Unmarshaler", func(t *testing.T) {
		t.Parallel()

		raw := "name: example\npackagepath: \"\"\n"
		var got strictCustom
		err := decodeWithKnownFields(raw, &got)
		if err == nil {
			t.Fatalf("expected unknown field error, got nil")
		}
		if !strings.Contains(err.Error(), "line 2: field packagepath not found") {
			t.Errorf("expected error mentioning line 2 and packagepath, got: %v", err)
		}
	})

	t.Run("rejects yaml:\"-\" field supplied in YAML", func(t *testing.T) {
		t.Parallel()

		raw := "name: example\nignored: should_fail\n"
		var got strictCustom
		err := decodeWithKnownFields(raw, &got)
		if err == nil {
			t.Fatalf("expected unknown field error for yaml:\"-\" field, got nil")
		}
		if !strings.Contains(err.Error(), "field ignored not found") {
			t.Errorf("expected error mentioning ignored, got: %v", err)
		}
	})

	t.Run("rejects unknown field inside nested struct, slice element, or map value", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name    string
			raw     string
			wantSub string
		}{
			{
				name: "nested struct",
				raw: `name: example
inner:
  allowed: ok
  result:
    bad: true
`,
				wantSub: "line 4: field result not found",
			},
			{
				name: "slice element",
				raw: `name: example
items:
  - allowed: ok
    typo_field: bad
`,
				wantSub: "line 4: field typo_field not found",
			},
			{
				name: "map value",
				raw: `name: example
named_items:
  k:
    bad_map_val: true
`,
				wantSub: "line 4: field bad_map_val not found",
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				var got strictCustom
				err := decodeWithKnownFields(tc.raw, &got)
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantSub)
				}
				if !strings.Contains(err.Error(), tc.wantSub) {
					t.Errorf("expected error containing %q, got: %v", tc.wantSub, err)
				}
			})
		}
	})

	t.Run("allows arbitrary keys when struct has inline map", func(t *testing.T) {
		t.Parallel()

		raw := "name: example\ndynamic_key: 123\n"
		var got strictInlineMap
		if err := decodeWithKnownFields(raw, &got); err != nil {
			t.Fatalf("expected inline map struct to allow extra keys, got: %v", err)
		}
		if got.Extra["dynamic_key"] != 123 {
			t.Errorf("expected Extra[dynamic_key]=123, got %+v", got.Extra)
		}
	})
}
