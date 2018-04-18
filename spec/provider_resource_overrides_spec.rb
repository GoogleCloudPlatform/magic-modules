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

class File
  class << self
    alias real_open open
    alias real_read read
  end
end

describe Provider::ResourceOverrides do
  context 'good file product' do
    let(:product) { Api::Compiler.new('spec/data/good-file.yaml').run }

    before(:each) do
      allow_open 'spec/data/good-file.yaml'

      product.validate
    end

    context 'with overrides' do
      let(:overrides) do
        Provider::ResourceOverride.parse(
          IO.read('spec/data/good-resource-overrides.yaml')
        )
      end

      before(:each) do
        allow_open 'spec/data/good-resource-overrides.yaml'

        overrides.consume_api(product)
        overrides.validate
      end

      context 'for resource' do
        let(:resource) do
          product.objects.find { |o| o.name == 'AnotherResource' }
        end

        context 'overrides resource description' do
          subject { resource.description }
          it { is_expected.to eq 'blah blah bar' }
        end

        context 'overrides property description' do
          subject do
            resource.properties.find { |p| p.name == 'property1' }.description
          end
          it { is_expected.to eq 'foo' }
        end

        context 'overrides nested property description' do
          subject do
            property = resource.properties.find do |p|
              p.name == 'nested-property'
            end

            nested_property = property.properties.find do |p|
              p.name == 'property1'
            end
            nested_property.description
          end
          it { is_expected.to eq 'bar' }
        end

        context 'overrides array of nested property description' do
          subject do
            property = resource.properties.find do |p|
              p.name == 'array-property'
            end
            nested_property = property.item_type.properties.find do |p|
              p.name == 'property1'
            end
            nested_property.description
          end
          it { is_expected.to eq 'baz' }
        end
      end
    end

    context 'with resource with invalid property path' do
      context 'referring to missing top-level property' do
        subject { -> { create_overrides(product, 'missing-property') } }

        it { is_expected.to raise_error StandardError, /missing-property/ }
      end

      context 'referring to missing nested property' do
        subject do
          -> { create_overrides(product, 'nested-property.missing-property') }
        end

        it do
          is_expected.to raise_error(
            StandardError,
            /nested-property.missing-property/
          )
        end
      end

      context 'referring to missing array nested property' do
        subject do
          -> { create_overrides(product, 'array-property.missing-property') }
        end

        it do
          is_expected.to raise_error(
            StandardError,
            /array-property.missing-property/
          )
        end
      end

      context 'referring to a nested property in non-nested type' do
        subject do
          -> { create_overrides(product, 'property1.missing-property') }
        end

        it do
          is_expected.to raise_error(
            StandardError,
            /property1.missing-property/
          )
        end
      end
    end
  end

  def create_overrides(product, property_path)
    overrides = Google::YamlValidator.parse(
      %(
      !ruby/object:Provider::ResourceOverrides
      AnotherResource: !ruby/object:Provider::Terraform::ResourceOverride
        properties:
          #{property_path}: !ruby/object:Provider::Terraform::PropertyOverride
            description: foobar
      )
    )

    overrides.consume_api(product)
    overrides.validate

    overrides
  end

  def allow_open(file_name)
    IO.expects(:read).with(file_name).returns(File.real_read(file_name))
      .at_least(1)
  end
end
