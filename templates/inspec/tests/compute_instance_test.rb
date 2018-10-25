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

class InstanceTest < Instance
  def initialize(data)
    @fetched = data
  end
end

instance_fixture = {"kind"=>"compute#instance",
 "id"=>"1154794430415066980",
 "creationTimestamp"=>"2018-10-24T14:52:15.794-07:00",
 "name"=>"gcp-inspec-app-mig3-4pp8",
 "tags"=>
  {"items"=>["allow-gcp-inspec-app-mig3", "allow-ssh"],
   "fingerprint"=>"rOcaWmVHbAQ="},
 "machineType"=>
  "https://www.googleapis.com/compute/v1/projects/sam-inspec/zones/europe-west2-c/machineTypes/f1-micro",
 "status"=>"RUNNING",
 "zone"=>
  "https://www.googleapis.com/compute/v1/projects/sam-inspec/zones/europe-west2-c",
 "networkInterfaces"=>
  [{"kind"=>"compute#networkInterface",
    "network"=>
     "https://www.googleapis.com/compute/v1/projects/sam-inspec/global/networks/default",
    "subnetwork"=>
     "https://www.googleapis.com/compute/v1/projects/sam-inspec/regions/europe-west2/subnetworks/default",
    "networkIP"=>"10.154.0.7",
    "name"=>"nic0",
    "accessConfigs"=>
     [{"kind"=>"compute#accessConfig",
       "type"=>"ONE_TO_ONE_NAT",
       "name"=>"external-nat",
       "natIP"=>"35.242.153.92",
       "networkTier"=>"PREMIUM"}],
    "fingerprint"=>"gqa1nAlsW2g="}],
 "disks"=>
  [{"kind"=>"compute#attachedDisk",
    "type"=>"PERSISTENT",
    "mode"=>"READ_WRITE",
    "source"=>
     "https://www.googleapis.com/compute/v1/projects/sam-inspec/zones/europe-west2-c/disks/gcp-inspec-app-mig3-4pp8",
    "deviceName"=>"persistent-disk-0",
    "index"=>0,
    "boot"=>true,
    "autoDelete"=>true,
    "licenses"=>
     ["https://www.googleapis.com/compute/v1/projects/debian-cloud/global/licenses/debian-9-stretch"],
    "interface"=>"SCSI",
    "guestOsFeatures"=>[{"type"=>"VIRTIO_SCSI_MULTIQUEUE"}]}],
 "metadata"=>
  {"kind"=>"compute#metadata",
   "fingerprint"=>"8sKIPG5hwJ8=",
   "items"=>
    [{"key"=>"instance-template",
      "value"=>
       "projects/577278241961/global/instanceTemplates/default-20181011161843755000000001"},
     {"key"=>"created-by",
      "value"=>
       "projects/577278241961/zones/europe-west2-c/instanceGroupManagers/gcp-inspec-app-mig3"},
     {"key"=>"tf_depends_id", "value"=>""},
     {"key"=>"startup-script",
      "value"=>"val"}]},
 "serviceAccounts"=>
  [{"email"=>"577278241961-compute@developer.gserviceaccount.com",
    "scopes"=>
     ["https://www.googleapis.com/auth/devstorage.full_control",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/monitoring.write"]}],
 "selfLink"=>
  "https://www.googleapis.com/compute/v1/projects/sam-inspec/zones/europe-west2-c/instances/gcp-inspec-app-mig3-4pp8",
 "scheduling"=>
  {"onHostMaintenance"=>"MIGRATE",
   "automaticRestart"=>true,
   "preemptible"=>false},
 "cpuPlatform"=>"Intel Broadwell",
 "labelFingerprint"=>"42WmSpB8rSM=",
 "startRestricted"=>false,
 "deletionProtection"=>false}

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