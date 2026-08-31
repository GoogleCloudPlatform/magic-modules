package provider

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

// Synthetic properties covering every branch (scalar, nested object,
// array-of-object, map-of-object, arbitrary-depth nesting) exercised by both
// the schema and field-mapping ephemeral templates below.
func newEphemeralTestResource() (res *api.Resource, effectiveIntelligenceConfig *api.Type) {
	deepest := &api.Type{Name: "innerValue", Type: "String"}
	innerObject := &api.Type{
		Name:       "inner",
		Type:       "NestedObject",
		Properties: []*api.Type{deepest},
	}

	effectiveIntelligenceConfig = &api.Type{
		Name:   "effectiveIntelligenceConfig",
		Type:   "NestedObject",
		Output: true,
		Properties: []*api.Type{
			{Name: "intelligenceConfig", Type: "String", Output: true},
			{Name: "effectiveEdition", Type: "String", Output: true},
			innerObject,
		},
	}

	props := []*api.Type{
		{Name: "name", Type: "String", Required: true},
		{Name: "editionConfig", Type: "String", DefaultFromApi: true},
		{
			Name:   "updateTime",
			Type:   "String",
			Output: true,
			// Multi-line, matching real MMv1 YAML ("description: |" blocks) —
			// regression test for a real bug where this broke Go string
			// literal syntax when naively interpolated (see
			// ephemeral-resource-generator-phase0-writeup.md's Phase 2 findings).
			Description: "The time at which the resource was last updated.\nMultiple lines, like real YAML descriptions almost always are.",
		},
		effectiveIntelligenceConfig,
		{
			Name:   "items",
			Type:   "Array",
			Output: true,
			ItemType: &api.Type{
				Type: "NestedObject",
				Properties: []*api.Type{
					{Name: "itemName", Type: "String", Output: true},
				},
			},
		},
		{
			Name:   "labelsLike",
			Type:   "Map",
			Output: true,
			ValueType: &api.Type{
				Type: "NestedObject",
				Properties: []*api.Type{
					{Name: "value", Type: "String", Output: true},
				},
			},
		},
	}

	res = &api.Resource{
		Name:       "TestResource",
		Properties: props,
		// Mirrors storagecontrol/ProjectIntelligenceConfig.yaml's self_link shape:
		// "name" is the only URL identifier, everything else is read-only output.
		SelfLink: "projects/{{name}}/locations/global/intelligenceConfig",
	}
	p := &api.Product{
		Name:        "TestProduct",
		Versions:    []*product.Version{{Name: "ga"}},
		Objects:     []*api.Resource{res},
		PackagePath: "products/testproduct",
	}
	p.Version = p.Versions[0]
	res.SetDefault(p)

	return res, effectiveIntelligenceConfig
}

// Parses the given wrapper template alongside schema_property_ephemeral_fw.go.tmpl.
// Parses/executes directly via text/template rather than TemplateData.GenerateFile,
// because GenerateFile calls glog.Exit (fatal, kills the whole test binary) on
// template errors instead of returning a normal error.
func parseEphemeralTemplate(t *testing.T, wrapperFile string, extraFiles ...string) *template.Template {
	t.Helper()

	// go test's working directory is this package's directory (mmv1/provider),
	// but template paths are rooted at mmv1/.
	templateFS := os.DirFS("..")
	funcMap := template.FuncMap{
		"TemplatePath": func() string { return "test" },
	}
	for k, v := range google.TemplateFunctions(templateFS) {
		funcMap[k] = v
	}

	files := append([]string{"templates/terraform/" + wrapperFile}, extraFiles...)
	tmpl, err := template.New("wrapper").Funcs(funcMap).ParseFS(templateFS, files...)
	if err != nil {
		t.Fatalf("template failed to parse: %v", err)
	}
	return tmpl
}

func renderEphemeral(tmpl *template.Template, wrapperFile string, input any) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, wrapperFile, input)
	return buf.String(), err
}

