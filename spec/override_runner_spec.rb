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

class TestResourceOverride < Provider::Overrides::ResourceOverride
  def self.attributes
    [:test_field]
  end
end

describe Provider::Overrides::Runner do
  context 'simple overrides' do
    describe 'should be able to override a product field' do
      let(:overrides) do
        Provider::Overrides::ResourceOverrides.new(
          'product' => Provider::Overrides::ResourceOverride.new(
            'name' => 'My Test Product'
          )
        )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::Overrides::Runner.new(api, overrides)
        new_api = runner.build
        expect(new_api.name).to eq(overrides['product']['name'])
      }
    end

    describe 'should be able to override a resource field' do
      let(:overrides) do
        Provider::Overrides::ResourceOverrides.new(
          'MyResource' => Provider::Overrides::ResourceOverride.new(
            'description' => 'A description'
          )
        )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::Overrides::Runner.new(api, overrides)
        new_api = runner.build
        resource = new_api.objects.select { |p| p.name == 'MyResource' }.first
        expect(resource.description).to eq(overrides['MyResource']['description'])
      }
    end

    describe 'should be able to override a property field' do
      let(:overrides) do
        Provider::Overrides::ResourceOverrides.new(
          'ReferencedResource' => Provider::Overrides::ResourceOverride.new(
            'properties' => {
              'name' => Provider::Overrides::PropertyOverride.new(
                'description' => 'My overriden description'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::Overrides::Runner.new(api, overrides)
        new_api = runner.build
        resource = new_api.objects.select { |p| p.name == 'ReferencedResource' }.first
        prop = resource.properties.select { |p| p.name == 'name' }.first
        expect(prop.description).to eq(
          overrides['ReferencedResource']['properties']['name']['description']
        )
      }
    end

    describe 'should be able to override a property type' do
      let(:overrides) do
        Provider::Overrides::ResourceOverrides.new(
          'ReferencedResource' => Provider::Overrides::ResourceOverride.new(
            'properties' => {
              'name' => Provider::Overrides::PropertyOverride.new(
                'type' => 'Api::Type::Integer'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::Overrides::Runner.new(api, overrides)
        new_api = runner.build
        resource = new_api.objects.select { |p| p.name == 'ReferencedResource' }.first
        prop = resource.properties.select { |p| p.name == 'name' }.first
        expect(prop.class).to eq(Api::Type::Integer)
      }
    end

    describe 'should be able to override a nested-nested property' do
      let(:overrides) do
        Provider::Overrides::ResourceOverrides.new(
          'AnotherResource' => Provider::Overrides::ResourceOverride.new(
            'properties' => {
              'nested-property2.property1.property1-nested' =>
              Provider::Overrides::PropertyOverride.new(
                'type' => 'Api::Type::Integer'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::Overrides::Runner.new(api, overrides)
        new_api = runner.build
        resource = new_api.objects.select { |p| p.name == 'AnotherResource' }.first
        prop = resource.properties.select { |p| p.name == 'nested-property2' }.first
        expect(prop.properties[0].properties[0].class).to eq(Api::Type::Integer)
      }
    end

    describe 'should be able to override a nested array property' do
      let(:overrides) do
        Provider::Overrides::ResourceOverrides.new(
          'AnotherResource' => Provider::Overrides::ResourceOverride.new(
            'properties' => {
              'array-property[].property1' => Provider::Overrides::PropertyOverride.new(
                'type' => 'Api::Type::Integer'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::Overrides::Runner.new(api, overrides)
        new_api = runner.build
        resource = new_api.objects.select { |p| p.name == 'AnotherResource' }.first
        prop = resource.properties.select { |p| p.name == 'array-property' }.first
        expect(prop.item_type.properties[0].class).to eq(Api::Type::Integer)
      }
    end

    describe 'should be able to override a namevalue -> object map' do
      let(:overrides) do
        Provider::Overrides::ResourceOverrides.new(
          'AnotherResource' => TestResourceOverride.new(
            'properties' => {
              'namevalue-property.nv-prop1' => Provider::Overrides::PropertyOverride.new(
                'description' => 'overriden'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new('spec/data/good-file.yaml').run }

      it {
        runner = Provider::Overrides::Runner.new(api, overrides, TestResourceOverride)
        new_api = runner.build
        resource = new_api.objects.select { |p| p.name == 'AnotherResource' }.first
        prop = resource.properties.select { |p| p.name == 'namevalue-property' }
                       .first
                       .value_type
                       .properties.first
        expect(prop.description).to eq('overriden')
      }
    end
  end
end
