package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFixDocsRun(t *testing.T) {
	nested := &schema.Resource{Schema: map[string]*schema.Schema{
		"top": {Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"plain": {Type: schema.TypeString, Optional: true},
		}}},
	}}
	resources := map[string]*schema.Resource{
		"google_handwritten": nested,
		"google_generated":   nested, // no doc in the directory
	}
	const doc = "## Argument Reference\n\n* `top` - (Optional) Top.\n\nThe `top` block supports:\n\n* `plain` - (Optional) Plain.\n"
	const fixed = "## Argument Reference\n\n* `top` - (Optional) Top. Structure is [documented below](#nested_top).\n\n<a name=\"nested_top\"></a>The `top` block supports:\n\n* `plain` - (Optional) Plain.\n"

	cases := []struct {
		name    string
		write   bool
		args    []string
		wantDoc string
		wantOut string
		wantErr bool
	}{
		{name: "dry run on all docs", wantDoc: doc, wantOut: "would change 1 docs"},
		{name: "write all docs", write: true, wantDoc: fixed, wantOut: "changed 1 docs"},
		{name: "write named resource", write: true, args: []string{"google_handwritten"}, wantDoc: fixed, wantOut: "changed 1 docs"},
		{name: "named resource without a doc", args: []string{"google_generated"}, wantDoc: doc, wantErr: true},
		{name: "unknown resource", args: []string{"google_unknown"}, wantDoc: doc, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "handwritten.html.markdown")
			if err := os.WriteFile(path, []byte(doc), 0644); err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			o := fixDocsOptions{
				resources: func() map[string]*schema.Resource { return resources },
				stdout:    &out,
				write:     tc.write,
			}
			err := o.run(dir, tc.args)
			if (err != nil) != tc.wantErr {
				t.Fatalf("run() error = %v, wantErr %v", err, tc.wantErr)
			}
			if got, _ := os.ReadFile(path); string(got) != tc.wantDoc {
				t.Errorf("doc after run() = %q, want %q", got, tc.wantDoc)
			}
			if !strings.Contains(out.String(), tc.wantOut) {
				t.Errorf("run() output = %q, want it to contain %q", out.String(), tc.wantOut)
			}
		})
	}
}
