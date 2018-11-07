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
################
# Discovery Doc Builder
#
# This script takes in a yaml file with a Docs object that
# describes which Discovery APIs are being taken in.
#
# The script will then build api.yaml files using
# the Discovery API

# Load everything from MM root.
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '../../')
Dir.chdir(File.join(File.dirname(__FILE__), '../../'))

require 'tools/linter/builder/discovery'
require 'tools/linter/builder/docs'
require 'tools/linter/builder/override'

require 'optparse'

module Api
  class Object
    # Create a setter if the setter doesn't exist
    # Yes, this isn't pretty and I apologize
    def method_missing(method_name, *args)
      matches = /([a-z_]*)=/.match(method_name)
      super unless matches
      create_setter(matches[1])
      method(method_name.to_sym).call(*args)
    end

    def create_setter(variable)
      self.class.define_method("#{variable}=") { |val| instance_variable_set("@#{variable}", val) }
    end

    def validate
    end
  end
end

doc_file = 'tools/linter/docs.yaml'

OptionParser.new do |opts|
  opts.banner = "Discovery doc runner. Usage: run.rb [docs.yaml]"
  opts.on("-f", "--file [file]") { |file| doc_file = file }
end.parse!

docs = YAML::load(File.read(doc_file))

docs.each do |doc|
  product = DiscoveryProduct.new(doc)
  product_obj = product.get_product
  (doc.overrides || []).each do |override|
    override = DiscoveryOverride::Runner.new(product_obj, override)
    product_obj = override.run
  end
  File.write(doc.output, product_obj.to_yaml)
end
