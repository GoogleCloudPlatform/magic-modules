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
    let(:config) { Provider::Config.parse('spec/data/terraform-config.yaml') }
    let(:product) { Api::Compiler.new('spec/data/good-file.yaml').run }
    let(:provider) { Provider::Terraform.new(config, product) }

    before do
      allow_open 'spec/data/good-file.yaml'
      allow_open 'spec/data/terraform-config.yaml'
      product.validate
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
      subject { provider.collection_url(product.objects[0]) }
      it do
        is_expected.to eq 'http://myproduct.google.com/api/referencedresource'
      end
    end

    describe '#self_link_url' do
      subject { provider.self_link_url(product.objects[0]) }
      it do
        is_expected.to eq(
          'http://myproduct.google.com/api/referencedresource/{{name}}'
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
end
