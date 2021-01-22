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

require 'api/async'
require 'api/product'
require 'api/resource'
require 'api/type'
require 'compile/core'
require 'google/yaml_validator'

module Api
  # Process <product>.yaml and produces output module
  class Compiler
    include Compile::Core

    attr_reader :product

    def initialize(catalog)
      @catalog = catalog
    end

    def run
      # Compile step #1: compile with generic class to instantiate target class
      config = Google::YamlValidator.parse(@catalog)
      unless config.class <= Api::Product
        raise StandardError, "#{@catalog} is #{config.class}"\
          ' instead of Api::Product' \
      end
      config
    end
  end
end
