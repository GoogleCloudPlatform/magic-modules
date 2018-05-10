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

describe Provider::Chef::Config do
  it 'returns Provider::Chef as provider' do
    expect(Provider::Chef::Config.new.provider).to be Provider::Chef
  end
end

describe Provider::Chef do
  context 'one type generated' do
    let(:product) { Api::Compiler.new('spec/data/good-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/chef-config.yaml', product)
    end
    let(:provider) { Provider::Chef.new(config, product) }

    before do
      allow_open 'spec/data/good-file.yaml'
      allow_open 'spec/data/chef-config.yaml'
      allow_open_license
      allow_open_properties
      allow_open_libraries
      allow_open_typed_array
      allow_open_spec_templates(true)
      product.validate
    end

    it do
      out = dummy_writer
      output_expectations kind: 'myproduct', name: 'another_resource',
                          provider: { writer: out, tester: out },
                          recipe: out,
                          resourceref: {
                            name: 'referenced_resource',
                            imports: 'name',
                            writer: out
                          },
                          named_prop: 'anotherresource',
                          arrays: [:string_array]

      provider.generate 'blah', []
    end
  end

  context 'exports generated' do
    before do
      allow_open 'spec/data/good-export-file.yaml'
      allow_open 'spec/data/chef-config.yaml'
      allow_open_license
      allow_open_properties
      allow_open_libraries
      allow_open_typed_array
      allow_open_real_templates
      product.validate
    end

    let(:config) { Provider::Config.parse('spec/data/chef-config.yaml') }
    let(:product) { Api::Compiler.new('spec/data/good-export-file.yaml').run }
    let(:provider) { Provider::Chef.new(config, product) }
    let(:expected) do
      ["self_link: __fetched['selfLink']",
       'property1: property1',
       'super_long_name: super_long_name'].to_set
    end

    let(:matched) do
      matched = Set.new
      out = dummy_writer
      provider_tester = mock('File')
      provider_tester.stubs(:write).with(anything) do |arg|
        match = expected.select { |e| arg.include?(e) }
        matched << match[0] unless match.empty?
        arg
      end
      output_expectations kind: 'myproduct', name: 'another_resource',
                          provider: { writer: provider_tester, tester: out },
                          recipe: out,
                          resourceref: {
                            name: 'referenced_resource',
                            imports: 'name',
                            writer: out
                          },
                          named_prop: 'anotherresource',
                          arrays: [:string_array]

      provider.generate 'blah', []
      matched
    end

    subject { matched }
    it { is_expected.to eq expected.to_set }
  end

  private

  def allow_open(file_name)
    IO.expects(:read).with(file_name).returns(File.real_read(file_name))
      .at_least(0)
  end

  def allow_read(file_name)
    File.expects(:read).at_least(0).with(file_name)
        .returns(File.real_read(file_name))
  end

  def file_read_map(target, source)
    File.expects(:read).with(target).returns(File.real_read(source))
  end

  # rubocop:disable Metrics/AbcSize
  def output_expectations(data)
    output_file "blah/resources/#{data[:name]}.rb",
                data[:provider][:writer]
    output_file "blah/spec/#{data[:name]}_spec.rb", data[:provider][:tester]
    property_expectations(data)
    output_expectations_libraries(data)
    output_expectations_string_array data
    output_expectations_named_properties(data)

    return unless data[:resourceref]

    # ResourceRef
    output_file "blah/resources/#{data[:resourceref][:name]}.rb",
                data[:resourceref][:writer]
    output_file "blah/spec/#{data[:resourceref][:name]}_spec.rb",
                data[:resourceref][:writer]
    output_real_network_data("#{data[:kind]}_#{data[:resourceref][:name]}")
    allow_open_real_network_data(data)
  end
  # rubocop:enable Metrics/AbcSize

  def output_expectations_named_properties(data)
    return unless data[:named_prop]
    output_file \
      ["blah/libraries/google/#{data[:kind][1..-1]}",
       "property/#{data[:named_prop]}_array_property.rb"].join('/'),
      dummy_writer
    output_file \
      ["blah/libraries/google/#{data[:kind][1..-1]}",
       "property/#{data[:named_prop]}_nested_property.rb"].join('/'),
      dummy_writer
  end

  def output_expectations_string_array(data)
    return if data[:arrays].nil? || !data[:arrays].include?(:string_array)
    output_file \
      "blah/libraries/google/#{data[:kind][1..-1]}/property/string_array.rb",
      dummy_writer
    output_file \
      "blah/libraries/google/#{data[:kind][1..-1]}/property/array.rb",
      dummy_writer
  end

  # rubocop:disable Metrics/AbcSize
  # rubocop:disable Metrics/MethodLength
  def output_expectations_libraries(data)
    dw = dummy_writer
    output_file \
      "blah/libraries/google/#{data[:kind][1..-1]}/network/base.rb", dw
    output_file \
      "blah/libraries/google/#{data[:kind][1..-1]}/network/delete.rb", dw
    output_file \
      "blah/libraries/google/#{data[:kind][1..-1]}/network/get.rb", dw
    output_file \
      "blah/libraries/google/#{data[:kind][1..-1]}/network/post.rb", dw
    output_file \
      "blah/libraries/google/#{data[:kind][1..-1]}/network/put.rb", dw
    output_file 'blah/spec/network_delete_spec.rb', dw
    output_file 'blah/spec/network_get_spec.rb', dw
    output_file 'blah/spec/network_post_spec.rb', dw
    output_file 'blah/spec/network_put_spec.rb', dw
    output_file 'blah/spec/network_blocker.rb', dw
    output_file 'blah/spec/network_blocker_spec.rb', dw
  end
  # rubocop:enable Metrics/MethodLength
  # rubocop:enable Metrics/AbcSize

  def property_expectations(data)
    dw = dummy_writer
    output_file \
      "blah/libraries/google/#{data[:kind][1..-1]}/property/enum.rb", dw
    output_file \
      "blah/libraries/google/#{data[:kind][1..-1]}/property/string.rb", dw
    return unless data[:resourceref]

    resourceref_name = [data[:resourceref][:name].delete('_'),
                        data[:resourceref][:imports]].join('_')
    output_file \
      ["blah/libraries/google/#{data[:kind][1..-1]}/property",
       "#{resourceref_name}.rb"].join('/'), dw
  end

  def allow_open_spec_templates(resourceref = false)
    file_read_map 'templates/chef/resource.erb', 'spec/data/prov_template.erb'
    file_read_map 'templates/chef/resource_spec.erb',
                  'spec/data/prov_spec_template.erb'
    6.times.each do
      file_read_map 'templates/network_spec.yaml.erb',
                    'spec/data/network_spec_template.erb'
    end

    allow_open_spec_templates(false) if resourceref
  end

  def allow_open_real_templates
    allow_read 'templates/chef/resource_spec.erb'
    allow_read 'templates/chef/resource.erb'

    allow_open 'templates/autogen_notice.erb'
    allow_open 'templates/expand_variables.erb'
    allow_open 'templates/network_mocks.erb'
    allow_open 'templates/network_spec.yaml.erb'
    allow_open 'templates/provider_helpers.erb'
    allow_open 'templates/chef/resourceref_expandvars.erb'
    allow_open 'templates/resourceref_mocks.erb'
    allow_open 'templates/return_if_object.erb'
    allow_open 'templates/transport.erb'

    allow_open_real_tests

    allow_open_license
  end

  def allow_open_real_tests
    allow_open 'templates/chef/tests/present~create.erb'
    allow_open 'templates/chef/tests/absent~no_action.erb'
    allow_open 'templates/chef/tests/absent~delete.erb'
    allow_open 'templates/chef/tests/present~no_changes.erb'
  end

  def allow_open_libraries
    allow_read 'templates/network/base.rb.erb'
    allow_read 'templates/network/delete.rb.erb'
    allow_read 'templates/network/delete_spec.rb.erb'
    allow_read 'templates/network/get.rb.erb'
    allow_read 'templates/network/get_spec.rb.erb'
    allow_read 'templates/network/post.rb.erb'
    allow_read 'templates/network/post_spec.rb.erb'
    allow_read 'templates/network/put.rb.erb'
    allow_read 'templates/network/put_spec.rb.erb'
    allow_read 'templates/network/network_blocker.rb.erb'
    allow_read 'templates/network/network_blocker_spec.rb.erb'
  end

  def allow_open_typed_array
    allow_open 'templates/chef/property/array_typed.rb.erb'
    allow_open 'templates/chef/property/array.rb.erb'
  end

  def allow_open_properties
    allow_read 'templates/chef/property/array.rb.erb'
    allow_read 'templates/chef/property/enum.rb.erb'
    allow_read 'templates/chef/property/nested_object.rb.erb'
    allow_read 'templates/chef/property/string.rb.erb'
    allow_read 'templates/chef/property/resourceref.rb.erb'
  end

  def allow_open_real_network_data(data)
    name = "#{data[:kind]}_#{data[:name]}"
    output_real_network_data(name)
  end

  def output_real_network_data(name)
    dw = dummy_writer
    3.times.each do |id|
      %w[name title].each do |title|
        output_file ['blah/spec/data/network',
                     "#{name}/success#{id + 1}~#{title}.yaml"].join('/'),
                    dw
      end
    end
  end

  def allow_open_license
    allow_open 'templates/license.erb'
    allow_open 'templates/autogen_notice.erb'
  end

  def output_file(name, output)
    File.expects(:open).with(name, 'w').yields(output)
        .at_least(0)
  end

  def dummy_writer
    out = mock('File')
    out.expects(:write).at_least(0)
    out
  end
end
