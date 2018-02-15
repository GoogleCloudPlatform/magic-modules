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
  context 'static' do
    let(:config) { Provider::Config.parse('spec/data/terraform-config.yaml') }
    let(:product) { Api::Compiler.new('spec/data/good-file.yaml').run }
    let(:provider) { Provider::Terraform.new(config, product) }

    before do
      allow_open 'spec/data/good-file.yaml'
      allow_open 'spec/data/terraform-config.yaml'
      product.validate
    end

    it 'should generate all accepted import id formats' do
      a_resource = resource(
        'base_url: "projects/{{project}}/regions/{{region}}/subnetworks"'
      )
      formats = provider.import_id_formats(a_resource)

      expect(formats).to contain_exactly(
        'projects/{{project}}/regions/{{region}}/subnetworks/{{name}}',
        '{{project}}/{{region}}/{{name}}',
        '{{name}}'
      )
    end

    it 'should transform id formats to a regex' do
      r = provider.format2regex 'projects/{{project}}/global/networks/{{name}}'
      expect(r).to eq(
        'projects/(?P<project>[^/]+)/global/networks/(?P<name>[^/]+)'
      )
    end
  end

  def allow_open(file_name)
    IO.expects(:read).with(file_name).returns(File.real_read(file_name))
      .at_least(0)
  end

  def resource(*data)
    Google::YamlValidator.parse(['--- !ruby/object:Api::Resource'].concat(data)
                                                                  .join("\n"))
  end
end
