# Copyright 2023 Google Inc.
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

describe 'Provider::Core.expected_output_folder?' do
  # Inputs for tests
  api = Api::Compiler.new(File.read('spec/data/good-file.yaml')).run
  version_name = 'ga'
  start_time = Time.now

  provider = Provider::Terraform.new(api, version_name, start_time)

  # rubocop:disable Layout/LineLength
  it 'should identify `terraform-provider-google` as an expected output folder' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/terraform-provider-google'
    expect(provider.expected_output_folder?(path)).to eq true
  end

  it 'should identify `terraform-provider-google-beta` as an expected output folder' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/terraform-provider-google-beta'
    expect(provider.expected_output_folder?(path)).to eq true
  end

  it 'should identify `terraform-next` as an expected output folder' do
    path = '/User/PersonsName/go/src/github.com/GoogleCloudPlatform/terraform-next'
    expect(provider.expected_output_folder?(path)).to eq true
  end

  it 'should identify `terraform-google-conversion` as an expected output folder' do
    path = '/User/PersonsName/go/src/github.com/GoogleCloudPlatform/terraform-google-conversion'
    expect(provider.expected_output_folder?(path)).to eq true
  end

  it 'should identify suffixed versions of expected folder names as unexpected' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/terraform-provider-google-unexpected-suffix'
    expect(provider.expected_output_folder?(path)).to eq false
  end

  it 'should identify prefixed versions of expected folder names as unexpected' do
    path = '/User/PersonsName/go/src/github.com/hashicorp/unexpected-prefix-terraform-provider-google'
    expect(provider.expected_output_folder?(path)).to eq false
  end
  # rubocop:enable Layout/LineLength
end
