package chromepolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkdiag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"google.golang.org/api/chromepolicy/v1"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// policies uses DynamicAttribute (not SetNestedBlock) because the Plugin Framework doesn't support dynamic types in collections (terraform-plugin-framework#973).
// This allows native HCL types in value without jsonencode(). We use a list (not set) to avoid dirty-state issues (terraform-plugin-framework#1008); order is normalized during Read by matching existing state order.

var (
	_ resource.Resource                = &googleChromePoliciesResource{}
	_ resource.ResourceWithConfigure   = &googleChromePoliciesResource{}
	_ resource.ResourceWithImportState = &googleChromePoliciesResource{}
)

func NewGoogleChromePoliciesResource() resource.Resource {
	return &googleChromePoliciesResource{}
}

type googleChromePoliciesResource struct {
	config *transport_tpg.Config
}

type chromePoliciesModel struct {
	Id           types.String  `tfsdk:"id"`
	CustomerId   types.String  `tfsdk:"customer_id"`
	OrgUnitId    types.String  `tfsdk:"org_unit_id"`
	GroupId      types.String  `tfsdk:"group_id"`
	SchemaFilter types.String  `tfsdk:"schema_filter"`
	Policies     types.Dynamic `tfsdk:"policies"`
}

func (r *googleChromePoliciesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chrome_policies"
}

func (r *googleChromePoliciesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T.", req.ProviderData))
		return
	}
	r.config = config
}

func (r *googleChromePoliciesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages all Chrome policies matching a schema filter for an org unit or group. Policies not in config are removed (inherited). Requires the chrome.management.policy OAuth scope.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"customer_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The customer ID. Defaults to my_customer.",
				Default:     stringdefault.StaticString("my_customer"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org_unit_id": schema.StringAttribute{
				Optional:    true,
				Description: "The target org unit. Exactly one of org_unit_id or group_id must be set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "The target group. Exactly one of org_unit_id or group_id must be set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema_filter": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The schema filter defining the authoritative scope.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policies": schema.DynamicAttribute{
				Optional:    true,
				Description: "List of policies to enforce. Each entry has schema, value, and optional additional_target_keys.",
			},
		},
	}
}

//CRUD

