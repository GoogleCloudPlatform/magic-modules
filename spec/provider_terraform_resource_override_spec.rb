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

require 'spec_helper'
require 'google/yaml_validator'
require 'provider/terraform/resource_override'

class File
  class << self
    alias real_open open
    alias real_read read
  end
end

describe Provider::Terraform::ResourceOverride do
  context 'good resource' do
    let(:resource) do
      Provider::Terraform::ResourceOverride.parse(
        IO.read('spec/data/good-resource.yaml')
      )
    end

    before do
      allow_open 'spec/data/good-resource.yaml'

      resource.validate
    end

    context 'with empty override' do
      let(:override) { Provider::Terraform::ResourceOverride.new }
      before { override.apply resource }

      it { expect(resource.description).to eq 'foo' }
    end

    it 'extends description' do
      create_override('description', '{{description}}bar').apply resource

      expect(resource.description).to eq 'foobar'
    end

    it 'overrides description' do
      create_override('description', 'bar').apply resource

      expect(resource.description).to eq 'bar'
    end
  end

  def allow_open(file_name)
    IO.expects(:read).with(file_name).returns(File.real_read(file_name))
      .at_least(0)
  end

  def create_override(property_name, override_val)
    Google::YamlValidator.parse(
      format("--- !ruby/object:Provider::Terraform::ResourceOverride\n" \
             "%<k>s: '%<v>s'",
             k: property_name,
             v: override_val)
    )
  end
end
