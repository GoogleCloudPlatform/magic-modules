// Copyright 2019 Google LLC
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
package tfdata

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Must be set to the same value as the internal typeObject const
const typeObject schema.ValueType = 8

// This is more or less equivalent to the internal getResult struct
// used by schema.ResourceData
type getResult struct {
	Value    interface{}
	Computed bool
	Exists   bool
	Schema   *schema.Schema
}

// Compare to https://github.com/hashicorp/terraform-plugin-sdk/blob/97b4465/helper/schema/resource_data.go#L15
type FakeResourceData struct {
	reader schema.FieldReader
	kind   string
	schema map[string]*schema.Schema
}

// Kind returns the type of resource (i.e. "google_storage_bucket").
func (d *FakeResourceData) Kind() string {
	return d.kind
}

// Id returns the ID of the resource from state.
func (d *FakeResourceData) Id() string {
	return ""
}

func (d *FakeResourceData) getRaw(key string) getResult {
	var parts []string
	if key != "" {
		parts = strings.Split(key, ".")
	}
	return d.get(parts)
}

func (d *FakeResourceData) get(addr []string) getResult {
	r, err := d.reader.ReadField(addr)
	if err != nil {
		panic(err)
	}

	var s *schema.Schema
	if schemaPath := addrToSchema(addr, d.schema); len(schemaPath) > 0 {
		s = schemaPath[len(schemaPath)-1]
	}
	if r.Value == nil && s != nil {
		r.Value = r.ValueOrZero(s)
	}

	return getResult{
		Value:    r.Value,
		Computed: r.Computed,
		Exists:   r.Exists,
		Schema:   s,
	}
}

// Get reads a single field by key.
func (d *FakeResourceData) Get(name string) interface{} {
	val, _ := d.GetOk(name)
	return val
}

// Get reads a single field by key and returns a boolean indicating
// whether the field exists.
func (d *FakeResourceData) GetOk(name string) (interface{}, bool) {
	r := d.getRaw(name)
	exists := r.Exists && !r.Computed

	if exists {
		// Verify that it's not the zero-value
		value := r.Value
		zero := r.Schema.Type.Zero()

		if eq, ok := value.(schema.Equal); ok {
			exists = !eq.Equal(zero)
		} else {
			exists = !reflect.DeepEqual(value, zero)
		}
	}

	return r.Value, exists
}

func (d *FakeResourceData) GetOkExists(key string) (interface{}, bool) {
	r := d.getRaw(key)
	exists := r.Exists && !r.Computed
	return r.Value, exists
}

// These methods are required by some mappers but we don't actually have (or need)
// implementations for them.
func (d *FakeResourceData) HasChange(string) bool             { return false }
func (d *FakeResourceData) Set(string, interface{}) error     { return nil }
func (d *FakeResourceData) SetId(string)                      {}
func (d *FakeResourceData) GetProviderMeta(interface{}) error { return nil }
func (d *FakeResourceData) Timeout(key string) time.Duration  { return time.Duration(1) }

func NewFakeResourceData(kind string, resourceSchema map[string]*schema.Schema, values map[string]interface{}) *FakeResourceData {
	state := map[string]string{}
	var address []string
	attributes(values, address, state, resourceSchema)
	reader := &schema.MapFieldReader{
		Map:    schema.BasicMapReader(state),
		Schema: resourceSchema,
	}
	return &FakeResourceData{
		kind:   kind,
		schema: resourceSchema,
		reader: reader,
	}
}