func (r *googleChromePoliciesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan chromePoliciesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.validateTarget(&plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	customerID := plan.CustomerId.ValueString()
	kind, targetID, targetResource := r.resolveTarget(&plan)

	policies := r.expandPolicies(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := plan.SchemaFilter.ValueString()
	if filter == "" {
		if len(policies) != 1 {
			resp.Diagnostics.AddError("Missing schema_filter",
				"schema_filter is required when zero or multiple policies are defined.")
			return
		}
		filter = policies[0].(map[string]interface{})["schema"].(string)
		plan.SchemaFilter = types.StringValue(filter)
	}

	// Validate all policy schemas match the filter.
	r.validateSchemasMatchFilter(policies, filter, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, diags := r.getClient()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(policies) > 0 {
		sc := newSchemaCache(svc.Customers.PolicySchemas, customerID)
		r.appendSdkDiags(validatePolicies(ctx, policies, sc), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		r.appendSdkDiags(batchModifyPolicies(ctx, svc.Customers.Policies, customerID, targetResource, policies), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	plan.Id = types.StringValue(chromePoliciesResourceID(customerID, kind, targetID, filter))
	tflog.Info(ctx, "Created Chrome Policies", map[string]interface{}{"id": plan.Id.ValueString()})

	// DynamicAttribute requires the provider to return the same type structure as the plan after apply.
	// We preserve plan values here; Read/refresh reconciles with the API.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *googleChromePoliciesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state chromePoliciesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoState(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *googleChromePoliciesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state chromePoliciesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	customerID := plan.CustomerId.ValueString()
	kind, targetID, targetResource := r.resolveTarget(&plan)
	filter := plan.SchemaFilter.ValueString()

	svc, diags := r.getClient()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newPolicies := r.expandPolicies(ctx, &plan, &resp.Diagnostics)
	oldPolicies := r.expandPolicies(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	r.validateSchemasMatchFilter(newPolicies, filter, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compare as sets (by identity key), ignoring list order. DynamicAttribute may trigger Update on reorder alone — we diff here to skip API calls when nothing actually changed.
	oldByKey := policyMapByKey(oldPolicies)
	newByKey := policyMapByKey(newPolicies)

	if !policySetsEqual(oldByKey, newByKey) {
		if len(newPolicies) > 0 {
			sc := newSchemaCache(svc.Customers.PolicySchemas, customerID)
			r.appendSdkDiags(validatePolicies(ctx, newPolicies, sc), &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		// Removed policies get inherited (reverted to parent).
		var inheritRequests []*chromepolicy.GoogleChromePolicyVersionsV1InheritOrgUnitPolicyRequest
		for key, pol := range oldByKey {
			if _, exists := newByKey[key]; !exists {
				inheritRequests = append(inheritRequests, &chromepolicy.GoogleChromePolicyVersionsV1InheritOrgUnitPolicyRequest{
					PolicyTargetKey: buildPolicyTargetKey(targetResource, pol),
					PolicySchema:    identityFromPolicy(pol).SchemaName,
				})
			}
		}
		r.appendSdkDiags(inheritPolicies(ctx, svc.Customers.Policies, customerID, inheritRequests), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		// New or changed policies get applied.
		var modifiedPolicies []interface{}
		for key, newPol := range newByKey {
			oldPol, existed := oldByKey[key]
			if !existed || !policyValuesEqual(oldPol, newPol) {
				modifiedPolicies = append(modifiedPolicies, newPol)
			}
		}
		if len(modifiedPolicies) > 0 {
			r.appendSdkDiags(batchModifyPolicies(ctx, svc.Customers.Policies, customerID, targetResource, modifiedPolicies), &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	// Preserve plan values for DynamicAttribute type consistency (same as Create).
	plan.Id = types.StringValue(chromePoliciesResourceID(customerID, kind, targetID, filter))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *googleChromePoliciesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state chromePoliciesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	customerID := state.CustomerId.ValueString()
	_, _, targetResource := r.resolveTarget(&state)
	filter := state.SchemaFilter.ValueString()

	svc, diags := r.getClient()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entries, sdkDiags := resolveDirectlySetPolicies(ctx, svc.Customers.Policies, customerID, filter, targetResource)
	r.appendSdkDiags(sdkDiags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(entries) == 0 {
		return
	}

	var requests []*chromepolicy.GoogleChromePolicyVersionsV1InheritOrgUnitPolicyRequest
	for _, entry := range entries {
		requests = append(requests, &chromepolicy.GoogleChromePolicyVersionsV1InheritOrgUnitPolicyRequest{
			PolicyTargetKey: &chromepolicy.GoogleChromePolicyVersionsV1PolicyTargetKey{
				TargetResource:       targetResource,
				AdditionalTargetKeys: entry.Identity.AdditionalTargetKeys,
			},
			PolicySchema: entry.Identity.SchemaName,
		})
	}

	r.appendSdkDiags(inheritPolicies(ctx, svc.Customers.Policies, customerID, requests), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted Chrome Policies for %s: inherited %d policies", targetResource, len(requests)))
}

func (r *googleChromePoliciesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Format: {customerId}/{orgunits|groups}/{targetId}/{schemaFilter}
	idParts := strings.SplitN(req.ID, "/", 4)
	if len(idParts) != 4 {
		resp.Diagnostics.AddError("Invalid import ID",
			fmt.Sprintf("Expected format: {customerId}/{orgunits|groups}/{targetId}/{schemaFilter}, got: %s", req.ID))
		return
	}

	customerID := idParts[0]
	kind := chromePolicyTargetKind(idParts[1])
	targetID := idParts[2]
	filter := idParts[3]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("customer_id"), customerID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("schema_filter"), filter)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

	switch kind {
	case targetOrgUnit:
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_unit_id"), targetID)...)
	case targetGroup:
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), targetID)...)
	default:
		resp.Diagnostics.AddError("Invalid target kind",
			fmt.Sprintf("Expected \"orgunits\" or \"groups\", got: %s", idParts[1]))
	}
}

//Helpers

func (r *googleChromePoliciesResource) getClient() (*chromepolicy.Service, diag.Diagnostics) {
	var diags diag.Diagnostics
	svc := r.config.NewChromePolicyClient(
		fwtransport.GenerateFrameworkUserAgentString(&fwmodels.ProviderMetaModel{}, r.config.UserAgent),
	)
	if svc == nil {
		diags.AddError("Client Error", "Failed to create Chrome Policy client")
	}
	return svc, diags
}

// appendSdkDiags converts SDKv2 diag.Diagnostics to Framework diag.Diagnostics.
func (r *googleChromePoliciesResource) appendSdkDiags(sdkDiags sdkdiag.Diagnostics, fwDiags *diag.Diagnostics) {
	for _, d := range sdkDiags {
		if d.Severity == sdkdiag.Error {
			fwDiags.AddError(d.Summary, d.Detail)
		} else {
			fwDiags.AddWarning(d.Summary, d.Detail)
		}
	}
}

func (r *googleChromePoliciesResource) validateTarget(model *chromePoliciesModel, diags *diag.Diagnostics) {
	hasOU := !model.OrgUnitId.IsNull() && model.OrgUnitId.ValueString() != ""
	hasGroup := !model.GroupId.IsNull() && model.GroupId.ValueString() != ""
	if !hasOU && !hasGroup {
		diags.AddError("Missing target", "Exactly one of org_unit_id or group_id must be set.")
	}
	if hasOU && hasGroup {
		diags.AddError("Conflicting targets", "Only one of org_unit_id or group_id can be set.")
	}
}

func (r *googleChromePoliciesResource) validateSchemasMatchFilter(policies []interface{}, filter string, diags *diag.Diagnostics) {
	for _, p := range policies {
		name := p.(map[string]interface{})["schema"].(string)
		if !schemaNameMatchesFilter(name, filter) {
			diags.AddError("Schema mismatch",
				fmt.Sprintf("Policy schema %q does not match schema_filter %q", name, filter))
			return
		}
	}
}

func (r *googleChromePoliciesResource) resolveTarget(model *chromePoliciesModel) (chromePolicyTargetKind, string, string) {
	if !model.OrgUnitId.IsNull() && model.OrgUnitId.ValueString() != "" {
		targetID := strings.TrimPrefix(model.OrgUnitId.ValueString(), "id:")
		return targetOrgUnit, targetID, chromePolicyTargetResource(targetOrgUnit, targetID)
	}
	if !model.GroupId.IsNull() && model.GroupId.ValueString() != "" {
		return targetGroup, model.GroupId.ValueString(), chromePolicyTargetResource(targetGroup, model.GroupId.ValueString())
	}
	return "", "", ""
}

//Dynamic <-> native Go conversion

// expandPolicies converts the policies DynamicAttribute into []interface{} for the shared helpers.
func (r *googleChromePoliciesResource) expandPolicies(ctx context.Context, model *chromePoliciesModel, diags *diag.Diagnostics) []interface{} {
	if model.Policies.IsNull() || model.Policies.IsUnknown() {
		return nil
	}

	underlying := model.Policies.UnderlyingValue()
	tupleVal, ok := underlying.(basetypes.TupleValue)
	if !ok {
		diags.AddError("Invalid policies type", fmt.Sprintf("Expected a list/tuple, got: %T", underlying))
		return nil
	}

	var result []interface{}
	for _, elem := range tupleVal.Elements() {
		objVal, ok := elem.(basetypes.ObjectValue)
		if !ok {
			diags.AddError("Invalid policy entry", fmt.Sprintf("Expected an object, got: %T", elem))
			return nil
		}

		attrs := objVal.Attributes()
		pol := map[string]interface{}{}

		if sn, ok := attrs["schema"].(basetypes.StringValue); ok {
			pol["schema"] = sn.ValueString()
		} else {
			diags.AddError("Missing schema", "Each policy must have a schema string.")
			return nil
		}

		if sv, ok := attrs["value"]; ok && sv != nil {
			elems := attrElements(sv)
			if elems == nil {
				diags.AddError("Invalid value", fmt.Sprintf("Expected object or map, got: %T", sv))
				return nil
			}
			schemaValues := make(map[string]interface{}, len(elems))
			for k, attrVal := range elems {
				schemaValues[k] = attrValueToJSON(attrVal)
			}
			pol["value"] = schemaValues
		}

		if atk, ok := attrs["additional_target_keys"]; ok && atk != nil && !atk.IsNull() {
			if elems := attrElements(atk); len(elems) > 0 {
				atkMap := make(map[string]interface{}, len(elems))
				for k, av := range elems {
					if s, ok := av.(basetypes.StringValue); ok {
						atkMap[k] = s.ValueString()
					}
				}
				pol["additional_target_keys"] = atkMap
			}
		}

		result = append(result, pol)
	}

	return result
}

// attrValueToJSON converts a Framework attr.Value to a JSON string for the API.
func attrValueToJSON(v attr.Value) string {
	b, _ := json.Marshal(attrValueToNative(v))
	return string(b)
}

// attrValueToNative converts a Framework attr.Value to a native Go value.
func attrValueToNative(v attr.Value) interface{} {
	switch val := v.(type) {
	case basetypes.StringValue:
		return val.ValueString()
	case basetypes.BoolValue:
		return val.ValueBool()
	case basetypes.NumberValue:
		f, _ := val.ValueBigFloat().Float64()
		return f
	case basetypes.Int64Value:
		return val.ValueInt64()
	case basetypes.Float64Value:
		return val.ValueFloat64()
	case basetypes.ListValue:
		return elementsToNative(val.Elements())
	case basetypes.TupleValue:
		return elementsToNative(val.Elements())
	case basetypes.SetValue:
		return elementsToNative(val.Elements())
	default:
		if elems := attrElements(v); elems != nil {
			m := make(map[string]interface{}, len(elems))
			for k, av := range elems {
				m[k] = attrValueToNative(av)
			}
			return m
		}
		return fmt.Sprintf("%v", v)
	}
}

func elementsToNative(elements []attr.Value) []interface{} {
	items := make([]interface{}, len(elements))
	for i, el := range elements {
		items[i] = attrValueToNative(el)
	}
	return items
}

// attrElements extracts key-value pairs from an ObjectValue or MapValue.
func attrElements(v attr.Value) map[string]attr.Value {
	switch val := v.(type) {
	case basetypes.ObjectValue:
		return val.Attributes()
	case basetypes.MapValue:
		return val.Elements()
	default:
		return nil
	}
}

// nativeToAttrValue converts a native Go value (from JSON unmarshal) to a Framework attr.Value.
func nativeToAttrValue(ctx context.Context, v interface{}) attr.Value {
	switch val := v.(type) {
	case string:
		return types.StringValue(val)
	case bool:
		return types.BoolValue(val)
	case float64:
		return types.NumberValue(big.NewFloat(val))
	case []interface{}:
		elements := make([]attr.Value, len(val))
		elementTypes := make([]attr.Type, len(val))
		for i, el := range val {
			fwVal := nativeToAttrValue(ctx, el)
			elements[i] = fwVal
			elementTypes[i] = fwVal.Type(ctx)
		}
		tv, _ := types.TupleValue(elementTypes, elements)
		return tv
	case map[string]interface{}:
		attrs := make(map[string]attr.Value, len(val))
		attrTypes := make(map[string]attr.Type, len(val))
		for k, av := range val {
			fwVal := nativeToAttrValue(ctx, av)
			attrs[k] = fwVal
			attrTypes[k] = fwVal.Type(ctx)
		}
		ov, _ := types.ObjectValue(attrTypes, attrs)
		return ov
	case nil:
		return types.StringNull()
	default:
		return types.StringValue(fmt.Sprintf("%v", v))
	}
}

// readIntoState reads policies from the API and populates the model.
func (r *googleChromePoliciesResource) readIntoState(ctx context.Context, model *chromePoliciesModel, diags *diag.Diagnostics) {
	customerID := model.CustomerId.ValueString()
	kind, targetID, targetResource := r.resolveTarget(model)
	filter := model.SchemaFilter.ValueString()

	model.Id = types.StringValue(chromePoliciesResourceID(customerID, kind, targetID, filter))

	svc, d := r.getClient()
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	entries, sdkDiags := resolveDirectlySetPolicies(ctx, svc.Customers.Policies, customerID, filter, targetResource)
	r.appendSdkDiags(sdkDiags, diags)
	if diags.HasError() {
		return
	}

	// Sort to match existing state order so reordering in config doesn't cause spurious diffs.
	existingOrder := make(map[string]int)
	existingPolicies := r.expandPolicies(ctx, model, diags)
	for i, p := range existingPolicies {
		pol := p.(map[string]interface{})
		existingOrder[identityFromPolicy(pol).key()] = i
	}
	sort.SliceStable(entries, func(i, j int) bool {
		oi, okI := existingOrder[entries[i].Identity.key()]
		oj, okJ := existingOrder[entries[j].Identity.key()]
		switch {
		case okI && okJ:
			return oi < oj
		case okI:
			return true
		case okJ:
			return false
		default:
			return entries[i].Identity.key() < entries[j].Identity.key()
		}
	})

	// Build the policies list as Framework dynamic values.
	policyElements := make([]attr.Value, 0, len(entries))
	policyElementTypes := make([]attr.Type, 0, len(entries))

	for _, entry := range entries {
		var rawValues map[string]interface{}
		if err := json.Unmarshal(entry.Value.Value, &rawValues); err != nil {
			diags.AddError("Failed to parse policy values", err.Error())
			return
		}

		// Build value as an object with native types.
		svAttrs := make(map[string]attr.Value, len(rawValues))
		svAttrTypes := make(map[string]attr.Type, len(rawValues))
		for k, v := range rawValues {
			fwVal := nativeToAttrValue(ctx, v)
			svAttrs[k] = fwVal
			svAttrTypes[k] = fwVal.Type(ctx)
		}
		schemaValuesObj, objDiags := types.ObjectValue(svAttrTypes, svAttrs)
		diags.Append(objDiags...)
		if diags.HasError() {
			return
		}

		// Build the policy object.
		polAttrs := map[string]attr.Value{
			"schema": types.StringValue(entry.Value.PolicySchema),
			"value":  schemaValuesObj,
		}
		polAttrTypes := map[string]attr.Type{
			"schema": types.StringType,
			"value":  schemaValuesObj.Type(ctx),
		}

		if len(entry.Identity.AdditionalTargetKeys) > 0 {
			atkAttrs := make(map[string]attr.Value, len(entry.Identity.AdditionalTargetKeys))
			atkAttrTypes := make(map[string]attr.Type, len(entry.Identity.AdditionalTargetKeys))
			for k, v := range entry.Identity.AdditionalTargetKeys {
				atkAttrs[k] = types.StringValue(v)
				atkAttrTypes[k] = types.StringType
			}
			atkObj, d := types.ObjectValue(atkAttrTypes, atkAttrs)
			diags.Append(d...)
			polAttrs["additional_target_keys"] = atkObj
			polAttrTypes["additional_target_keys"] = atkObj.Type(ctx)
		}

		polObj, d := types.ObjectValue(polAttrTypes, polAttrs)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		policyElements = append(policyElements, polObj)
		policyElementTypes = append(policyElementTypes, polObj.Type(ctx))
	}

	policiesTuple, d := types.TupleValue(policyElementTypes, policyElements)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	model.Policies = types.DynamicValue(policiesTuple)
}
