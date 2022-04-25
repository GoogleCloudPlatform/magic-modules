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
require 'api/object'

describe Api::Resource do
  context 'uses the correct API version' do
    let(:product) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }

    before { product.validate }

    subject do
      product.objects.find { |o| o.name == 'AnotherResource' }
    end

    context 'ga' do
      it do
        version = product.version_obj('ga')
        subject.exclude_if_not_in_version!(version)
        is_expected.not_to(contain_property_with_name('beta-property'))
        is_expected.to(contain_property_with_name('property1'))
      end
    end

    context 'beta' do
      it do
        version = product.version_obj('beta')
        subject.exclude_if_not_in_version!(version)
        is_expected.to(contain_property_with_name('beta-property'))
        is_expected.to(contain_property_with_name('property1'))
      end
    end
  end

  # TODO: Fill in these tests or get rid of them completely
  it 'uses product base_url if missing' do
  end

  it 'ignores product base_url if absolute' do
  end

  it 'combines base_url with product if relative' do
  end
end

RSpec::Matchers.define :contain_property_with_name do |expected|
  match do |actual|
    actual.properties.find { |p| p.name == expected }
  end
  failure_message do |actual|
    "expected #{actual.properties.map(&:name)} to contain #{expected}"
  end
  failure_message_when_negated do |actual|
    "expected #{actual.properties.map(&:name)} not to contain #{expected}"
  end
end
