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

require 'google_compute_zones'
require 'json'

class ZonesTest < Zones
  def fetch_resource(data)
    return data
  end
end

zones_fixture = JSON.parse(File.read('fixtures/zones_fixture.json'))

RSpec.describe Zones, "zones" do
  it "plural test" do
    zones_mock_resource = ZonesTest.new([zones_fixture])
    expect(zones_mock_resource.names.size).to eq 3
    expect(zones_mock_resource.names).to include 'us-east1-d'
    expect(zones_mock_resource.names).to include 'us-east1-b'
    expect(zones_mock_resource.names).to include 'us-east1-c'
    expect(zones_mock_resource.statuses).to include 'UP'
    expect(zones_mock_resource.statuses.size).to eq 3
    expect(zones_mock_resource.ids).to include '2231'
    expect(zones_mock_resource.ids).to include '2234'
    expect(zones_mock_resource.ids).to include '2233'
  end
end

no_zones_fixture = JSON.parse(File.read('fixtures/zones_fixture.json'))
no_zones_fixture['items'] = []
no_zones = ZonesTest.new([no_zones_fixture])
RSpec.describe Zones, "#fetch_resource" do
  it "no zones" do
    expect(no_zones.names.size).to eq 0
  end
end