require 'spec_helper'
require 'provider/core'

describe 'Provider::Core.expected_output_folder?' do
  # Inputs
  config = 'foo'
  api = Api::Compiler.new(File.read('spec/data/good-file.yaml')).run
  version_name = 'ga'
  start_time = Time.now

  provider = Provider::Core.new(config, api, version_name, start_time)

  it 'should identify `terraform-provider-google` as an expected output folder' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/terraform-provider-google'
    expect(provider.expected_output_folder?(path)).to eq true
  end

  it 'should identify `terraform-provider-google-beta` as an expected output folder' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/terraform-provider-google-beta'
    expect(provider.expected_output_folder?(path)).to eq true
  end

  it 'should identify `terraform-next` as an expected output folder' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/terraform-next'
    expect(provider.expected_output_folder?(path)).to eq true
  end

  it 'should identify `terraform-google-conversion` as an expected output folder' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/terraform-google-conversion'
    expect(provider.expected_output_folder?(path)).to eq true
  end

  it 'should identify `terraform-provider-google-unexpected-suffix` as an UNexpected output folder' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/terraform-provider-google-unexpected-suffix'
    expect(provider.expected_output_folder?(path)).to eq false
  end

  it 'should identify `unexpected-prefix-terraform-provider-google` as an UNexpected output folder' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/unexpected-prefix-terraform-provider-google'
    expect(provider.expected_output_folder?(path)).to eq false
  end
end
