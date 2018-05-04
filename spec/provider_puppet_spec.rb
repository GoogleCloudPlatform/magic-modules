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

require 'set'
require 'spec_helper'

class File
  class << self
    alias real_open open
    alias real_read read
  end
end

describe Provider::Puppet::Config do
  it 'returns Provider::Puppet as provider' do
    expect(Provider::Puppet::Config.new.provider).to be Provider::Puppet
  end
end

describe Provider::Puppet do
  context 'one type with resourceref generated' do
    let(:product) { Api::Compiler.new('spec/data/good-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-config.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }

    before do
      allow_open 'spec/data/good-file.yaml'
      allow_open 'spec/data/puppet-config.yaml'
      allow_open_license
      allow_open_typed_array
      allow_open_properties
      allow_open_libraries
      allow_open_spec_templates(true)
      product.validate
    end

    it do
      out = dummy_writer
      output_expectations kind: 'myproduct', name: 'another_resource',
                          type: { writer: out, tester: out },
                          provider: { writer: out, tester: out },
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

  context 'multiple types generated' do
    let(:product) { Api::Compiler.new('spec/data/good-multi-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-config.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }

    before do
      allow_open 'spec/data/good-multi-file.yaml'
      allow_open 'spec/data/puppet-config.yaml'
      allow_open_typed_array
      allow_open_properties
      allow_open_libraries
      allow_open_real_templates
      product.validate
    end

    it do
      out1 = dummy_writer
      out2 = dummy_writer
      output_expectations kind: 'multiproduct', name: 'my_resource',
                          type: { writer: out1, tester: out1 },
                          provider: { writer: out1, tester: out1 },
                          resourceref: {
                            name: 'referenced_resource',
                            imports: 'name',
                            writer: out1
                          }
      output_expectations kind: 'multiproduct', name: 'another_resource',
                          type: { writer: out2, tester: out2 },
                          provider: { writer: out2, tester: out2 },
                          arrays: [:string_array]
      provider.generate 'blah', []
    end
  end

  context 'test template' do
    let(:product) { Api::Compiler.new('spec/data/good-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-config.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }

    before do
      allow_open 'spec/data/good-file.yaml'
      allow_open 'spec/data/puppet-config.yaml'
      allow_open_license
      allow_open_typed_array
      allow_open_properties
      allow_open_libraries
      allow_open_spec_templates(true)
      product.validate
    end

    it do
      out = dummy_writer
      type_writer = mock('File')
      type_writer.expects(:write).with("type: myproduct_another_resource\n")
      type_writer.expects(:write).with("property: property1 (String)\n")
      type_writer.expects(:write).with("property: property2 (String)\n")
      type_writer.expects(:write).with("property: property3 (Array)\n")
      type_writer.expects(:write).with("property: property4 (Enum)\n")
      type_writer.expects(:write).with("parameter: property5 (ResourceRef)\n")
      type_writer.expects(:write).with(
        "property: nested_property (NestedObject)\n"
      )
      type_writer.expects(:write).with(
        "property: array_property (Array)\n"
      )
      provider_writer = mock('File')
      provider_writer.expects(:write).with("type: myproduct_another_resource\n")
      output_expectations kind: 'myproduct', name: 'another_resource',
                          type: { writer: type_writer, tester: out },
                          provider: { writer: provider_writer, tester: out },
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

  context 'long collection uri' do
    let(:product) { Api::Compiler.new('spec/data/good-longuri.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-longuri-config.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }
    let(:expected) do
      ['URI.join(', "'http://myproduct.google.com/myapi/v123',",
       'expand_variables(', "'some/long/path/over/80chars',", 'data', ')',
       ')']
    end

    before do
      allow_open 'spec/data/good-longuri.yaml'
      allow_open 'spec/data/puppet-longuri-config.yaml'
      allow_open_real_templates
      allow_open_properties
      allow_open_libraries
      product.validate
    end

    it do
      collected = []
      out = dummy_writer
      provider_writer = mock('File')
      provider_writer.stubs(:write).with(anything) do |arg|
        collected << arg.strip
      end
      output_expectations kind: 'myproduct', name: 'my_resource',
                          type: { writer: out, tester: out },
                          provider: { writer: provider_writer, tester: out }
      provider.generate 'blah', []
      expect(collected).to contain_array(expected)
    end
  end

  context 'multiline code format' do
    let(:product) { Api::Compiler.new('spec/data/good-single-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-multicode.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }
    let(:expected) do
      ["    longline-create-line1\n",
       "    longline-create-line2\n",
       "    longline-create-line3\n",
       "    longline-create-line4\n"]
    end

    before do
      allow_open 'spec/data/good-single-file.yaml'
      allow_open 'spec/data/puppet-multicode.yaml'
      allow_open_real_templates
      allow_open_properties
      allow_open_libraries
      product.validate
    end

    it do
      matched = []
      out = dummy_writer
      provider_writer = mock('File')
      provider_writer.stubs(:write).with(anything) do |arg|
        match = expected.select { |e| arg.include?(e) }
        matched << match[0] unless match.empty?
        arg
      end
      output_expectations kind: 'myproduct', name: 'my_resource',
                          type: { writer: out, tester: out },
                          provider: { writer: provider_writer, tester: out }
      provider.generate 'blah', []
      expect(matched).to eq expected
    end
  end

  context 'readonly false sets all require requests' do
    let(:product) { Api::Compiler.new('spec/data/good-single-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-multicode.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }
    let(:expected) do
      ["require 'google/yproduct/network/delete'\n",
       "require 'google/yproduct/network/post'\n",
       "require 'google/yproduct/network/put'\n"]
    end

    before do
      allow_open 'spec/data/puppet-multicode.yaml'
      allow_open 'spec/data/good-single-file.yaml'
      allow_open_real_templates
      allow_open_properties
      allow_open_libraries
      product.validate
    end

    it 'blank or false readonly should have all requests required' do
      matched = []
      out = dummy_writer
      provider_writer = mock('File')
      provider_writer.stubs(:write).with(anything) do |arg|
        match = expected.select { |e| arg.include?(e) }
        matched << match[0] unless match.empty?
        arg
      end
      output_expectations kind: 'myproduct', name: 'my_resource',
                          type: { writer: out, tester: out },
                          provider: { writer: provider_writer, tester: out }
      provider.generate 'blah', []
      expect(matched).to eq expected
    end
  end

  context 'readonly true sets only get require requests' do
    let(:product) do
      Api::Compiler.new('spec/data/good-single-readonly-file.yaml').run
    end
    let(:config) do
      Provider::Config.parse('spec/data/puppet-multicode.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }
    let(:not_expected) do
      ["require 'google/request/post'\n",
       "require 'google/request/delete'\n"]
    end

    before do
      allow_open 'spec/data/puppet-multicode.yaml'
      allow_open 'spec/data/good-single-readonly-file.yaml'
      allow_open_real_templates
      allow_open_properties
      allow_open_libraries
      product.validate
    end

    it 'blank or false readonly should have all requests required' do
      matched = []
      out = dummy_writer
      provider_writer = mock('File')
      provider_writer.stubs(:write).with(anything) do |arg|
        match = not_expected.select { |e| arg.include?(e) }
        matched << match[0] unless match.empty?
        arg
      end
      output_expectations kind: 'myproduct', name: 'my_resource',
                          type: { writer: out, tester: out },
                          provider: { writer: provider_writer, tester: out }
      provider.generate 'blah', []
      expect(matched).to eq []
    end
  end

  context 'filters creation based on command line list' do
    let(:product) { Api::Compiler.new('spec/data/good-multi2-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-config.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }
    let(:type) { 'YetAnotherResource' }

    before do
      allow_open 'spec/data/good-multi2-file.yaml'
      allow_open 'spec/data/puppet-config.yaml'
      allow_open_spec_templates
      allow_open_properties
      allow_open_libraries
      allow_open_license
      product.validate
    end

    it do
      out = dummy_writer
      output_expectations kind: 'multiproduct', name: 'yet_another_resource',
                          type: { writer: out, tester: out },
                          provider: { writer: out, tester: out }

      provider.generate 'blah', [type]
    end
  end

  context 'exports proper provider info when specified' do
    before do
      allow_open 'spec/data/good-export-file.yaml'
      allow_open 'spec/data/puppet-config.yaml'
      allow_open_typed_array
      allow_open_real_templates
      allow_open_properties
      allow_open_libraries
      product.validate
    end

    let(:product) { Api::Compiler.new('spec/data/good-export-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-config.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }
    let(:expected) do
      ['property1: resource[:property1]',
       'self_link: @fetched[\'selfLink\']',
       'super_long_name: resource[:super_long_name]'].to_set
    end

    let(:matched) do
      matched = []
      out = dummy_writer
      provider_tester = mock('File')
      provider_tester.stubs(:write).with(anything) do |arg|
        match = expected.select { |e| arg.include?(e) }
        matched << match[0] unless match.empty?
        arg
      end
      output_expectations kind: 'myproduct', name: 'another_resource',
                          type: { writer: out, tester: out },
                          provider: { writer: provider_tester, tester: out },
                          resourceref: {
                            name: 'referenced_resource',
                            imports: 'name',
                            writer: out
                          },
                          arrays: [:string_array]
      provider.generate 'blah', []
      matched
    end

    subject { matched.to_set }

    it { is_expected.to eq expected }
  end

  context 'exports proper type info when specified' do
    before do
      allow_open 'spec/data/good-export-file.yaml'
      allow_open 'spec/data/puppet-config.yaml'
      allow_open_typed_array
      allow_open_real_templates
      allow_open_properties
      allow_open_libraries
      product.validate
    end

    let(:product) { Api::Compiler.new('spec/data/good-export-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-config.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }
    let(:expected) do
      ['def exports',
       'provider.exports']
    end

    let(:matched) do
      matched = []
      out = dummy_writer
      type_tester = mock('File')
      type_tester.stubs(:write).with(anything) do |arg|
        match = expected.select { |e| arg.include?(e) }
        matched << match[0] unless match.empty?
        arg
      end
      output_expectations kind: 'myproduct', name: 'another_resource',
                          type: { writer: type_tester, tester: out },
                          provider: { writer: out, tester: out },
                          resourceref: {
                            name: 'referenced_resource',
                            imports: 'name',
                            writer: out
                          },
                          arrays: [:string_array]
      provider.generate 'blah', []
      matched
    end

    subject { matched }

    it { is_expected.to eq expected }
  end

  context 'does not exports proper provider info when not specified' do
    before do
      allow_open 'spec/data/good-file.yaml'
      allow_open 'spec/data/puppet-config.yaml'
      allow_open_typed_array
      allow_open_real_templates
      allow_open_properties
      allow_open_libraries
      product.validate
    end

    let(:product) { Api::Compiler.new('spec/data/good-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-config.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }
    let(:expected) do
      ['self_link: @fetched[\'selfLink\']',
       'property1: @fetched[\'property1\']',
       'super_long_name: @fetched[\'superLongName\']']
    end

    let(:matched) do
      matched = []
      out = dummy_writer
      provider_tester = mock('File')
      provider_tester.stubs(:write).with(anything) do |arg|
        match = expected.select { |e| arg.include?(e) }
        matched << match[0] unless match.empty?
        arg
      end
      output_expectations kind: 'myproduct', name: 'another_resource',
                          type: { writer: out, tester: out },
                          provider: { writer: provider_tester, tester: out },
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

    it { is_expected.to eq [] }
  end

  context 'does not exports proper type info when not specified' do
    before do
      allow_open 'spec/data/good-file.yaml'
      allow_open 'spec/data/puppet-config.yaml'
      allow_open_typed_array
      allow_open_real_templates
      allow_open_properties
      allow_open_libraries
      product.validate
    end

    let(:product) { Api::Compiler.new('spec/data/good-file.yaml').run }
    let(:config) do
      Provider::Config.parse('spec/data/puppet-config.yaml', product)
    end
    let(:provider) { Provider::Puppet.new(config, product) }
    let(:expected) do
      ['def exports',
       'provider.exports']
    end

    let(:matched) do
      matched = []
      out = dummy_writer
      type_tester = mock('File')
      type_tester.stubs(:write).with(anything) do |arg|
        match = expected.select { |e| arg.include?(e) }
        matched << match[0] unless match.empty?
        arg
      end
      output_expectations kind: 'myproduct', name: 'another_resource',
                          type: { writer: type_tester, tester: out },
                          provider: { writer: out, tester: out },
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

    it { is_expected.to eq [] }
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

  def output_expectations(data)
    name = "#{data[:kind]}_#{data[:name]}"
    output_file "blah/lib/puppet/type/#{name}.rb", data[:type][:writer]
    output_file "blah/lib/puppet/provider/#{name}/google.rb",
                data[:provider][:writer]
    output_file "blah/spec/#{name}_provider_spec.rb", data[:provider][:tester]

    output_expectations_properties data
    output_expectations_libraries data
    output_expectations_string_array data
    output_expectations_resource_ref data unless data[:resourceref].nil?
    output_expectations_network_data data
    output_expectations_named_properties data
  end

  def output_expectations_network_data(data)
    dw = dummy_writer
    name = "#{data[:kind]}_#{data[:name]}"
    3.times.each do |id|
      %w[name title].each do |title|
        output_file ['blah/spec/data/network',
                     "#{name}/success#{id + 1}~#{title}.yaml"].join('/'),
                    dw
      end
    end
  end

  # rubocop:disable Metrics/AbcSize
  # rubocop:disable Metrics/MethodLength
  def output_expectations_resource_ref(data)
    rref_name = "#{data[:kind]}_#{data[:resourceref][:name]}"
    rref_file = [data[:resourceref][:name].tr('_', ''),
                 data[:resourceref][:imports]].join('_')

    output_file "blah/lib/puppet/type/#{rref_name}.rb",
                data[:resourceref][:writer]
    output_file "blah/lib/puppet/provider/#{rref_name}/google.rb",
                data[:resourceref][:writer]
    output_file "blah/spec/#{rref_name}_provider_spec.rb",
                data[:resourceref][:writer]

    output_file ["blah/lib/google/#{data[:kind][1..-1]}/property/",
                 "#{rref_file}.rb"].join,
                data[:resourceref][:writer]

    3.times.each do |id|
      %w[name title].each do |title|
        output_file ['blah/spec/data/network',
                     "#{rref_name}/success#{id + 1}~#{title}.yaml"].join('/'),
                    data[:resourceref][:writer]
      end
    end
  end
  # rubocop:enable Metrics/AbcSize
  # rubocop:enable Metrics/MethodLength

  def output_expectations_named_properties(data)
    return unless data[:named_prop]
    output_file \
      ["blah/lib/google/#{data[:kind][1..-1]}",
       "property/#{data[:named_prop]}_array_property.rb"].join('/'),
      dummy_writer
    output_file \
      ["blah/lib/google/#{data[:kind][1..-1]}",
       "property/#{data[:named_prop]}_nested_property.rb"].join('/'),
      dummy_writer
  end

  # rubocop:disable Metrics/AbcSize
  def output_expectations_libraries(data)
    dw = dummy_writer
    output_file "blah/lib/google/#{data[:kind][1..-1]}/network/base.rb", dw
    output_file "blah/lib/google/#{data[:kind][1..-1]}/network/delete.rb", dw
    output_file "blah/lib/google/#{data[:kind][1..-1]}/network/get.rb", dw
    output_file "blah/lib/google/#{data[:kind][1..-1]}/network/post.rb", dw
    output_file "blah/lib/google/#{data[:kind][1..-1]}/network/put.rb", dw
    output_file 'blah/spec/network_delete_spec.rb', dw
    output_file 'blah/spec/network_get_spec.rb', dw
    output_file 'blah/spec/network_post_spec.rb', dw
    output_file 'blah/spec/network_put_spec.rb', dw
    output_file 'blah/spec/network_blocker.rb', dw
    output_file 'blah/spec/network_blocker_spec.rb', dw
  end
  # rubocop:enable Metrics/AbcSize

  def output_expectations_properties(data)
    dw = dummy_writer
    output_file "blah/lib/google/#{data[:kind][1..-1]}/property/base.rb", dw
    output_file "blah/lib/google/#{data[:kind][1..-1]}/property/enum.rb", dw
    output_file "blah/lib/google/#{data[:kind][1..-1]}/property/resourceref.rb",
                dw
    output_file "blah/lib/google/#{data[:kind][1..-1]}/property/string.rb", dw
  end

  def output_expectations_string_array(data)
    return if data[:arrays].nil? || !data[:arrays].include?(:string_array)
    output_file \
      "blah/lib/google/#{data[:kind][1..-1]}/property/string_array.rb",
      dummy_writer
    output_file \
      "blah/lib/google/#{data[:kind][1..-1]}/property/array.rb",
      dummy_writer
    output_file \
      "blah/lib/google/#{data[:kind][1..-1]}/property/base.rb",
      dummy_writer
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

  def allow_open_properties
    allow_read 'templates/puppet/property/base.rb.erb'
    allow_read 'templates/puppet/property/enum.rb.erb'
    allow_read 'templates/puppet/property/nested_object.rb.erb'
    allow_read 'templates/puppet/property/resourceref.rb.erb'
    allow_read 'templates/puppet/property/string.rb.erb'
  end

  def allow_open_spec_templates(rref = false)
    file_read_map 'templates/puppet/type.erb', 'spec/data/type_template.erb'
    file_read_map 'templates/puppet/provider.erb', 'spec/data/prov_template.erb'
    file_read_map 'templates/puppet/provider_spec.erb',
                  'spec/data/prov_spec_template.erb'
    # 6 combinations of success1-3 and name/title
    6.times.each do
      file_read_map 'templates/network_spec.yaml.erb',
                    'spec/data/network_spec_template.erb'
    end

    # Templates will be opened a second time to create the resourceref
    allow_open_spec_templates if rref
  end

  def allow_open_real_templates
    allow_read 'templates/puppet/type.erb'
    allow_read 'templates/puppet/provider.erb'
    allow_read 'templates/puppet/provider_spec.erb'

    allow_open 'templates/expand_variables.erb'
    allow_open 'templates/provider_helpers.erb'
    allow_open 'templates/network_mocks.erb'
    allow_open 'templates/network_spec.yaml.erb'
    allow_open 'templates/return_if_object.erb'
    allow_open 'templates/transport.erb'

    allow_open_real_tests
    allow_open_real_resourceref

    allow_open_license
  end

  def allow_open_real_resourceref
    allow_open 'templates/resourceref_mocks.erb'
    allow_open 'templates/puppet/resourceref_expandvars.erb'
  end

  def allow_open_real_tests
    allow_open 'templates/puppet/test~absent~delete.erb'
    allow_open 'templates/puppet/test~absent~no_action.erb'
    allow_open 'templates/puppet/test~present~create.erb'
    allow_open 'templates/puppet/test~present~no_changes.erb'
  end

  def allow_open_license
    allow_open 'templates/autogen_notice.erb'
    allow_open 'templates/license.erb'
  end

  def allow_open_typed_array
    allow_open 'templates/puppet/property/array_typed.rb.erb'
    allow_read 'templates/puppet/property/array.rb.erb'
  end

  def output_file(name, output)
    File.expects(:open).with(name, 'w').yields(output).at_least(0)
  end

  def dummy_writer
    out = mock('File')
    out.expects(:write).at_least(0)
    out
  end
end
