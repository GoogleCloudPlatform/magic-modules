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

require 'api/async'
require 'api/compiler'
require 'api/product'
require 'api/resource'
require 'api/type'
require 'google/yaml_validator'

# Responsbible for grabbing api.yaml and getting resources from it
class ProductApi
  attr_reader :api

  def initialize(product_name)
    @api = get_api(product_name)
  end

  def resource(name)
    @api.objects.select { |obj| obj.name == name }.first
  end

  def all_resource_names
    @api.objects.map(&:name)
  end

  private

  def get_api(product_name)
    Api::Compiler.new("products/#{product_name}/api.yaml").run
  end
end
