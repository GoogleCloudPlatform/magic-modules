package docfixer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func nested(fields map[string]*schema.Schema) *schema.Schema {
	return &schema.Schema{Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: fields}}
}

func str() *schema.Schema {
	return &schema.Schema{Type: schema.TypeString, Optional: true}
}

func TestFix(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{
		"name":  {Type: schema.TypeString, Required: true},
		"top":   nested(map[string]*schema.Schema{"plain": str(), "inner": nested(map[string]*schema.Schema{"x": str(), "y": str()})}),
		"other": nested(map[string]*schema.Schema{"inner": nested(map[string]*schema.Schema{"x": str(), "y": str()})}),
	}}
	const (
		args = "## Argument Reference\n\n"
		// The `top` block with and without an anchor, documenting `plain`.
		topBlock      = "\n\nThe `top` block supports:\n\n* `plain` - (Optional) Plain.\n"
		topBlockFixed = "\n\n<a name=\"nested_top\"></a>The `top` block supports:\n\n* `plain` - (Optional) Plain.\n"
	)
	cases := []struct {
		name string
		in   string
		want string // "" means unchanged
	}{
		{
			name: "anchor and append link",
			in:   args + "* `name` - (Required) Name.\n\n* `top` - (Optional) The top block." + topBlock,
			want: args + "* `name` - (Required) Name.\n\n* `top` - (Optional) The top block. Structure is [documented below](#nested_top)." + topBlockFixed,
		},
		{
			name: "links existing plain text",
			in:   args + "* `top` - (Optional) The top block. Structure is documented below." + topBlock,
			want: args + "* `top` - (Optional) The top block. Structure is [documented below](#nested_top)." + topBlockFixed,
		},
		{
			name: "links existing plain text without is",
			in:   args + "* `top` - (Optional) structure documented below" + topBlock,
			want: args + "* `top` - (Optional) structure [documented below](#nested_top)" + topBlockFixed,
		},
		{
			name: "fixes broken link",
			in:   args + "* `top` - (Optional) Structure is [documented below](#nested_wrong)." + topBlock,
			want: args + "* `top` - (Optional) Structure is [documented below](#nested_top)." + topBlockFixed,
		},
		{
			name: "nested and shared blocks",
			in: args + "* `top` - (Optional) Top.\n\n* `other` - (Optional) Other.\n\n" +
				"The `top` block supports:\n\n* `inner` - (Optional) Inner.\n\n" +
				"The `other` block supports:\n\n* `inner` - (Optional) Inner.\n\n" +
				"The `inner` block supports:\n\n* `x` - (Optional) X.\n",
			want: args + "* `top` - (Optional) Top. Structure is [documented below](#nested_top).\n\n* `other` - (Optional) Other. Structure is [documented below](#nested_other).\n\n" +
				"<a name=\"nested_top\"></a>The `top` block supports:\n\n* `inner` - (Optional) Inner. Structure is [documented below](#nested_other_inner).\n\n" +
				"<a name=\"nested_other\"></a>The `other` block supports:\n\n* `inner` - (Optional) Inner. Structure is [documented below](#nested_other_inner).\n\n" +
				"<a name=\"nested_other_inner\"></a>The `inner` block supports:\n\n* `x` - (Optional) X.\n",
		},
		{
			name: "existing link with other wording is kept",
			in:   args + "* `top` - (Optional) Top. See [structure below](#nested_top)." + topBlockFixed,
		},
		{
			name: "existing link split across lines is kept",
			in:   args + "* `top` - (Optional) Top. Structure is [documented\nbelow](#nested_top)." + topBlock,
			want: args + "* `top` - (Optional) Top. Structure is [documented\nbelow](#nested_top)." + topBlockFixed,
		},
		{
			name: "already fixed is unchanged",
			in:   args + "* `top` - (Optional) Top. Structure is [documented below](#nested_top)." + topBlockFixed,
		},
		{
			name: "unknown block is unchanged",
			in:   args + "* `top` - (Optional) Top.\n\nThe `top` block supports:\n\n* `missing` - (Optional) Not in the schema.\n",
		},
		{
			name: "multiple argument reference sections are unchanged",
			in:   "## Argument Reference for A\n\n* `top` - (Optional) Top." + topBlock + "\n## Argument Reference for B\n\n* `top` - (Optional) Top." + topBlock,
		},
		{
			name: "CRLF is preserved",
			in:   "## Argument Reference\r\n\r\n* `top` - (Optional) Top.\r\n\r\nThe `top` block supports:\r\n\r\n* `plain` - (Optional) Plain.\r\n",
			want: "## Argument Reference\r\n\r\n* `top` - (Optional) Top. Structure is [documented below](#nested_top).\r\n\r\n<a name=\"nested_top\"></a>The `top` block supports:\r\n\r\n* `plain` - (Optional) Plain.\r\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			want := tc.want
			if want == "" {
				want = tc.in
			}
			got := Fix(tc.in, r)
			if diff := cmp.Diff(want, got.Text); diff != "" {
				t.Errorf("Fix() diff (-want +got):\n%s\nchanges: %v\nnotes: %v", diff, got.Changes, got.Notes)
			}
			// Fixing is idempotent.
			if again := Fix(got.Text, r); again.Text != got.Text {
				t.Errorf("second Fix() changed the doc again: %v", again.Changes)
			}
		})
	}
}
