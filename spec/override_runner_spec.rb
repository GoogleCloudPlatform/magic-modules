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

describe Provider::OverrideRunner do
  context 'simple overrides' do
    describe 'should be able to override a product field' do
      let(:overrides) do
        Provider::ResourceOverrides.new(
          'product' => Provider::ResourceOverride.new(
            'name' => 'My Test Product'
          )
        )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::OverrideRunner.new(api, overrides)
        new_api = runner.build
        expect(new_api.name).to eq(overrides['product']['name'])
      }
    end

    describe 'should be able to override a resource field' do
      let(:overrides) do
        Provider::ResourceOverrides.new(
          'MyResource' => Provider::ResourceOverride.new(
            'description' => 'A description'
          )
        )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::OverrideRunner.new(api, overrides)
        new_api = runner.build
        resource = new_api.objects.select { |p| p.name == 'MyResource' }.first
        expect(resource.description).to eq(overrides['MyResource']['description'])
      }
    end

    describe 'should be able to override a property field' do
      let(:overrides) do
        Provider::ResourceOverrides.new(
          'ReferencedResource' => Provider::ResourceOverride.new(
            'properties' => {
              'name' => Provider::PropertyOverride.new(
                'description' => 'My overriden description'
              )
            }
            )
          )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::OverrideRunner.new(api, overrides)
        new_api = runner.build
        resource = new_api.objects.select { |p| p.name == 'ReferencedResource' }.first
        prop = resource.properties.select { |p| p.name == 'name' }.first
        expect(prop.description).to eq(
          overrides['ReferencedResource']['properties']['name']['description']
        )
      }
    end
  end
end
