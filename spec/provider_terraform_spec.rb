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

describe Provider::Terraform do
  context 'good file product' do
    let(:product) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }
    let(:config) do
      Provider::Config.parse('spec/data/terraform-config.yaml', product)[1]
    end
    let(:provider) { Provider::Terraform.new(config, product, 'ga', Time.now) }
    let(:resource) { product.objects[0] }

    before do
      allow_open 'spec/data/good-file.yaml'
      allow_open 'spec/data/terraform-config.yaml'
      product.validate
      config.validate
    end

    describe '#format2regex' do
      subject do
        provider.format2regex 'projects/{{project}}/global/networks/{{name}}'
      end

      it do
        is_expected.to eq(
          'projects/(?P<project>[^/]+)/global/networks/(?P<name>[^/]+)'
        )
      end
    end

    context '#titlelize_property' do
      describe 'short property name' do
        subject { provider.titlelize_property(named_property('fooBar')) }
        it { is_expected.to eq 'FooBar' }
      end

      describe 'titlelizes long property name' do
        subject do
          provider.titlelize_property(named_property('fooBarBazFooBar'))
        end
        it { is_expected.to eq 'FooBarBazFooBar' }
      end
    end

    describe '#collection_url' do
      subject { resource.collection_url }
      it do
        version = product.version_obj_or_closest(nil)
        product.set_properties_based_on_version(version)
        is_expected.to eq 'http://myproduct.google.com/api/referencedresource'
      end
    end

    describe '#collection_url beta' do
      subject { resource.collection_url }
      it do
        version = product.version_obj_or_closest('beta')
        product.set_properties_based_on_version(version)
        is_expected.to eq 'http://myproduct.google.com/api/beta/referencedresource'
      end
    end

    describe '#self_link_url' do
      subject { resource.self_link_url }
      it do
        version = product.version_obj_or_closest(nil)
        product.set_properties_based_on_version(version)
        is_expected.to eq(
          'http://myproduct.google.com/api/referencedresource/{{name}}'
        )
      end
    end

    describe '#self_link_url beta' do
      subject { resource.self_link_url }
      it do
        version = product.version_obj_or_closest('beta')
        product.set_properties_based_on_version(version)
        is_expected.to eq(
          'http://myproduct.google.com/api/beta/referencedresource/{{name}}'
        )
      end
    end

    describe '#properties_by_custom_update' do
      let(:postUrl1) { custom_update_property('p1', 'url1', :POST) }
      let(:otherPostUrl1) { custom_update_property('p2', 'url1', :POST) }
      let(:postUrl2) { custom_update_property('p3', 'url2', :POST) }
      let(:putUrl2) { custom_update_property('p4', 'url2', :PUT) }
      let(:props) do
        [
          custom_update_property('no-custom-update'),
          postUrl1, otherPostUrl1, postUrl2, putUrl2
        ]
      end
      subject { provider.properties_by_custom_update(props) }

      it do
        is_expected.to eq(
          {
            update_url: 'url1',
            update_verb: :POST,
            update_id: nil,
            fingerprint_name: nil
          } =>
            [postUrl1, otherPostUrl1],
          {
            update_url: 'url2',
            update_verb: :POST,
            update_id: nil,
            fingerprint_name: nil
          } => [postUrl2],
          {
            update_url: 'url2',
            update_verb: :PUT,
            update_id: nil,
            fingerprint_name: nil
          } => [putUrl2]
        )
      end
    end
  end

  def allow_open(file_name)
    IO.expects(:read).with(file_name).returns(File.real_read(file_name))
      .at_least(0)
  end

  def named_property(name)
    Google::YamlValidator.parse(
      format("--- !ruby/object:Api::Object::Named\nname: '%<name>s'",
             name: name)
    )
  end

  def custom_update_property(name, update_url = nil, update_verb = nil)
    lines = []
    lines.push '--- !ruby/object:Api::Type::String'
    lines.push "name: '#{name}'"
    lines.push "update_url: '#{update_url}'" unless update_url.nil?
    lines.push "update_verb: :#{update_verb}" unless update_verb.nil?

    Google::YamlValidator.parse(lines.join("\n"))
  end
end
