# Copyright 2018 Google Inc.
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

require 'rspec'

RSpec.shared_examples 'property_tests' do |_disc_prop, api_prop, tags|
  it 'should exist', property: true, **tags do
    expect(api_prop).to be_truthy
  end
end

RSpec.shared_examples 'resource_tests' do |disc_res, api_res, tags|
  # This test will be skipped if the Discovery Doc doesn't have a kind listed.
  it 'should have kind', skip: !disc_res.schema.dig('properties', 'kind', 'default'),
                         resource: true, **tags do
    expect(disc_res.schema.dig('properties', 'kind', 'default')).to eq(api_res.kind)
  end
end
