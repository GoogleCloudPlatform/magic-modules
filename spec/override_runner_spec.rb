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

class TestResourceOverride < Overrides::ResourceOverride
  def self.attributes
    [:test_field]
  end
end

describe Overrides::Runner do
  context 'simple overrides' do
    describe 'should be able to override a product field' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'product' => Overrides::ResourceOverride.new(
            'name' => 'My Test Product'
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        new_api = Overrides::Runner.build(api, overrides)
        expect(new_api.name).to eq(overrides['product']['name'])
      }
    end

    describe 'should be able to override a resource field' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'MyResource' => Overrides::ResourceOverride.new(
            'description' => 'A description'
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        new_api = Overrides::Runner.build(api, overrides)
        resource = new_api.objects.select { |p| p.name == 'MyResource' }.first
        expect(resource.description).to eq(overrides['MyResource']['description'])
      }
    end

    describe 'should be able to override a property field' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'ReferencedResource' => Overrides::ResourceOverride.new(
            'properties' => {
              'name' => Overrides::PropertyOverride.new(
                'description' => 'My overridden description'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        new_api = Overrides::Runner.build(api, overrides)
        resource = new_api.objects.select { |p| p.name == 'ReferencedResource' }.first
        prop = resource.properties.select { |p| p.name == 'name' }.first
        expect(prop.description).to eq(
          overrides['ReferencedResource']['properties']['name']['description']
        )
      }
    end

    describe 'should be able to override a property type' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'ReferencedResource' => Overrides::ResourceOverride.new(
            'properties' => {
              'name' => Overrides::PropertyOverride.new(
                'type' => 'Api::Type::Integer'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        new_api = Overrides::Runner.build(api, overrides)
        resource = new_api.objects.select { |p| p.name == 'ReferencedResource' }.first
        prop = resource.properties.select { |p| p.name == 'name' }.first
        expect(prop.class).to eq(Api::Type::Integer)
      }
    end

    describe 'should be able to override a nested-nested property' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'AnotherResource' => Overrides::ResourceOverride.new(
            'properties' => {
              'nested-property2.property1.property1-nested' =>
              Overrides::PropertyOverride.new(
                'type' => 'Api::Type::Integer'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        new_api = Overrides::Runner.build(api, overrides)
        resource = new_api.objects.select { |p| p.name == 'AnotherResource' }.first
        prop = resource.properties.select { |p| p.name == 'nested-property2' }.first
        expect(prop.properties[0].properties[0].class).to eq(Api::Type::Integer)
      }
    end

    describe 'should be able to override a nested array property' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'AnotherResource' => Overrides::ResourceOverride.new(
            'properties' => {
              'array-property.property1' => Overrides::PropertyOverride.new(
                'type' => 'Api::Type::Integer'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        new_api = Overrides::Runner.build(api, overrides)
        resource = new_api.objects.select { |p| p.name == 'AnotherResource' }.first
        prop = resource.properties.select { |p| p.name == 'array-property' }.first
        expect(prop.item_type.properties[0].class).to eq(Api::Type::Integer)
      }
    end

    describe 'should be able to override a namevalue -> object map' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'AnotherResource' => TestResourceOverride.new(
            'properties' => {
              'namevalue-property.nv-prop1' => Overrides::PropertyOverride.new(
                'description' => 'overridden'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        new_api = Overrides::Runner.build(api, overrides, TestResourceOverride)
        resource = new_api.objects.select { |p| p.name == 'AnotherResource' }.first
        prop = resource.properties.select { |p| p.name == 'namevalue-property' }
                       .first
                       .value_type
                       .properties.first
        expect(prop.description).to eq('overridden')
      }
    end
  end
end
