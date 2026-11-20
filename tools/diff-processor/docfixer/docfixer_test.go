package docfixer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func block(children ...string) *schema.Schema {
	s := map[string]*schema.Schema{}
	for _, c := range children {
		s[c] = &schema.Schema{Type: schema.TypeString, Optional: true}
	}
	return &schema.Schema{Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: s}}
}

func TestFix(t *testing.T) {
	inner := block("x", "y")
	r := &schema.Resource{Schema: map[string]*schema.Schema{
		"top": {Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"inner": inner,
			"plain": {Type: schema.TypeString, Optional: true},
		}}},
		"other": {Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"inner": block("x", "y"),
		}}},
		"name": {Type: schema.TypeString, Required: true},
	}}
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "anchor and append link",
			in:   "## Argument Reference\n\n* `name` - (Required) Name.\n\n* `top` - (Optional) The top block.\n\nThe `top` block supports:\n\n* `plain` - (Optional) Plain.\n",
			want: "## Argument Reference\n\n* `name` - (Required) Name.\n\n* `top` - (Optional) The top block. Structure is [documented below](#nested_top).\n\n<a name=\"nested_top\"></a>The `top` block supports:\n\n* `plain` - (Optional) Plain.\n",
		},
		{
			name: "links existing plain text",
			in:   "## Argument Reference\n\n* `top` - (Optional) The top block. Structure is documented below.\n\nThe `top` block supports:\n\n* `plain` - (Optional) Plain.\n",
			want: "## Argument Reference\n\n* `top` - (Optional) The top block. Structure is [documented below](#nested_top).\n\n<a name=\"nested_top\"></a>The `top` block supports:\n\n* `plain` - (Optional) Plain.\n",
		},
		{
			name: "nested and shared blocks",
			in:   "## Argument Reference\n\n* `top` - (Optional) Top.\n\n* `other` - (Optional) Other.\n\nThe `top` block supports:\n\n* `inner` - (Optional) Inner.\n\nThe `other` block supports:\n\n* `inner` - (Optional) Inner.\n\nThe `inner` block supports:\n\n* `x` - (Optional) X.\n",
			want: "## Argument Reference\n\n* `top` - (Optional) Top. Structure is [documented below](#nested_top).\n\n* `other` - (Optional) Other. Structure is [documented below](#nested_other).\n\n<a name=\"nested_top\"></a>The `top` block supports:\n\n* `inner` - (Optional) Inner. Structure is [documented below](#nested_other_inner).\n\n<a name=\"nested_other\"></a>The `other` block supports:\n\n* `inner` - (Optional) Inner. Structure is [documented below](#nested_other_inner).\n\n<a name=\"nested_other_inner\"></a>The `inner` block supports:\n\n* `x` - (Optional) X.\n",
		},
		{
			name: "already fixed is unchanged",
			in:   "## Argument Reference\n\n* `top` - (Optional) Top. Structure is [documented below](#nested_top).\n\n<a name=\"nested_top\"></a>The `top` block supports:\n\n* `plain` - (Optional) Plain.\n",
			want: "## Argument Reference\n\n* `top` - (Optional) Top. Structure is [documented below](#nested_top).\n\n<a name=\"nested_top\"></a>The `top` block supports:\n\n* `plain` - (Optional) Plain.\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Fix(tc.in, r)
			if diff := cmp.Diff(tc.want, got.Text); diff != "" {
				t.Errorf("Fix() diff (-want +got):\n%s\nchanges: %v\nnotes: %v", diff, got.Changes, got.Notes)
			}
		})
	}
}
