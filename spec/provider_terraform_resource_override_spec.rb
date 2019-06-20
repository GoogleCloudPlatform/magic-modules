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
require 'overrides/terraform/resource_override'

class File
  class << self
    alias real_open open
    alias real_read read
  end
end

describe Overrides::Terraform::ResourceOverride do
  context 'good resource' do
    let(:resource) do
      Overrides::Terraform::ResourceOverride.parse(
        IO.read('spec/data/good-resource.yaml')
      )
    end

    before(:each) do
      allow_open 'spec/data/good-resource.yaml'

      resource.validate
    end

    # The ResourceOverride object will get the new description.
    # During the application phase, if the ResourceOverride object
    # has a description, it'll be applied to the new Api Object.
    subject { override.description }

    context 'with extend description' do
      let(:override) { create_override('description', '{{description}}bar') }
      subject { override.description }
      before(:each) { override.apply resource }
      it { is_expected.to eq 'foobar' }
    end

    context 'with override description' do
      let(:override) { create_override('description', 'bar') }
      subject { override.description }
      before(:each) { override.apply resource }
      it { is_expected.to eq 'bar' }
    end
  end

  private

  def allow_open(file_name)
    IO.expects(:read).with(file_name).returns(File.real_read(file_name))
      .at_least(1)
  end

  def create_override(property_name, override_val)
    Google::YamlValidator.parse(
      ['--- !ruby/object:Overrides::Terraform::ResourceOverride',
       "#{property_name}: '#{override_val}'"].join("\n")
    )
  end
end
