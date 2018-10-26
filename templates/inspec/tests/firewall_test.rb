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

require 'google_compute_firewall'
require 'json'

class FirewallTest < Firewall
  def initialize(data)
    @fetched = data
  end
end

firewall_fixture = JSON.parse(File.read('fixtures/firewall_fixture.json'))

RSpec.describe Firewall, '#parse' do
  before do 
    @firewall_mock = FirewallTest.new(firewall_fixture)
    @firewall_mock.parse
  end
  context 'firewall attributes' do
    it { expect(@firewall_mock.exists?).to be true }
    it { expect(@firewall_mock.creation_timestamp).to eq Time.at(1539322976).to_datetime }
    it { expect(@firewall_mock.description).to eq 'description' }
    it { expect(@firewall_mock.allowed.size).to be 1 }
    it { expect(@firewall_mock.allowed[0].ip_protocol).to eq 'tcp' }
    it { expect(@firewall_mock.allowed[0].ports).to include "443" }
    it { expect(@firewall_mock.denied.size).to be 1 }
    it { expect(@firewall_mock.denied[0].ip_protocol).to eq 'udp' }
    it { expect(@firewall_mock.denied[0].ports).to include "555" }
    it { expect(@firewall_mock.direction).to eq 'INGRESS' }
    it { expect(@firewall_mock.network).to match('/default$') }
    it { expect(@firewall_mock.source_ranges).to include "113.197.104.0/22" }

  end
end


no_firewall = FirewallTest.new(nil)
RSpec.describe Firewall, "#parse" do
  it "does not exist" do
    expect(no_firewall.exists?).to be false
  end
end