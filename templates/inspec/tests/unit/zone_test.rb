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

require 'google_compute_zone'
require 'json'

class ZoneTest < Zone
  def initialize(data)
    @fetched = data
  end
end

zone_fixture = JSON.parse(File.read('fixtures/zone_fixture.json'))

RSpec.describe Zone, "parse" do
  it "zone attributes" do
    zone_mock = ZoneTest.new(zone_fixture)
    zone_mock.parse
    expect(zone_mock.exists?).to be true
    expect(zone_mock.name).to eq 'us-east1-b'
    expect(zone_mock.status).to eq 'UP'
    expect(zone_mock.deprecated.obsolete).to eq nil
    time = Time.at(628232400).to_datetime
    expect(zone_mock.creation_timestamp).to eq time
  end

end

RSpec.describe Zone, "#parse" do
  it "no result" do
    no_zone_resource = ZoneTest.new(nil)
    expect(no_zone_resource.exists?).to be false
  end
end
