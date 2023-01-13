# Copyright 2023 Google Inc.
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
  # Magic Modules Provider for KCC ServiceMappings and TF samples.
  class TerraformKCC < Provider::Terraform
    def generate(output_folder, types, _product_path, _dump_yaml, generate_code, generate_docs)
      @base_url = @version.base_url
      generate_objects(output_folder, types, generate_code, generate_docs)
      compile_product_files(output_folder)
    end

    # Create a directory of sample per test case.
    # Filter out samples that have no test and that don't match the current
    # product version.
    def generate_resource(pwd, data, _generate_code, _generate_docs)
      kind = data.product.name + data.name
      # skip_test examples and examples with test_env_vars should also be
      # included. Whether and how to convert them into KCC examples will be
      # handled separately.
      examples = data.object.examples
                     .reject { |e| @version < @api.version_obj_or_closest(e.min_version) }

      examples.each do |example|
        folder_name = "#{data.product.name}-#{kind}-#{example.name}"
        folder_name += '-skipped' if example.skip_test
        target_folder = File.join('samples', folder_name)

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

    def generate_iam_policy(pwd, data, generate_code, generate_docs) end

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

    # Generate the metadata mapping name based on the last import format. It's
    # usually the field name in the placeholder in the last section of the last
    # import format.
    # E.g. "{{name}}" -> "name";
    # "projects/{{project}}/locations/{{location}}/buckets/{{bucket_id}}" ->
    # "bucket_id".
    def guess_metadata_mapping_name(object)
      # Split the last import format by '/' and take the last part. Then use
      # the regex to verify if it is a value field in the format of {{value}}.
      last_import_part =
        import_id_formats_from_resource(object)[-1].split('/')[-1].scan(/{{[[:word:]]+}}/)
      # If it is a value field, the length of last_import_part will be 1;
      # otherwise it'll be 0.
      # Remove '{{' and '}}' and only return the field name.
      last_import_part.first.gsub('{{', '').gsub('}}', '') if last_import_part.length == 1
    end

    def server_generated_name?(name, object)
      if object.custom_code.post_create
        has_computed_name_configured =
          object.custom_code.post_create == 'templates/terraform/post_create/set_computed_name.erb'
      end
      camel_case_name = name.camelize(:lower)
      has_output_only_name =
        object.all_properties.any? { |p| p.name == camel_case_name && p.output }
      has_computed_name_configured || has_output_only_name
    end

    def supports_conditions(iam_policy)
      request_type = iam_policy.iam_conditions_request_type
      valid_request_types = %w[QUERY_PARAM QUERY_PARAM_NESTED REQUEST_BODY]
      valid_request_types.include?(request_type.to_s)
    end

    def get_resource_id_value_template(id_template, is_server_generated_name, object)
      return nil unless is_server_generated_name

      if id_template.split('/').length == 1 && object.base_url != id_template
        raw_value_template = object.base_url
      end

      return nil if raw_value_template.nil?

      "#{raw_value_template}/{{value}}"
    end

    def get_container(id_template, is_server_generated_name, object)
      container = get_container_from_template(id_template)
      raise 'error having more than one container' if container.length > 2

      if container.empty? && is_server_generated_name
        value_template =
          get_resource_id_value_template(id_template, is_server_generated_name, object)
        container = get_container_from_template(value_template) unless value_template.nil?
        raise 'error having more than one container' if container.length > 2
      end
      container
    end

    def get_container_from_template(template)
      container = []
      id_template_parts = template.split('/')

      projects_field_index = id_template_parts.find_index('projects')
      unless projects_field_index.nil?
        project_field_name =
          id_template_parts[projects_field_index + 1].gsub('{{', '').gsub('}}', '')
      end
      if !project_field_name.nil? && project_field_name != 'name'
        container += ['project', project_field_name]
      end

      folders_field_index = id_template_parts.find_index('folders')
      unless folders_field_index.nil?
        folder_field_name = id_template_parts[folders_field_index + 1].gsub('{{', '').gsub('}}', '')
      end
      if !folder_field_name.nil? && folder_field_name != 'name'
        container += ['folder', folder_field_name]
      end

      organizations_field_index = id_template_parts.find_index('organizations')
      unless organizations_field_index.nil?
        organization_field_name =
          id_template_parts[organizations_field_index + 1].gsub('{{', '').gsub('}}', '')
      end
      if !organization_field_name.nil? && organization_field_name != 'name'
        container += ['organization', organization_field_name]
      end

      billing_accounts_field_index = id_template_parts.find_index('billingAccounts')
      unless billing_accounts_field_index.nil?
        billing_account_field_name =
          id_template_parts[billing_accounts_field_index + 1].gsub('{{', '').gsub('}}', '')
      end
      if !billing_account_field_name.nil? && billing_account_field_name != 'name'
        container += ['billingAccount', billing_account_field_name]
      end

      container
    end

    def get_hierarchical_reference(container)
      hierarchical_reference = []
      hierarchical_reference += [container[0], "#{container[0]}Ref"] if container.length == 2
      hierarchical_reference
    end
  end
end
