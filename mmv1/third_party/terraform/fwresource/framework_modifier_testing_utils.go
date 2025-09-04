package fwresource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func GenerateModifyPlanRequestWithSetValue(s schema.Schema, field string, value interface{}) *resource.ModifyPlanRequest{
	request := GenerateModifyPlanRequest(s)
	request.Plan.SetAttribute(context.Background(), path.Root(field), value)
	return request
}

func GenerateModifyPlanRequest(s schema.Schema) *resource.ModifyPlanRequest{
	state := tfsdk.State{
		Schema: s,
	}
	plan := tfsdk.Plan{
		Schema: s,
	}
	return resource.ModifyPlanRequest{
		State: state,
		Plan: plan,
	}
}

func GenerateModifyPlanResponse(s schema.Schema) *resource.ModifyPlanResponse{
	state := tfsdk.State{
		Schema: s,
	}
	plan := tfsdk.Plan{
		Schema: s,
	}
	return resource.ModifyPlanRequest{
		State: state,
		Plan: plan,
	}
}