func TestSchemaFieldsEphemeralFWTemplate(t *testing.T) {
	res, effectiveIntelligenceConfig := newEphemeralTestResource()
	wrapperFile := "ephemeral_schema_field_test_wrapper.go.tmpl"
	tmpl := parseEphemeralTemplate(t, wrapperFile, "templates/terraform/schema_property_ephemeral_fw.go.tmpl")

	rendered := map[string]string{}
	for _, prop := range res.OrderProperties(res.AllUserProperties()) {
		out, err := renderEphemeral(tmpl, wrapperFile, prop)
		if err != nil {
			t.Fatalf("rendering property %q failed: %v", prop.Name, err)
		}
		t.Logf("property %q rendered:\n%s", prop.Name, out)
		rendered[prop.Name] = out
	}

	// "name" is the only self_link URL variable — must be Required. Everything
	// else is Computed-only regardless of its role in the base resource: this
	// is what distinguishes an ephemeral resource's schema from a regular
	// resource's — there's no create/update lifecycle, only a read.
	if !strings.Contains(rendered["name"], "Required: true,") {
		t.Errorf("expected \"name\" (the only URL var) to be Required, got:\n%s", rendered["name"])
	}
	// Regression test: "updateTime" has a multi-line Description, matching real
	// MMv1 YAML. Description must never be emitted at all (see the template's
	// own comment for why) — this used to break Go string literal syntax when
	// naively interpolated, a real bug only caught by generating the actual
	// resource (whose YAML descriptions are multi-line), not by this test suite,
	// since this synthetic property was originally single-line.
	if strings.Contains(rendered["updateTime"], "Description:") {
		t.Errorf("expected no Description field to be emitted at all, got:\n%s", rendered["updateTime"])
	}
	// editionConfig is `default_from_api: true` in the base resource (would be
	// Optional+Computed for a real resource), but it's not a URL var, so for
	// the ephemeral case it must be Computed-only — confirmed against the
	// Phase 0 hand-written spike, which used Computed-only for this exact field.
	if !strings.Contains(rendered["editionConfig"], "Computed: true,") || strings.Contains(rendered["editionConfig"], "Optional: true,") {
		t.Errorf("expected \"editionConfig\" to be Computed-only (not Optional), got:\n%s", rendered["editionConfig"])
	}

	out, err := renderEphemeral(tmpl, wrapperFile, effectiveIntelligenceConfig)
	if err != nil {
		t.Fatalf("rendering nested property failed: %v", err)
	}
	for _, want := range []string{
		"schema.ObjectAttribute{",
		"AttributeTypes: map[string]attr.Type{",
		`"intelligence_config": types.StringType`,
		`"inner": types.ObjectType{AttrTypes: map[string]attr.Type{`,
		`"inner_value": types.StringType`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected rendered output to contain %q, got:\n%s", want, out)
		}
	}
	for _, mustNotContain := range []string{"SingleNestedAttribute", "NestedAttributeObject{", "schema.Block"} {
		if strings.Contains(out, mustNotContain) {
			t.Errorf("expected rendered output NOT to contain %q (NestedAttribute/Block are known-broken for ephemeral resources), got:\n%s", mustNotContain, out)
		}
	}
}

func TestFieldMappingEphemeralFWTemplate(t *testing.T) {
	res, effectiveIntelligenceConfig := newEphemeralTestResource()
	wrapperFile := "ephemeral_field_mapping_test_wrapper.go.tmpl"
	tmpl := parseEphemeralTemplate(t, wrapperFile, "templates/terraform/schema_property_ephemeral_fw.go.tmpl")

	for _, prop := range res.OrderProperties(res.AllUserProperties()) {
		out, err := renderEphemeral(tmpl, wrapperFile, prop)
		if err != nil {
			t.Fatalf("rendering property %q failed: %v", prop.Name, err)
		}
		t.Logf("property %q rendered:\n%s", prop.Name, out)
	}

	out, err := renderEphemeral(tmpl, wrapperFile, effectiveIntelligenceConfig)
	if err != nil {
		t.Fatalf("rendering nested property failed: %v", err)
	}
	for _, want := range []string{
		`fwutils.EphemeralFieldMapping{`,
		`ApiName: "effectiveIntelligenceConfig"`,
		`AttrType: types.ObjectType{AttrTypes: map[string]attr.Type{`,
		`Fields: map[string]fwutils.EphemeralFieldMapping{`,
		`"intelligence_config": fwutils.EphemeralFieldMapping{`,
		`ApiName: "intelligenceConfig"`,
		`"inner": fwutils.EphemeralFieldMapping{`,
		`ApiName: "inner"`,
		`"inner_value": fwutils.EphemeralFieldMapping{`,
		`ApiName: "innerValue"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected rendered output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestEphemeralUrlTemplate(t *testing.T) {
	wrapperFile := "ephemeral_url_test_wrapper.go.tmpl"
	tmpl := parseEphemeralTemplate(t, wrapperFile, "templates/terraform/ephemeral_url_fw.go.tmpl")

	// Mirrors the action-based spike (publicca/ExternalAccountKey.yaml):
	// create_url with one well-known var (project) and one custom var (location).
	actionRes := &api.Resource{
		Name:      "ExternalAccountKey",
		CreateUrl: "projects/{{project}}/locations/{{location}}/externalAccountKeys",
	}
	out, err := renderEphemeral(tmpl, wrapperFile, map[string]any{
		"Resource": actionRes,
		"Url":      actionRes.CreateUrl,
	})
	if err != nil {
		t.Fatalf("rendering action-based URL failed: %v", err)
	}
	t.Logf("action-based URL rendered:\n%s", out)
	wantAction := `fmt.Sprintf("%sprojects/%s/locations/%s/externalAccountKeys", transport_tpg.BaseUrl(Product, config), project, data.Location.ValueString())`
	if strings.TrimSpace(out) != wantAction {
		t.Errorf("action-based URL:\ngot:  %s\nwant: %s", strings.TrimSpace(out), wantAction)
	}

	// Mirrors the GET-based spike (storagecontrol/ProjectIntelligenceConfig):
	// self_link with one custom var (name), no well-known project/region/zone.
	getRes := &api.Resource{
		Name:     "ProjectIntelligenceConfig",
		SelfLink: "projects/{{name}}/locations/global/intelligenceConfig",
	}
	getUrl := getRes.SelfLinkUri()
	out, err = renderEphemeral(tmpl, wrapperFile, map[string]any{
		"Resource": getRes,
		"Url":      getUrl,
	})
	if err != nil {
		t.Fatalf("rendering GET-based URL failed: %v", err)
	}
	t.Logf("GET-based URL rendered:\n%s", out)
	wantGet := `fmt.Sprintf("%sprojects/%s/locations/global/intelligenceConfig", transport_tpg.BaseUrl(Product, config), data.Name.ValueString())`
	if strings.TrimSpace(out) != wantGet {
		t.Errorf("GET-based URL:\ngot:  %s\nwant: %s", strings.TrimSpace(out), wantGet)
	}
}
