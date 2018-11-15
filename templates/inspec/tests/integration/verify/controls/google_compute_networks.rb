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

title 'Test GCP plural compute networks'

control 'gcp-compute-networks-1.0' do

  impact 1.0
  title 'GCP compute networks plural.'
  # TODO(slevenick): remove only_if once we generate this again
  only_if { false }
  VCR.use_cassette('gcp-compute-networks') do
    resource = google_compute_networks(project: attribute('project_name'))
    
    describe resource do
      it { should exist }

      its ('names.size') { should eq 2 }
      its ('names') { should include attribute('network')['name'] }
      its ('names') { should include 'default' }
      
    end
  end
end