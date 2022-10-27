# Copyright 2019 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

require 'provider/terraform'
require 'provider/terraform/import'

module Provider
  # Magic Modules Provider for KCC ServiceMappings and other related templates.
  # Instead of generating KCC directly, this provider generates a KCC-compatible
  # library to be consumed by KCC.
  class TerraformKCC < Provider::Terraform
    def generate(output_folder, types, _product_path, _dump_yaml, generate_code, generate_docs)
      @base_url = @version.base_url
      generate_objects(output_folder, types, generate_code, generate_docs)
      compile_product_files(output_folder)
    end

    # Create a directory of samples per resource
    # Filter out samples that have no test and don't necessarily run, use
    # externally injected values (env vars), and that don't match the current
    # product version.
    def generate_resource(pwd, data, _generate_code, _generate_docs)
      kind = data.product.name + data.name
      # TODO: support skip_test tests in a separate output subfolder.
      examples = data.object.examples
                     .reject(&:skip_test)
                     .reject { |e| !e.test_env_vars.nil? && e.test_env_vars.any? }
                     .reject { |e| @version < @api.version_obj_or_closest(e.min_version) }

      examples.each do |example|
        target_folder = File.join('samples', data.product.name + "-" + kind + '-' + example.name)
        FileUtils.mkpath target_folder
        data.example = example
        data.generate(
          pwd,
          'templates/kcc/samples/sample.tf.erb',
          File.join(target_folder, 'main.tf'),
          self
        )
      end
    end

    def generate_resource_tests(pwd, data) end

    def generate_resource_sweepers(pwd, data) end

    def generate_iam_policy(pwd, data, generate_code, generate_docs)end

    def compile_product_files(output_folder)
      file_template = ProductFileTemplate.new(
        output_folder,
        nil,
        @api,
        @target_version_name,
        build_env
      )
      compile_file_list(output_folder,
                        [
                          [
                            "servicemappings/#{@api.name.downcase}.yaml",
                            'templates/kcc/product/service_mapping.yaml.erb'
                          ]
                        ],
                        file_template)
    end

    def compile_common_files(output_folder, products, _common_compile_file) end

    def copy_common_files(output_folder, generate_code, generate_docs) end

    # A strict mapping from K8S name -> Terraform resource "name" doesn't make
    # sense for some resources but we can approximate this well enough for most
    # of them. This is often the `name` field, but it could be named something
    # else.
    # If {{name}} is part of a resource id, it should be the last import format.
    # Otherwise, {{value}} or values/{{value}} are also valid. If the final id
    # has multiple terms, we can reject it (by returning nil) as we can't create
    # a 1:1 mapping from K8S resource name : Terraform field.
    def guess_metadata_mapping_name(object)
      # Split the last import format by '/' and take the last part. Then use
      # the regex to verify if it is a value field in the format of {{value}}.
      last_import_part = import_id_formats_from_resource(object)[-1].split('/')[-1].scan(/{{[[:word:]]+}}/)
      # If it is a value field, the length of last_import_part will be 1;
      # otherwise it'll be 0.
      # Remove '{{' and '}}' and only return the field name.
      last_import_part.first.gsub('{{', '').gsub('}}', '') if last_import_part.length == 1
    end

    # TODO: Incrementally cover all the server generated ID cases.
    def is_server_generated_name(name, object)
      has_computed_name_configured = object.custom_code.post_create == 'templates/terraform/post_create/set_computed_name.erb' if object.custom_code.post_create
      camel_case_name = name.camelize(:lower)
      has_output_only_name = object.all_properties.any?{ |p| p.name == camel_case_name && p.output }
      has_computed_name_configured || has_output_only_name
    end

    def supports_conditions(iam_policy)
      request_type = iam_policy.iam_conditions_request_type
      valid_request_types = ['QUERY_PARAM', 'QUERY_PARAM_NESTED', 'REQUEST_BODY']
      return valid_request_types.include?(request_type.to_s)
    end

    def get_resource_id_value_template(id_template, is_server_generated_name, object)
      return nil if !is_server_generated_name

      if id_template.split('/').length == 1 && object.base_url != id_template
        raw_value_template = object.base_url
      end
      return nil if raw_value_template == nil

      value_template = raw_value_template + '/{{value}}'
    end

    def get_container(id_template, is_server_generated_name, object)
      container = get_container_from_template(id_template)
      raise 'error having more than one container' if container.length > 2
      if container.length == 0 && is_server_generated_name
        value_template = get_resource_id_value_template(id_template, is_server_generated_name, object)
        container = get_container_from_template(value_template) if value_template != nil
        raise 'error having more than one container' if container.length > 2
      end
      return container
    end

    def get_container_from_template(template)
      container = Array.new()
      id_template_parts = template.split('/')

      projects_field_index = id_template_parts.find_index('projects')
      project_field_name = id_template_parts[projects_field_index + 1].gsub('{{', '').gsub('}}', '') if !projects_field_index.nil?
      container += ['project', project_field_name] if !project_field_name.nil? && project_field_name != 'name'

      folders_field_index = id_template_parts.find_index('folders')
      folder_field_name = id_template_parts[folders_field_index + 1].gsub('{{', '').gsub('}}', '') if !folders_field_index.nil?
      container += ['folder', folder_field_name] if !folder_field_name.nil? && folder_field_name != 'name'

      organizations_field_index = id_template_parts.find_index('organizations')
      organization_field_name = id_template_parts[organizations_field_index + 1].gsub('{{', '').gsub('}}', '') if !organizations_field_index.nil?
      container += ['organization', organization_field_name] if !organization_field_name.nil? && organization_field_name != 'name'

      billing_accounts_field_index = id_template_parts.find_index('billingAccounts')
      billing_account_field_name = id_template_parts[billing_accounts_field_index + 1].gsub('{{', '').gsub('}}', '') if !billing_accounts_field_index.nil?
      container += ['billingAccount', billing_account_field_name] if !billing_account_field_name.nil? && billing_account_field_name != 'name'

      return container
    end

    def get_hierarchical_reference(container)
      hierarchical_reference = Array.new()
      hierarchical_reference += [container[0], container[0] + 'Ref'] if container.length == 2
      return hierarchical_reference
    end

    def format_id_template(id_template, object)
      # transform from buckets/{{bucket}} to {{bucket}}
      id_template_parts = id_template.scan(/{{[[:word:]]+}}/)
      id_template_parts -= ['{{project}}', '{{region}}']
      id_template_formatted = id_template_parts.join('/')

      # transform refs from {{bucket}} to {{bucketRef.name}} form
      prop_names = id_template.scan(/{{[[:word:]]+}}/).map { |p| p.gsub('{{', '').gsub('}}', '') }
      # probably won't catch overridden names
      object.all_properties
            .reject { |p| p.name == 'zone' } # exclude special fields
            .select { |p| p.is_a?(Api::Type::ResourceRef) } # select resource refs
            .select { |p| prop_names.include?(p.name.camelize(:lower)) } # canonical name
            .each do |prop|
        id_template_formatted = id_template_formatted
                                .gsub(
                                  "{{#{prop.name.camelize(:lower)}}}",
                                  "{{#{prop.name}Ref.name}}"
                                )
      end
      id_template_formatted
    end
  end
end
