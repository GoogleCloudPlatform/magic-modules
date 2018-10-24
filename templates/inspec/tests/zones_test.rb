require 'google_compute_zones'

class ZonesTest < Zones
  def fetch_resource(data)
    return data
  end
end
    
zones_fixture = {"kind"=>"compute#zoneList",
 "id"=>"projects/sam-inspec/zones",
 "items"=>
  [{"kind"=>"compute#zone",
    "id"=>"2231",
    "creationTimestamp"=>"1969-12-31T16:00:00.000-08:00",
    "name"=>"us-east1-b",
    "description"=>"us-east1-b",
    "status"=>"UP",
    "region"=>
     "https://www.googleapis.com/compute/v1/projects/sam-inspec/regions/us-east1",
    "selfLink"=>
     "https://www.googleapis.com/compute/v1/projects/sam-inspec/zones/us-east1-b",
    "availableCpuPlatforms"=>
     ["Intel Skylake", "Intel Broadwell", "Intel Haswell"]},
   {"kind"=>"compute#zone",
    "id"=>"2233",
    "creationTimestamp"=>"1969-12-31T16:00:00.000-08:00",
    "name"=>"us-east1-c",
    "description"=>"us-east1-c",
    "status"=>"UP",
    "region"=>
     "https://www.googleapis.com/compute/v1/projects/sam-inspec/regions/us-east1",
    "selfLink"=>
     "https://www.googleapis.com/compute/v1/projects/sam-inspec/zones/us-east1-c",
    "availableCpuPlatforms"=>
     ["Intel Skylake", "Intel Broadwell", "Intel Haswell"]},
   {"kind"=>"compute#zone",
    "id"=>"2234",
    "creationTimestamp"=>"1969-12-31T16:00:00.000-08:00",
    "name"=>"us-east1-d",
    "description"=>"us-east1-d",
    "status"=>"UP",
    "region"=>
     "https://www.googleapis.com/compute/v1/projects/sam-inspec/regions/us-east1",
    "selfLink"=>
     "https://www.googleapis.com/compute/v1/projects/sam-inspec/zones/us-east1-d",
    "availableCpuPlatforms"=>
     ["Intel Skylake", "Intel Broadwell", "Intel Haswell"]}]}

RSpec.describe Zones, "zones" do
  it "plural test" do
    zones_mock_resource = ZonesTest.new(zones_fixture)
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