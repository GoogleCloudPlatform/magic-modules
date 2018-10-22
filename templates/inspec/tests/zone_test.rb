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

RSpec.describe Zone, "zone resource" do
	it "parse test" do
		t = ZoneTest.new(zone_fixture)
		t.parse
		expect(t.exists?).to be true
		expect(t.name).to eq 'us-east1-b'
		expect(t.status).to eq 'UP'
		expect(t.deprecated.obsolete).to eq nil
		time = Time.at(628232400).to_datetime
		expect(t.creation_timestamp).to eq time
	end
	it "no response" do
		t = ZoneTest.new(nil)
		expect(t.exists?).to be false
	end
end

