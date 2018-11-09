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

require 'google_compute_instances'
require 'json'

class InstancesTest < Instances
  def fetch_resource(data)
    return data
  end
end

instances_fixture = JSON.parse(File.read('fixtures/instances_fixture.json'))

instances_mock = InstancesTest.new([instances_fixture])
RSpec.describe Instances, '#fetch_resource' do
  context 'instances plural' do
    it { expect(instances_mock.names.size).to eq 5 }
    it { expect(instances_mock.ids).to include '4819361437611903243' }
    it { expect(instances_mock.names).to include 'gcp-inspec-ext-linux-vm' }
  end
end


no_instances_fixture = JSON.parse(File.read('fixtures/instances_fixture.json'))
no_instances_fixture['items'] = []

RSpec.describe Instance, "none" do
  it "no result" do
    no_instances_mock = InstancesTest.new([no_instances_fixture])
    expect(no_instances_mock.names.size).to eq 0
  end
end