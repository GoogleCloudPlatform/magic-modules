# Copyright 2018 Google Inc.
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

module Provider
  # Handles all generation logic for a provider.
  # Providers must call generation_steps to specify which generation functions
  # are being called and in what order.
  module Generation
    # Main entry point for the compiler. As this method is simply invoking other
    # generators, it is okay to ignore Rubocop warnings about method size and
    # complexity.
    def generate(output_folder, types, version_name)
      steps = self.class.instance_variable_get(:@steps)
      raise 'Use `generation_steps` to specify which steps and ordering' unless steps

      steps.each { |s| raise "#{s} function does not exist" unless method(s) }
      steps.each { |s| method(s).call(output_folder, types, version_name) }
    end

    def generate_objects(output_folder, types, version_name)
      version = @api.version_obj_or_default(version_name)
      @api.set_properties_based_on_version(version)
      (@api.objects || []).each do |object|
        if !types.empty? && !types.include?(object.name)
          Google::LOGGER.info "Excluding #{object.name} per user request"
        elsif types.empty? && object.exclude
          Google::LOGGER.info "Excluding #{object.name} per API catalog"
        elsif types.empty? && object.exclude_if_not_in_version(version)
          Google::LOGGER.info "Excluding #{object.name} per API version"
        else
          # version_name will differ from version.name if the resource is being
          # generated at its default version instead of the one that was passed
          # in to the compiler. Terraform needs to know which version was passed
          # in so it can name its output directories correctly.
          generate_object object, output_folder, version_name
        end
      end
    end

    def copy_files(output_folder, _types, _version_name)
      return if @config.files.nil? || @config.files.copy.nil?

      copy_file_list(output_folder, @config.files.copy)
    end

    def compile_files(output_folder, _types, version_name)
      return if @config.files.nil? || @config.files.compile.nil?

      compile_file_list(output_folder, @config.files.compile, version: version_name)
    end

    # version_name is actually used because all of the variables in scope in this method
    # are made available within the templates by the compile call. This means that version_name
    # is exposed to the templating logic and version_name is used in other places in the same
    # way so it needs to be named consistently
    # rubocop:disable Lint/UnusedMethodArgument
    def copy_common_files(output_folder, version_name = 'ga')
      provider_name = self.class.name.split('::').last.downcase
      return unless File.exist?("provider/#{provider_name}/common~copy.yaml")

      Google::LOGGER.info "Copying common files for #{provider_name}"
      files = YAML.safe_load(compile("provider/#{provider_name}/common~copy.yaml"))
      copy_file_list(output_folder, files)
    end
    # rubocop:enable Lint/UnusedMethodArgument

    def compile_common_files(output_folder, _types, version_name = nil)
      provider_name = self.class.name.split('::').last.downcase
      return unless File.exist?("provider/#{provider_name}/common~compile.yaml")

      Google::LOGGER.info "Compiling common files for #{provider_name}"
      files = YAML.safe_load(compile("provider/#{provider_name}/common~compile.yaml"))
      compile_file_list(output_folder, files, version: version_name)
    end

    private

    def copy_file_list(output_folder, files)
      files.each do |target, source|
        target_file = File.join(output_folder, target)
        target_dir = File.dirname(target_file)
        Google::LOGGER.debug "Copying #{source} => #{target}"
        FileUtils.mkpath target_dir unless Dir.exist?(target_dir)
        FileUtils.copy_entry source, target_file
      end
    end

    def compile_examples(output_folder)
      compile_file_map(
        output_folder,
        @config.examples,
        lambda do |_object, file|
          ["examples/#{file}",
           "products/#{@api.prefix[1..-1]}/files/examples~#{file}"]
        end
      )
    end

    def compile_file_list(output_folder, files, data = {})
      files.each do |target, source|
        Google::LOGGER.debug "Compiling #{source} => #{target}"
        target_file = File.join(output_folder, target)
                          .gsub('{{product_name}}', @api.prefix[1..-1])

        manifest = @config.respond_to?(:manifest) ? @config.manifest : {}
        generate_file(
          data.clone.merge(
            name: target,
            product: @api,
            object: {},
            config: {},
            scopes: @api.scopes,
            manifest: manifest,
            tests: '',
            template: source,
            compiler: compiler,
            output_folder: output_folder,
            out_file: target_file,
            prop_ns_dir: @api.prefix[1..-1].downcase,
            product_ns: @api.prefix[1..-1].camelize(:upper)
          )
        )

        %x(goimports -w #{target_file}) if File.extname(target_file) == '.go'
      end
    end

    def generate_object(object, output_folder, version_name)
      data = build_object_data(object, output_folder, version_name)

      generate_resource data
      generate_resource_tests data
    end

    def build_object_data(object, output_folder, version)
      {
        name: object.out_name,
        object: object,
        tests: (@config.tests || {}).select { |o, _v| o == object.name }
                                    .fetch(object.name, {}),
        output_folder: output_folder,
        product_name: object.__product.prefix[1..-1],
        version: version
      }
    end

    def generate_resource_file(data)
      product_ns = if @config.name.nil?
                     data[:object].__product.prefix[1..-1].camelize(:upper)
                   else
                     @config.name
                   end
      generate_file(data.clone.merge(
        # Override with provider specific template for this object, if needed
        template: data[:default_template],
        product_ns: product_ns
      ))
    end
  end
end
