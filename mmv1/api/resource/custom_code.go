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

package resource

// Inserts custom code into terraform resources.
type CustomCode struct {
	// Collection of fields allowed in the CustomCode section for
	// Terraform.

	// All custom code attributes are string-typed.  The string should
	// be the name of a template file which will be compiled in the
	// specified / described place.
	//
	// ======================
	// schema.Resource stuff
	// ======================
	// Extra Schema Entries go below all other schema entries in the
	// resource's Resource.Schema map.  They should be formatted as
	// entries in the map, e.g. `"foo": &schema.Schema{ ... },`.
	ExtraSchemaEntry string `yaml:"extra_schema_entry"`

	// ====================
	// Encoders & Decoders
	// ====================
	// The encoders are functions which take the `obj` map after it
	// has been assembled in either "Create" or "Update" and mutate it
	// before it is sent to the server.  There are lots of reasons you
	// might want to use these - any differences between local schema
	// and remote schema will be placed here.
	// Because the call signature of this function cannot be changed,
	// the template will place the function header and closing } for
	// you, and your custom code template should *not* include them.
	Encoder string

	// The update encoder is the encoder used in Update - if one is
	// not provided, the regular encoder is used.  If neither is
	// provided, of course, neither is used.  Similarly, the custom
	// code should *not* include the function header or closing }.
	// Update encoders are only used if object.input is false,
	// because when object.input is true, only individual fields
	// can be updated - in that case, use a custom expander.
	UpdateEncoder string `yaml:"update_encoder"`

	// The decoder is the opposite of the encoder - it's called
	// after the Read succeeds, rather than before Create / Update
	// are called.  Like with encoders, the decoder should not
	// include the function header or closing }.
	Decoder string

	// =====================
	// Simple customizations
	// =====================
	// Constants go above everything else in the file, and include
	// things like methods that will be referred to by name elsewhere
	// (e.g. "fooBarDiffSuppress") and regexes that are necessarily
	// exported (e.g. "fooBarValidationRegex").
	Constants string

	// This code is run before the Create call happens.  It's placed
	// in the Create function, just before the Create call is made.
	PreCreate string `yaml:"pre_create"`

	// This code is run after the Create call succeeds.  It's placed
	// in the Create function directly without modification.
	PostCreate string `yaml:"post_create"`

	// This code is run after the Create call fails before the error is
	// returned. It's placed in the Create function directly without
	// modification.
	PostCreateFailure string `yaml:"post_create_failure"`

	// This code replaces the entire contents of the Create call. It
	// should be used for resources that don't have normal creation
	// semantics that cannot be supported well by other MM features.
	CustomCreate string `yaml:"custom_create"`

	// This code is run before the Read call happens.  It's placed
	// in the Read function.
	PreRead string `yaml:"pre_read"`

	// This code is run before the Update call happens.  It's placed
	// in the Update function, just after the encoder call, before
	// the Update call.  Just like the encoder, it is only used if
	// object.input is false.
	PreUpdate string `yaml:"pre_update"`

	// This code is run after the Update call happens.  It's placed
	// in the Update function, just after the call succeeds.
	// Just like the encoder, it is only used if object.input is
	// false.
	PostUpdate string `yaml:"post_update"`

	// This code replaces the entire contents of the Update call. It
	// should be used for resources that don't have normal update
	// semantics that cannot be supported well by other MM features.
	CustomUpdate string `yaml:"custom_update"`

	// This code is run just before the Delete call happens.  It's
	// useful to prepare an object for deletion, e.g. by detaching
	// a disk before deleting it.
	PreDelete string `yaml:"pre_delete"`

	// This code is run just after the Delete call happens.
	PostDelete string `yaml:"post_delete"`

	// This code replaces the entire delete method.  Since the delete
	// method's function header can't be changed, the template
	// inserts that for you - do not include it in your custom code.
	CustomDelete string `yaml:"custom_delete"`

	// This code replaces the entire import method.  Since the import
	// method's function header can't be changed, the template
	// inserts that for you - do not include it in your custom code.
	CustomImport string `yaml:"custom_import"`

	// This code is run just after the import method succeeds - it
	// is useful for parsing attributes that are necessary for
	// the Read() method to succeed.
	PostImport string `yaml:"post_import"`

	// This code is run in the generated test file to check that the
	// resource was successfully deleted. Use this if the API responds
	// with a success HTTP code for deleted resources
	TestCheckDestroy string `yaml:"test_check_destroy"`
}
