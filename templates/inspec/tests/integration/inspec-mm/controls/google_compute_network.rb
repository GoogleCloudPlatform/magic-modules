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

require 'vcr_config'

title 'Test single GCP compute network'

control 'gcp-compute-network-1.0' do

  impact 1.0
  title 'Ensure GCP compute network has the correct properties.'
  # TODO(slevenick): remove only_if once we generate this again
  only_if { false }
  VCR.use_cassette('gcp-compute-network') do
    resource = google_compute_network({project: attribute('project_name'), name: attribute('network')['name']})
    describe resource do
      it { should exist }

      its ('subnetworks.count') { should eq 1 }
      its ('creation_timestamp') { should be > (Time.now - 365*60*60*24*1).to_datetime }
      its ('routing_config.routing_mode') { should eq attribute('network')['routing_mode'] }
      its ('auto_create_subnetworks'){ should be false }
    end

    subnetwork_name = attribute('subnetwork')['name']
    describe.one do
      resource.subnetworks.each do |subnetwork|
        describe subnetwork do
          # using attribute within this block seems to cause InSpec issues.
          it { should match '/' + subnetwork_name + '$' }
        end
      end
    end
  end
end