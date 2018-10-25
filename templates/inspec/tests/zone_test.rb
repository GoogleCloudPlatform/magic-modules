require 'google_compute_zone'

class ZoneTest < Zone
  def initialize(data)
    @fetched = data
  end
end

zone_fixture = {"kind"=>"compute#zone",
 "id"=>"2231",
 "creationTimestamp"=>"1989-11-28T00:00:00-05:00",
 "name"=>"us-east1-b",
 "description"=>"us-east1-b",
 "status"=>"UP",
 "region"=>
  "https://www.googleapis.com/compute/v1/projects/sam-inspec/regions/us-east1",
 "selfLink"=>
  "https://www.googleapis.com/compute/v1/projects/sam-inspec/zones/us-east1-b",
 "availableCpuPlatforms"=>
  ["Intel Skylake", "Intel Broadwell", "Intel Haswell"]}

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
  it "no response" do
    no_zone_resource = ZoneTest.new(nil)
    expect(no_zone_resource.exists?).to be false
  end
end
