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

require 'google_compute_instance'
require 'json'

class InstanceTest < Instance
  def initialize(data)
    @fetched = data
  end
end

instance_fixture = JSON.parse(File.read('fixtures/instance_fixture.json'))

RSpec.describe Instance, "#parse" do
  it "compute instance parse" do
    instance_mock = InstanceTest.new(instance_fixture)
    instance_mock.parse
    expect(instance_mock.exists?).to be true
    expect(instance_mock.disks.size).to eq 1
    expect(instance_mock.disks[0].mode).to eq 'READ_WRITE'
    expect(instance_mock.disks[0].auto_delete).to be true
    expect(instance_mock.scheduling.preemptible).to be false
    expect(instance_mock.scheduling.automatic_restart).to be true
    expect(instance_mock.service_accounts.size).to eq 1
    expect(instance_mock.service_accounts[0].email).to eq "577278241961-compute@developer.gserviceaccount.com"
    expect(instance_mock.service_accounts[0].scopes).to include "https://www.googleapis.com/auth/compute"
    
  end
end

RSpec.describe Instance, "none" do
  it "no result" do
    instance_mock = InstanceTest.new(nil)
    expect(instance_mock.exists?).to be false
  end
end