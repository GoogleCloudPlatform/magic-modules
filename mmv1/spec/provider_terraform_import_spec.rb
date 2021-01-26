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
require 'provider/terraform'

class File
  class << self
    alias real_open open
    alias real_read read
  end
end

describe Provider::Terraform do
  context 'static' do
    let(:product) { Api::Compiler.new(File.read('spec/data/good-file.yaml')).run }
    let(:config) do
      Provider::Config.parse('spec/data/terraform-config.yaml', product)[1]
    end
    let(:provider) { Provider::Terraform.new(config, product, 'ga', Time.now) }

    before do
      allow_open 'spec/data/good-file.yaml'
      allow_open 'spec/data/terraform-config.yaml'
      product.validate
      config.validate
    end

    describe '#import_id_formats_from_resource' do
      subject do
        provider.import_id_formats_from_resource(
          resource(
            'base_url: "projects/{{project}}/regions/{{region}}/subnetworks"'
          )
        )
      end

      it do
        is_expected.to contain_exactly(
          'projects/{{project}}/regions/{{region}}/subnetworks/{{name}}',
          '{{project}}/{{region}}/{{name}}',
          '{{region}}/{{name}}',
          '{{name}}'
        )
      end
    end

    def allow_open(file_name)
      IO.expects(:read).with(file_name).returns(File.real_read(file_name))
        .at_least(0)
    end

    def resource(*data)
      res = Google::YamlValidator.parse(['--- !ruby/object:Api::Resource']
                                 .concat(["name: 'testobject'"])
                                 .concat(data)
                                 .join("\n"))
      product.objects.append(res)
      new_product = Overrides::Runner.build(product, config.overrides,
                                            config.resource_override,
                                            config.property_override)
      new_product.objects.last
    end
  end
end
