// Copyright 2024 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package google

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// A helper class to validate contents coming from YAML files.
type YamlValidator struct{}

func (v *YamlValidator) Parse(content []byte, obj interface{}, yamlPath string) {
	// Create a new decoder to enable strict validation with KnownFields(true)
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)

	if err := decoder.Decode(obj); err != nil {
		log.Fatalf("Cannot unmarshal data from file %s: %v", yamlPath, err)
	}
}

// DecodeYAMLStrict validates that value contains only fields known to target
// (preserving KnownFields(true) enforcement inside custom yaml.Unmarshaler
// implementations, where yaml.Node.Decode otherwise resets knownFields to
// false) and then decodes value into target.
func DecodeYAMLStrict(value *yaml.Node, target interface{}) error {
	if err := CheckKnownFields(value, target); err != nil {
		return err
	}
	return value.Decode(target)
}

// CheckKnownFields recursively verifies that all mapping keys in node match
// known YAML fields on target's Go type, including nested structs, slices,
// maps, inline structs/maps, and YAML merge keys.
func CheckKnownFields(node *yaml.Node, target interface{}) error {
	if node == nil || target == nil {
		return nil
	}
	errs := checkKnownFieldsNode(node, reflect.TypeOf(target))
	if len(errs) > 0 {
		return fmt.Errorf("yaml: unmarshal errors:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}

func checkKnownFieldsNode(node *yaml.Node, t reflect.Type) []string {
	if node == nil || t == nil {
		return nil
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if node.Kind == yaml.DocumentNode {
		var errs []string
		for _, c := range node.Content {
			errs = append(errs, checkKnownFieldsNode(c, t)...)
		}
		return errs
	}
	if node.Kind == yaml.AliasNode {
		return checkKnownFieldsNode(node.Alias, t)
	}

	switch t.Kind() {
	case reflect.Struct:
		if node.Kind != yaml.MappingNode {
			return nil
		}
		fields, hasInlineMap := collectStructYAMLFields(t)
		var errs []string
		for i := 0; i+1 < len(node.Content); i += 2 {
			kNode := node.Content[i]
			vNode := node.Content[i+1]
			if kNode.Tag == "!!merge" || kNode.Value == "<<" {
				if vNode.Kind == yaml.SequenceNode {
					for _, elem := range vNode.Content {
						errs = append(errs, checkKnownFieldsNode(elem, t)...)
					}
				} else {
					errs = append(errs, checkKnownFieldsNode(vNode, t)...)
				}
				continue
			}
			ft, ok := fields[kNode.Value]
			if !ok {
				if !hasInlineMap {
					errs = append(errs, fmt.Sprintf("line %d: field %s not found in type %s", kNode.Line, kNode.Value, t.String()))
				}
				continue
			}
			errs = append(errs, checkKnownFieldsNode(vNode, ft)...)
		}
		return errs
	case reflect.Slice, reflect.Array:
		if node.Kind != yaml.SequenceNode {
			return nil
		}
		var errs []string
		for _, c := range node.Content {
			errs = append(errs, checkKnownFieldsNode(c, t.Elem())...)
		}
		return errs
	case reflect.Map:
		if node.Kind != yaml.MappingNode {
			return nil
		}
		var errs []string
		for i := 0; i+1 < len(node.Content); i += 2 {
			vNode := node.Content[i+1]
			errs = append(errs, checkKnownFieldsNode(vNode, t.Elem())...)
		}
		return errs
	}
	return nil
}

func collectStructYAMLFields(t reflect.Type) (map[string]reflect.Type, bool) {
	fields := map[string]reflect.Type{}
	hasInlineMap := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}
		tag := f.Tag.Get("yaml")
		if tag == "-" {
			continue
		}
		parts := strings.Split(tag, ",")
		name := parts[0]
		inline := slices.Contains(parts[1:], "inline")
		if inline {
			ft := f.Type
			for ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				sub, subMap := collectStructYAMLFields(ft)
				for k, v := range sub {
					fields[k] = v
				}
				if subMap {
					hasInlineMap = true
				}
			} else if ft.Kind() == reflect.Map {
				hasInlineMap = true
			}
			continue
		}
		if f.PkgPath != "" {
			continue
		}
		if name == "" {
			name = strings.ToLower(f.Name)
		}
		fields[name] = f.Type
	}
	return fields, hasInlineMap
}
