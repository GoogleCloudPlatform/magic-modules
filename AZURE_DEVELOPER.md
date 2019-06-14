# Developer Guide for Azure Resources

We extended magic-modules to support Azure SDKs and resources. Our design principle is to share as much code as we can, but if not, we will put Azure specific code and templates in `azure` folder, and under `Azure` namespace. For example, default data type definitions are under `Api::Type` namespace in `api/type.rb` file, while Azure specific type definitions are under `Api::Azure::Type` namespace in `api/azure/type.rb`.

For the original magic-modules development documentation, please refer to [DEVELOPER.md](DEVELOPER.md). Before reading the documentation, please make sure you know the basic concepts and coding technologies of Ruby and Ruby template (erb).

## Folder Structure

We reused most of the folder structure defined by the Google's magic-modules, but extended some Azure specific folders. Here is the big-picture of some important folders:

```
magic-modules
  |- api
  |  |- azure
  |  |  |- *.rb         [all Azure specific types could be used in api.yaml]
  |  |- *.rb            [all types could be used in api.yaml]
  |- google
  |  |- *.rb            [all utility functions by Google]
  |- provider
  |  |- ansible         [Ansible specific type and helper function definitions]
  |  |- terraform       [Terraform specific type and helper function definitions]
  |  |- azure
  |  |  |- ansible      [Azure-ansible specific type and helper function definitions]
  |  |  |- terraform    [Azure-terraform specific type and helper function definitions]
  |  |  |- example      [Shared example related type definitions]
  |  |  |- ansible.rb   [Root object to include all sub-modules in ansible folder]
  |  |  |- terraform.rb [Root object to include all sub-modules in terraform folder]
  |  |  |- core.rb      [Helper functions to parse example types]
  |  |- core.rb, abstract_core.rb, terraform.rb, ansible.rb, config.rb [See below]
  |  |- *.rb [all types for <product>.yaml and helper functions for templates]
  |- templates
  |  |- ansible
  |  |  |- facts.erb            [Ansible info module template]
  |  |  |- resource.erb         [Ansible module template]
  |  |  |- integration_test.erb [Ansible test template]
  |  |- terraform
  |  |  |- datasource.erb               [Terraform data source template]
  |  |  |- resource.erb                 [Terraform resource template]
  |  |  |- resource.html.markdown.erb   [Terraform resource documentation template]
  |  |  |- datasource.html.markdown.erb [Terraform data source documentation template]
  |  |  |- schemas
  |  |  |  |- *.erb       [Terraform schema sub-templates, including definition, d.Get, d.Set, etc.]
  |  |  |- *.erb          [Terraform sub-template, e.g. expand, flatten, etc.]
  |  |- azure
  |  |  |- ansible
  |  |  |  |- example  [Ansible test yaml and documentation yaml templates]
  |  |  |  |- module   [Sub-tempaltes to generate Ansible modules or info modules]
  |  |  |  |- sdk      [Sub-templates to generate Python SDK related code like method call]
  |  |  |  |- sdktypes [Sub-templates to generate schema<->SDK marshalling code]
  |  |  |  |- test     [Ansible test helper templates]
  |  |  |- terraform
  |  |  |  |- acctest  [Helper sub-templates to generate Terraform tests]
  |  |  |  |- example  [Terraform test HCL and documentation HCL templates]
  |  |  |  |- schemas  [Azure-specific schema sub-temapltes, including definition, d.Get, d.Set, etc.]
  |  |  |  |- sdk      [Sub-templates to generate Go SDK related code like method call, ID parse, fmt.Errorf]
  |  |  |  |- sdktypes [Sub-templates to generate schema<->SDK marshalling code]
  |  |- *.erb    [sharable templates like auto-gen comment]
  |- compiler.rb [entry point]
```

## Compiler

The entry point of magic module is `compiler.rb`. It will parse the command line, try to read `api.yaml`, `terraform.yaml` and `ansible.yaml` in the input directory, and load the corresponding provider. For Terraform, the provider is located in `provider/terraform.rb`; while for Ansible is `provider/ansible.rb`. When generating code for a specific product (let's say Terraform), all code templates (erb files) will only see functions defined in the corresponding provider.

Class inheritance structure is illustrated below:

```
Provider::Core (provider/core.rb)
  |- Provider::AbstractCore (provider/abstract_core.rb)
  |    |- Provider::Terraform (provider/terraform.rb)
  |- Provider::Ansible::Core (provider/ansible.rb)
```

These providers are root object, besides defining some common helper functions, they will `include` all submodules of the provider (for example, the definition of all Azure specific data types `provider/azure/terraform.rb` and `provider/azure/ansible.rb`). Together with the configuration definitions (`provider/terraform/config.rb` or `provider/ansible/config.rb`), we will be able to use all defined data types and properties in `api.yaml`, `terraform.yaml`, `ansible.yaml` and all ERB templates. As a developer, you need to make sure your types are eventually included in these root objects, otherwise magic-modules will raise errors.

## Code Templates

The overall structure of all code templates are listed in the Folder Structure section. Typically for each template, we introduced a helper function with it, and throughout the code base, we should call that helper function to actually apply the code template.

For example, for a template `property_to_sdkobject.erb` in Terraform, we have the following functions/files:

* Actual template file: `templates/azure/terraform/sdktypes/property_to_sdkobject.erb`
* Helper function: `build_property_to_sdk_object` in `provider/azure/terraform/sdk/sub_template.rb`

We should always call the helper functions throughout the code base like the one in `templates/azure/terraform/sdktypes/nested_object_field_assign.erb`.

## Logic

We handled both Terraform and Ansible resource code generation in a similar way.

1. Magic-modules core handles all definitions (`api.yaml`) and overrides (`terraform.yaml`/`ansible.yaml`) for us
2. Generate schema definition structure
3. Generate code to marshal data from schema to SDK
4. Generate code to call CRUDL SDK APIs
5. Generate code to marshal data from SDK to schema

Since magic-modules core handles all overrides, I will not talk too much about it, please reference to the original magic-modules documentation. As a developer for Azure, we only need to define the overridable attributes (which will be used in `terraform.yaml`/`ansible.yaml`) in one of the `resource_override.rb` or `property_override.rb`, and then we are able to use them throughout the code base including templates and helper functions.

Generating schema definition is simple in Ansible since we only need to call `to_yaml` helper function. It requires some additional efforts for Terraform. That's the reason why `templates/terraform/schemas` and `templates/azure/terraform/schemas` exist.

Marshalling is another tough job to do because Azure SDK objects are deeply hierarchical objects. Both Ansible and Terraform handle it in a recursive way (recursive code template). We put both-way marshalling code templates in `templates/azure/[terraform|ansible]/sdktypes` folder. You will be able to find the corresponding helper functions by using the names of the templates.

It is not too difficult to generate Azure SDK API calls since they are all defined in `azure_sdk_definitions` section of `api.yaml`. Typically API calls are put directly in the resource/module code template, but with an reusable method call sub-template.