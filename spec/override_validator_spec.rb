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

describe Overrides::Validator do
  context 'simple overrides' do
    describe 'should be able identify a bad field' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'AnotherResource' => Overrides::ResourceOverride.new(
            'blahbad' => 'A description'
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        runner = Overrides::Validator.new(api, overrides)
        expect { runner.run }.to raise_error(RuntimeError,
                                             /@blahbad does not exist on AnotherResource/)
      }
    end
    describe 'should be able identify a bad property' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'AnotherResource' => Overrides::ResourceOverride.new(
            'properties' => {
              'blahbad' => TestResourceOverride.new(
                'description' => 'A description'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        runner = Overrides::Validator.new(api, overrides)
        expect { runner.run }.to raise_error(RuntimeError,
                                             /blahbad does not exist on AnotherResource/)
      }
    end

    describe 'should be able validate a namevalues nestedobject properly' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'AnotherResource' => Overrides::ResourceOverride.new(
            'properties' => {
              'namevalue-property.nv-prop1' => TestResourceOverride.new(
                'description' => 'A description'
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        runner = Overrides::Validator.new(api, overrides)
        expect { runner.run }.not_to raise_error
      }
    end

    describe 'should be able validate a changed type with new properties' do
      let(:overrides) do
        Overrides::ResourceOverrides.new(
          'AnotherResource' => Overrides::ResourceOverride.new(
            'properties' => {
              'namevalue-property.nv-prop1' => TestResourceOverride.new(
                'type' => 'Api::Type::Enum',
                'values' => %i[test1 test2]
              )
            }
          )
        )
      end
      let(:api) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

      it {
        runner = Overrides::Validator.new(api, overrides)
        expect { runner.run }.not_to raise_error
      }
    end
  end
end
