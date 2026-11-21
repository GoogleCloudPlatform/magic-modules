// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = CloudNatDynamicPortFromNumber{}
)

func NewCloudNatDynamicPortFromNumber() function.Function {
	return CloudNatDynamicPortFromNumber{}
}

type CloudNatDynamicPortFromNumber struct{}

func (r CloudNatDynamicPortFromNumber) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "cloud_nat_dynamic_port_from_number"
}

func (r CloudNatDynamicPortFromNumber) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Function that give you as return a valid gcp dynamic port value",
		MarkdownDescription: "Return the closest pow of 2 from the input parameter to be able to pass it as dynmic port value for GCP Cloud Nat, minimum is 32",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "port_type",
				MarkdownDescription: "The dynamic port type you need, only accept : min or max",
			},
			function.NumberParameter{
				Name:                "input",
				MarkdownDescription: "Input number port you want that will be use to get the closest inferior pow of 2 that can be use for gcp cloud nat dynamic port value.",
			},
		},
		Return: function.NumberReturn{},
	}
}

func (r CloudNatDynamicPortFromNumber) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var portType string
	var inputNumber int64
	var allowed_dynamic_ports []int64
	var correctType bool = false
	var output int64 = 32

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &portType, &inputNumber))

	// check port_type is only min or max
	for _, v := range []string{"min", "max"} {
		if portType == v {
			correctType = true
		}
	}
	if !correctType {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(0, "Invalid port_type : must be either min or max string"))
	}

	if resp.Error != nil {
		return
	}
	// Code of the function bellow
	allowed_dynamic_ports = []int64{
		32,
		64,
		128,
		256,
		512,
		1024,
		2048,
		4096,
		8192,
		16384,
		32768,
	}
	// all pow of 2 allowed for min
	if "max" == portType {
		allowed_dynamic_ports = append(allowed_dynamic_ports, 65536) // add specific pow of 2 available for max
	}
	for _, v := range allowed_dynamic_ports {
		if inputNumber >= v {
			output = v
		}
	}
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, output))
}
