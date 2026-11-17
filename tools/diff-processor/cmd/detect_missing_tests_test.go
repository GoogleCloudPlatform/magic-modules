package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDetectMissingTestsCmdRun(t *testing.T) {
	sd := diff.SchemaDiff{
		"covered_resource": {
			Fields: map[string]diff.FieldDiff{
				"field_one": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
			},
		},
	}

	t.Run("default-non-exhaustive", func(t *testing.T) {
		var buf bytes.Buffer
		o := &detectMissingTestsOptions{
			computeSchemaDiff: func() diff.SchemaDiff { return sd },
			stdout:            &buf,
		}
		if err := o.run([]string{"../../test-reader/reader/testdata"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var got map[string]*detector.MissingTestInfo
		if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 missing tests in non-exhaustive mode, got %v", got)
		}
	})

	t.Run("exhaustive-json", func(t *testing.T) {
		var buf bytes.Buffer
		o := &detectMissingTestsOptions{
			computeSchemaDiff: func() diff.SchemaDiff { return sd },
			stdout:            &buf,
			exhaustive:        true,
			format:            "json",
		}
		if err := o.run([]string{"../../test-reader/reader/testdata"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var got map[string]*detector.ExhaustiveMissingTestInfo
		if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		info, ok := got["covered_resource"]
		if !ok {
			t.Fatalf("expected covered_resource in exhaustive output (missing MODIFY on field_one)")
		}
		if info.Fields["field_one"].Set.Status != detector.StatusCovered {
			t.Errorf("expected field_one SET to be COVERED, got %+v", info.Fields["field_one"].Set)
		}
		if info.Fields["field_one"].Modify.Status != detector.StatusMissing {
			t.Errorf("expected field_one MODIFY to be MISSING, got %+v", info.Fields["field_one"].Modify)
		}
		if info.Fields["field_one"].Unset.Status != detector.StatusCovered {
			t.Errorf("expected field_one UNSET to be COVERED, got %+v", info.Fields["field_one"].Unset)
		}
	})

	t.Run("exhaustive-markdown", func(t *testing.T) {
		var buf bytes.Buffer
		o := &detectMissingTestsOptions{
			computeSchemaDiff: func() diff.SchemaDiff { return sd },
			stdout:            &buf,
			exhaustive:        true,
			format:            "markdown",
		}
		if err := o.run([]string{"../../test-reader/reader/testdata"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "### `covered_resource`") || !strings.Contains(out, "| `field_one` | Leaf |") {
			t.Errorf("unexpected markdown output:\n%s", out)
		}
	})
}
