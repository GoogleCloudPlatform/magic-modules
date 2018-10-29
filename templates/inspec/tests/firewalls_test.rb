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

require 'google_compute_firewalls'
require 'json'

class FirewallsTest < Firewalls
  def fetch_resource(data)
    return data
  end
end

firewalls_fixture = JSON.parse(File.read('fixtures/firewalls_fixture.json'))

RSpec.describe Firewalls, '#fetch_resource' do
  before do 
    @firewalls_mock = FirewallsTest.new([firewalls_fixture])
  end
  context 'firewalls plural' do
    it { expect(@firewalls_mock.names.size).to eq 3 }
    it { expect(@firewalls_mock.names).to include 'default-knsku4qwwbtr3bhcf3y6vcmu' }
    it { expect(@firewalls_mock.names).to include 'default-7mzjmae3tlidh4yoidvnpe53' }
    it { expect(@firewalls_mock.names).to include 'default-2wnao3jebww7ldrn463stwke' }
  end
end

no_firewalls_fixture = JSON.parse(File.read('fixtures/firewalls_fixture.json'))

no_firewalls_fixture['items'] = []
no_firewalls = FirewallsTest.new([no_firewalls_fixture])
RSpec.describe Firewalls, "#fetch_resource" do
  it "no firewalls" do
    expect(no_firewalls.names.size).to eq 0
  end
end