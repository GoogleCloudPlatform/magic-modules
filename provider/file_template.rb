# Copyright 2017 Google Inc.
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

require 'fileutils'
require 'google/logger'
require 'pathname'

module Provider
  # Parent class for specific types of files. Contains methods to generate files
  class FileTemplate
    include Compile::Core
    # The root folder we're outputting to.
    attr_accessor :output_folder
    # Information about the local environment
    # (which formatters are enabled, start-time)
    attr_accessor :env
    # The API version
    attr_accessor :version

    # Given the data object for a file, write that file and verify that it
    # passes these conditions:
    #
    # - The file has not been generated already this run
    # - The file has an autogen exception or an autogen notice defined
    #
    # Once the file's contents are written, set the proper [chmod] mode and
    # format the file with a language-appropriate formatter.
    def generate(template, path, provider)
      folder = File.dirname(path)
      FileUtils.mkpath folder unless Dir.exist?(folder)

      # If we've modified a file since starting an MM run, it's a reasonable
      # assumption that it was this run that modified it.
      if File.exist?(path) && File.mtime(path) > @env[:start_time]
        raise "#{path} was already modified during this run"
      end

      # You're looking at some magic here!
      # This is how variables are made available in templates; we iterate
      # through each key:value pair in this object, and we set them
      # in the scope of the provider.
      #
      # The templates get access to everything in the provider +
      # all of the variables in this object.
      ctx = provider.provider_binding
      instance_variables.each do |name|
        ctx.local_variable_set(name[1..-1], instance_variable_get(name))
      end

      # This variable is used in ansible/resource.erb
      ctx.local_variable_set('file_relative', relative_path(path, @output_folder).to_s)

      Google::LOGGER.debug "Generating #{path}"
      File.open(path, 'w') { |f| f.puts compile_file(ctx, template) }

      # Files are often generated in parallel.
      # We can use thread-local variables to ensure that autogen checking
      # stays specific to the file each thred represents.
      raise "#{path} missing autogen" unless Thread.current[:autogen]

      old_file_chmod_mode = File.stat(template).mode
      FileUtils.chmod(old_file_chmod_mode, path)

      format_output_file(path)
    end

    private

    def format_output_file(path)
      if path.end_with?('.py') && @env[:pyformat_enabled]
        run_formatter("python3 -m black --line-length 160 -S #{path}")
      elsif path.end_with?('.go') && @env[:goformat_enabled]
        run_formatter("gofmt -w -s #{path}")
        run_formatter("goimports -w #{path}")
      end
    end

    def run_formatter(command)
      output = %x(#{command} 2>&1)
      Google::LOGGER.error output unless $CHILD_STATUS.to_i.zero?
    end

    def relative_path(target, base)
      Pathname.new(target).relative_path_from(Pathname.new(base))
    end
  end

  # Responsible for compiling provider-level files, rather than product-specific ones
  class ProviderFileTemplate < Provider::FileTemplate
    # All the products that are being compiled with the provider on this run
    attr_accessor :products

    # Optional path to the directory where overrides reside. Used to locate files
    # outside of the MM root directory
    attr_accessor :override_path

    def initialize(output_folder, version, env, products, override_path = nil)
      @output_folder = output_folder
      @version = version
      @env = env
      @products = products
      @override_path = override_path
    end
  end

  # Responsible for generating a file in the context of a product
  # with a given set of parameters.
  class ProductFileTemplate < Provider::FileTemplate
    # The name of the resource
    attr_accessor :name
    # The provider-specific configuration.
    attr_accessor :config
    # The namespace of the product.
    attr_accessor :product_ns
    # The resource itself.
    attr_accessor :object
    # The entire API object.
    attr_accessor :product

    class << self
      # Construct a new ProductFileTemplate based on a resource object
      def file_for_resource(output_folder, object, version, config, env)
        file_template = new(output_folder, object.name, object.__product, version, env)
        file_template.object = object
        file_template.config = config
        file_template
      end
    end

    def initialize(output_folder, name, product, version, env)
      @name = name
      @product = product
      @product_ns = product.name
      @output_folder = output_folder
      @version = version
      @env = env
    end
  end
end
