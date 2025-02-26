// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// In order to interact with resource converters, we need to be able to create
// "terraform resource data" that supports a very limited subset of the API actually
// used during the conversion process.
package models

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FakeResourceDataWithMeta struct {
	FakeResourceData
	FakeResourceMeta
}

type FakeResourceMeta struct {
	kind      string
	address   string
	isDeleted bool
}

// Kind returns the type of resource (i.e. "google_storage_bucket").
func (d *FakeResourceMeta) Kind() string {
	return d.kind
}

func (d *FakeResourceMeta) Address() string {
	return d.address
}

func (d *FakeResourceMeta) IsDeleted() bool {
	return d.isDeleted
}

func NewFakeResourceDataWithMeta(kind string, resourceSchema map[string]*schema.Schema, values map[string]interface{}, isDeleted bool, tfplanAddress string) *FakeResourceDataWithMeta {
	state := map[string]string{}
	var address []string
	attributes(values, address, state, resourceSchema)
	reader := &schema.MapFieldReader{
		Map:    schema.BasicMapReader(state),
		Schema: resourceSchema,
	}
	return &FakeResourceDataWithMeta{
		FakeResourceData{
			schema: resourceSchema,
			reader: reader,
		},
		FakeResourceMeta{
			kind:      kind,
			isDeleted: isDeleted,
			address:   tfplanAddress,
		},
	}
}
