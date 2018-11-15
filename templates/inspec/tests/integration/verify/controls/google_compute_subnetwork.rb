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

title 'Test Google compute subnetwork resource'

control 'gcp-compute-subnetwork-1.0' do

  impact 1.0
  title 'Ensure GCP compute subnetwork resource works.'
  # TODO(slevenick): remove only_if once we generate this again
  only_if { false }
  VCR.use_cassette('gcp-compute-subnetwork') do
    resource = google_compute_subnetwork({project: attribute('project_name'), region: attribute('region'), name: attribute('subnetwork')['name']})
    describe resource do
      it { should exist }
      its('region') { should match attribute('region') }
      its('creation_timestamp') { should be > (Time.now - 365*60*60*24*1).to_datetime }
      its('ip_cidr_range') { should eq attribute('subnetwork')['ip_range'] }
      its('network') { should match attribute('network')['name'] }
      its('private_ip_google_access') { should be false }
    end
  end
end