// addrToSchema finds the final element schema for the given address
// and the given schema. It returns all the schemas that led to the final
// schema. These are in order of the address (out to in).
// NOTE: This function was copied from the terraform library:
// github.com/hashicorp/terraform/helper/schema/field_reader.go
func addrToSchema(addr []string, schemaMap map[string]*schema.Schema) []*schema.Schema {
	current := &schema.Schema{
		Type: typeObject,
		Elem: schemaMap,
	}

	// If we aren't given an address, then the user is requesting the
	// full object, so we return the special value which is the full object.
	if len(addr) == 0 {
		return []*schema.Schema{current}
	}

	result := make([]*schema.Schema, 0, len(addr))
	for len(addr) > 0 {
		k := addr[0]
		addr = addr[1:]

	REPEAT:
		// We want to trim off the first "typeObject" since its not a
		// real lookup that people do. i.e. []string{"foo"} in a structure
		// isn't {typeObject, typeString}, its just a {typeString}.
		if len(result) > 0 || current.Type != typeObject {
			result = append(result, current)
		}

		switch t := current.Type; t {
		case schema.TypeBool, schema.TypeInt, schema.TypeFloat, schema.TypeString:
			if len(addr) > 0 {
				return nil
			}
		case schema.TypeList, schema.TypeSet:
			isIndex := len(addr) > 0 && addr[0] == "#"

			switch v := current.Elem.(type) {
			case *schema.Resource:
				current = &schema.Schema{
					Type: typeObject,
					Elem: v.Schema,
				}
			case *schema.Schema:
				current = v
			case schema.ValueType:
				current = &schema.Schema{Type: v}
			default:
				// we may not know the Elem type and are just looking for the
				// index
				if isIndex {
					break
				}

				if len(addr) == 0 {
					// we've processed the address, so return what we've
					// collected
					return result
				}

				if len(addr) == 1 {
					if _, err := strconv.Atoi(addr[0]); err == nil {
						// we're indexing a value without a schema. This can
						// happen if the list is nested in another schema type.
						// Default to a TypeString like we do with a map
						current = &schema.Schema{Type: schema.TypeString}
						break
					}
				}

				return nil
			}

			// If we only have one more thing and the next thing
			// is a #, then we're accessing the index which is always
			// an int.
			if isIndex {
				current = &schema.Schema{Type: schema.TypeInt}
				break
			}

		case schema.TypeMap:
			if len(addr) > 0 {
				switch v := current.Elem.(type) {
				case schema.ValueType:
					current = &schema.Schema{Type: v}
				default:
					// maps default to string values. This is all we can have
					// if this is nested in another list or map.
					current = &schema.Schema{Type: schema.TypeString}
				}
			}
		case typeObject:
			// If we're already in the object, then we want to handle Sets
			// and Lists specially. Basically, their next key is the lookup
			// key (the set value or the list element). For these scenarios,
			// we just want to skip it and move to the next element if there
			// is one.
			if len(result) > 0 {
				lastType := result[len(result)-2].Type
				if lastType == schema.TypeSet || lastType == schema.TypeList {
					if len(addr) == 0 {
						break
					}

					k = addr[0]
					addr = addr[1:]
				}
			}

			m := current.Elem.(map[string]*schema.Schema)
			val, ok := m[k]
			if !ok {
				return nil
			}

			current = val
			goto REPEAT
		}
	}

	return result
}

// attributes function takes json parsed JSON object (value param) and fill map[string]string with it's
// content (state param) for example JSON:
//
//	{
//		"foo": {
//			"name" : "value"
//		},
//	  "list": ["item1", "item2"]
//	}
//
// will be translated to map with following key/value set:
//
//	foo.name => "value"
//	list.# => 2
//	list.0 => "item1"
//	list.1 => "item2"
//
// Map above will be passed to schema.BasicMapReader that have all appropriate logic to read fields
// correctly during conversion to CAI.
func attributes(value interface{}, address []string, state map[string]string, schemas map[string]*schema.Schema) {
	schemaArr := addrToSchema(address, schemas)
	if len(schemaArr) == 0 {
		return
	}
	sch := schemaArr[len(schemaArr)-1]
	addr := strings.Join(address, ".")
	// int is special case, can't use handle it in main switch because number will be always parsed from JSON as float
	// need to identify it by schema.TypeInt and convert to int from int or float
	if sch.Type == schema.TypeInt && value != nil {
		switch value.(type) {
		case int:
			state[addr] = strconv.Itoa(value.(int))
		case float64:
			state[addr] = strconv.Itoa(int(value.(float64)))
		case float32:
			state[addr] = strconv.Itoa(int(value.(float32)))
		}
		return
	}

	switch value.(type) {
	case nil:
		defaultValue, err := sch.DefaultValue()
		if err != nil {
			panic(fmt.Sprintf("error getting default value for %v", address))
		}
		if defaultValue == nil {
			defaultValue = sch.ZeroValue()
		}
		attributes(defaultValue, address, state, schemas)
	case float64:
		state[addr] = strconv.FormatFloat(value.(float64), 'f', 6, 64)
	case float32:
		state[addr] = strconv.FormatFloat(value.(float64), 'f', 6, 32)
	case string:
		state[addr] = value.(string)
	case bool:
		state[addr] = strconv.FormatBool(value.(bool))
	case int:
		state[addr] = strconv.Itoa(value.(int))
	case []interface{}:
		arr := value.([]interface{})
		countAddr := addr + ".#"
		state[countAddr] = strconv.Itoa(len(arr))
		for i, e := range arr {
			addr := append(address, strconv.Itoa(i))
			attributes(e, addr, state, schemas)
		}
	case map[string]interface{}:
		m := value.(map[string]interface{})
		for k, v := range m {
			addr := append(address, k)
			attributes(v, addr, state, schemas)
		}
	case *schema.Set:
		set := value.(*schema.Set)
		attributes(set.List(), address, state, schemas)
	default:
		panic(fmt.Sprintf("unrecognized type %T", value))
	}
}
