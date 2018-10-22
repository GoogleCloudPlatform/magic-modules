require 'google_compute_zone'

RSpec.describe Zone, "zone" do
	it "first test" do
		t = Zone.new
		expect(t.thing).to eq 1

	end
end